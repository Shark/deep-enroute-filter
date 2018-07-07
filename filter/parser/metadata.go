package parser

import (
  "encoding/hex"
  "errors"

  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/types"

  "github.com/google/gopacket/layers"
  "github.com/zubairhamed/canopus"
)

func extractCOAPMetadata(ipv6Layer *layers.IPv6, udpLayer *layers.UDP, message canopus.Message) (metadata *types.COAPMessageMetadata, err error) {
  srcIP := ipv6Layer.SrcIP.String()
  dstIP := ipv6Layer.DstIP.String()
  srcPort := int(udpLayer.SrcPort)
  dstPort := int(udpLayer.DstPort)

  if(dstPort != 5683) {
    return nil, errors.New("Packet is not a COAP message")
  }

  coapMsgToken := hex.EncodeToString(message.GetToken())

  authTokenOption := message.GetOption(65000)
  var authTokenValue *string
  if(authTokenOption != nil) {
    if byteValue, ok := authTokenOption.GetValue().([]byte); ok {
      stringValue := string(byteValue)
      authTokenValue = &stringValue
    }
  }

  return &types.COAPMessageMetadata{
    srcIP,
    dstIP,
    srcPort,
    dstPort,
    coapMsgToken,
    authTokenValue,
  }, nil
}
