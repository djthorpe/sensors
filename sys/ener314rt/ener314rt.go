/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package ener314rt

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// STRUCTS

// Configuration
type ENER314RT struct {
	GPIO     gopi.GPIO     // GPIO interface
	Radio    sensors.RFM69 // Radio interface
	PinReset gopi.GPIOPin  // Reset pin
	PinLED1  gopi.GPIOPin  // LED1 (Green, Rx) pin
	PinLED2  gopi.GPIOPin  // LED2 (Red, Tx) pin
}

// ener314rt driver
type ener314rt struct {
	log   gopi.Logger
	gpio  gopi.GPIO
	radio sensors.RFM69
	reset gopi.GPIOPin
	led1  gopi.GPIOPin
	led2  gopi.GPIOPin
	ledrx gopi.GPIOPin
	ledtx gopi.GPIOPin
	mode  sensors.MiHomeMode

	// Locker
	sync.Mutex
}

type LED uint

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS, GLOBAL VARIABLES

const (
	LED_ALL LED = iota
	LED_1
	LED_2
	LED_RX
	LED_TX
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config ENER314RT) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sensors.ener314rt>Open{ reset=%v led1=%v led2=%v }", config.PinReset, config.PinLED1, config.PinLED2)

	if config.GPIO == nil || config.Radio == nil {
		// Fail when either GPIO or Radio is nil
		return nil, gopi.ErrBadParameter
	}

	this := new(ener314rt)
	this.log = log
	this.gpio = config.GPIO
	this.radio = config.Radio
	this.reset = config.PinReset

	// Set LED's
	this.led1 = config.PinLED1
	this.led2 = config.PinLED2
	this.ledrx = config.PinLED1
	this.ledtx = config.PinLED2
	if this.ledtx == gopi.GPIO_PIN_NONE {
		// Where the second LED doesn't exist, make it the first LED
		this.ledtx = this.led1
	} else if this.ledrx == gopi.GPIO_PIN_NONE {
		// Where the first LED doesn't exist, make it the second LED
		this.ledrx = this.led2
	}

	// Set mode to undefined
	this.mode = sensors.MIHOME_MODE_NONE

	// Return success
	return this, nil
}

func (this *ener314rt) Close() error {
	this.log.Debug("<sensors.ener314rt>Close{}")

	// Lock until finished
	this.Lock()
	defer this.Unlock()

	// Free resources
	this.gpio = nil
	this.radio = nil

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *ener314rt) String() string {
	return fmt.Sprintf("<sensors.energenie.MiHome>{ mode=%v reset=%v led1=%v led2=%v ledrx=%v ledtx=%v gpio=%v radio=%v }", this.mode, this.reset, this.led1, this.led2, this.ledrx, this.ledtx, this.gpio, this.radio)
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *ener314rt) ResetRadio() error {
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

	// Set undefined mode
	this.mode = sensors.MIHOME_MODE_NONE

	return nil
}

// SetLED switches an LED on or off
func (this *ener314rt) SetLED(led LED, state gopi.GPIOState) error {
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
	case LED_RX:
		if this.ledrx == gopi.GPIO_PIN_NONE {
			// Allow to silently do nothing where device does have RX indicator
			return nil
		} else {
			this.gpio.SetPinMode(this.ledrx, gopi.GPIO_OUTPUT)
			this.gpio.WritePin(this.ledrx, state)
		}
	case LED_TX:
		if this.ledtx == gopi.GPIO_PIN_NONE {
			// Allow to silently do nothing where device does have RX indicator
			return nil
		} else {
			this.gpio.SetPinMode(this.ledtx, gopi.GPIO_OUTPUT)
			this.gpio.WritePin(this.ledtx, state)
		}
	default:
		return gopi.ErrBadParameter
	}
	return nil
}

// Receive payloads until context is cancelled or timeout
func (this *ener314rt) Receive(ctx context.Context, mode sensors.MiHomeMode, payload chan<- []byte) error {
	this.log.Debug2("<sensors.ener314rt>Receive{ mode=%v }", mode)

	// Check incoming parameters
	if ctx == nil || payload == nil || mode == sensors.MIHOME_MODE_NONE {
		return gopi.ErrBadParameter
	}

	// Lock until finished
	this.Lock()
	defer this.Unlock()

	// Switch into correct mode
	if this.mode != mode {
		this.log.Debug("<sensors.ener314rt>Receive: Switch radio mode=%v", mode)
		if err := this.SetMode(mode); err != nil {
			return err
		} else {
			this.mode = mode
		}
	}

	this.log.Debug("<sensors.ener314rt>Receive: Set RX Mode")

	// Set RX mode
	defer this.radio.SetMode(sensors.RFM_MODE_STDBY)
	if err := this.radio.SetMode(sensors.RFM_MODE_RX); err != nil {
		return err
	} else if err := this.radio.SetSequencer(true); err != nil {
		return err
	}

	// Repeatedly read until context is done
	this.log.Debug("<sensors.ener314rt>Receive: Start Receive Loop")
FOR_LOOP:
	for {
		select {
		case <-ctx.Done():
			break FOR_LOOP
		default:
			this.log.Debug("<sensors.ener314rt>Receive: ReadPayload")
			if data, _, err := this.radio.ReadPayload(ctx); err != nil {
				return err
			} else if data != nil {
				// RX light on
				this.SetLED(LED_RX, gopi.GPIO_HIGH)
				defer this.SetLED(LED_RX, gopi.GPIO_LOW)

				// Emit payload
				payload <- data

				// Clear FIFO
				if err := this.radio.ClearFIFO(); err != nil {
					this.log.Error("<sensors.ener314rt>Receive: ClearFIFO: %v", err)
				}
			}
		}
	}

	// Return context error
	return ctx.Err()
}

// Send a raw payload
func (this *ener314rt) Send(payload []byte, repeat uint, mode sensors.MiHomeMode) error {
	this.log.Debug2("<sensors.ener314rt>Send{ mode=%v payload=%v repeat=%v }", mode, strings.ToUpper(hex.EncodeToString(payload)), repeat)

	// Lock until finished
	this.Lock()
	defer this.Unlock()

	// Check parameters
	if len(payload) == 0 || repeat == 0 {
		return gopi.ErrBadParameter
	}

	// Switch into correct mode
	if this.mode != mode {
		if err := this.SetMode(mode); err != nil {
			return err
		} else {
			this.mode = mode
		}
	}

	// Set TX Mode
	defer this.radio.SetMode(sensors.RFM_MODE_STDBY)
	if err := this.radio.SetMode(sensors.RFM_MODE_TX); err != nil {
		return err
	} else if err := this.radio.SetSequencer(true); err != nil {
		return err
	}

	// TX light on
	this.SetLED(LED_TX, gopi.GPIO_HIGH)
	defer this.SetLED(LED_TX, gopi.GPIO_LOW)

	// Write payload
	if err := this.radio.WritePayload(payload, repeat, 100*time.Millisecond); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *ener314rt) MeasureTemperature(tempoffset float32) (float32, error) {
	this.log.Debug2("<sensors.ener314rt>MeasureTemperature{}")

	// Lock until finished
	this.Lock()
	defer this.Unlock()

	// Need to put into standby mode to measure the temperature
	old_mode := this.radio.Mode()
	if old_mode != sensors.RFM_MODE_STDBY {
		if err := this.radio.SetMode(sensors.RFM_MODE_STDBY); err != nil {
			return 0, err
		}
	}

	// Perform the measurement
	value, err := this.radio.MeasureTemperature(tempoffset)

	// Return to previous mode of operation
	if old_mode != sensors.RFM_MODE_STDBY {
		if err := this.radio.SetMode(old_mode); err != nil {
			return 0, err
		}
	}

	// Return the value and error condition
	return value, err
}

func (this *ener314rt) Mode() sensors.MiHomeMode {
	return this.mode
}

func (this *ener314rt) SetMode(mode sensors.MiHomeMode) error {
	this.log.Debug2("<sensors.ener314rt>SetMode{ mode=%v }", mode)
	switch mode {
	case sensors.MIHOME_MODE_MONITOR:
		if err := this.setFSKMode(); err != nil {
			return err
		}
	case sensors.MIHOME_MODE_CONTROL:
		if err := this.setOOKMode(); err != nil {
			return err
		}
	default:
		return gopi.ErrBadParameter
	}

	// Success
	this.mode = mode
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *ener314rt) setFSKMode() error {

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
	} else if err := this.radio.SetAFCMode(sensors.RFM_AFCMODE_OFF); err != nil {
		return err
	} else if err := this.radio.SetAFCRoutine(sensors.RFM_AFCROUTINE_STANDARD); err != nil {
		return err
	} else if err := this.radio.SetLNA(sensors.RFM_LNA_IMPEDANCE_50, sensors.RFM_LNA_GAIN_AUTO); err != nil {
		return err
	} else if err := this.radio.SetRXFilter(sensors.RFM_RXBW_FREQUENCY_FSK_62P5, sensors.RFM_RXBW_CUTOFF_4); err != nil {
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
	} else if err := this.radio.SetPayloadSize(0x40); err != nil {
		return err
	} else if err := this.radio.SetSyncWord([]byte{0x2D, 0xD4}); err != nil {
		return err
	} else if err := this.radio.SetSyncTolerance(0); err != nil {
		return err
	} else if err := this.radio.SetNodeAddress(0x04); err != nil {
		return err
	} else if err := this.radio.SetBroadcastAddress(0x00); err != nil {
		return err
	} else if err := this.radio.SetAESKey(nil); err != nil {
		return err
	} else if err := this.radio.SetFIFOThreshold(1); err != nil {
		return err
	}

	// Success
	return nil
}

func (this *ener314rt) setOOKMode() error {
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
	} else if err := this.radio.SetPacketFormat(sensors.RFM_PACKET_FORMAT_FIXED); err != nil {
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
	} else if err := this.radio.SetSyncTolerance(0); err != nil {
		return err
	} else if err := this.radio.SetAESKey(nil); err != nil {
		return err
	} else if err := this.radio.SetFIFOThreshold(1); err != nil {
		return err
	}

	// Success
	return nil
}
