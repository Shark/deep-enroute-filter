package pipeline

import (
  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/types"

  "github.com/google/gopacket"
)

func Consume(incomingMessages <-chan *types.COAPMessage, outgoingPackets chan<- gopacket.Packet, whitelistedMessageHashes *map[string]bool) {
  for message := range incomingMessages {
    packetHash := message.Metadata.Hash()

    (*whitelistedMessageHashes)[packetHash] = true

    outgoingPackets <- message.Packet
  }
}
