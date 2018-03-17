/*
   Go Language Raspberry Pi Interface
   (c) Copyright David Thorpe 2016-2017
   All Rights Reserved
   Documentation http://djthorpe.github.io/gopi/
   For Licensing and Usage information, please see LICENSE.md
*/

// Interacts with the BME280 sensor over the I2C bus
package main

import (
	"errors"
	"fmt"
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"

	// Register modules
	_ "github.com/djthorpe/gopi/sys/hw/linux"
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/sensors/hw/bme280"

	// Module flavours for i2c and spi
	"github.com/djthorpe/sensors/cmd/bme280"

	// Tablewriter
	"github.com/olekukonko/tablewriter"
)

const (
	COMMAND_MEASURE = iota
	COMMAND_RESET
	COMMAND_STATUS
	COMMAND_HELP
)

////////////////////////////////////////////////////////////////////////////////

func GetModeFromString(mode string) (sensors.BME280Mode, error) {
	switch mode {
	case "normal":
		return sensors.BME280_MODE_NORMAL, nil
	case "forced":
		return sensors.BME280_MODE_FORCED, nil
	case "sleep":
		return sensors.BME280_MODE_SLEEP, nil
	default:
		return sensors.BME280_MODE_NORMAL, fmt.Errorf("Invalid mode: %v", mode)
	}
}

func GetOversampleFromUint(value uint) (sensors.BME280Oversample, error) {
	switch value {
	case 0:
		return sensors.BME280_OVERSAMPLE_SKIP, nil
	case 1:
		return sensors.BME280_OVERSAMPLE_1, nil
	case 2:
		return sensors.BME280_OVERSAMPLE_2, nil
	case 4:
		return sensors.BME280_OVERSAMPLE_4, nil
	case 8:
		return sensors.BME280_OVERSAMPLE_8, nil
	case 16:
		return sensors.BME280_OVERSAMPLE_16, nil
	default:
		return sensors.BME280_OVERSAMPLE_SKIP, fmt.Errorf("Invalid oversample value: %v", value)
	}
}

func GetFilterFromUint(value uint) (sensors.BME280Filter, error) {
	switch value {
	case 0:
		return sensors.BME280_FILTER_OFF, nil
	case 2:
		return sensors.BME280_FILTER_2, nil
	case 4:
		return sensors.BME280_FILTER_4, nil
	case 8:
		return sensors.BME280_FILTER_8, nil
	case 16:
		return sensors.BME280_FILTER_16, nil
	default:
		return sensors.BME280_FILTER_OFF, fmt.Errorf("Invalid filter value: %v", value)
	}
}

func GetStandbyFromFloat(value float64) (sensors.BME280Standby, error) {
	switch value {
	case 0.5:
		return sensors.BME280_STANDBY_0P5MS, nil
	case 62.5:
		return sensors.BME280_STANDBY_62P5MS, nil
	case 125:
		return sensors.BME280_STANDBY_125MS, nil
	case 250:
		return sensors.BME280_STANDBY_250MS, nil
	case 500:
		return sensors.BME280_STANDBY_500MS, nil
	case 1000:
		return sensors.BME280_STANDBY_1000MS, nil
	case 10:
		return sensors.BME280_STANDBY_10MS, nil
	case 20:
		return sensors.BME280_STANDBY_20MS, nil
	default:
		return sensors.BME280_STANDBY_MAX, fmt.Errorf("Invalid standby value: %v", value)
	}
}

////////////////////////////////////////////////////////////////////////////////

func status(device sensors.BME280) error {

	// Output register information
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Register", "Value"})

	chip_id, chip_version := device.ChipIDVersion()
	table.Append([]string{"chip_id", fmt.Sprintf("0x%02X", chip_id)})
	table.Append([]string{"chip_version", fmt.Sprintf("0x%02X", chip_version)})
	table.Append([]string{"mode", fmt.Sprint(device.Mode())})
	table.Append([]string{"filter", fmt.Sprint(device.Filter())})
	table.Append([]string{"standby", fmt.Sprint(device.Standby())})
	table.Append([]string{"duty cycle", fmt.Sprint(device.DutyCycle())})

	t, p, h := device.Oversample()
	table.Append([]string{"oversample temperature", fmt.Sprint(t)})
	table.Append([]string{"oversample pressure", fmt.Sprint(p)})
	table.Append([]string{"oversample humidity", fmt.Sprint(h)})

	if measuring, updating, err := device.Status(); err != nil {
		return err
	} else {
		table.Append([]string{"measuring", fmt.Sprint(measuring)})
		table.Append([]string{"updating", fmt.Sprint(updating)})
	}

	table.Render()
	return nil
}

func measure(device sensors.BME280) error {

	// If oversample temperature is in skip mode, then set to 1 to ensure
	// we have at least one thing to measure!
	if t_os, p_os, h_os := device.Oversample(); t_os == sensors.BME280_OVERSAMPLE_SKIP {
		if err := device.SetOversample(sensors.BME280_OVERSAMPLE_1, p_os, h_os); err != nil {
			return err
		}
	}

	// If sensor is in sleep mode then change to forced mode,
	// which will return it to sleep mode once the sample
	// has been read
	if device.Mode() == sensors.BME280_MODE_SLEEP {
		if err := device.SetMode(sensors.BME280_MODE_FORCED); err != nil {
			return err
		}
	}

	if t, p, h, err := device.ReadSample(); err != nil {
		return err
	} else {
		a := device.AltitudeForPressure(p, sensors.BME280_PRESSURE_SEALEVEL)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAlignment(tablewriter.ALIGN_RIGHT)
		table.SetHeader([]string{"Measurement", "Value"})
		if t > 0 {
			table.Append([]string{"temperature", fmt.Sprintf("%.2f \u00B0C", t)})
		}
		if p > 0.0 {
			table.Append([]string{"pressure", fmt.Sprintf("%.2f hPa", p)})
			table.Append([]string{"altitude", fmt.Sprintf("%.2f m", a)})
		}
		if h > 0 {
			table.Append([]string{"humidity", fmt.Sprintf("%.2f %%RH", h)})
		}
		table.Render()
	}
	return nil
}

func reset(device sensors.BME280) error {
	if err := device.SoftReset(); err != nil {
		return err
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// SET

func set_mode(device sensors.BME280, mode string) error {
	if mode_value, err := GetModeFromString(mode); err != nil {
		return err
	} else if err := device.SetMode(mode_value); err != nil {
		return err
	} else {
		return nil
	}
}

func set_filter(device sensors.BME280, filter uint) error {
	if filter_value, err := GetFilterFromUint(filter); err != nil {
		return err
	} else if err := device.SetFilter(filter_value); err != nil {
		return err
	} else {
		return nil
	}
}

func set_oversample(device sensors.BME280, oversample uint) error {
	if oversample_value, err := GetOversampleFromUint(oversample); err != nil {
		return err
	} else if err := device.SetOversample(oversample_value, oversample_value, oversample_value); err != nil {
		return err
	} else {
		return nil
	}
}

func set_standby(device sensors.BME280, standby float64) error {
	if standby_value, err := GetStandbyFromFloat(standby); err != nil {
		return err
	} else if err := device.SetStandby(standby_value); err != nil {
		return err
	} else {
		return nil
	}
}

func set(app *gopi.AppInstance, device sensors.BME280) error {
	if mode, exists := app.AppFlags.GetString("mode"); exists {
		if err := set_mode(device, mode); err != nil {
			return err
		}
	}
	if filter, exists := app.AppFlags.GetUint("filter"); exists {
		if err := set_filter(device, filter); err != nil {
			return err
		}
	}
	if oversample, exists := app.AppFlags.GetUint("oversample"); exists {
		if err := set_oversample(device, oversample); err != nil {
			return err
		}
	}
	if standby, exists := app.AppFlags.GetFloat64("standby"); exists {
		if err := set_standby(device, standby); err != nil {
			return err
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// MAIN FUNCTION

func MainLoop(app *gopi.AppInstance, done chan<- struct{}) error {

	// Determine the command to run
	command := COMMAND_HELP
	if args := app.AppFlags.Args(); len(args) == 1 && args[0] == "reset" {
		command = COMMAND_RESET
	} else if len(args) == 1 && args[0] == "status" {
		command = COMMAND_STATUS
	} else if len(args) == 0 || len(args) == 1 && args[0] == "measure" {
		command = COMMAND_MEASURE
	}

	// Run the command
	if device := app.ModuleInstance(bme280.MODULE_NAME).(sensors.BME280); device == nil {
		return errors.New("BME280 module not found")
	} else {
		switch command {
		case COMMAND_MEASURE:
			if err := set(app, device); err != nil {
				return err
			}
			if err := measure(device); err != nil {
				return err
			}
		case COMMAND_RESET:
			if err := reset(device); err != nil {
				return err
			}
			if err := set(app, device); err != nil {
				return err
			}
			if err := status(device); err != nil {
				return err
			}
		case COMMAND_STATUS:
			if err := set(app, device); err != nil {
				return err
			}
			if err := status(device); err != nil {
				return err
			}
		default:
			return gopi.ErrHelp
		}
	}

	// Exit
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig(bme280.MODULE_NAME)

	// Parameters
	config.AppFlags.FlagString("mode", "", "Sensor mode (normal,forced,sleep)")
	config.AppFlags.FlagUint("filter", 0, "Filter co-efficient (0,2,4,8,16)")
	config.AppFlags.FlagUint("oversample", 0, "Oversampling (0,1,2,4,8,16)")
	config.AppFlags.FlagFloat64("standby", 0, "Standby time, ms (0.5,10,20,62.5,125,250,500,1000)")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, MainLoop))
}
