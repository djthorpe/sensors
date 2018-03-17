// +build spi,!i2c

/*
   Go Language Raspberry Pi Interface
   (c) Copyright David Thorpe 2016-2017
   All Rights Reserved
   Documentation http://djthorpe.github.io/gopi/
   For Licensing and Usage information, please see LICENSE.md
*/

package iface

// Declare communicating with the BME280 sensor over the SPI bus
const (
	MODULE_NAME = "sensors/bme280:spi"
)
