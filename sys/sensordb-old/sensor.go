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
	"regexp"
	"strconv"
	"time"
)

var (
	regexp_key = regexp.MustCompile("^([0-9A-Fa-f]+):([0-9A-Fa-f]+)$")
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

func (this *sensor) Timestamp() time.Time {
	if this.TimeSeen.IsZero() == false {
		return this.TimeSeen
	} else if this.TimeCreated.IsZero() == false {
		return this.TimeCreated
	} else {
		return time.Time{}
	}
}

func (this *sensor) Product() uint8 {
	if parts := regexp_key.FindStringSubmatch(this.Key_); len(parts) == 3 {
		if product, err := strconv.ParseUint(parts[1], 16, 32); err == nil && product <= 0xFF {
			return uint8(product)
		}
	}
	// Invalid product so return 0
	return 0
}

func (this *sensor) Sensor() uint32 {
	if parts := regexp_key.FindStringSubmatch(this.Key_); len(parts) == 3 {
		if sensor, err := strconv.ParseUint(parts[2], 16, 64); err == nil && sensor <= 0xFFFFFFFF {
			return uint32(sensor)
		}
	}
	// Invalid sensor so return 0
	return 0
}

func (this *sensor) String() string {
	return fmt.Sprintf("Sensor<%v:%v>{ description='%v' ts=%v }", this.Namespace_, this.Key_, this.Description_, this.Timestamp().Format(time.Kitchen))
}
