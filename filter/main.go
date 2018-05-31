package main

import (
	"fmt"
	"os"
  "net"
  "syscall"
  "github.com/AkihiroSuda/go-netfilter-queue"
  "github.com/google/gopacket"
	"github.com/google/gopacket/layers"
  "github.com/zubairhamed/canopus"
)

func main() {
	var err error

  fd, err:= syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW,
  				syscall.ETH_P_ALL)
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

	for true {
		select {
		case p := <-packets:
      verdict := netfilter.NF_DROP

      ipv6 := p.Packet.Layer(layers.LayerTypeIPv6).(*layers.IPv6)
      if udpLayer := p.Packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
        udp, _ := udpLayer.(*layers.UDP)
        if(udp.DstPort == 5683) { // CoAP packet
          fmt.Println("This is a COAP packet!")
          msg, err := canopus.BytesToMessage(udp.LayerPayload())
          coapMsg := msg.(*canopus.CoapMessage)

          if(err != nil) {
            fmt.Printf("Error parsing CoAP: %v\n", err)
          } else {
            // check if message has been filtered
            checked := false
            option := coapMsg.GetOption(35)
            if(option != nil) {
              if value, ok := option.GetValue().(string); ok {
                if(value == "checked") {
                  checked = true
                }
              }
            }

            if(checked) {
              fmt.Println("Allowed packet");
              verdict = netfilter.NF_ACCEPT
            } else {
              coapMsg.AddOption(35, "checked")
              canopus.PrintMessage(msg)

              buf := gopacket.NewSerializeBuffer()
              opts := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}

              msgBytes, err := canopus.MessageToBytes(coapMsg)
              if(err != nil) {
                fmt.Printf("Error serializing message: %v", err)
              }
              bytes, _ := buf.PrependBytes(len(msgBytes))
              copy(bytes, msgBytes)

              udp.SetNetworkLayerForChecksum(ipv6)
              udp.SerializeTo(buf, opts)
              ipv6.SerializeTo(buf, opts)

              iface, err := net.InterfaceByName("lo")
              if(err != nil) {
                fmt.Println("Did not find iface")
              }
              lb, _ := net.ParseMAC("00:00:00:00:00:00")
              frame := layers.Ethernet{
                SrcMAC: lb,
                DstMAC: lb,
                EthernetType: 0x86DD, // IPv6
              }
              err = frame.SerializeTo(buf, opts)
              if(err != nil) {
                fmt.Printf("Error serialize eth: %v\n", err)
              }

              var addr syscall.SockaddrLinklayer
              addr.Protocol = syscall.ETH_P_IPV6
              addr.Ifindex = iface.Index
              addr.Hatype = syscall.ARPHRD_LOOPBACK

              decoded := gopacket.NewPacket(buf.Bytes(), layers.LayerTypeEthernet, gopacket.Default)
              for _, layer := range decoded.Layers() {
    fmt.Println("PACKET LAYER:", layer.LayerType())
  }

              // Send the packet
              err = syscall.Sendto(fd, buf.Bytes(), 0, &addr)

              if(err != nil) {
                fmt.Printf("Error sending ethernet packet: %v", err)
              }
            }
          }
        }
      }

      fmt.Println("verdict")
			p.SetVerdict(verdict)
		}
	}
}
