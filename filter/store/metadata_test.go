package store

import (
  "encoding/hex"
  "reflect"
	"testing"

  "github.com/google/gopacket"
  "github.com/google/gopacket/layers"
)


func TestExtractCOAPMetadataFromPacket(t *testing.T) {
  packetBytes, err := hex.DecodeString("600618270021113ffd1b2211c9b6000054a78239fdb8f6c7fd00000000000000000000fffe0053c08a64163300217b1e44019df9e91f89b5bb2e77656c6c2d6b6e6f776e04636f7265")
	if err != nil {
		t.Error("Byte decoding failed")
	}
  packet := gopacket.NewPacket(packetBytes, layers.LayerTypeIPv6, gopacket.Default)

  metadata, err := ExtractCOAPMetadataFromPacket(packet)

  if(err != nil) {
    t.Errorf("GenerateHashForPacket failed: %v", err)
  }

  expected := &CoapPackageMetadata{
    "fd1b:2211:c9b6:0:54a7:8239:fdb8:f6c7",
    "fd00::ff:fe00:53c0",
    35428,
    5683,
    "e91f89b5",
  }

  if(!reflect.DeepEqual(metadata, expected)) {
    t.Errorf("Expected metadata: %v, actual: %v", expected, metadata)
  }
}

func TestHashCOAPMetadata(t *testing.T) {
  metadata := &CoapPackageMetadata{
    "fd1b:2211:c9b6:0:54a7:8239:fdb8:f6c7",
    "fd00::ff:fe00:53c0",
    35428,
    5683,
    "e91f89b5",
  }

  hash := metadata.Hash()
  expected := "2df46c62c553821b4a705848c5aba873c9407df59f3095c2e2d83bd1398c1a77"

  if(hash != expected) {
    t.Errorf("Expected hash: %s != %s", expected, hash)
  }
}
