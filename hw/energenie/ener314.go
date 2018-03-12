/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package energenie

import (
	"sync"
	"time"

	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// STRUCTS

// ENER314 Configuration
type ENER314 struct {
	// the GPIO interface
	GPIO gopi.GPIO
}

// ENER314 Driver
type ener314 struct {
	log  gopi.Logger
	gpio gopi.GPIO
	lock sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	ENER314_K0         = gopi.GPIOPin(17)
	ENER314_K1         = gopi.GPIOPin(22)
	ENER314_K2         = gopi.GPIOPin(23)
	ENER314_K3         = gopi.GPIOPin(27)
	ENER314_MOD_SEL    = gopi.GPIOPin(24) // (low OOK high FSK)
	ENER314_MOD_EN     = gopi.GPIOPin(25) // (low off high on)
	ENER314_SOCKET_MIN = 1
	ENER314_SOCKET_MAX = 4
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config ENER314) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug2("<sensors.energenie.ENER314>Open{ gopi=%v }", config.GPIO)

	this := new(ener314)
	this.gpio = config.GPIO
	this.log = log

	// set output pins low
	for _, pin := range []gopi.GPIOPin{ENER314_K0, ENER314_K1, ENER314_K2, ENER314_K3, ENER314_MOD_SEL, ENER314_MOD_EN} {
		this.gpio.SetPinMode(pin, gopi.GPIO_OUTPUT)
		this.gpio.WritePin(pin, gopi.GPIO_LOW)
	}

	// Return success
	return this, nil
}

func (this *ener314) Close() error {
	this.log.Debug2("<sensors.energenie.ENER314>Close{ }")

	// set output pins low
	for _, pin := range []gopi.GPIOPin{ENER314_K0, ENER314_K1, ENER314_K2, ENER314_K3, ENER314_MOD_SEL, ENER314_MOD_EN} {
		this.gpio.WritePin(pin, gopi.GPIO_LOW)
	}

	this.gpio = nil

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// ON AND OFF

func (this *ener314) On(sockets ...uint) error {
	this.log.Debug2("<sensors.energenie.ENER314>On{ sockets=%v }", sockets)
	// Check for all sockets
	if len(sockets) == 0 {
		return this.send(0, true)
	}
	// Write socket
	for _, socket := range sockets {
		if socket < ENER314_SOCKET_MIN || socket > ENER314_SOCKET_MAX {
			return gopi.ErrBadParameter
		} else if err := this.send(socket, true); err != nil {
			return err
		}
	}
	return nil
}

func (this *ener314) Off(sockets ...uint) error {
	this.log.Debug2("<sensors.energenie.ENER314>Off{ sockets=%v }", sockets)
	// Check for all sockets
	if len(sockets) == 0 {
		return this.send(0, false)
	}
	// Write socket
	for _, socket := range sockets {
		if socket < ENER314_SOCKET_MIN || socket > ENER314_SOCKET_MAX {
			return gopi.ErrBadParameter
		} else if err := this.send(socket, false); err != nil {
			return err
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// SEND

func (this *ener314) write(reg byte, state bool) {
	if state {
		reg = reg | 8
	}
	// output to pins K0
	if (reg & 0x01) == 0 {
		this.gpio.WritePin(ENER314_K0, gopi.GPIO_LOW)
	} else {
		this.gpio.WritePin(ENER314_K0, gopi.GPIO_HIGH)
	}
	// output to pins K1
	if (reg & 0x02) == 0 {
		this.gpio.WritePin(ENER314_K1, gopi.GPIO_LOW)
	} else {
		this.gpio.WritePin(ENER314_K1, gopi.GPIO_HIGH)
	}
	// output to pins K2
	if (reg & 0x04) == 0 {
		this.gpio.WritePin(ENER314_K2, gopi.GPIO_LOW)
	} else {
		this.gpio.WritePin(ENER314_K2, gopi.GPIO_HIGH)
	}
	// output to pins K3
	if (reg & 0x08) == 0 {
		this.gpio.WritePin(ENER314_K3, gopi.GPIO_LOW)
	} else {
		this.gpio.WritePin(ENER314_K3, gopi.GPIO_HIGH)
	}

	// Let it settle, encoder requires this
	time.Sleep(100 * time.Millisecond)

	// Enable the modulator
	this.gpio.WritePin(ENER314_MOD_EN, gopi.GPIO_HIGH)

	// Keep enabled for a period
	time.Sleep(250 * time.Millisecond)

	// Disable the modulator
	this.gpio.WritePin(ENER314_MOD_EN, gopi.GPIO_LOW)

	// Let it settle
	time.Sleep(100 * time.Millisecond)
}

func (this *ener314) send(socket uint, state bool) error {
	switch {
	case socket == 0:
		this.write(0x3, state)
		break
	case socket == 1:
		this.write(0x7, state)
		break
	case socket == 2:
		this.write(0x6, state)
		break
	case socket == 3:
		this.write(0x5, state)
		break
	case socket == 4:
		this.write(0x4, state)
		break
	default:
		return gopi.ErrBadParameter
	}
	return nil
}
