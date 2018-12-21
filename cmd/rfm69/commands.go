/*
   Go Language Raspberry Pi Interface
   (c) Copyright David Thorpe 2016-2018
   All Rights Reserved
   Documentation http://djthorpe.github.io/gopi/
   For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"
	"github.com/olekukonko/tablewriter"
)

var (
	command_map = map[string]func(app *gopi.AppInstance, device sensors.RFM69) error{
		"TriggerAFC":      TriggerAFC,
		"ClearFIFO":       ClearFIFO,
		"ReadFIFO":        ReadFIFO,
		"ReadPayload":     ReadPayload,
		"WritePayload":    WritePayload,
		"Status":          Status,
		"ReadTemperature": ReadTemperature,
	}
)

/////////////////////////////////////////////////////////////////////
// LIST OF COMMANDS

func ListCommands() []string {
	commands := make([]string, 0)
	for k := range command_map {
		commands = append(commands, k)
	}
	return commands
}

/////////////////////////////////////////////////////////////////////
// COMMANDS

func TriggerAFC(app *gopi.AppInstance, device sensors.RFM69) error {
	return device.TriggerAFC()
}

func ClearFIFO(app *gopi.AppInstance, device sensors.RFM69) error {
	return device.ClearFIFO()
}

func ReadFIFO(app *gopi.AppInstance, device sensors.RFM69) error {

	// Read FIFO
	timeout, _ := app.AppFlags.GetDuration("timeout")
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	if data, err := device.ReadFIFO(ctx); err != nil {
		return err
	} else if data == nil {
		return fmt.Errorf("FIFO Empty")
	} else {
		// Output register information
		table := tablewriter.NewWriter(os.Stdout)

		table.SetHeader([]string{"FIFO", "Value"})
		table.Append([]string{"data", fmt.Sprintf("%v", strings.ToUpper(hex.EncodeToString(data)))})

		table.Render()
	}

	// Success
	return nil
}

func ReadTemperature(app *gopi.AppInstance, device sensors.RFM69) error {
	// Put into Standby mode
	if device.Mode() != sensors.RFM_MODE_STDBY {
		if err := device.SetMode(sensors.RFM_MODE_STDBY); err != nil {
			return err
		}
	}

	// calibration value
	calibration, _ := app.AppFlags.GetFloat64("temp_calibration")
	if value, err := device.MeasureTemperature(float32(calibration)); err != nil {
		return err
	} else {
		// Output register information
		table := tablewriter.NewWriter(os.Stdout)

		table.SetHeader([]string{"Parameter", "Value"})
		table.Append([]string{"temp_calibration", fmt.Sprintf("%vC", calibration)})
		table.Append([]string{"temperature", fmt.Sprintf("%vC", value)})

		table.Render()
	}

	// Success
	return nil
}

func ReadPayload(app *gopi.AppInstance, device sensors.RFM69) error {

	// Put into RX mode
	if device.Mode() != sensors.RFM_MODE_RX {
		if err := device.SetMode(sensors.RFM_MODE_RX); err != nil {
			return err
		}
	}

	// Read payload
	timeout, _ := app.AppFlags.GetDuration("timeout")
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	if data, crc_ok, err := device.ReadPayload(ctx); err != nil {
		return err
	} else if data == nil {
		return fmt.Errorf("Timeout waiting for payload")
	} else {
		// Output register information
		table := tablewriter.NewWriter(os.Stdout)

		table.SetHeader([]string{"Payload", "Value"})
		table.Append([]string{"payload", fmt.Sprintf("%v", strings.ToUpper(hex.EncodeToString(data)))})
		table.Append([]string{"crc_ok", fmt.Sprintf("%v", crc_ok)})

		table.Render()
	}

	// Success
	return nil
}

func WritePayload(app *gopi.AppInstance, device sensors.RFM69) error {

	// Put into TX mode
	if device.Mode() != sensors.RFM_MODE_TX {
		if err := device.SetMode(sensors.RFM_MODE_TX); err != nil {
			return err
		}
	}

	// Get data
	if data, _ := app.AppFlags.GetString("data"); data == "" {
		return fmt.Errorf("Missing -data argument")
	} else if data_, err := hex.DecodeString(data); err != nil {
		return err
	} else if err := device.WritePayload(data_, 1, 0); err != nil {
		return err
	} else {
		// Output register information
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Payload", "Value"})
		table.Append([]string{"payload", fmt.Sprintf("%v", strings.ToUpper(hex.EncodeToString(data_)))})
		table.Render()
	}

	// Success
	return nil
}

func Status(app *gopi.AppInstance, device sensors.RFM69) error {

	// Output register information
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Register", "Value"})

	// Output mode, data mode and modulation
	table.Append([]string{"mode", modeToString(device.Mode())})
	table.Append([]string{"listen", listenOnToString(device.ListenOn())})
	table.Append([]string{"sequencer", sequencerEnabledToString(device.SequencerEnabled())})
	table.Append([]string{"modulation", modulationToString(device.Modulation())})
	table.Append([]string{"bitrate", bitrateToString(device.Bitrate())})
	table.Append([]string{"freq_carrier", freqToString(device.FreqCarrier())})
	table.Append([]string{"freq_dev", freqToString(device.FreqDeviation())})

	// Automatic Frequency Correction
	table.Append([]string{"afc", fmt.Sprintf("%v Hz", device.AFC())})
	table.Append([]string{"afc_mode", fmt.Sprint(device.AFCMode())})
	table.Append([]string{"afc_routine", fmt.Sprint(device.AFCRoutine())})

	// Low Noise Amplifier Settings
	table.Append([]string{"lna_impedance", fmt.Sprintf("%v", device.LNAImpedance())})
	if gain, err := device.LNACurrentGain(); err != nil {
		return err
	} else {
		table.Append([]string{"lna_gain", fmt.Sprintf("%v (%v)", device.LNAGain(), gain)})
	}

	// RX Channel Filyer Parameters
	table.Append([]string{"rxbw_frequency", fmt.Sprintf("%v", device.RXFilterFrequency())})
	table.Append([]string{"rxbw_cutoff", fmt.Sprintf("%v", device.RXFilterCutoff())})

	// Packet parameters
	table.Append([]string{"datamode", dataModeToString(device.DataMode())})

	// Format, coding and filtering
	table.Append([]string{"packet_format", packetFormatString(device.PacketFormat())})
	table.Append([]string{"packet_coding", packetCodingString(device.PacketCoding())})
	table.Append([]string{"packet_filter", packetFilterString(device.PacketFilter())})
	table.Append([]string{"packet_crc", packetCRCString(device.PacketCRC())})

	// Node and Broadcast addresses
	table.Append([]string{"node_addr", fmt.Sprintf("0x%02X", device.NodeAddress())})
	table.Append([]string{"broadcast_addr", fmt.Sprintf("0x%02X", device.BroadcastAddress())})

	// Payload and Preamble sizes
	table.Append([]string{"preamble_size", fmt.Sprintf("%v bytes", device.PreambleSize())})
	table.Append([]string{"payload_size", fmt.Sprintf("%v bytes", device.PayloadSize())})

	// Sync Word
	if device.SyncWord() == nil {
		table.Append([]string{"sync_word", "disabled"})
	} else {
		table.Append([]string{"sync_word", fmt.Sprintf("%v", strings.ToUpper(hex.EncodeToString(device.SyncWord())))})
		table.Append([]string{"sync_tol", fmt.Sprintf("%v bits", device.SyncTolerance())})
	}

	// AES Key (or off)
	if device.AESKey() == nil {
		table.Append([]string{"aes_key", "disabled"})
	} else {
		table.Append([]string{"aes_key", fmt.Sprintf("%v", strings.ToUpper(hex.EncodeToString(device.AESKey())))})
	}

	// FIFO
	table.Append([]string{"fifo_threshold", fmt.Sprintf("%v bytes", device.FIFOThreshold())})

	table.Render()
	return nil
}
