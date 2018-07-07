package pipeline

import (
  "sync"

  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/rules"
  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/types"
)

func Consume(incomingMessages <-chan *types.COAPMessage, processedMessages chan<- types.ProcessedMessage, outgoingMessages chan<- *types.COAPMessage, whitelistedMessageHashes map[string]bool, whitelistedMessagesHashesMutex sync.RWMutex) {
  for message := range incomingMessages {
    packetHash := message.Metadata.Hash()

    rule := rules.MethodRule{AllowedMethods: []string{"GET"}}
    result := rule.Process(message)

    whitelistedMessagesHashesMutex.Lock()
    whitelistedMessageHashes[packetHash] = result.Allowed
    whitelistedMessagesHashesMutex.Unlock()

    processedMessages <- types.ProcessedMessage{message, []types.RuleProcessingResult{result}}

    outgoingMessages <- message
  }
}
