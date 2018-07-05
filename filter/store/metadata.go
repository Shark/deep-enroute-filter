package store

import (
  "crypto/sha256"
  "encoding/hex"
  "errors"
  "fmt"

  "github.com/google/gopacket"
  "github.com/google/gopacket/layers"
  "github.com/zubairhamed/canopus"
)

type CoapPacketMetadata struct {
  srcIP        string
  dstIP        string
  srcPort      int
  dstPort      int
  coapMsgToken string
}

func (m *CoapPacketMetadata) Hash() string {
  str := fmt.Sprintf("%s%s%d%d%s", m.srcIP, m.dstIP, m.srcPort, m.dstPort, m.coapMsgToken)
  sum := sha256.Sum256([]byte(str))
  return fmt.Sprintf("%x", sum)
}

func ExtractCOAPMetadataFromPacket(packet gopacket.Packet) (metadata *CoapPacketMetadata, err error) {
  inIPv6 := packet.Layer(layers.LayerTypeIPv6)
  inUDP := packet.Layer(layers.LayerTypeUDP)

  if(inIPv6 != nil && inUDP != nil) {
    ipv6Layer := inIPv6.(*layers.IPv6)
    udpLayer := inUDP.(*layers.UDP)

    srcIP := ipv6Layer.SrcIP.String()
    dstIP := ipv6Layer.DstIP.String()
    srcPort := int(udpLayer.SrcPort)
    dstPort := int(udpLayer.DstPort)

    if(dstPort != 5683) {
      return nil, errors.New("Packet is not a COAP message")
    }

    coapMsg, err := canopus.BytesToMessage(udpLayer.LayerPayload())
    if(err != nil) {
      return nil, errors.New("Failed to parse COAP message")
    }

    coapMsgToken := hex.EncodeToString(coapMsg.GetToken())

    return &CoapPacketMetadata{
      srcIP,
      dstIP,
      srcPort,
      dstPort,
      coapMsgToken,
    }, nil
  } else {
    return nil, errors.New("Packet does not have IPv6 or UDP layer")
  }
}
