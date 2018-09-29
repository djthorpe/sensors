/*
   Go Language Raspberry Pi Interface
   (c) Copyright David Thorpe 2016-2018
   All Rights Reserved
   Documentation http://djthorpe.github.io/gopi/
   For Licensing and Usage information, please see LICENSE.md
*/

// Interacts with the RFM69 device
package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func setParametersMode(app *gopi.AppInstance, device sensors.RFM69) error {
	if value, exists := app.AppFlags.GetString("mode"); exists {
		if mode, err := stringToMode(value); err != nil {
			return err
		} else if err := device.SetMode(mode); err != nil {
			return err
		} else if value == "listen" {
			if err := device.SetListenOn(true); err != nil {
				return err
			}
		}
	}

	if enabled, exists := app.AppFlags.GetBool("sequencer"); exists {
		if err := device.SetSequencer(enabled); err != nil {
			return err
		}
	}

	if enabled, exists := app.AppFlags.GetBool("listen"); exists {
		if err := device.SetListenOn(enabled); err != nil {
			return err
		}
	}

	// Success
	return nil
}

func setParametersDataMode(app *gopi.AppInstance, device sensors.RFM69) error {
	if value, exists := app.AppFlags.GetString("datamode"); exists == false {
		return nil
	} else if mode, err := stringToDataMode(value); err != nil {
		return err
	} else if err := device.SetDataMode(mode); err != nil {
		return err
	}

	// Success
	return nil
}

func setParametersBitrate(app *gopi.AppInstance, device sensors.RFM69) error {
	if value, exists := app.AppFlags.GetFloat64("bitrate"); exists {
		if err := device.SetBitrate(uint(value * 1000)); err != nil {
			return err
		}
	}

	if value, exists := app.AppFlags.GetFloat64("freq_carrier"); exists {
		if err := device.SetFreqCarrier(uint(value * 1000)); err != nil {
			return err
		}
	}

	if value, exists := app.AppFlags.GetFloat64("freq_dev"); exists {
		if err := device.SetFreqDeviation(uint(value * 1000)); err != nil {
			return err
		}
	}

	// Success
	return nil
}

func setParametersPacket(app *gopi.AppInstance, device sensors.RFM69) error {
	if value, exists := app.AppFlags.GetString("packet_format"); exists {
		if format, err := stringToPacketFormat(value); err != nil {
			return err
		} else if err := device.SetPacketFormat(format); err != nil {
			return err
		}
	}

	if value, exists := app.AppFlags.GetString("packet_coding"); exists {
		if format, err := stringToPacketCoding(value); err != nil {
			return err
		} else if err := device.SetPacketCoding(format); err != nil {
			return err
		}
	}

	if value, exists := app.AppFlags.GetString("packet_filter"); exists {
		if format, err := stringToPacketFilter(value); err != nil {
			return err
		} else if err := device.SetPacketFilter(format); err != nil {
			return err
		}
	}

	if value, exists := app.AppFlags.GetString("packet_crc"); exists {
		if crc, err := stringToPacketCRC(value); err != nil {
			return err
		} else if err := device.SetPacketCRC(crc); err != nil {
			return err
		}
	}

	if value, exists := app.AppFlags.GetUint("preamble_size"); exists {
		if err := device.SetPreambleSize(uint16(value)); err != nil {
			return err
		}
	}

	if value, exists := app.AppFlags.GetUint("payload_size"); exists {
		if err := device.SetPayloadSize(uint8(value)); err != nil {
			return err
		}
	}

	// Success
	return nil
}

func setParametersModulation(app *gopi.AppInstance, device sensors.RFM69) error {
	if value, exists := app.AppFlags.GetString("modulation"); exists == false {
		return nil
	} else if modulation, err := stringToModulation(value); err != nil {
		return err
	} else if err := device.SetModulation(modulation); err != nil {
		return err
	}

	// Success
	return nil
}

func setParametersAFC(app *gopi.AppInstance, device sensors.RFM69) error {
	if value, exists := app.AppFlags.GetString("afc_mode"); exists {
		if mode, err := stringToAFCMode(value); err != nil {
			return err
		} else if err := device.SetAFCMode(mode); err != nil {
			return err
		}
	}

	if value, exists := app.AppFlags.GetString("afc_routine"); exists {
		if routine, err := stringToAFCRoutine(value); err != nil {
			return err
		} else if err := device.SetAFCRoutine(routine); err != nil {
			return err
		}
	}

	// Success
	return nil
}

func setParametersNodeBroadcastAddr(app *gopi.AppInstance, device sensors.RFM69) error {
	if value, exists := app.AppFlags.GetString("node_addr"); exists {
		if addr, err := hex.DecodeString(value); err != nil || len(addr) != 1 {
			return fmt.Errorf("Invalid node_addr: %v", value)
		} else if err := device.SetNodeAddress(addr[0]); err != nil {
			return err
		}
	}

	if value, exists := app.AppFlags.GetString("broadcast_addr"); exists {
		if addr, err := hex.DecodeString(value); err != nil || len(addr) != 1 {
			return fmt.Errorf("Invalid broadcast_addr: %v", value)
		} else if err := device.SetBroadcastAddress(addr[0]); err != nil {
			return err
		}
	}

	// Success
	return nil
}

func setParametersAESKey(app *gopi.AppInstance, device sensors.RFM69) error {
	if value, exists := app.AppFlags.GetString("aes_key"); exists {
		if value == "" {
			// Disable AES
			if err := device.SetAESKey(nil); err != nil {
				return err
			}
		} else if key, err := hex.DecodeString(value); err != nil {
			return err
		} else if err := device.SetAESKey(key); err != nil {
			return err
		}
	}

	// Success
	return nil
}

func setParametersSyncWord(app *gopi.AppInstance, device sensors.RFM69) error {
	if value, exists := app.AppFlags.GetString("sync_word"); exists {
		if value == "" {
			// Disable Sync Word
			if err := device.SetSyncWord(nil); err != nil {
				return err
			}
		} else if word, err := hex.DecodeString(value); err != nil {
			return err
		} else if err := device.SetSyncWord(word); err != nil {
			return err
		}
	}

	if value, exists := app.AppFlags.GetUint("sync_tol"); exists {
		if value > 7 {
			return fmt.Errorf("Invalid sync_tol value: %v", value)
		} else if err := device.SetSyncTolerance(uint8(value)); err != nil {
			return err
		}
	}

	// Success
	return nil
}

func setParametersFIFO(app *gopi.AppInstance, device sensors.RFM69) error {
	if value, exists := app.AppFlags.GetUint("fifo_threshold"); exists {
		if value > 0xFF {
			return gopi.ErrBadParameter
		}
		if err := device.SetFIFOThreshold(uint8(value)); err != nil {
			return err
		}
	}

	// Success
	return nil
}

func setParameters(app *gopi.AppInstance, device sensors.RFM69) error {
	if err := setParametersMode(app, device); err != nil {
		return err
	}
	if err := setParametersDataMode(app, device); err != nil {
		return err
	}
	if err := setParametersBitrate(app, device); err != nil {
		return err
	}
	if err := setParametersPacket(app, device); err != nil {
		return err
	}
	if err := setParametersModulation(app, device); err != nil {
		return err
	}
	if err := setParametersAFC(app, device); err != nil {
		return err
	}
	if err := setParametersNodeBroadcastAddr(app, device); err != nil {
		return err
	}
	if err := setParametersAESKey(app, device); err != nil {
		return err
	}
	if err := setParametersSyncWord(app, device); err != nil {
		return err
	}
	if err := setParametersFIFO(app, device); err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RUN ARGS

func RunArgs(app *gopi.AppInstance, device sensors.RFM69) error {

	// Read the arguments, or use default 'Status'
	args := app.AppFlags.Args()
	if len(args) == 0 {
		args = []string{"Status"}
	}

	// Run the commands
	for _, arg := range args {
		if cmd, exists := command_map[arg]; exists == false {
			return fmt.Errorf("Invalid command: %v", arg)
		} else if err := cmd(app, device); err != nil {
			return err
		}
	}

	// success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// MAIN FUNCTION

func MainLoop(app *gopi.AppInstance, done chan<- struct{}) error {

	// Run the command
	if device := app.ModuleInstance(MODULE_NAME).(sensors.RFM69); device == nil {
		return errors.New("Module not found: " + MODULE_NAME)
	} else if err := setParameters(app, device); err != nil {
		return err
	} else if err := RunArgs(app, device); err != nil {
		return err
	}

	// Exit
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration, load the RFM69 instance
	config := gopi.NewAppConfig(MODULE_NAME)

	// Parameters
	config.AppFlags.FlagString("mode", "", "Device Mode (sleep,standby,fs,tx,rx,listen)")
	config.AppFlags.FlagBool("sequencer", false, "Enable sequencer")
	config.AppFlags.FlagBool("listen", false, "Enable listen mode")
	config.AppFlags.FlagString("datamode", "", "Data Mode (packet,nosync,sync)")
	config.AppFlags.FlagString("modulation", "", "Modulation (fsk,fsk_1.0,fsk_0.5,fsk_0.3,ook,ook_br,ook_2br)")
	config.AppFlags.FlagFloat64("bitrate", 0, "Bitrate (kbps)")
	config.AppFlags.FlagFloat64("freq_carrier", 0, "Carrier Frequency (kbps)")
	config.AppFlags.FlagFloat64("freq_dev", 0, "Frequency Deviation (kbps)")
	config.AppFlags.FlagString("node_addr", "", "Node Address (byte)")
	config.AppFlags.FlagString("broadcast_addr", "", "Broadcast Address (byte)")
	config.AppFlags.FlagUint("payload_size", 0, "Payload Size (bytes)")
	config.AppFlags.FlagUint("preamble_size", 0, "Preamble Size (bytes)")
	config.AppFlags.FlagString("aes_key", "", "AES Key (16 bytes) or empty")
	config.AppFlags.FlagString("sync_word", "", "Sync Word (1-8 bytes) or empty")
	config.AppFlags.FlagUint("sync_tol", 0, "Sync Word Tolerance (0-7 bits)")
	config.AppFlags.FlagString("packet_format", "", "Packet Format (fixed, variable)")
	config.AppFlags.FlagString("packet_coding", "", "Packet Coding (off, manchester, whitening)")
	config.AppFlags.FlagString("packet_filter", "", "Packet Filtering (off, node, broadcast)")
	config.AppFlags.FlagString("packet_crc", "", "Packet CRC (off, autoclear_off, autoclear_on)")
	config.AppFlags.FlagString("afc_mode", "", "AFC Mode (off, on, autoclear), ")
	config.AppFlags.FlagString("afc_routine", "", "AFC Routine (standard, improved)")
	config.AppFlags.FlagUint("fifo_threshold", 0, "FIFO Threshold (bytes)")
	config.AppFlags.FlagDuration("timeout", 5*time.Second, "FIFO and Payload read timeout")
	config.AppFlags.FlagFloat64("temp_calibration", 0, "Temperature Calibration Offset")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, MainLoop))
}
