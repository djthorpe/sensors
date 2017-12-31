/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package tsl2561

import (
	gopi "github.com/djthorpe/gopi"
	sensors "github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// I2C Configuration
type TSL2561 struct {
	// the I2C driver
	I2C gopi.I2C

	// The slave address, usually 0x77 or 0x76
	Slave uint8
}

type tsl2561 struct {
	i2c   gopi.I2C
	slave uint8
	log   gopi.Logger
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	TSL2561_I2CSLAVE_DEFAULT = 0x39
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config TSL2561) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sensors.TSL2561.Open>{ slave=0x%02X bus=%v }", config.Slave, config.I2C)

	this := new(tsl2561)
	this.i2c = config.I2C
	this.log = log
	this.slave = TSL2561_I2CSLAVE_DEFAULT

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

	// Return success
	return this, nil
}

func (this *tsl2561) Close() error {
	this.log.Debug2("<sensors.TSL2561.Close>{ }")

	return nil
}
