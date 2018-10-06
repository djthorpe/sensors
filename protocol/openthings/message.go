/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package openthings

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// MESSAGE IMPLEMENTATION

func (this *message) Name() string {
	return this.source.Name()
}

func (this *message) Source() gopi.Driver {
	return this.source
}

func (this *message) Timestamp() time.Time {
	return this.ts
}

func (this *message) Payload() []byte {
	return this.payload
}

func (this *message) Size() uint8 {
	if len(this.payload) > 0 {
		return this.payload[0]
	} else {
		return 0
	}
}

func (this *message) Manufacturer() sensors.OTManufacturer {
	if len(this.payload) >= 2 {
		m := sensors.OTManufacturer(this.payload[1])
		if m <= sensors.OT_MANUFACTURER_MAX {
			return m
		}
	}
	return sensors.OT_MANUFACTURER_NONE
}

func (this *message) ProductID() uint8 {
	if len(this.payload) >= 3 {
		return this.payload[2]
	} else {
		return 0
	}
}

func (this *message) SensorID() uint32 {
	return this.sensor_id
}

func (this *message) CRC() uint16 {
	return this.crc
}

func (this *message) Records() []sensors.OTRecord {
	return this.records
}

func (this *message) IsDuplicate(other sensors.Message) bool {
	if this == other {
		return true
	}
	if other == nil || this.Name() != other.Name() {
		return false
	}
	// TODO
	return false
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *message) String() string {
	var params []string
	if this.Size() > 0 {
		params = append(params, fmt.Sprintf("payload_size=%v", this.Size()))
	}
	if this.Manufacturer() != sensors.OT_MANUFACTURER_NONE {
		params = append(params, fmt.Sprintf("manufacturer=%v", this.Manufacturer()))
		params = append(params, fmt.Sprintf("product_id=0x%02X", this.ProductID()))
	} else {
		params = append(params, fmt.Sprintf("payload=%v", strings.ToUpper(hex.EncodeToString(this.payload))))
	}
	return fmt.Sprintf("<protocol.openthings.Message>{ %v }", strings.Join(params, " "))
}
