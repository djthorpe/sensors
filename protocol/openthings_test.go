package protocol_test

import (
	"encoding/hex"
	"strings"
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/sensors/protocol/openthings"
)

func Test_OT_000(t *testing.T) {
	// Create an OOK module
	if app, err := gopi.NewAppInstance(gopi.NewAppConfig("sensors/protocol/openthings")); err != nil {
		t.Fatal(err)
	} else if _, ok := app.ModuleInstance("sensors/protocol/openthings").(sensors.Proto); ok == false {
		t.Fatal("Does not comply to Proto interface")
	} else if _, ok := app.ModuleInstance("sensors/protocol/openthings").(sensors.OTProto); ok == false {
		t.Fatal("Does not comply to OTProto interface")
	}
}

func Test_OT_001(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else if msg_, err := proto.New(sensors.OT_MANUFACTURER_ENERGENIE, 0, 0); err != nil {
		t.Fatal(err)
	} else if msg, ok := msg_.(sensors.OTMessage); ok == false {
		t.Fatal("Not an OTMessage")
	} else if msg.Manufacturer() != sensors.OT_MANUFACTURER_ENERGENIE {
		t.Error("Unexpected manufacturer")
	} else if msg.Product() != 0 {
		t.Error("Unexpected product")
	} else if msg.Sensor() != 0 {
		t.Error("Unexpected sensor")
	}
}

func Test_OT_002(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		for sensor := uint32(0); sensor < 0xFFFFFF; sensor += 0x1234 {
			product := uint8(sensor % 0x7F)
			if msg, err := proto.New(sensors.OT_MANUFACTURER_ENERGENIE, product, sensor); err != nil {
				t.Error(err)
			} else {
				t.Log(msg)
			}
		}
	}
}
func Test_OT_003(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		for sensor := uint32(0); sensor < 0xFFFFFF; sensor += 0x1234 {
			product := uint8(sensor % 0x7F)
			if msg, err := proto.New(sensors.OT_MANUFACTURER_ENERGENIE, product, sensor); err != nil {
				t.Error(err)
			} else if payload := proto.Encode(msg); len(payload) == 0 {
				t.Error("Empty payload")
			} else {
				t.Log(msg)
				t.Logf("...%v", strings.ToUpper(hex.EncodeToString(payload)))
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// OT

func OTProto() sensors.OTProto {
	if app, err := gopi.NewAppInstance(gopi.NewAppConfig("sensors/protocol/openthings")); err != nil {
		return nil
	} else if proto, ok := app.ModuleInstance("sensors/protocol/openthings").(sensors.OTProto); ok == false {
		return nil
	} else {
		return proto
	}
}
