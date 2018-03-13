/*
   Go Language Raspberry Pi Interface
   (c) Copyright David Thorpe 2016-2018
   All Rights Reserved
   Documentation http://djthorpe.github.io/gopi/
   For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"github.com/djthorpe/sensors"
)

func SetFSKMode(device sensors.RFM69) error {
	if err := device.SetMode(sensors.RFM_MODE_STDBY); err != nil {
		return err
	} else if err := device.SetModulation(sensors.RFM_MODULATION_FSK); err != nil {
		return err
	} else if err := device.SetSequencer(true); err != nil {
		return err
	} else if err := device.SetBitrate(4800); err != nil {
		return err
	} else if err := device.SetFreqCarrier(434300000); err != nil {
		return err
	} else if err := device.SetFreqDeviation(30000); err != nil {
		return err
	} else if err := device.SetAFCMode(sensors.RFM_AFCMODE_ON); err != nil {
		return err
	} else if err := device.SetAFCRoutine(sensors.RFM_AFCROUTINE_STANDARD); err != nil {
		return err
	} else if err := device.SetDataMode(sensors.RFM_DATAMODE_PACKET); err != nil {
		return err
	} else if err := device.SetPacketFormat(sensors.RFM_PACKET_FORMAT_VARIABLE); err != nil {
		return err
	} else if err := device.SetPacketCoding(sensors.RFM_PACKET_CODING_MANCHESTER); err != nil {
		return err
	} else if err := device.SetPacketFilter(sensors.RFM_PACKET_FILTER_NONE); err != nil {
		return err
	} else if err := device.SetPacketCRC(sensors.RFM_PACKET_CRC_OFF); err != nil {
		return err
	} else if err := device.SetPreambleSize(3); err != nil {
		return err
	} else if err := device.SetPayloadSize(66); err != nil {
		return err
	} else if err := device.SetSyncWord([]byte{0xD4, 0x2D}); err != nil {
		return err
	} else if err := device.SetSyncTolerance(3); err != nil {
		return err
	} else if err := device.SetAESKey(nil); err != nil {
		return err
	} else if err := device.SetFIFOThreshold(1); err != nil {
		return err
	}

	// Success
	return nil
}
