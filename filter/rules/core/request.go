package core

import (
  "fmt"

  "github.com/zubairhamed/canopus"
)

func fetchCore(dstIP string) (*[]string, error) {
  conn, err := canopus.Dial(fmt.Sprintf("[%s]:5683", dstIP))
  if err != nil {
    return nil, err
  }

  req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Get).(*canopus.CoapRequest)
  req.SetRequestURI("/.well-known/core")

  resp, err := conn.Send(req)
  if err != nil {
    return nil, err
  }

  payload := resp.GetMessage().GetPayload().String()
  definition := parseDefinition(payload)

  return &definition, nil
}
