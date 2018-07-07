package pipeline

import (
  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/rules"
  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/types"
)

func Consume(incomingMessages <-chan *types.COAPMessage, processedMessages chan<- types.ProcessedMessage, outgoingMessages chan<- *types.COAPMessage, authenticityToken string) {
  for message := range incomingMessages {
    rule := rules.MethodRule{AllowedMethods: []string{"GET"}}
    result := rule.Process(message)

    processedMessages <- types.ProcessedMessage{message, []types.RuleProcessingResult{result}}

    if result.Allowed {
      message.Message.AddOption(65000, authenticityToken)
      outgoingMessages <- message
    }
  }
}
