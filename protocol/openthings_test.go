package protocol_test

import (
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
	} else if _, err := proto.New(sensors.OT_MANUFACTURER_ENERGENIE, 0, 0); err != nil {
		t.Fatal(err)
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
