package pipeline

import (
  "sync"

  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/rules"
  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/types"

  "github.com/google/gopacket"
)

func Consume(incomingMessages <-chan *types.COAPMessage, processedMessages chan<- types.ProcessedMessage, outgoingPackets chan<- gopacket.Packet, whitelistedMessageHashes map[string]bool, whitelistedMessagesHashesMutex sync.RWMutex) {
  for message := range incomingMessages {
    packetHash := message.Metadata.Hash()

    rule := rules.MethodRule{AllowedMethods: []string{"GET"}}
    result := rule.Process(message)

    if result.Allowed {
      whitelistedMessagesHashesMutex.Lock()
      whitelistedMessageHashes[packetHash] = true
      whitelistedMessagesHashesMutex.Unlock()
    }

    processedMessages <- types.ProcessedMessage{message, []types.RuleProcessingResult{result}}

    outgoingPackets <- message.Packet
  }
}
