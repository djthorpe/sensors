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
	PIMOTE_K0         = gopi.GPIOPin(17)
	PIMOTE_K1         = gopi.GPIOPin(22)
	PIMOTE_K2         = gopi.GPIOPin(23)
	PIMOTE_K3         = gopi.GPIOPin(27)
	PIMOTE_MOD_SEL    = gopi.GPIOPin(24)
	PIMOTE_MOD_EN     = gopi.GPIOPin(25)
	PIMOTE_SOCKET_MIN = 1
	PIMOTE_SOCKET_MAX = 4
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config ENER314) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug2("<sensors.energenie.ENER314>Open{ gopi=%v }", config.GPIO)

	this := new(ener314)
	this.gpio = config.GPIO
	this.log = log

	// set output pins low
	for _, pin := range []gopi.GPIOPin{PIMOTE_K0, PIMOTE_K1, PIMOTE_K2, PIMOTE_K3, PIMOTE_MOD_SEL, PIMOTE_MOD_EN} {
		this.gpio.SetPinMode(pin, gopi.GPIO_OUTPUT)
		this.gpio.WritePin(pin, gopi.GPIO_LOW)
	}

	// Return success
	return this, nil
}

func (this *ener314) Close() error {
	this.log.Debug2("<sensors.energenie.ENER314>Close{ }")

	// set output pins low
	for _, pin := range []gopi.GPIOPin{PIMOTE_K0, PIMOTE_K1, PIMOTE_K2, PIMOTE_K3, PIMOTE_MOD_SEL, PIMOTE_MOD_EN} {
		this.gpio.WritePin(pin, gopi.GPIO_LOW)
	}

	this.gpio = nil

	return nil
}
