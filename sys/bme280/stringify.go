/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package bme280

import (
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *bme280) String() string {
	var bus string
	if this.i2c != nil {
		bus = fmt.Sprintf("%v", this.i2c)
	}
	if this.spi != nil {
		bus = fmt.Sprintf("%v", this.spi)
	}
	return fmt.Sprintf("<sensors.BME280>{ chipid=0x%02X version=0x%02X mode=%v filter=%v t_sb=%v spi3w_en=%v osrs_t=%v osrs_p=%v osrs_h=%v bus=%v calibration=%v }", this.chipid, this.version, this.mode, this.filter, this.t_sb, this.spi3w_en, this.osrs_t, this.osrs_p, this.osrs_h, bus, this.calibration)
}

func (this *calibation) String() string {
	return fmt.Sprintf("<Calibration>{ T1=%v T2=%v T3=%v P1=%v P2=%v P3=%v P4=%v P5=%v P6=%v P7=%v P8=%v P9=%v H1=%v H2=%v H3=%v H4=%v H5=%v H6=%v }", this.T1, this.T2, this.T3, this.P1, this.P2, this.P3, this.P4, this.P5, this.P6, this.P7, this.P8, this.P9, this.H1, this.H2, this.H3, this.H4, this.H5, this.H6)
}
