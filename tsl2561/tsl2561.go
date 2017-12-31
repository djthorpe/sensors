/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package tsl2561

import (
	"errors"

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

	chipid         uint8
	version        uint8
	integrate_time sensors.TSL2561IntegrateTime
	gain           sensors.TSL2561Gain
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

	// Chip and version
	if chip_id, revision, err := this.readChipVersion(); err != nil {
		return nil, err
	} else {
		this.chipid = chip_id
		this.version = revision
	}

	// Obtain gain and integrate_time
	if gain, integrate_time, err := this.readTiming(); err != nil {
		return nil, err
	} else {
		this.gain = gain
		this.integrate_time = integrate_time
	}

	// Return success
	return this, nil
}

func (this *tsl2561) Close() error {
	this.log.Debug2("<sensors.TSL2561.Close>{ }")

	// Zero out fields
	this.i2c = nil

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE - GET

// Return ChipID and Version
func (this *tsl2561) ChipIDVersion() (uint8, uint8) {
	return this.chipid, this.version
}

// Return Gain
func (this *tsl2561) Gain() sensors.TSL2561Gain {
	return this.gain
}

// Return Integrate Time
func (this *tsl2561) IntegrateTime() sensors.TSL2561IntegrateTime {
	return this.integrate_time
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE - SET

func (this *tsl2561) SetGain(value sensors.TSL2561Gain) error {
	if err := this.writeTiming(value, this.integrate_time); err != nil {
		return err
	} else if gain_read, _, err := this.readTiming(); err != nil {
		return err
	} else if gain_read != value {
		return errors.New("Unexpected gain value %v, expected %v", gain_read, value)
	} else {
		this.gain = gain_read
		return nil
	}
}

func (this *tsl2561) SetIntegrateTime(value sensors.TSL2561IntegrateTime) error {
	if err := this.writeTiming(this.gain, value); err != nil {
		return err
	} else if _, integrate_time_read, err := this.readTiming(); err != nil {
		return err
	} else if integrate_time_read != value {
		return errors.New("Unexpected integrate_time value %v, expected %v", integrate_time_read, value)
	} else {
		this.integrate_time = integrate_time_read
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE - METHODS

func (this *tsl2561) PowerOn() error {
	return this.powerOn()
}

func (this *tsl2561) PowerOff() error {
	return this.powerOff()
}
