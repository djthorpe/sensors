/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package energenie

import (
	"fmt"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// STRUCTS

// Configuration
type MiHome struct {
	GPIO     gopi.GPIO     // the GPIO interface
	Radio    sensors.RFM69 // the Radio interface
	PinReset gopi.GPIOPin  // the reset pin
	PinLED1  gopi.GPIOPin  // the LED1 (Green) pin
	PinLED2  gopi.GPIOPin  // the LED2 (Red) pin
}

// mihome driver
type mihome struct {
	log   gopi.Logger
	gpio  gopi.GPIO
	radio sensors.RFM69

	reset gopi.GPIOPin
	led1  gopi.GPIOPin
	led2  gopi.GPIOPin
}

type LED uint

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	LED_ALL LED = iota
	LED_1
	LED_2
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config MiHome) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug2("<sensors.energenie.MiHome>Open{ reset=%v led1=%v led2=%v }", config.PinReset, config.PinLED1, config.PinLED2)

	this := new(mihome)
	this.log = log
	this.gpio = config.GPIO
	this.radio = config.Radio
	this.reset = config.PinReset
	this.led1 = config.PinLED1
	this.led2 = config.PinLED2

	// Return success
	return this, nil
}

func (this *mihome) Close() error {
	this.log.Debug2("<sensors.energenie.MiHome>Close{ }")

	this.gpio = nil
	this.radio = nil

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *mihome) String() string {
	return fmt.Sprintf("<sensors.energenie.MiHome>{ gpio=%v radio=%v reset=%v led1=%v led2=%v }", this.gpio, this.radio, this.reset, this.led1, this.led2)
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *mihome) ResetRadio() error {
	// If reset is not defined, then return not implemented
	if this.reset == gopi.GPIO_PIN_NONE {
		return gopi.ErrNotImplemented
	}

	// Ensure pin is output
	this.gpio.SetPinMode(this.reset, gopi.GPIO_OUTPUT)

	// Turn all LED's on
	if err := this.SetLED(LED_ALL, gopi.GPIO_HIGH); err != nil {
		return err
	}

	// Pull reset high for 100ms and then low for 5ms
	this.gpio.WritePin(this.reset, gopi.GPIO_HIGH)
	time.Sleep(time.Millisecond * 100)
	this.gpio.WritePin(this.reset, gopi.GPIO_LOW)
	time.Sleep(time.Millisecond * 5)

	// Turn all LED's off
	if err := this.SetLED(LED_ALL, gopi.GPIO_LOW); err != nil {
		return err
	}

	return nil
}

func (this *mihome) SetLED(led LED, state gopi.GPIOState) error {
	switch led {
	case LED_ALL:
		if this.led1 != gopi.GPIO_PIN_NONE {
			if err := this.SetLED(LED_1, state); err != nil {
				return err
			}
		}
		if this.led2 != gopi.GPIO_PIN_NONE {
			if err := this.SetLED(LED_2, state); err != nil {
				return err
			}
		}
	case LED_1:
		if this.led1 == gopi.GPIO_PIN_NONE {
			return gopi.ErrNotImplemented
		} else {
			this.gpio.SetPinMode(this.led1, gopi.GPIO_OUTPUT)
			this.gpio.WritePin(this.led1, state)
		}
	case LED_2:
		if this.led2 == gopi.GPIO_PIN_NONE {
			return gopi.ErrNotImplemented
		} else {
			this.gpio.SetPinMode(this.led2, gopi.GPIO_OUTPUT)
			this.gpio.WritePin(this.led2, state)
		}
	default:
		return gopi.ErrBadParameter
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *mihome) setFSKMode() error {
	if err := this.radio.SetMode(sensors.RFM_MODE_STDBY); err != nil {
		return err
	} else if err := this.radio.SetModulation(sensors.RFM_MODULATION_FSK); err != nil {
		return err
	} else if err := this.radio.SetSequencer(true); err != nil {
		return err
	} else if err := this.radio.SetBitrate(4800); err != nil {
		return err
	} else if err := this.radio.SetFreqCarrier(434300000); err != nil {
		return err
	} else if err := this.radio.SetFreqDeviation(30000); err != nil {
		return err
	} else if err := this.radio.SetAFCMode(sensors.RFM_AFCMODE_ON); err != nil {
		return err
	} else if err := this.radio.SetAFCRoutine(sensors.RFM_AFCROUTINE_STANDARD); err != nil {
		return err
	} else if err := this.radio.SetDataMode(sensors.RFM_DATAMODE_PACKET); err != nil {
		return err
	} else if err := this.radio.SetPacketFormat(sensors.RFM_PACKET_FORMAT_VARIABLE); err != nil {
		return err
	} else if err := this.radio.SetPacketCoding(sensors.RFM_PACKET_CODING_MANCHESTER); err != nil {
		return err
	} else if err := this.radio.SetPacketFilter(sensors.RFM_PACKET_FILTER_NONE); err != nil {
		return err
	} else if err := this.radio.SetPacketCRC(sensors.RFM_PACKET_CRC_OFF); err != nil {
		return err
	} else if err := this.radio.SetPreambleSize(3); err != nil {
		return err
	} else if err := this.radio.SetPayloadSize(66); err != nil {
		return err
	} else if err := this.radio.SetSyncWord([]byte{0xD4, 0x2D}); err != nil {
		return err
	} else if err := this.radio.SetSyncTolerance(3); err != nil {
		return err
	} else if err := this.radio.SetAESKey(nil); err != nil {
		return err
	} else if err := this.radio.SetFIFOThreshold(1); err != nil {
		return err
	}

	// Success
	return nil
}

func (this *mihome) setOOKMode() error {
	if err := this.radio.SetMode(sensors.RFM_MODE_STDBY); err != nil {
		return err
	} else if err := this.radio.SetModulation(sensors.RFM_MODULATION_OOK); err != nil {
		return err
	} else if err := this.radio.SetSequencer(true); err != nil {
		return err
	} else if err := this.radio.SetBitrate(4800); err != nil {
		return err
	} else if err := this.radio.SetFreqCarrier(433920000); err != nil {
		return err
	} else if err := this.radio.SetFreqDeviation(0); err != nil {
		return err
	} else if err := this.radio.SetAFCMode(sensors.RFM_AFCMODE_OFF); err != nil {
		return err
	} else if err := this.radio.SetDataMode(sensors.RFM_DATAMODE_PACKET); err != nil {
		return err
	} else if err := this.radio.SetPacketFormat(sensors.RFM_PACKET_FORMAT_VARIABLE); err != nil {
		return err
	} else if err := this.radio.SetPacketCoding(sensors.RFM_PACKET_CODING_NONE); err != nil {
		return err
	} else if err := this.radio.SetPacketFilter(sensors.RFM_PACKET_FILTER_NONE); err != nil {
		return err
	} else if err := this.radio.SetPacketCRC(sensors.RFM_PACKET_CRC_OFF); err != nil {
		return err
	} else if err := this.radio.SetPreambleSize(0); err != nil {
		return err
	} else if err := this.radio.SetPayloadSize(0); err != nil {
		return err
	} else if err := this.radio.SetSyncWord(nil); err != nil {
		return err
	} else if err := this.radio.SetAESKey(nil); err != nil {
		return err
	} else if err := this.radio.SetFIFOThreshold(1); err != nil {
		return err
	}

	// Success
	return nil
}
