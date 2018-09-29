// +build i2c,!spi

package main

import (
	_ "github.com/djthorpe/gopi-hw/sys/i2c"
)

const (
	MODULE_NAME = "sensors/bme680/i2c"
)
