// +build spi,!i2c

package main

import (
	// Modules
	_ "github.com/djthorpe/gopi-hw/sys/spi"
)

const (
	MODULE_NAME = "sensors/bme680/spi"
)
