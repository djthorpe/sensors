// +build i2c,!spi

package main

import (
	// Modules
	_ "github.com/djthorpe/gopi-hw/sys/i2c"
)

const (
	MODULE_NAME = "sensors/bme280/i2c"
)
