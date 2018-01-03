/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package rfm69

import (
	gopi "github.com/djthorpe/gopi"
	sensors "github.com/djthorpe/sensors"
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

	// Get version - and check against expected value
	if version, err := this.getVersion(); err != nil {
		return nil, sensors.ErrNoDevice
	} else if version != RFM_VERSION_VALUE {
		return nil, sensors.ErrNoDevice
	} else {
		this.version = version
	}

	// Get operational mode
	if mode, listen_on, sequencer_off, err := this.getOpMode(); err != nil {
		return nil, err
	} else {
		this.mode = mode
		this.listen_on = listen_on
		this.sequencer_off = sequencer_off
	}

	// Get data mode and modulation
	if data_mode, modulation, err := this.getDataModul(); err != nil {
		return nil, err
	} else {
		this.data_mode = data_mode
		this.modulation = modulation
	}

	// Automatic frequency correction
	if afc, err := this.getAfc(); err != nil {
		return nil, err
	} else {
		this.afc = afc
	}

	// Return success
	return this, nil
}

func (this *rfm69) Close() error {
	this.log.Debug("<sensors.RFM69.Close>{ }")

	// Lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Blank out SPI value
	this.spi = nil

	return nil
}
