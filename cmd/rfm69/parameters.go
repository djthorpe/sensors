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
