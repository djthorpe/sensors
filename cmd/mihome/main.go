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
	"regexp"
	"strings"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"
	"github.com/olekukonko/tablewriter"
)

type CommandFunc func() error
type ActionFunc func(app *gopi.AppInstance, sensor sensors.Sensor, message sensors.Message) error

var (
	actions = map[string]ActionFunc{
		"0C:000C70": ActionMotionSensorInLivingRoom, // Turn on/off christmas tree
		"13:000F25": ActionClicker,                  // Turn on/off light
	}
	command_queue = make(chan CommandFunc, 100)
	regexp_sensor = regexp.MustCompile("^(\\w+):([0-9A-Fa-f]+:[0-9A-Fa-f]+)$")
)

////////////////////////////////////////////////////////////////////////////////
// ACTIONS

func ActionMotionSensorInLivingRoom(app *gopi.AppInstance, sensor sensors.Sensor, message sensors.Message) error {
	if db := app.ModuleInstance("sensors/db").(sensors.Database); db == nil {
		return fmt.Errorf("Missing or invalid sensors database")
	} else if sensor = db.Lookup("openthings", "02:002ED3"); sensor == nil {
		return fmt.Errorf("Unknown sensor device")
	}

	if message_, ok := message.(sensors.OTMessage); ok {
		for _, record := range message_.Records() {
			if record.Name() == sensors.OT_PARAM_MOTION_DETECTOR {
				if value, err := record.BoolValue(); err != nil {
					return err
				} else if value == false {
					// Switch on socket one
					command_queue <- func() error {
						return CommandOn(app, sensor)
					}

				} else if value == true {
					// Switch off socket one
					command_queue <- func() error {
						return CommandOff(app, sensor)
					}
				}
			}
		}
	}
	return nil
}

func ActionClicker(app *gopi.AppInstance, sensor sensors.Sensor, message sensors.Message) error {
	if db := app.ModuleInstance("sensors/db").(sensors.Database); db == nil {
		return fmt.Errorf("Missing or invalid sensors database")
	} else if sensor = db.Lookup("ook", "F1:6C6C6"); sensor == nil {
		return fmt.Errorf("Unknown sensor device")
	}

	if message_, ok := message.(sensors.OTMessage); ok {
		for _, record := range message_.Records() {
			if record.Name() == sensors.OT_PARAM_CLICK {
				if value, err := record.UintValue(); err != nil {
					return err
				} else if value == 1 {
					// Switch on socket one
					command_queue <- func() error {
						return CommandOn(app, sensor)
					}
				} else if value == 2 {
					// Switch off socket one
					command_queue <- func() error {
						return CommandOff(app, sensor)
					}
				}
			}
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// ON AND OFF

/*
func CommandOnOld(app *gopi.AppInstance, sensor sensors.Sensor) error {
	if sensor_flag, _ := app.AppFlags.GetString("sensor"); sensor_flag == "" {
		return fmt.Errorf("Missing -sensor argument")
	} else if db := app.ModuleInstance("sensors/db").(sensors.Database); db == nil {
		return fmt.Errorf("Missing or invalid sensors database")
	} else if parts := regexp_sensor.FindStringSubmatch(sensor_flag); len(parts) != 3 {
		return fmt.Errorf("Invalid -sensor argument")
	} else if mihome := app.ModuleInstance("sensors/mihome").(sensors.MiHome); mihome == nil {
		return fmt.Errorf("Missing mihome device")
	} else if sensor := db.Lookup(parts[1], parts[2]); sensor == nil {
		return fmt.Errorf("Unknown sensor device")
	} else if err := mihome.RequestSwitchOn(sensors.MiHomeProduct(sensor.Product()), sensor.Sensor()); err != nil {
		return err
	}

	return nil
}
*/

// Switch a sensor on
func CommandOn(app *gopi.AppInstance, sensor sensors.Sensor) error {
	if mihome := app.ModuleInstance("sensors/mihome").(sensors.MiHome); mihome == nil {
		return fmt.Errorf("Missing mihome device")
	} else if err := mihome.RequestSwitchOn(sensors.MiHomeProduct(sensor.Product()), sensor.Sensor()); err != nil {
		return err
	}
	return nil
}

// Switch a sensor off
func CommandOff(app *gopi.AppInstance, sensor sensors.Sensor) error {
	if mihome := app.ModuleInstance("sensors/mihome").(sensors.MiHome); mihome == nil {
		return fmt.Errorf("Missing mihome device")
	} else if err := mihome.RequestSwitchOff(sensors.MiHomeProduct(sensor.Product()), sensor.Sensor()); err != nil {
		return err
	}
	return nil
}

// CommandList will list all the current sensors
func CommandList(app *gopi.AppInstance, sensor sensors.Sensor) error {
	if db := app.ModuleInstance("sensors/db").(sensors.Database); db == nil {
		return fmt.Errorf("Missing or invalid sensors database")
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAlignment(tablewriter.ALIGN_RIGHT)
		table.SetHeader([]string{"Sensor", "Description", "Last Seen"})
		for _, sensor := range db.Sensors() {
			ts_ := time.Since(sensor.Timestamp()).Truncate(time.Minute)
			table.Append([]string{
				fmt.Sprintf("%v:%v", sensor.Namespace(), sensor.Key()),
				sensor.Description(),
				fmt.Sprint(ts_),
			})
		}
		table.Render()
	}
	return nil
}

/*
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
		} else {
			// Print out the message
			if message_, ok := message.(sensors.OTMessage); ok {
				fmt.Printf("%9s %30s | %s\n", sensor.Key(), sensor.Description(), RecordString(message_))
			} else {
				fmt.Printf("%9s %30s | %s\n", sensor.Key(), sensor.Description(), message)
			}
			// Perform action on the message
			if action, exists := actions[sensor.Key()]; exists {
				action(app, sensor, message)
			}
		}
	}
	return nil
}

func Receive(app *gopi.AppInstance, start chan<- struct{}, stop <-chan struct{}) error {
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
		case fn := <-command_queue:
			if err := fn(); err != nil {
				app.Logger.Error("ProcessCommand: %v", err)
			}
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

	// Run the command line tool
	os.Exit(gopi.CommandLineTool2(config, Main, Receive))
}
