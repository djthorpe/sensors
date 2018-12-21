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

type CommandFunc func(app *gopi.AppInstance) error

var (
	commands = map[string]CommandFunc{
		"on":   CommandOn,
		"off":  CommandOff,
		"list": CommandList,
	}
	regexp_sensor = regexp.MustCompile("^(\\w+):([0-9A-Fa-f]+:[0-9A-Fa-f]+)$")
)

////////////////////////////////////////////////////////////////////////////////
// ON AND OFF

func CommandOn(app *gopi.AppInstance) error {
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

func CommandOff(app *gopi.AppInstance) error {
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
	} else if err := mihome.RequestSwitchOff(sensors.MiHomeProduct(sensor.Product()), sensor.Sensor()); err != nil {
		return err
	}

	return nil
}

// CommandList will list all the current sensors
func CommandList(app *gopi.AppInstance) error {
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

func SendRequestSwitchState(app *gopi.AppInstance, state bool) error {

	if mihome := app.ModuleInstance("sensors/mihome").(sensors.MiHome); mihome == nil {
		return gopi.ErrAppError
	} else if db := app.ModuleInstance("sensors/db").(sensors.Database); db == nil {
		return fmt.Errorf("Missing or invalid sensors database")
	} else if sensor_flag, _ := app.AppFlags.GetString("sensor"); sensor_flag == "" {
		return fmt.Errorf("Missing -sensor argument")
	} else if parts := regexp_sensor.FindStringSubmatch(sensor_flag); len(parts) != 3 {
		return fmt.Errorf("Invalid -sensor argument")
	} else if sensor := db.Lookup(parts[1], parts[2]); sensor == nil {
		return fmt.Errorf("Unknown sensor device")
	} else {
		switch state {
		case false:
			app.Logger.Info("off=%v", sensor)
			return mihome.RequestSwitchOff(sensors.MiHomeProduct(sensor.Product()), sensor.Sensor())
		case true:
			app.Logger.Info("on=%v", sensor)
			return mihome.RequestSwitchOn(sensors.MiHomeProduct(sensor.Product()), sensor.Sensor())
		}
	}

	// Success
	return nil
}

func Send(app *gopi.AppInstance, start chan<- struct{}, stop <-chan struct{}) error {
	app.Logger.Info("Send started")

	command_queue := make(chan CommandFunc, 1)
	for _, command := range app.AppFlags.Args() {
		command_ := strings.ToLower(command)
		if f, exists := commands[command_]; exists == false {
			start <- gopi.DONE
			app.SendSignal()
			return fmt.Errorf("Invalid command: '%v'", command)
		} else {
			command_queue <- f
		}
	}

	// Event loop
	start <- gopi.DONE
FOR_LOOP:
	for {
		select {
		case <-stop:
			break FOR_LOOP
		case fn := <-command_queue:
			if err := fn(app); err != nil {
				app.Logger.Error("Error: %v", err)
			}
		}
	}

	app.Logger.Info("Send stopped")
	close(command_queue)
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
	config.AppFlags.FlagString("sensor", "ook:F0:6C6C6", "Sensor to control")

	// Usage
	/*config.AppFlags.SetUsageFunc(func(flags *gopi.Flags) {
		//fmt.Fprintf(os.Stderr, "Usage:\n  %v <flags...> <%v>\n\n", flags.Name(), strings.Join(CommandNames(), "|"))
		//PrintCommands(os.Stderr)
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		flags.PrintDefaults()
	})*/

	// Run the command line tool
	os.Exit(gopi.CommandLineTool2(config, Main, Receive, Send))
}
