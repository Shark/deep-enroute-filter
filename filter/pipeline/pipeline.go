package pipeline

import (
  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/rules"
  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/types"
)

type ProcessedMessageEvent struct {
  processedMessage types.ProcessedMessage
}

func (e *ProcessedMessageEvent) Type() string {
  return "ProcessedMessageEvent"
}

func (e *ProcessedMessageEvent) Payload() interface{} {
  return e.processedMessage
}

func Consume(incomingMessages <-chan *types.COAPMessage, outgoingMessages chan<- *types.COAPMessage, authenticityToken string, events chan types.Event) {
  for message := range incomingMessages {
    rule := rules.MethodRule{AllowedMethods: []string{"GET"}}
    result := rule.Process(message)

    events <- &ProcessedMessageEvent{
      types.ProcessedMessage{message, []types.RuleProcessingResult{result}},
    }

    if result.Allowed {
      message.Message.AddOption(65000, authenticityToken)
      outgoingMessages <- message
    }
  }
}
