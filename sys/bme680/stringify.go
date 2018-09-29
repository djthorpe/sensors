/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package bme680

import (
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *bme680) String() string {
	var bus string
	if this.i2c != nil {
		bus = fmt.Sprintf("%v", this.i2c)
	}
	if this.spi != nil {
		bus = fmt.Sprintf("%v", this.spi)
	}
	return fmt.Sprintf("<sensors.BME680>{ chipid=0x%02X bus=%v }", this.chipid, bus)
}
