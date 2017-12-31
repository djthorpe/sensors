/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package bme280

import (
	"fmt"

	gopi "github.com/djthorpe/gopi"
	sensors "github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config BME280_I2C) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sensors.BME280.Open>{ slave=0x%02X bus=%v }", config.Slave, config.I2C)

	this := new(bme280)
	this.i2c = config.I2C
	this.log = log
	this.slave = BME280_I2CSLAVE_DEFAULT

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

	// Now perform additional setup
	if err := this.setup(); err != nil {
		return nil, err
	}

	// Return success
	return this, nil
}

func (config BME280_SPI) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sensors.BME280.Open>{ speed=%vHz bus=%v }", config.Speed, config.SPI)

	this := new(bme280)
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
		if err := this.spi.SetMaxSpeedHz(BME280_SPI_MAXSPEEDHZ); err != nil {
			return nil, err
		}
	}

	// Now perform additional setup
	if err := this.setup(); err != nil {
		return nil, err
	}

	// Return success
	return this, nil
}

func (this *bme280) Close() error {
	this.log.Debug2("<sensors.BME280.Close>{ }")

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *bme280) setup() error {
	// Read Chip ID and Version
	if chip_id, version, err := this.readChipVersion(); err != nil {
		return err
	} else if chip_id != BME280_CHIPID_DEFAULT {
		return fmt.Errorf("Unexpected chip_id: 0x%02X (expected 0x%02X)", chip_id, BME280_CHIPID_DEFAULT)
	} else {
		this.chipid = chip_id
		this.version = version
	}

	return this.read_registers()
}

func (this *bme280) read_registers() error {
	// Read calibration values
	if calibration, err := this.readCalibration(); err != nil {
		return err
	} else {
		this.calibration = calibration
	}

	// Read control registers
	if osrs_t, osrs_p, osrs_h, mode, err := this.readControl(); err != nil {
		return err
	} else {
		this.osrs_t = osrs_t
		this.osrs_p = osrs_p
		this.osrs_h = osrs_h
		this.mode = mode
	}

	// Read config registers
	if t_sb, filter, spi3w_en, err := this.readConfig(); err != nil {
		return err
	} else {
		this.t_sb = t_sb
		this.filter = filter
		this.spi3w_en = spi3w_en
	}

	return nil
}
