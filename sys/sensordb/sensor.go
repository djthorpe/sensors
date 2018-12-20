/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package sensordb

import (
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////
// SENSOR

func (this *sensor) Namespace() string {
	return this.Namespace_
}

func (this *sensor) Key() string {
	return this.Key_
}

func (this *sensor) Description() string {
	return this.Description_
}

func (this *sensor) String() string {
	return fmt.Sprintf("%v<%v>{ description='%v' }", this.Namespace_, this.Key_, this.Description_)
}
