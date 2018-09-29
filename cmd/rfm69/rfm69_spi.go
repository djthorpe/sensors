// +build spi,!i2c

package main

import (
	// Modules
	_ "github.com/djthorpe/gopi-hw/sys/spi"
	_ "github.com/djthorpe/sensors/sys/rfm69"
)

const (
	MODULE_NAME = "sensors/rfm69/spi"
)
