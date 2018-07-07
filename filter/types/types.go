package types

import (
  "fmt"
  "crypto/sha256"
  "github.com/zubairhamed/canopus"
  "github.com/google/gopacket"
)

type COAPMessage struct {
  Metadata COAPMessageMetadata
  NetworkLayer gopacket.NetworkLayer
  TransportLayer gopacket.TransportLayer
  Message canopus.Message
}

type COAPMessageMetadata struct {
  SrcIP           string
  DstIP           string
  SrcPort         int
  DstPort         int
  CoapMsgToken    string
  AuthOptionValue *string
}

func (m *COAPMessageMetadata) Hash() string {
  str := fmt.Sprintf("%s%s%d%d%s", m.SrcIP, m.DstIP, m.SrcPort, m.DstPort, m.CoapMsgToken)
  sum := sha256.Sum256([]byte(str))
  return fmt.Sprintf("%x", sum)
}

type RuleProcessingResult struct {
  Allowed bool
  Rule Rule
  RuleMessage *string
}

type Rule interface {
  Process(message *COAPMessage) RuleProcessingResult
  Name() string
}

type ProcessedMessage struct {
  Message *COAPMessage
  RuleProcessingResults []RuleProcessingResult
}
