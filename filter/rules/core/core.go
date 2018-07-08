package core

import (
  "fmt"
  "strings"

  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/types"
)

type CoreRule struct {
  endpoints map[string][]string
}

func (c CoreRule) Process(message *types.COAPMessage) types.RuleProcessingResult {
  if message.Metadata.UriPath == "/.well-known/core" {
    return types.RuleProcessingResult{true, c, nil}
  }

  var endpoints []string
  if cachedEndpoints, ok := c.endpoints[message.Metadata.DstIP]; ok {
    endpoints = cachedEndpoints
  } else {
    fetchedEndpoints, err := fetchCore(message.Metadata.DstIP)
    if(err != nil) {
      fmt.Printf("Error fetching core: %v", err)
      return types.RuleProcessingResult{true, c, nil}
    }
    c.endpoints[message.Metadata.DstIP] = *fetchedEndpoints
    endpoints = *fetchedEndpoints
  }

  allowed := false
  for _, endpoint := range endpoints {
    if strings.HasPrefix(message.Metadata.UriPath, endpoint) {
      allowed = true
    }
  }

  if allowed {
    return types.RuleProcessingResult{
      true,
      c,
      nil,
    }
  } else {
    message := fmt.Sprintf("%s is not included in .well-known/core: %q", message.Metadata.UriPath, endpoints)
    return types.RuleProcessingResult{
      false,
      c,
      &message,
    }
  }
}

func (r CoreRule) Name() string {
  return "CoreRule"
}

func NewCoreRule() CoreRule {
  return CoreRule{make(map[string][]string)}
}
