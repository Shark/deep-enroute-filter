package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"syscall"

	"github.com/AkihiroSuda/go-netfilter-queue"

	"gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/network"
	"gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/parser"
	"gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/pipeline"
	"gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/types"
	"gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/web"
)

func generateAuthenticityToken() (*string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	token := hex.EncodeToString(bytes)
	return &token, nil
}

func main() {
	var err error
	authenticityToken, err := generateAuthenticityToken()
	if err != nil {
		fmt.Printf("Error generating authenticityToken: %v", err)
		return
	}

	fd, err:= syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, syscall.ETH_P_ALL)
	if (err != nil) {
		fmt.Println("Error: " + err.Error())
		return
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
	processedMessages := make(chan types.ProcessedMessage, 10)
	outgoingMessages := make(chan *types.COAPMessage, 10)

	go func() {
		pipeline.Consume(incomingMessages, processedMessages, outgoingMessages, *authenticityToken)
	}()

	go func() {
		network.ReinjectPackets(outgoingMessages, fd)
	}()

	go func() {
		web.ListenAndServe(processedMessages)
	}()

	for true {
		select {
		case p := <-packets:
			verdict := netfilter.NF_DROP
			var packet *[]byte

			message, err := parser.ParseCOAPMessageFromPacket(p.Packet)
			if(err != nil) {
				fmt.Printf("Error parsing packet: %v\n", err)
				continue
			}

			if message.Metadata.AuthenticityToken != nil && *message.Metadata.AuthenticityToken == *authenticityToken {
				message.Message.RemoveOptions(65000)
				packet, err = network.SerializeMessage(message, false)
				if err != nil {
					fmt.Printf("Error serializing packet: %v\n", err)
				} else {
					verdict = netfilter.NF_ACCEPT
				}
			} else {
				incomingMessages <- message
			}

			if packet != nil {
				p.SetVerdictWithPacket(verdict, *packet)
			} else {
				p.SetVerdict(verdict)
			}
		}
	}
}
