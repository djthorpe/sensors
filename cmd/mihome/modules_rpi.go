// +build rpi

package main

import (
	// Modules
	_ "github.com/djthorpe/gopi-hw/sys/gpio"
	_ "github.com/djthorpe/gopi-hw/sys/hw"
	_ "github.com/djthorpe/gopi-hw/sys/metrics"
	_ "github.com/djthorpe/gopi-hw/sys/spi"
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/sensors/protocol/ook"
	_ "github.com/djthorpe/sensors/protocol/openthings"
	_ "github.com/djthorpe/sensors/sys/ener314rt"
	_ "github.com/djthorpe/sensors/sys/mihome"
	_ "github.com/djthorpe/sensors/sys/rfm69"
)
