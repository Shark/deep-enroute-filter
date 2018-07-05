package parser

import (
  "errors"
  "fmt"

  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/types"

  "github.com/google/gopacket"
  "github.com/google/gopacket/layers"
  "github.com/zubairhamed/canopus"
)

func parsePacketPayloadAsCOAPMessage(packet gopacket.Packet) (message canopus.Message, err error) {
  defer func() {
    if r := recover(); r != nil {
      message = nil
      err = errors.New("unspecified error")
    }
  }()

  if payloadLayer, ok := packet.ApplicationLayer().(*gopacket.Payload); ok {
    return canopus.BytesToMessage(payloadLayer.LayerContents())
  } else {
    return nil, errors.New("Packet does not have a payload layer")
  }
}

func ParseCOAPMessageFromPacket(packet gopacket.Packet) (*types.COAPMessage, error) {
  if udpLayer, ok := packet.Layer(layers.LayerTypeUDP).(*layers.UDP); ok {
    dstPort := int(udpLayer.DstPort)

    if(dstPort != 5683) {
      return nil, errors.New("Packet is not a COAP message")
    }

    coapMsg, err := parsePacketPayloadAsCOAPMessage(packet)
    if(err != nil) {
      return nil, fmt.Errorf("Failed to parse COAP message: %v", err)
    }

    metadata, err := extractCOAPMetadata(packet, coapMsg)
    if(err != nil) {
      return nil, fmt.Errorf("Failed to extract COAP metadata: %v", err)
    }

    return &types.COAPMessage{
      *metadata,
      packet,
      coapMsg,
    }, nil
  } else {
    return nil, errors.New("Packet does not have IPv6 or UDP layer")
  }
}
