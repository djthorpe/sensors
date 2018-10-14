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

func (this *message) Manufacturer() sensors.OTManufacturer {
	return this.manufacturer
}

func (this *message) Product() uint8 {
	return this.product
}

func (this *message) Sensor() uint32 {
	return this.sensor
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

		// Check records
		if this_records, other_records := this.Records(), other_.Records(); len(this_records) != len(other_records) {
			return false
		} else {
			// Records should be in the same order
			for i := 0; i < len(this_records); i++ {
				if this_records[i].IsDuplicate(other_records[i]) == false {
					return false
				}
			}
		}
	}

	// Return success
	return true
}

// Append a record
func (this *message) Append(record sensors.OTRecord) {
	if record != nil {
		this.records = append(this.records, record)
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *message) String() string {
	var params []string
	manufacturer := strings.TrimPrefix(fmt.Sprint(this.Manufacturer()), "OT_MANUFACTURER_")
	params = append(params, fmt.Sprintf("name='%v'", this.Name()))
	params = append(params, fmt.Sprintf("manufacturer=%v", manufacturer))
	params = append(params, fmt.Sprintf("product=0x%02X", this.Product()))
	params = append(params, fmt.Sprintf("sensor=0x%05X", this.Sensor()))
	for _, record := range this.records {
		params = append(params, fmt.Sprint(record))
	}
	if this.ts.IsZero() == false {
		params = append(params, fmt.Sprintf("ts=%v", this.Timestamp().Format(time.Kitchen)))
	}
	return fmt.Sprintf("<protocol.openthings.Message>{ %v }", strings.Join(params, " "))
}

////////////////////////////////////////////////////////////////////////////////
// ENCODE PARTS OF THE MESSAGE

func (this *message) encode_header(pip uint16) []byte {
	header := make([]byte, OT_MESSAGE_HEADER_SIZE, 0xFF)
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

func (this *message) encode_records(payload []byte) ([]byte, error) {
	// Append record data
	for _, r := range this.records {
		if data, err := r.Data(); err != nil {
			return nil, err
		} else {
			payload = append(payload, data...)
		}
	}
	// Return payload
	return payload, nil
}

func (this *message) encode_footer(payload []byte, crc uint16) ([]byte, error) {
	return append(payload, 0, uint8(crc&0xFF00>>8), uint8(crc&0x00FF>>0)), nil
}
