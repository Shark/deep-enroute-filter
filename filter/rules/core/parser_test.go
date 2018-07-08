package core

import (
	"testing"
)

func TestParseDefinition(t *testing.T) {
  definition := `
    </.well-known/core>;ct=40,
    </actuators/leds>;title="LEDs: ?color=r|g|b, POST/PUT mode=on|off";rt="Control",
    </sensors/sht21>;title="Temperature and Humidity";rt="Sht21",
    </sensors/max44009>;title="Light";rt="max44009"
  `

  result := parseDefinition(definition)

  if(len(result) != 4) {
    t.Errorf("expected length 4, actual length %d", len(result))
  }

  if(result[0] != "/.well-known/core") {
    t.Errorf("expected /.well-known/core, got: %s", result[0])
  }

  if(result[3] != "/sensors/max44009") {
    t.Errorf("expected /sensors/max44009, got: %s", result[3])
  }
}
