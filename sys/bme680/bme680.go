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
	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// I2C Configuration
type BME680_I2C struct {
	// the I2C driver
	I2C gopi.I2C

	// The slave address, usually 0x77 or 0x76
	Slave uint8
}

// SPI Configuration
type BME680_SPI struct {
	// the SPI driver
	SPI gopi.SPI

	// SPI Device speed in Hertz
	Speed uint32
}

// Concrete driver
type bme680 struct {
	spi    gopi.SPI
	i2c    gopi.I2C
	slave  uint8
	chipid uint8
	log    gopi.Logger
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	BME680_I2CSLAVE_DEFAULT uint8  = 0x76
	BME680_SPI_MAXSPEEDHZ   uint32 = 5000
	BME680_CHIPID_DEFAULT   uint8  = 0x61
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACE - GET

// Return ChipID and Version
func (this *bme680) ChipID() uint8 {
	return this.chipid
}
