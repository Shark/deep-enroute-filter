package pipeline

import (
  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/rules"
  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/rules/core"
  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/rules/ratelimit"
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
  methodRule := rules.MethodRule{AllowedMethods: []string{"GET"}}
  coreRule := core.NewCoreRule(events)
  rateLimitRule := ratelimit.NewRateLimit(events)

  for message := range incomingMessages {
    methodRuleResult := methodRule.Process(message)
    coreRuleResult := coreRule.Process(message)
    rateLimitResult := rateLimitRule.Process(message)

    events <- &ProcessedMessageEvent{
      types.ProcessedMessage{message, []types.RuleProcessingResult{methodRuleResult, coreRuleResult, rateLimitResult}},
    }

    if methodRuleResult.Allowed && coreRuleResult.Allowed && rateLimitResult.Allowed {
      message.Message.AddOption(65000, authenticityToken)
      outgoingMessages <- message
    }
  }
}
