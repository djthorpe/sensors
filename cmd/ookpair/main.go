// +build rpi

package main

import (
	"encoding/hex"
	"os"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"

	// Modules
	_ "github.com/djthorpe/gopi-hw/sys/gpio"
	_ "github.com/djthorpe/gopi-hw/sys/hw"
	_ "github.com/djthorpe/gopi-hw/sys/metrics"
	_ "github.com/djthorpe/gopi-hw/sys/spi"
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/sensors/protocol/ook"
	_ "github.com/djthorpe/sensors/sys/ener314rt"
	_ "github.com/djthorpe/sensors/sys/mihome"
	_ "github.com/djthorpe/sensors/sys/rfm69"
)

////////////////////////////////////////////////////////////////////////////////

func Main(app *gopi.AppInstance, done chan<- struct{}) error {
	ener314rt := app.ModuleInstance("sensors/ener314rt").(sensors.ENER314RT)
	if err := ener314rt.ResetRadio(); err != nil {
		return err
	} else if payload_off, err := hex.DecodeString("800000008EE8EE888EE8EE888EE8EEEE"); err != nil {
		return err
	} else if payload_on, err := hex.DecodeString("800000008EE8EE888EE8EE888EE8EEE8"); err != nil {
		return err
	} else {
		for i := 0; i < 5; i++ {
			if err := ener314rt.Send(payload_on, 1, sensors.MIHOME_MODE_CONTROL); err != nil {
				return err
			}
			time.Sleep(time.Second)
			if err := ener314rt.Send(payload_off, 1, sensors.MIHOME_MODE_CONTROL); err != nil {
				return err
			}
		}
	}
	/*
		mihome := app.ModuleInstance("sensors/mihome").(sensors.MiHome)

		if err := mihome.RequestSwitchOn(sensors.MIHOME_PRODUCT_CONTROL_ALL, 0x6C6C6); err != nil {
			return err
		}

		time.Sleep(1 * time.Second)

		if err := mihome.RequestSwitchOff(sensors.MIHOME_PRODUCT_CONTROL_ALL, 0x6C6C6); err != nil {
			return err
		}
	*/

	// Return success
	return nil
}

func main() {
	// Create the configuration
	//	config := gopi.NewAppConfig("sensors/mihome", "sensors/protocol/ook")
	config := gopi.NewAppConfig("sensors/ener314rt", "sensors/protocol/ook")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool2(config, Main))
}
