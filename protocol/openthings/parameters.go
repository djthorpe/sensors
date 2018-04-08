/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package openthings

import (
	"encoding/binary"
	"fmt"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"
)

type ot_record struct {
	name     sensors.OTParameter
	request  bool
	datatype sensors.OTDataType
	datasize uint8
	data     []byte
}

const (
	ot_state_start = iota
	ot_state_length
	ot_state_data
)

////////////////////////////////////////////////////////////////////////////////
// READ RECORDS

func read_records(data []byte) ([]sensors.OTRecord, error) {
	records := make([]sensors.OTRecord, 0)
	state := ot_state_start
	record := &ot_record{}
	for _, v := range data {
		switch state {
		case ot_state_start:
			record.name = sensors.OTParameter(v & 0x7F)
			record.request = to_uint8_bool(v & 0x80)
			state = ot_state_length
		case ot_state_length:
			record.datatype = sensors.OTDataType((v >> 4) & 0x0F)
			record.datasize = v & 0x0F
			record.data = make([]byte, 0, record.datasize)
			state = ot_state_data
		case ot_state_data:
			record.data = append(record.data, v)
			if len(record.data) == int(record.datasize) {
				state = ot_state_start
				records = append(records, record)
				record = &ot_record{}
			}
		}
	}
	// Add on the last record
	if record.name != sensors.OT_PARAM_NONE {
		records = append(records, record)
	}
	// Return the records
	return records, nil
}

////////////////////////////////////////////////////////////////////////////////
// OTRECORD INTERFACE

func (this *ot_record) Name() sensors.OTParameter {
	return this.name
}

func (this *ot_record) Type() sensors.OTDataType {
	return this.datatype
}

func (this *ot_record) StringValue() (string, error) {
	switch this.datatype {
	case sensors.OT_DATATYPE_UDEC_0:
		if value, err := this.UIntValue(); err != nil {
			return "", err
		} else {
			return fmt.Sprint(value), nil
		}
	case sensors.OT_DATATYPE_DEC_0:
		if value, err := this.IntValue(); err != nil {
			return "", err
		} else {
			return fmt.Sprint(value), nil
		}
	case sensors.OT_DATATYPE_UDEC_8, sensors.OT_DATATYPE_DEC_8:
		if value, err := this.FloatValue(); err != nil {
			return "", err
		} else {
			return fmt.Sprintf("%.8f", value), nil
		}
	case sensors.OT_DATATYPE_UDEC_4,
		sensors.OT_DATATYPE_UDEC_12, sensors.OT_DATATYPE_UDEC_16,
		sensors.OT_DATATYPE_UDEC_20, sensors.OT_DATATYPE_UDEC_24:
		if value, err := this.FloatValue(); err != nil {
			return "", err
		} else {
			return fmt.Sprint(value), nil
		}
	case sensors.OT_DATATYPE_DEC_16, sensors.OT_DATATYPE_DEC_24:
		if value, err := this.FloatValue(); err != nil {
			return "", err
		} else {
			return fmt.Sprint(value), nil
		}
	default:
		return "", fmt.Errorf("StringValue: Not Implemented: %v", this.datatype)
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *ot_record) String() string {
	if string_value, err := this.StringValue(); err != nil {
		return fmt.Sprintf("%v<req=%v err=%v type=%v>", this.name, this.request, err, this.datatype)
	} else {
		return fmt.Sprintf("%v<req=%v value=%v>", this.name, this.request, string_value)
	}
}

////////////////////////////////////////////////////////////////////////////////
// OTHER METHODS: TODO

// Type OT_DATATYPE_UDEC_0
func (this *ot_record) UIntValue() (uint64, error) {
	// Check data type
	if this.datatype != sensors.OT_DATATYPE_UDEC_0 {
		return 0, gopi.ErrBadParameter
	}
	// Check data length
	if int(this.datasize) != len(this.data) {
		return 0, gopi.ErrOutOfOrder
	}
	// Return uint
	return this.uintValue()
}

// Type OT_DATATYPE_DEC_0
func (this *ot_record) IntValue() (int64, error) {
	// Check data type
	if this.datatype != sensors.OT_DATATYPE_DEC_0 {
		return 0, gopi.ErrBadParameter
	}
	// Check data length
	if int(this.datasize) != len(this.data) {
		return 0, gopi.ErrOutOfOrder
	}
	// Return int
	return this.intValue()
}

// Returns a float value with precision
func (this *ot_record) FloatValue() (float64, error) {
	// Check data length
	if int(this.datasize) != len(this.data) {
		return 0, gopi.ErrOutOfOrder
	}
	// Convert fixed point into floating point
	switch this.datatype {
	case sensors.OT_DATATYPE_UDEC_0:
		value, err := this.uintValue()
		return float64(value), err
	case sensors.OT_DATATYPE_UDEC_8:
		value, err := this.uintValue()
		return float64(value) / float64(1<<8), err
	case sensors.OT_DATATYPE_DEC_8:
		value, err := this.intValue()
		return float64(value) / float64(1<<8), err
	default:
		return 0, gopi.ErrBadParameter
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Returns an unsigned integer for any UDEC of length 1,2,4 or 8 bytes
func (this *ot_record) uintValue() (uint64, error) {
	// Allow 1, 2,4 and 8 byte unsigned integers
	switch len(this.data) {
	case 1:
		return uint64(this.data[0]), nil
	case 2:
		return uint64(binary.BigEndian.Uint16(this.data)), nil
	case 4:
		return uint64(binary.BigEndian.Uint32(this.data)), nil
	case 8:
		return uint64(binary.BigEndian.Uint64(this.data)), nil
	default:
		return 0, gopi.ErrBadParameter
	}
}

// Returns a signed integer for any DEC of length 1,2,4 or 8 bytes
func (this *ot_record) intValue() (int64, error) {
	// Allow 1, 2,4 and 8 byte signed integers
	switch len(this.data) {
	case 1:
		return int64(int8(this.data[0])), nil
	case 2:
		return int64(int16(binary.BigEndian.Uint16(this.data))), nil
	case 4:
		return int64(int32(binary.BigEndian.Uint32(this.data))), nil
	case 8:
		return int64(binary.BigEndian.Uint64(this.data)), nil
	default:
		return 0, gopi.ErrBadParameter
	}
}
