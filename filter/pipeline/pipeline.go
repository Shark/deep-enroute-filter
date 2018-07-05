package pipeline

import (
  "sync"

  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/types"

  "github.com/google/gopacket"
)

func Consume(incomingMessages <-chan *types.COAPMessage, outgoingPackets chan<- gopacket.Packet, whitelistedMessageHashes map[string]bool, whitelistedMessagesHashesMutex sync.RWMutex) {
  for message := range incomingMessages {
    packetHash := message.Metadata.Hash()

    whitelistedMessagesHashesMutex.Lock()
    whitelistedMessageHashes[packetHash] = true
    whitelistedMessagesHashesMutex.Unlock()

    outgoingPackets <- message.Packet
  }
}
