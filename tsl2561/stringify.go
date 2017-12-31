/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package tsl2561

import (
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *tsl2561) String() string {
	return fmt.Sprintf("<sensors.TSL2561>{ chipid=0x%02X version=0x%02X package_type=%v integrate_time=%v gain=%v bus=%v }", this.chipid, this.version, this.package_type, this.integrate_time, this.gain, this.i2c)
}
