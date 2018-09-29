/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package sensors

import (
	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES - BME680 AIR QUALITY SENSOR
// Note this driver is still in development

type BME680 interface {
	gopi.Driver

	// Get ChipID
	ChipID() uint8

	// Reset
	SoftReset() error
}
