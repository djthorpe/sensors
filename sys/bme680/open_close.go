/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package bme680

import (
	// Frameworks
	"fmt"

	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config BME680_I2C) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sensors.BME680.Open>{ slave=0x%02X bus=%v }", config.Slave, config.I2C)

	this := new(bme680)
	this.i2c = config.I2C
	this.log = log
	this.slave = BME680_I2CSLAVE_DEFAULT

	if config.Slave != 0 {
		this.slave = config.Slave
	}

	if this.i2c == nil {
		return nil, gopi.ErrBadParameter
	}

	// Detect slave
	if detected, err := this.i2c.DetectSlave(this.slave); err != nil {
		return nil, err
	} else if detected == false {
		return nil, sensors.ErrNoDevice
	}

	// Set slave
	if err := this.i2c.SetSlave(this.slave); err != nil {
		return nil, err
	}

	// Call setup
	if err := this.setup(); err != nil {
		return nil, err
	}

	// Return success
	return this, nil
}

func (config BME680_SPI) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sensors.BME680.Open>{ speed=%vHz bus=%v }", config.Speed, config.SPI)

	this := new(bme680)
	this.log = log

	if config.SPI != nil {
		this.spi = config.SPI
	} else {
		return nil, gopi.ErrBadParameter
	}

	// Set SPI bus speed
	if config.Speed != 0 {
		if err := this.spi.SetMaxSpeedHz(config.Speed); err != nil {
			return nil, err
		}
	} else {
		if err := this.spi.SetMaxSpeedHz(BME680_SPI_MAXSPEEDHZ); err != nil {
			return nil, err
		}
	}

	// Call setup
	if err := this.setup(); err != nil {
		return nil, err
	}

	// Return success
	return this, nil
}

func (this *bme680) Close() error {
	this.log.Debug2("<sensors.BME680.Close>{ }")

	// Zero out fields
	this.i2c = nil
	this.spi = nil

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *bme680) setup() error {
	// Read Chip ID and Version
	if chip_id, err := this.readChipId(); err != nil {
		return err
	} else if chip_id != BME680_CHIPID_DEFAULT {
		return fmt.Errorf("Unexpected chip_id: 0x%02X (expected 0x%02X)", chip_id, BME680_CHIPID_DEFAULT)
	} else {
		this.chipid = chip_id
	}

	// Return sucess
	return nil
}
