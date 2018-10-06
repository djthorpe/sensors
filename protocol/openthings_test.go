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
	} else if _, ok := app.ModuleInstance("sensors/protocol/openthings").(sensors.OTProto); ok == false {
		t.Fatal("OOK does not comply to OTProto interface")
	}
}
