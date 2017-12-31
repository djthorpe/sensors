// +build i2c,!spi

/*
   Go Language Raspberry Pi Interface
   (c) Copyright David Thorpe 2016-2017
   All Rights Reserved
   Documentation http://djthorpe.github.io/gopi/
   For Licensing and Usage information, please see LICENSE.md
*/

package bme280

// Declare communicating with the BME280 sensor over the I2C bus
const (
	MODULE_NAME = "sensors/bme280:i2c"
)
