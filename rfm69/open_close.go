/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package rfm69

import (
	"sync"

	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Configuration
type RFM69 struct {
	// the SPI driver
	SPI gopi.SPI

	// Device speed
	Speed uint32
}

// driver
type rfm69 struct {
	spi  gopi.SPI
	log  gopi.Logger
	lock sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	RFM_SPI_MODE    = gopi.SPI_MODE_0
	RFM_SPI_SPEEDHZ = 4000000 // 4MHz
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config RFM69) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sensors.RFM69.Open>{ spi=%v speed=%v }", config.SPI, config.Speed)

	this := new(rfm69)
	this.spi = config.SPI
	this.log = log

	if this.spi == nil {
		return nil, gopi.ErrBadParameter
	}

	// Set SPI mode
	if err := this.spi.SetMode(RFM_SPI_MODE); err != nil {
		return nil, err
	}

	// Set SPI speed
	if config.Speed > 0 {
		if err := this.spi.SetMaxSpeedHz(config.Speed); err != nil {
			return nil, err
		}
	} else {
		if err := this.spi.SetMaxSpeedHz(RFM_SPI_SPEEDHZ); err != nil {
			return nil, err
		}
	}

	// Lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Return success
	return this, nil
}

func (this *rfm69) Close() error {
	this.log.Debug("<sensors.RFM69.Close>{ }")

	// Lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Blank out
	this.spi = nil

	return nil
}
