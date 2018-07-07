package parser

import (
  "encoding/hex"
	"testing"

  "github.com/google/gopacket"
  "github.com/google/gopacket/layers"
)


func TestParseCOAPMessageFromPacket1(t *testing.T) {
  // valid IPv6-UDP-COAP packet without AuthenticityToken
  packetBytes, err := hex.DecodeString("600618270021113ffd1b2211c9b6000054a78239fdb8f6c7fd00000000000000000000fffe0053c08a64163300217b1e44019df9e91f89b5bb2e77656c6c2d6b6e6f776e04636f7265")
	if err != nil {
		t.Error("Byte decoding failed")
	}
  packet := gopacket.NewPacket(packetBytes, layers.LayerTypeIPv6, gopacket.Default)

  _, err = ParseCOAPMessageFromPacket(packet)

  if(err != nil) {
    t.Errorf("ParseCOAPMessageFromPacket failed: %v", err)
  }
}

func TestParseCOAPMessageFromPacket2(t *testing.T) {
  // valid IPv6-UDP-COAP packet with AuthenticityToken
  packetBytes, err := hex.DecodeString("600add9c00361140fdbacd7e4c8ea6fd0000000000000002fdf0a23f8cae5b9700000000000000029c4616330036ae8f40016838b56261736963edfcd0133231366366616261353464323966303366663465333161393831393363613862")
	if err != nil {
		t.Error("Byte decoding failed")
	}
  packet := gopacket.NewPacket(packetBytes, layers.LayerTypeIPv6, gopacket.Default)

  message, err := ParseCOAPMessageFromPacket(packet)

  if(err != nil) {
    t.Errorf("ParseCOAPMessageFromPacket failed: %v", err)
  }

  if(message.Metadata.AuthenticityToken == nil) {
    t.Errorf("ParseCOAPMessageFromPacket: expected AuthenticityToken to be present")
  }
}

func TestParseCOAPMessageFromPacket3(t *testing.T) {
  // valid IPv6-UDP-COAP packet
  packetBytes, err := hex.DecodeString("600551e00020113ffd1b2211c9b6000054a78239fdb8f6c7fd00000000000000464e6dfffe21b1a1ee0b00350020fb62f1d201000001000000000000027831036c616e0000010001")
	if err != nil {
		t.Error("Byte decoding failed")
	}
  packet := gopacket.NewPacket(packetBytes, layers.LayerTypeIPv6, gopacket.Default)

  _, err = ParseCOAPMessageFromPacket(packet)

  if(err == nil || err.Error() != "Packet is not a COAP message") {
    t.Error("ParseCOAPMessageFromPacket expected to fail, but did not")
  }
}

func TestParseCOAPMessageFromPacket4(t *testing.T) {
  // invalid IPv6-UDP-COAP packet
  packetBytes, err := hex.DecodeString("600618270021113ffd1b2211c9b6000054a78239fdb8f6c7fd00000000000000000000fffe0053c08a64163300217b1e40019df9e91f89b5bb2e77656c6c2d6b6e6f7704636f7265")
	if err != nil {
		t.Error("Byte decoding failed")
	}
  packet := gopacket.NewPacket(packetBytes, layers.LayerTypeIPv6, gopacket.Default)

  _, err = ParseCOAPMessageFromPacket(packet)

  if(err == nil || err.Error() != "Failed to parse COAP message: unspecified error") {
    t.Error("ParseCOAPMessageFromPacket expected to fail, but did not")
  }
}
