/*
   Go Language Raspberry Pi Interface
   (c) Copyright David Thorpe 2016-2018
   All Rights Reserved
   Documentation http://djthorpe.github.io/gopi/
   For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"fmt"
	"os"
	"strings"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
/*
func ProductName(message sensors.OTMessage) string {
	if message.Manufacturer() == sensors.OT_MANUFACTURER_ENERGENIE {
		product := strings.TrimPrefix(fmt.Sprint(sensors.MiHomeProduct(message.Product())), "MIHOME_PRODUCT_")
		return product
	} else {
		manufacturer := strings.TrimPrefix(fmt.Sprint(message.Manufacturer()), "OT_MANUFACTURER_")
		return fmt.Sprintf("%v<%02X>", manufacturer, message.Product())
	}
}

func Sensor(message sensors.OTMessage) string {
	return fmt.Sprintf("%s<%05X>", ProductName(message), message.Sensor())
}

func PrintEvent(evt gopi.Event) {
	if message, ok := evt.(sensors.OTMessage); ok {
		records := ""
		for _, record := range message.Records() {
			records += fmt.Sprint(record) + " "
		}
		fmt.Printf("%20s %s\n", Sensor(message), strings.TrimSpace(records))

		if sensors.MiHomeProduct(message.Product()) == sensors.MIHOME_PRODUCT_MIHO013 {

		}

	} else if message, ok := evt.(sensors.OOKMessage); ok {
		fmt.Println("OOKMessage", message)
	} else {
		fmt.Println("Other", evt)
	}
}

func RequestIdentify(mihome sensors.MiHome, evt gopi.Event) {
	if message, ok := evt.(sensors.OTMessage); ok == false {
		return
	} else if sensors.MiHomeProduct(message.Product()) != sensors.MIHOME_PRODUCT_MIHO013 {
		return
	} else if records := message.Records(); len(records) == 0 {
		return
	} else if records[0].Name() != sensors.OT_PARAM_TEMPERATURE {
		return
	} else if err := mihome.RequestIdentify(sensors.MiHomeProduct(message.Product()), message.Sensor()); err != nil {
		fmt.Println(err)
	}
}

func RequestDiagnostics(mihome sensors.MiHome, evt gopi.Event) {
	if message, ok := evt.(sensors.OTMessage); ok == false {
		return
	} else if sensors.MiHomeProduct(message.Product()) != sensors.MIHOME_PRODUCT_MIHO013 {
		return
	} else if records := message.Records(); len(records) == 0 {
		return
	} else if records[0].Name() != sensors.OT_PARAM_TEMPERATURE {
		return
	} else if err := mihome.RequestDiagnostics(sensors.MiHomeProduct(message.Product()), message.Sensor()); err != nil {
		fmt.Println(err)
	}
}

func RequestExercise(mihome sensors.MiHome, evt gopi.Event) {
	if message, ok := evt.(sensors.OTMessage); ok == false {
		return
	} else if sensors.MiHomeProduct(message.Product()) != sensors.MIHOME_PRODUCT_MIHO013 {
		return
	} else if records := message.Records(); len(records) == 0 {
		return
	} else if records[0].Name() != sensors.OT_PARAM_TEMPERATURE {
		return
	} else if err := mihome.RequestExercise(sensors.MiHomeProduct(message.Product()), message.Sensor()); err != nil {
		fmt.Println(err)
	}
}

func RequestBatteryLevel(mihome sensors.MiHome, evt gopi.Event) {
	if message, ok := evt.(sensors.OTMessage); ok == false {
		return
	} else if sensors.MiHomeProduct(message.Product()) != sensors.MIHOME_PRODUCT_MIHO013 {
		return
	} else if records := message.Records(); len(records) == 0 {
		return
	} else if records[0].Name() != sensors.OT_PARAM_TEMPERATURE {
		return
	} else if err := mihome.RequestBatteryLevel(sensors.MiHomeProduct(message.Product()), message.Sensor()); err != nil {
		fmt.Println(err)
	}
}

func RequestReportingInterval(mihome sensors.MiHome, evt gopi.Event) {
	if message, ok := evt.(sensors.OTMessage); ok == false {
		return
	} else if sensors.MiHomeProduct(message.Product()) != sensors.MIHOME_PRODUCT_MIHO013 {
		return
	} else if records := message.Records(); len(records) == 0 {
		return
	} else if records[0].Name() != sensors.OT_PARAM_TEMPERATURE {
		return
	} else if err := mihome.RequestReportInterval(sensors.MiHomeProduct(message.Product()), message.Sensor(), time.Second*300); err != nil {
		fmt.Println(err)
	}
}

func RequestJoin(mihome sensors.MiHome, evt gopi.Event) {
	if message, ok := evt.(sensors.OTMessage); ok == false {
		return
	} else if records := message.Records(); len(records) == 0 {
		return
	} else if records[0].Name() != sensors.OT_PARAM_JOIN {
		return
	} else if err := mihome.SendJoin(sensors.MiHomeProduct(message.Product()), message.Sensor()); err != nil {
		fmt.Println(err)
	}
}

func RequestTargetTemp(mihome sensors.MiHome, evt gopi.Event) {
	if message, ok := evt.(sensors.OTMessage); ok == false {
		return
	} else if records := message.Records(); len(records) == 0 {
		return
	} else if records[0].Name() != sensors.OT_PARAM_TEMPERATURE {
		return
	} else if err := mihome.RequestTargetTemperature(sensors.MiHomeProduct(message.Product()), message.Sensor(), 22); err != nil {
		fmt.Println(err)
	}
}
*/

func RecordString(message sensors.OTMessage) string {
	records := ""
	for _, record := range message.Records() {
		records += fmt.Sprint(record) + " "
	}
	return strings.TrimSpace(records)
}

func ProcessEvent(app *gopi.AppInstance, evt gopi.Event) error {
	if db := app.ModuleInstance("sensors/db").(sensors.Database); db == nil {
		return fmt.Errorf("Missing or invalid sensors database")
	} else if message, ok := evt.(sensors.Message); ok {
		if sensor, err := db.Register(message); err != nil {
			return err
		} else if message_, ok := message.(sensors.OTMessage); ok {
			fmt.Printf("%9s %30s | %s\n", sensor.Key(), sensor.Description(), RecordString(message_))
		} else {
			fmt.Printf("%9s %30s | %s\n", sensor.Key(), sensor.Description(), message)
		}
		return nil
	} else {
		return fmt.Errorf("Unknown event: %v", evt)
	}
}

func Receive(app *gopi.AppInstance, start chan<- struct{}, stop <-chan struct{}) error {

	// Get slave flag
	//slave, _ := app.AppFlags.GetBool("slave")

	// Reset the mihome device
	mihome := app.ModuleInstance("sensors/mihome").(sensors.MiHome)
	if mihome == nil {
		return gopi.ErrAppError
	}
	if err := mihome.Reset(); err != nil {
		return err
	}

	// Subscribe to events from the MiHome device
	evts := mihome.Subscribe()

	// Event loop
	start <- gopi.DONE
FOR_LOOP:
	for {
		select {
		case <-stop:
			break FOR_LOOP
		case evt := <-evts:
			if err := ProcessEvent(app, evt); err != nil {
				app.Logger.Error("ProcessEvent: %v", err)
			}
		}
	}

	// Unsubscribe
	mihome.Unsubscribe(evts)

	// Return success
	return nil
}

func Main(app *gopi.AppInstance, done chan<- struct{}) error {
	// Wait for signal
	app.WaitForSignal()
	// Send done signal to tasks
	done <- gopi.DONE
	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("sensors/mihome", "sensors/protocol/ook", "sensors/protocol/openthings", "sensors/db")
	config.AppFlags.FlagBool("slave", false, "Listen only, don't send")

	// Usage
	/*config.AppFlags.SetUsageFunc(func(flags *gopi.Flags) {
		//fmt.Fprintf(os.Stderr, "Usage:\n  %v <flags...> <%v>\n\n", flags.Name(), strings.Join(CommandNames(), "|"))
		//PrintCommands(os.Stderr)
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		flags.PrintDefaults()
	})*/

	// Run the command line tool
	os.Exit(gopi.CommandLineTool2(config, Main, Receive))
}
