package core

import (
  "fmt"
  "strings"
  "time"

  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/types"
)

type coreEndpointsEvent struct {
  endpoints map[string]int
}

func (e *coreEndpointsEvent) Type() string {
  return "CoreEndpointsEvent"
}

func (e *coreEndpointsEvent) Payload() interface{} {
  return e.endpoints
}

type CoreRule struct {
  endpoints map[string][]string
  events chan types.Event
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

func (c CoreRule) publishState() {
  endpointsForEvent := make(map[string]int)
  for dstIP, endpoints := range c.endpoints {
    endpointsForEvent[dstIP] = len(endpoints)
  }
  c.events <- &coreEndpointsEvent{endpointsForEvent}
}

func (r CoreRule) Name() string {
  return "CoreRule"
}

func NewCoreRule(events chan types.Event) *CoreRule {
  coreRule := CoreRule{
    endpoints: make(map[string][]string),
    events: events,
  }

  go func() {
    ticker := time.NewTicker(2 * time.Second)
    for _ = range ticker.C {
      coreRule.publishState()
    }
  }()

  return &coreRule
}
