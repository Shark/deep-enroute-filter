package parser

import (
  "encoding/hex"
  "errors"

  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/types"

  "github.com/google/gopacket"
  "github.com/google/gopacket/layers"
  "github.com/zubairhamed/canopus"
)

func extractCOAPMetadata(packet gopacket.Packet, message canopus.Message) (metadata *types.COAPMessageMetadata, err error) {
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

    coapMsgToken := hex.EncodeToString(message.GetToken())

    return &types.COAPMessageMetadata{
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
