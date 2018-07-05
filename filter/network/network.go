package network

import (
  "errors"
  "fmt"
  "net"
  "syscall"

  "github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func ReinjectPackets(outgoingPackets <-chan gopacket.Packet, rawSocketFd int) {
  for packet := range outgoingPackets {
    err := ReinjectPacket(rawSocketFd, packet)
    if err != nil {
      fmt.Printf("Error reinjecting packet: %v", err)
    }
  }
}

func ReinjectPacket(rawSocketFd int, packet gopacket.Packet) error {
  buf := gopacket.NewSerializeBuffer()
  opts := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}

  if payloadLayer, ok := packet.ApplicationLayer().(*gopacket.Payload); ok {
    err := payloadLayer.SerializeTo(buf, opts)
    if(err != nil) {
      return err;
    }
  } else {
    return errors.New("payload layer not found in packet")
  }

  if ipv6Layer, ok := packet.NetworkLayer().(*layers.IPv6); ok {
    if udpLayer, ok := packet.TransportLayer().(*layers.UDP); ok {
      udpLayer.SetNetworkLayerForChecksum(ipv6Layer)
      err := udpLayer.SerializeTo(buf, opts)
      if(err != nil) {
        return err;
      }
    } else {
      return errors.New("UDP layer not found in packet")
    }

    err := ipv6Layer.SerializeTo(buf, opts)
    if(err != nil) {
      return err;
    }
  } else {
    return errors.New("IPv6 layer not found in packet")
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
