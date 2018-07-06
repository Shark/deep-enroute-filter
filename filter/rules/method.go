package rules

import (
  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/types"

  "github.com/zubairhamed/canopus"
)

type MethodRule struct {
  AllowedMethods []string
}

func (r *MethodRule) Process(message *types.COAPMessage) bool {
  methodString := canopus.MethodString(message.Message.GetCode())

  for _, method := range r.AllowedMethods {
    if(method == methodString) {
      return true
    }
  }

  return false
}
