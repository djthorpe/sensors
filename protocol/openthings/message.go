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

func (this *Message) Name() string {
	return "OTMessage"
}

func (this *Message) Source() gopi.Driver {
	return nil
}

func (this *Message) Timestamp() time.Time {
	return time.Time{}
}

func (this *Message) Payload() []byte {
	return this.payload
}

func (this *Message) Size() uint8 {
	if len(this.payload) > 0 {
		return this.payload[0]
	} else {
		return 0
	}
}

func (this *Message) Manufacturer() sensors.OTManufacturer {
	if len(this.payload) >= 2 {
		m := sensors.OTManufacturer(this.payload[1])
		if m <= sensors.OT_MANUFACTURER_MAX {
			return m
		}
	}
	return sensors.OT_MANUFACTURER_NONE
}

func (this *Message) ProductID() uint8 {
	if len(this.payload) >= 3 {
		return this.payload[2]
	} else {
		return 0
	}
}

func (this *Message) SensorID() uint32 {
	return this.sensor_id
}

func (this *Message) CRC() uint16 {
	return this.crc
}

func (this *Message) Records() []sensors.OTRecord {
	return this.records
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Message) String() string {
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
