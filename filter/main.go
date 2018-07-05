package main

import (
	"fmt"
	"os"
	"syscall"
	"github.com/AkihiroSuda/go-netfilter-queue"
	"github.com/zubairhamed/canopus"

	"gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/parser"
	"gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/network"
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

	allowedTokens := make(map[string]bool)

	for true {
		select {
		case p := <-packets:
			verdict := netfilter.NF_DROP

			message, err := parser.ParseCOAPMessageFromPacket(p.Packet)
			if(err != nil) {
				fmt.Println("Error parsing packet: %v", err)
				continue
			}

			packetHash := message.Metadata.Hash()

			// check if message has been filtered
			checked := allowedTokens[packetHash] == true;

			if(checked) {
				fmt.Println("Allowed packet");
				verdict = netfilter.NF_ACCEPT
			} else {
				canopus.PrintMessage(message.Message)

				allowedTokens[packetHash] = true
				err := network.ReinjectPacket(fd, p.Packet)
				if(err != nil) {
					fmt.Printf("Error reinjecting packet: %v", err)
				}
			}

			p.SetVerdict(verdict)
		}
	}
}
