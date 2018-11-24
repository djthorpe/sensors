/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package ads1x15

import (
	"fmt"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	sensors "github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// ADS1015 Configuration
type ADS1015 struct {
	// the I2C driver
	I2C gopi.I2C

	// The slave address
	Slave uint8
}

// ADS1115 Configuration
type ADS1115 struct {
	// the I2C driver
	I2C gopi.I2C

	// The slave address
	Slave uint8
}

type ads1x15 struct {
	i2c     gopi.I2C
	slave   uint8
	product sensors.ADS1X15Product
	log     gopi.Logger
}

type register uint8

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	ADS1x15_SLAVE_DEFAULT      register = 0x48
	ADS1x15_REG_CONVERSION     register = 0x00
	ADS1x15_REG_CONFIG         register = 0x01
	ADS1x15_REG_LOW_THRESHOLD  register = 0x02
	ADS1x15_REG_HIGH_THRESHOLD register = 0x03
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config ADS1015) Open(log gopi.Logger) (gopi.Driver, error) {
	return open(config.I2C, config.Slave, sensors.ADS1X15_PRODUCT_1015, log)
}

func (config ADS1115) Open(log gopi.Logger) (gopi.Driver, error) {
	return open(config.I2C, config.Slave, sensors.ADS1X15_PRODUCT_1115, log)
}

func open(i2c gopi.I2C, slave uint8, product sensors.ADS1X15Product, log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sensors.ADS1X15.Open>{ product=%v slave=0x%02X bus=%v }", product, slave, i2c, log)

	this := new(ads1x15)
	this.product = product
	this.i2c = i2c
	this.log = log
	this.slave = slave

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

	return this, nil
}

func (this *ads1x15) Close() error {
	this.log.Debug("<sensors.ADS1015.Close>{ slave=0x%02X }", this.slave)

	// Release resources
	this.i2c = nil

	// Return success
	return nil
}

func (this *ads1x15) String() string {
	return fmt.Sprintf("<sensors.ADS1X15>{ product=%v slave=0x%02X bus=%v }", this.product, this.slave, this.i2c)
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE

func (this *ads1x15) Product() sensors.ADS1X15Product {
	return this.product
}

////////////////////////////////////////////////////////////////////////////////
// REGISTERS

func (this *ads1x15) read_conversion() (uint16, error) {
	return this.read_uint16(ADS1x15_REG_CONVERSION)
}

func (this *ads1x15) read_config() (uint16, error) {
	return this.read_uint16(ADS1x15_REG_CONFIG)
}

func (this *ads1x15) read_low_threshold() (uint16, error) {
	return this.read_uint16(ADS1x15_REG_LOW_THRESHOLD)
}

func (this *ads1x15) read_high_threshold() (uint16, error) {
	return this.read_uint16(ADS1x15_REG_HIGH_THRESHOLD)
}

func (this *ads1x15) read_uint16(reg register) (uint16, error) {
	if recv, err := this.i2c.ReadUint16(uint8(reg)); err != nil {
		return 0, err
	} else {
		return recv, nil
	}
}
