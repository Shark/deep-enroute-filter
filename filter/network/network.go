package network

import (
  "errors"
  "fmt"
  "net"
  "syscall"

  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/types"

  "github.com/zubairhamed/canopus"
  "github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func ReinjectPackets(outgoingMessages <-chan *types.COAPMessage, rawSocketFd int) {
  for packet := range outgoingMessages {
    err := ReinjectPacket(rawSocketFd, packet)
    if err != nil {
      fmt.Printf("Error reinjecting packet: %v", err)
    }
  }
}

func ReinjectPacket(rawSocketFd int, message *types.COAPMessage) error {
  buf := gopacket.NewSerializeBuffer()
  opts := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}

  messageBytes, err := canopus.MessageToBytes(message.Message)
  if err != nil {
    return fmt.Errorf("Error serializing COAP message: %v", err)
  }
  bufBytes, err := buf.PrependBytes(len(messageBytes))
  if err != nil {
    return fmt.Errorf("Error prepending bytes to buffer: %v", err)
  }
  copy(bufBytes, messageBytes)

  if udpLayer, ok := message.TransportLayer.(*layers.UDP); ok {
    udpLayer.SetNetworkLayerForChecksum(message.NetworkLayer)
    err = udpLayer.SerializeTo(buf, opts)
    if err != nil {
      return fmt.Errorf("Error serializing TransportLayer: %v", err)
    }
  } else {
    return errors.New("Can not use TransportLayer as UDP layer")
  }

  if ipv6Layer, ok := message.NetworkLayer.(*layers.IPv6); ok {
    err = ipv6Layer.SerializeTo(buf, opts)
    if err != nil {
      return fmt.Errorf("Error serializing NetworkLayer: %v", err)
    }
  } else {
    return errors.New("Can not use NetworkLayer as IPv6 layer")
  }

  loopbackIface, err := net.InterfaceByName("lo")
  if(err != nil) {
    return errors.New("unable to find loopback interface")
  }
  nullMAC, _ := net.ParseMAC("00:00:00:00:00:00")
  ethernetFrame := layers.Ethernet{
    SrcMAC: nullMAC,
    DstMAC: nullMAC,
    EthernetType: 0x86DD, // IPv6
  }
  err = ethernetFrame.SerializeTo(buf, opts)
  if(err != nil) {
    return err
  }

  var addr syscall.SockaddrLinklayer
  addr.Protocol = syscall.ETH_P_IPV6
  addr.Ifindex = loopbackIface.Index
  addr.Hatype = syscall.ARPHRD_LOOPBACK

  // Send the packet
  err = syscall.Sendto(rawSocketFd, buf.Bytes(), 0, &addr)

  if(err != nil) {
    return fmt.Errorf("Error sending ethernet packet: %v", err)
  }
  return nil
}
