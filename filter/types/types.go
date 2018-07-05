package types

import (
  "fmt"
  "crypto/sha256"
  "github.com/zubairhamed/canopus"
  "github.com/google/gopacket"
)

type COAPMessage struct {
  Metadata COAPMessageMetadata
  Packet gopacket.Packet
  Message canopus.Message
}

type COAPMessageMetadata struct {
  SrcIP        string
  DstIP        string
  SrcPort      int
  DstPort      int
  CoapMsgToken string
}

func (m *COAPMessageMetadata) Hash() string {
  str := fmt.Sprintf("%s%s%d%d%s", m.SrcIP, m.DstIP, m.SrcPort, m.DstPort, m.CoapMsgToken)
  sum := sha256.Sum256([]byte(str))
  return fmt.Sprintf("%x", sum)
}
