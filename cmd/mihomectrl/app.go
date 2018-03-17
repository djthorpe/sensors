/*
   Go Language Raspberry Pi Interface
   (c) Copyright David Thorpe 2016-2018
   All Rights Reserved
   Documentation http://djthorpe.github.io/gopi/
   For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type App struct {
	log        gopi.Logger
	gpio       gopi.GPIO
	rfm69      sensors.RFM69
	openthings gopi.Driver
	reset      gopi.GPIOPin
	led1       gopi.GPIOPin
	led2       gopi.GPIOPin
	address    []byte // OOK Address (10 bytes)
}

type LED uint8

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS & GLOBAL VARIABLES

const (
	// LED numbers
	LEDALL LED = iota
	LED1
	LED2
)

const (
	OOK_ON_ALL  = 0x0B
	OOK_OFF_ALL = 0x03
	OOK_ON_1    = 0x0F
	OOK_OFF_1   = 0x07
	OOK_ON_2    = 0x0E
	OOK_OFF_2   = 0x06
	OOK_ON_3    = 0x0D
	OOK_OFF_3   = 0x05
	OOK_ON_4    = 0x0C
	OOK_OFF_4   = 0x04
)

var (
	OOK_ZERO     = byte(0x08)
	OOK_ONE      = byte(0x0E)
	OOK_PREAMBLE = []byte{0x80, 0x00, 0x00, 0x00}
)

var (
	ErrInvalidAddress = errors.New("Invalid Address")
)

////////////////////////////////////////////////////////////////////////////////
// CONSTRUCTOR

func ConfigFlags(config gopi.AppConfig) {
	config.AppFlags.FlagUint("gpio.reset", 25, "Reset Pin (Logical)")
	config.AppFlags.FlagUint("gpio.led1", 27, "Green LED Pin (Logical)")
	config.AppFlags.FlagUint("gpio.led2", 22, "Red LED Pin (Logical)")
}

func NewApp(app *gopi.AppInstance) *App {
	this := new(App)
	this.log = app.Logger

	if rfm69 := app.ModuleInstance("sensors/rfm69").(sensors.RFM69); rfm69 == nil {
		this.log.Error("Missing sensors/rfm69 module")
		return nil
	} else if openthings := app.ModuleInstance("protocol/openthings"); openthings == nil {
		this.log.Error("Missing protocol/openthings module")
		return nil
	} else if gpio := app.ModuleInstance("linux/gpio").(gopi.GPIO); gpio == nil {
		this.log.Error("Missing linux/gpio module")
		return nil
	} else {
		this.gpio = gpio
		this.rfm69 = rfm69
		this.openthings = openthings
	}

	// Obtain the RESET, LED1 and LED2 logical pins
	this.reset = gopi.GPIO_PIN_NONE
	this.led1 = gopi.GPIO_PIN_NONE
	this.led2 = gopi.GPIO_PIN_NONE
	if reset, _ := app.AppFlags.GetUint("gpio.reset"); reset > 0 {
		this.reset = gopi.GPIOPin(reset)
	}
	if led1, _ := app.AppFlags.GetUint("gpio.led1"); led1 > 0 {
		this.led1 = gopi.GPIOPin(led1)
	}
	if led2, _ := app.AppFlags.GetUint("gpio.led2"); led2 > 0 {
		this.led2 = gopi.GPIOPin(led2)
	}

	return this
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *App) ResetRadio() error {
	// If reset is not defined, then return not implemented
	if this.reset == gopi.GPIO_PIN_NONE {
		return gopi.ErrNotImplemented
	}

	// Ensure pin is output
	this.gpio.SetPinMode(this.reset, gopi.GPIO_OUTPUT)

	// Turn all LED's on
	if err := this.SetLED(LEDALL, gopi.GPIO_HIGH); err != nil {
		return err
	}

	// Pull reset high for 100ms and then low for 5ms
	this.gpio.WritePin(this.reset, gopi.GPIO_HIGH)
	time.Sleep(time.Millisecond * 100)
	this.gpio.WritePin(this.reset, gopi.GPIO_LOW)
	time.Sleep(time.Millisecond * 5)

	// Turn all LED's off
	if err := this.SetLED(LEDALL, gopi.GPIO_LOW); err != nil {
		return err
	}

	return nil
}

func (this *App) SetLED(led LED, state gopi.GPIOState) error {
	switch led {
	case LEDALL:
		if this.led1 != gopi.GPIO_PIN_NONE {
			if err := this.SetLED(LED1, state); err != nil {
				return err
			}
		}
		if this.led2 != gopi.GPIO_PIN_NONE {
			if err := this.SetLED(LED2, state); err != nil {
				return err
			}
		}
	case LED1:
		if this.led1 == gopi.GPIO_PIN_NONE {
			return gopi.ErrNotImplemented
		} else {
			this.gpio.SetPinMode(this.led1, gopi.GPIO_OUTPUT)
			this.gpio.WritePin(this.led1, state)
		}
	case LED2:
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

func (this *App) SetFSKMode() error {
	if err := this.rfm69.SetMode(sensors.RFM_MODE_STDBY); err != nil {
		return err
	} else if err := this.rfm69.SetModulation(sensors.RFM_MODULATION_FSK); err != nil {
		return err
	} else if err := this.rfm69.SetSequencer(true); err != nil {
		return err
	} else if err := this.rfm69.SetBitrate(4800); err != nil {
		return err
	} else if err := this.rfm69.SetFreqCarrier(434300000); err != nil {
		return err
	} else if err := this.rfm69.SetFreqDeviation(30000); err != nil {
		return err
	} else if err := this.rfm69.SetAFCMode(sensors.RFM_AFCMODE_ON); err != nil {
		return err
	} else if err := this.rfm69.SetAFCRoutine(sensors.RFM_AFCROUTINE_STANDARD); err != nil {
		return err
	} else if err := this.rfm69.SetDataMode(sensors.RFM_DATAMODE_PACKET); err != nil {
		return err
	} else if err := this.rfm69.SetPacketFormat(sensors.RFM_PACKET_FORMAT_VARIABLE); err != nil {
		return err
	} else if err := this.rfm69.SetPacketCoding(sensors.RFM_PACKET_CODING_MANCHESTER); err != nil {
		return err
	} else if err := this.rfm69.SetPacketFilter(sensors.RFM_PACKET_FILTER_NONE); err != nil {
		return err
	} else if err := this.rfm69.SetPacketCRC(sensors.RFM_PACKET_CRC_OFF); err != nil {
		return err
	} else if err := this.rfm69.SetPreambleSize(3); err != nil {
		return err
	} else if err := this.rfm69.SetPayloadSize(66); err != nil {
		return err
	} else if err := this.rfm69.SetSyncWord([]byte{0xD4, 0x2D}); err != nil {
		return err
	} else if err := this.rfm69.SetSyncTolerance(3); err != nil {
		return err
	} else if err := this.rfm69.SetAESKey(nil); err != nil {
		return err
	} else if err := this.rfm69.SetFIFOThreshold(1); err != nil {
		return err
	}

	// Success
	return nil
}

func (this *App) SetOOKMode() error {
	if err := this.rfm69.SetMode(sensors.RFM_MODE_STDBY); err != nil {
		return err
	} else if err := this.rfm69.SetModulation(sensors.RFM_MODULATION_OOK); err != nil {
		return err
	} else if err := this.rfm69.SetSequencer(true); err != nil {
		return err
	} else if err := this.rfm69.SetBitrate(4800); err != nil {
		return err
	} else if err := this.rfm69.SetFreqCarrier(433920000); err != nil {
		return err
	} else if err := this.rfm69.SetFreqDeviation(0); err != nil {
		return err
	} else if err := this.rfm69.SetAFCMode(sensors.RFM_AFCMODE_OFF); err != nil {
		return err
	} else if err := this.rfm69.SetDataMode(sensors.RFM_DATAMODE_PACKET); err != nil {
		return err
	} else if err := this.rfm69.SetPacketFormat(sensors.RFM_PACKET_FORMAT_VARIABLE); err != nil {
		return err
	} else if err := this.rfm69.SetPacketCoding(sensors.RFM_PACKET_CODING_NONE); err != nil {
		return err
	} else if err := this.rfm69.SetPacketFilter(sensors.RFM_PACKET_FILTER_NONE); err != nil {
		return err
	} else if err := this.rfm69.SetPacketCRC(sensors.RFM_PACKET_CRC_OFF); err != nil {
		return err
	} else if err := this.rfm69.SetPreambleSize(0); err != nil {
		return err
	} else if err := this.rfm69.SetPayloadSize(0); err != nil {
		return err
	} else if err := this.rfm69.SetSyncWord(nil); err != nil {
		return err
	} else if err := this.rfm69.SetAESKey(nil); err != nil {
		return err
	} else if err := this.rfm69.SetFIFOThreshold(1); err != nil {
		return err
	}

	// Success
	return nil
}

func (this *App) SetAddress(addr uint32) error {
	this.log.Debug("SetAddress{ addr=%05X }", addr)
	// Check to make sure it's 20 bits
	if addr&0x0FFFFF != addr {
		return ErrInvalidAddress
	}
	// create new byte array, 10 bytes
	this.address = make([]byte, 10)
	// Iterate through bits
	for i := 9; i >= 0; i-- {
		by := byte(0)
		for j := 0; j < 2; j++ {
			by >>= 4
			if (addr & 1) == 0 {
				by |= (OOK_ZERO << 4)
			} else {
				by |= (OOK_ONE << 4)
			}
			addr >>= 1
		}
		this.address[i] = by
	}
	return nil
}

func (this *App) EncodeByte(value byte) []byte {
	this.log.Debug("EncodeByte{ value=%02X }", value)
	// A byte is encoded as 4 bytes
	encoded := make([]byte, 4)
	for i := 0; i < 4; i++ {
		by := byte(0)
		for j := 0; j < 2; j++ {
			by <<= 4
			if (value & 1) == 0 {
				by |= OOK_ZERO
			} else {
				by |= OOK_ONE
			}
			value >>= 1
		}
		encoded[i] = by
	}
	return encoded
}

func (this *App) OOKPayload(cmd byte) ([]byte, error) {
	// Check to make sure address has been set
	if this.address == nil {
		return nil, ErrInvalidAddress
	}
	if encoded := this.EncodeByte(cmd); len(encoded) != 4 {
		return nil, gopi.ErrAppError
	} else {
		this.log.Debug("cmd %X => %v", cmd, strings.ToUpper(hex.EncodeToString(encoded)))

		// The payload is 16 bytes (preamble, address, command)
		payload := make([]byte, 0, 16)
		payload = append(payload, OOK_PREAMBLE...)
		payload = append(payload, this.address...)
		payload = append(payload, encoded[0], encoded[1])
		return payload, nil
	}
}

func (this *App) WritePayload(data []byte, repeat uint) error {
	// Set radio into TX Mode
	old_mode := this.rfm69.Mode()
	if err := this.rfm69.SetMode(sensors.RFM_MODE_TX); err != nil {
		return err
	} else if err := this.rfm69.SetSequencer(true); err != nil {
		return err
	}

	// Send data
	send_err := this.rfm69.WritePayload(data, repeat)

	// Return radio to old mode
	if err := this.rfm69.SetMode(old_mode); err != nil {
		return err
	}

	// Return error on sending (or nil on success)
	return send_err
}

func (this *App) SendOOK(command byte) error {
	if payload, err := this.OOKPayload(command); err != nil {
		return err
	} else if err := this.WritePayload(payload, 16); err != nil {
		return err
	}
	return nil
}
