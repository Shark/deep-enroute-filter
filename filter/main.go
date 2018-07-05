package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/AkihiroSuda/go-netfilter-queue"
	"github.com/google/gopacket"

	"gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/network"
	"gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/parser"
	"gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/pipeline"
	"gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/types"
)

func main() {
	var err error

	fd, err:= syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, syscall.ETH_P_ALL)
	if (err != nil) {
		fmt.Println("Error: " + err.Error())
		return;
	}
	fmt.Println("Obtained fd ", fd)
	defer syscall.Close(fd)

	nfq, err := netfilter.NewNFQueue(0, 100, netfilter.NF_DEFAULT_PACKET_SIZE)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer nfq.Close()
	packets := nfq.GetPackets()

	incomingMessages := make(chan *types.COAPMessage, 10)
	outgoingPackets := make(chan gopacket.Packet, 10)
	whitelistedMessageHashes := make(map[string]bool)

	go func() {
		pipeline.Consume(incomingMessages, outgoingPackets, &whitelistedMessageHashes)
	}()

	go func() {
		network.ReinjectPackets(outgoingPackets, fd)
	}()

	for true {
		select {
		case p := <-packets:
			verdict := netfilter.NF_DROP

			message, err := parser.ParseCOAPMessageFromPacket(p.Packet)
			if(err != nil) {
				fmt.Println("Error parsing packet: %v", err)
				continue
			}

			incomingMessages <- message

			if whitelistedMessageHashes[message.Metadata.Hash()] {
				verdict = netfilter.NF_ACCEPT
			}

			p.SetVerdict(verdict)
		}
	}
}
