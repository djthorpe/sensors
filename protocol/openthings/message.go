/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package openthings

import (
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

/*
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
*/

func (this *message) Manufacturer() sensors.OTManufacturer {
	return this.manufacturer
}

/*
	if len(this.payload) >= 2 {
		m := sensors.OTManufacturer(this.payload[1])
		if m <= sensors.OT_MANUFACTURER_MAX {
			return m
		}
	}
	return sensors.OT_MANUFACTURER_NONE
}
*/
func (this *message) Product() uint8 {
	return this.product
}

/*
	if len(this.payload) >= 3 {
		return this.payload[2]
	} else {
		return 0
	}
}
*/
func (this *message) Sensor() uint32 {
	return this.sensor
}

/*
func (this *message) CRC() uint16 {
	return this.crc
}

func (this *message) Records() []sensors.OTRecord {
	return this.records
}
*/

func (this *message) IsDuplicate(other sensors.Message) bool {
	if this == other {
		return true
	}
	if other == nil || this.Name() != other.Name() {
		return false
	}
	if other_, ok := other.(*message); ok == false {
		return false
	} else {
		if this.manufacturer != other_.manufacturer {
			return false
		}
		if this.product != other_.product {
			return false
		}
		if this.sensor != other_.sensor {
			return false
		}
	}

	// TODO - records

	// Return success
	return true
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *message) String() string {
	var params []string
	params = append(params, fmt.Sprintf("name='%v'", this.Name()))
	params = append(params, fmt.Sprintf("manufacturer=%v", this.Manufacturer()))
	params = append(params, fmt.Sprintf("product=0x%02X", this.Product()))
	params = append(params, fmt.Sprintf("sensor=0x%02X", this.Sensor()))
	if this.ts.IsZero() == false {
		params = append(params, fmt.Sprintf("ts=%v", this.Timestamp()))
	}
	return fmt.Sprintf("<protocol.openthings.Message>{ %v }", strings.Join(params, " "))
}

////////////////////////////////////////////////////////////////////////////////
// ENCODE PARTS OF THE MESSAGE

func (this *message) encode_header(pip uint16) []byte {
	header := make([]byte, OT_MESSAGE_HEADER_SIZE)
	header[0] = 0 // Length not yet in place
	header[1] = uint8(this.manufacturer) & 0x7F
	header[2] = uint8(this.product)
	header[3] = uint8(pip & 0xFF00 >> 8)
	header[4] = uint8(pip & 0x00FF >> 0)
	header[5] = uint8(this.sensor & 0xFF0000 >> 16)
	header[6] = uint8(this.sensor & 0x00FF00 >> 8)
	header[7] = uint8(this.sensor & 0x0000FF >> 0)
	return header
}

func (this *message) encode_records() []byte {
	// TODO
	return []byte{}
}

func (this *message) encode_footer(crc uint16) []byte {
	// Return the footer
	footer := make([]byte, OT_MESSAGE_FOOTER_SIZE)
	footer[0] = 0 // End of data
	footer[1] = uint8(crc & 0xFF00 >> 8)
	footer[2] = uint8(crc & 0x00FF >> 0)
	return footer
}
