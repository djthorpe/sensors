/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package bme280

import (
	"errors"
)

////////////////////////////////////////////////////////////////////////////////
// VARIABLES

var (
	ErrNoDevice      = errors.New("Device not found")
	ErrSampleSkipped = errors.New("Temperature sampling skipped")
)
