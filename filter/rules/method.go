package rules

import (
  "fmt"

  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/types"

  "github.com/zubairhamed/canopus"
)

type MethodRule struct {
  AllowedMethods []string
}

func (r *MethodRule) Process(message *types.COAPMessage) types.RuleProcessingResult {
  methodString := canopus.MethodString(message.Message.GetCode())

  for _, method := range r.AllowedMethods {
    if(method == methodString) {
      return types.RuleProcessingResult{Allowed: true}
    }
  }

  return types.RuleProcessingResult{
    false,
    fmt.Sprintf("%s is not allowed: %v", methodString, r.AllowedMethods),
  }
}
