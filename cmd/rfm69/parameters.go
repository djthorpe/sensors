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

	"github.com/djthorpe/sensors"
)

var (
	mode_map = map[string]sensors.RFMMode{
		"sleep":   sensors.RFM_MODE_SLEEP,
		"standby": sensors.RFM_MODE_STDBY,
		"fs":      sensors.RFM_MODE_FS,
		"tx":      sensors.RFM_MODE_TX,
		"rx":      sensors.RFM_MODE_RX,
	}
	datamode_map = map[string]sensors.RFMDataMode{
		"packet": sensors.RFM_DATAMODE_PACKET,
		"nosync": sensors.RFM_DATAMODE_CONTINUOUS_NOSYNC,
		"sync":   sensors.RFM_DATAMODE_CONTINUOUS_SYNC,
	}
	modulation_map = map[string]sensors.RFMModulation{
		"fsk":     sensors.RFM_MODULATION_FSK,
		"fsk_1.0": sensors.RFM_MODULATION_FSK_BT_1P0,
		"fsk_0.5": sensors.RFM_MODULATION_FSK_BT_0P5,
		"fsk_0.3": sensors.RFM_MODULATION_FSK_BT_0P3,
		"ook":     sensors.RFM_MODULATION_OOK,
		"ook_br":  sensors.RFM_MODULATION_OOK_BR,
		"ook_2br": sensors.RFM_MODULATION_OOK_2BR,
	}
	packet_format_map = map[string]sensors.RFMPacketFormat{
		"fixed":    sensors.RFM_PACKET_FORMAT_FIXED,
		"variable": sensors.RFM_PACKET_FORMAT_VARIABLE,
	}
	packet_coding_map = map[string]sensors.RFMPacketCoding{
		"off":        sensors.RFM_PACKET_CODING_NONE,
		"manchester": sensors.RFM_PACKET_CODING_MANCHESTER,
		"whitening":  sensors.RFM_PACKET_CODING_WHITENING,
	}
	packet_filter_map = map[string]sensors.RFMPacketFilter{
		"off":       sensors.RFM_PACKET_FILTER_NONE,
		"node":      sensors.RFM_PACKET_FILTER_NODE,
		"broadcast": sensors.RFM_PACKET_FILTER_BROADCAST,
	}
	packet_crc_map = map[string]sensors.RFMPacketCRC{
		"off":           sensors.RFM_PACKET_CRC_OFF,
		"autoclear_off": sensors.RFM_PACKET_CRC_AUTOCLEAR_OFF,
		"autoclear_on":  sensors.RFM_PACKET_CRC_AUTOCLEAR_ON,
	}

	afc_mode_map = map[string]sensors.RFMAFCMode{
		"off":       sensors.RFM_AFCMODE_OFF,
		"on":        sensors.RFM_AFCMODE_ON,
		"autoclear": sensors.RFM_AFCMODE_AUTOCLEAR,
	}

	afc_routine_map = map[string]sensors.RFMAFCRoutine{
		"standard": sensors.RFM_AFCROUTINE_STANDARD,
		"improved": sensors.RFM_AFCROUTINE_IMPROVED,
	}
)

/////////////////////////////////////////////////////////////////////
// DEVICE MODE

func modeToString(value sensors.RFMMode) string {
	for k, v := range mode_map {
		if value == v {
			return k
		}
	}
	return fmt.Sprint(value)
}

func stringToMode(value string) (sensors.RFMMode, error) {
	// Listen mode is actually standby
	if value == "listen" {
		value = "standby"
	}
	if mode, ok := mode_map[value]; ok == false {
		return 0, fmt.Errorf("Invalid mode flag: %v", value)
	} else {
		return mode, nil
	}
}

/////////////////////////////////////////////////////////////////////
// DATA MODE

func dataModeToString(value sensors.RFMDataMode) string {
	for k, v := range datamode_map {
		if value == v {
			return k
		}
	}
	return fmt.Sprint(value)
}

func stringToDataMode(value string) (sensors.RFMDataMode, error) {
	if mode, ok := datamode_map[value]; ok == false {
		return 0, fmt.Errorf("Invalid datamode flag: %v", value)
	} else {
		return mode, nil
	}
}

/////////////////////////////////////////////////////////////////////
// LISTEN_ON AND SEQUENCER_ENABLED

func listenOnToString(listenOn bool) string {
	if listenOn {
		return "on"
	} else {
		return "off"
	}
}

func sequencerEnabledToString(sequencerEnabled bool) string {
	if sequencerEnabled {
		return "enabled"
	} else {
		return "off"
	}
}

/////////////////////////////////////////////////////////////////////
// MODULATION

func modulationToString(value sensors.RFMModulation) string {
	for k, v := range modulation_map {
		if value == v {
			return k
		}
	}
	return fmt.Sprint(value)
}

func stringToModulation(value string) (sensors.RFMModulation, error) {
	if modulation, ok := modulation_map[value]; ok == false {
		return 0, fmt.Errorf("Invalid modulation flag: %v", value)
	} else {
		return modulation, nil
	}
}

/////////////////////////////////////////////////////////////////////
// PACKET FORMAT, CODING, FILTER

func packetFormatString(value sensors.RFMPacketFormat) string {
	for k, v := range packet_format_map {
		if value == v {
			return k
		}
	}
	return fmt.Sprint(value)
}

func packetCodingString(value sensors.RFMPacketCoding) string {
	for k, v := range packet_coding_map {
		if value == v {
			return k
		}
	}
	return fmt.Sprint(value)
}

func packetFilterString(value sensors.RFMPacketFilter) string {
	for k, v := range packet_filter_map {
		if value == v {
			return k
		}
	}
	return fmt.Sprint(value)
}

func packetCRCString(value sensors.RFMPacketCRC) string {
	for k, v := range packet_crc_map {
		if value == v {
			return k
		}
	}
	return fmt.Sprint(value)
}

func stringToPacketFormat(value string) (sensors.RFMPacketFormat, error) {
	if format, ok := packet_format_map[value]; ok == false {
		return 0, fmt.Errorf("Invalid packet_format flag: %v", value)
	} else {
		return format, nil
	}
}

func stringToPacketCoding(value string) (sensors.RFMPacketCoding, error) {
	if coding, ok := packet_coding_map[value]; ok == false {
		return 0, fmt.Errorf("Invalid packet_coding flag: %v", value)
	} else {
		return coding, nil
	}
}

func stringToPacketFilter(value string) (sensors.RFMPacketFilter, error) {
	if filter, ok := packet_filter_map[value]; ok == false {
		return 0, fmt.Errorf("Invalid packet_filter flag: %v", value)
	} else {
		return filter, nil
	}
}

func stringToPacketCRC(value string) (sensors.RFMPacketCRC, error) {
	if crc, ok := packet_crc_map[value]; ok == false {
		return 0, fmt.Errorf("Invalid packet_crc flag: %v", value)
	} else {
		return crc, nil
	}
}

/////////////////////////////////////////////////////////////////////
// BITRATE & FREQUENCY TO STRING

func bitrateToString(bitrate uint) string {
	if bitrate < 1000 {
		return fmt.Sprintf("%v bps", bitrate)
	}
	kbitrate := float64(bitrate) / 1000.0
	return fmt.Sprintf("%v kbps", kbitrate)
}

func freqToString(hertz uint) string {
	if hertz < 1000 {
		return fmt.Sprintf("%v Hz", hertz)
	}
	khertz := float64(hertz) / 1000.0
	if khertz < 1000 {
		return fmt.Sprintf("%v KHz", khertz)
	}
	mhertz := khertz / 1000.0
	return fmt.Sprintf("%v MHz", mhertz)
}

/////////////////////////////////////////////////////////////////////
// AFC

func stringToAFCMode(value string) (sensors.RFMAFCMode, error) {
	if mode, ok := afc_mode_map[value]; ok == false {
		return 0, fmt.Errorf("Invalid afc_mode flag: %v", value)
	} else {
		return mode, nil
	}
}

func stringToAFCRoutine(value string) (sensors.RFMAFCRoutine, error) {
	if routine, ok := afc_routine_map[value]; ok == false {
		return 0, fmt.Errorf("Invalid afc_routine flag: %v", value)
	} else {
		return routine, nil
	}
}
