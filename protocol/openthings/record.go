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
	"encoding/hex"
	"fmt"
	"math"
	"strings"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"
)

const (
	ot_state_param = iota
	ot_state_length
	ot_state_data
	ot_state_end
)

////////////////////////////////////////////////////////////////////////////////
// DECRYPT PARAMETERS

func (this *openthings) decode_parameters(data []byte) ([]sensors.OTRecord, error) {
	this.log.Debug2("<protocol.openthings>DecodeParameters{ data=%v }", strings.ToUpper(hex.EncodeToString(data)))
	if data[len(data)-1] != byte(0x00) {
		this.log.Warn("<protocol.openthings>DecodeParameters: Parameters does not end with a zero byte")
		return nil, sensors.ErrMessageCorruption
	}
	state := ot_state_param
	r := &record{}
	parameters := make([]sensors.OTRecord, 0, 2)

	for i, v := range data {
		switch state {
		case ot_state_param:
			if v == 0 {
				state = ot_state_end
			} else {
				r._Name = sensors.OTParameter(v & 0x7F)
				r.report = (v & 0x80) != 0x00
				state = ot_state_length
			}
		case ot_state_length:
			r._Type = sensors.OTDataType((v >> 4) & 0x0F)
			r._Size = v & 0x0F

			// Check for unsupported data types
			if r._Type == sensors.OT_DATATYPE_ENUM || r._Type == sensors.OT_DATATYPE_FLOAT {
				this.log.Warn("<protocol.openthings>DecodeParameters: Type %v is not yet supported", r._Type)
				return nil, gopi.ErrNotImplemented
			}

			// For non-zero data sizes, make the data structure for storing data or else
			// move back into the start-of-record mode
			if r._Size > 0 {
				// Sanity check size
				if r._Size > uint8(len(data)-i-2) {
					return nil, sensors.ErrMessageCorruption
				}
				r._Data = make([]byte, 0, r._Size)
				state = ot_state_data
			} else if r._Size == 0 {
				parameters = append(parameters, r)
				state = ot_state_param
				r = &record{}
			}
		case ot_state_data:
			r._Data = append(r._Data, v)
			if len(r._Data) == int(r._Size) {
				parameters = append(parameters, r)
				state = ot_state_param
				r = &record{}
			}
		default:
			// Unknown state
			return nil, sensors.ErrMessageCorruption
		}
	}
	if state != ot_state_end {
		this.log.Warn("<protocol.openthings>DecodeParameters: Missing records terminator")
		return nil, sensors.ErrMessageCorruption
	}
	return parameters, nil
}

////////////////////////////////////////////////////////////////////////////////
// CREATE RECORDS

func (this *openthings) NewFloat(sensors.OTParameter, sensors.OTDataType, float64, bool) (sensors.OTRecord, error) {
	return nil, gopi.ErrNotImplemented
}

func (this *openthings) NewString(name sensors.OTParameter, value string, report bool) (sensors.OTRecord, error) {
	// Check incoming parameters
	if name == sensors.OT_PARAM_NONE || name > sensors.OT_PARAM_MAX {
		return nil, gopi.ErrBadParameter
	}
	data := []byte(value)
	if len(data) > int(0x0F) {
		return nil, gopi.ErrBadParameter
	}

	// Create the record
	record := new(record)
	record._Name = name
	record._Type = sensors.OT_DATATYPE_STRING
	record._Size = uint8(len(data))
	record._Data = data
	record.report = report

	// Success
	return record, nil
}

func (this *openthings) NewNull(name sensors.OTParameter, report bool) (sensors.OTRecord, error) {
	// Check incoming parameters
	if name == sensors.OT_PARAM_NONE || name > sensors.OT_PARAM_MAX {
		return nil, gopi.ErrBadParameter
	}

	// Create the record
	record := new(record)
	record._Name = name
	record._Type = sensors.OT_DATATYPE_UDEC_0
	record._Size = 0
	record.report = report

	// Success
	return record, nil
}

func (this *openthings) NewInt(name sensors.OTParameter, value int64, report bool) (sensors.OTRecord, error) {
	// Check incoming parameters
	if name == sensors.OT_PARAM_NONE || name > sensors.OT_PARAM_MAX {
		return nil, gopi.ErrBadParameter
	}

	// Create the record
	record := new(record)
	record._Name = name
	record._Type = sensors.OT_DATATYPE_DEC_0

	// Populate data
	if value <= math.MaxInt8 && value >= math.MinInt8 {
		record._Data = make([]byte, 1)
		record._Data[0] = uint8(value)
	} else if value <= math.MaxInt16 && value >= math.MinInt16 {
		record._Data = make([]byte, 2)
		binary.BigEndian.PutUint16(record._Data, uint16(value))
	} else if value <= math.MaxInt32 && value >= math.MinInt32 {
		record._Data = make([]byte, 4)
		binary.BigEndian.PutUint32(record._Data, uint32(value))
	} else {
		record._Data = make([]byte, 8)
		binary.BigEndian.PutUint64(record._Data, uint64(value))
	}
	record._Size = uint8(len(record._Data))
	record.report = report

	fmt.Printf("%16X => %v\n", value, strings.ToUpper(hex.EncodeToString(record._Data)))

	// Success
	return record, nil
}

func (this *openthings) NewUint(name sensors.OTParameter, value uint64, report bool) (sensors.OTRecord, error) {
	// Check incoming parameters
	if name == sensors.OT_PARAM_NONE || name > sensors.OT_PARAM_MAX {
		return nil, gopi.ErrBadParameter
	}

	// Create the record
	record := new(record)
	record._Name = name
	record._Type = sensors.OT_DATATYPE_UDEC_0
	record.report = report

	// Populate data
	if value <= math.MaxUint8 {
		record._Data = make([]byte, 1)
		record._Data[0] = uint8(value)
	} else if value <= math.MaxUint16 {
		record._Data = make([]byte, 2)
		binary.BigEndian.PutUint16(record._Data, uint16(value))
	} else if value <= math.MaxUint32 {
		record._Data = make([]byte, 4)
		binary.BigEndian.PutUint32(record._Data, uint32(value))
	} else {
		record._Data = make([]byte, 8)
		binary.BigEndian.PutUint64(record._Data, uint64(value))
	}
	record._Size = uint8(len(record._Data))
	record.report = report

	// Success
	return record, nil
}

func (this *openthings) NewBool(name sensors.OTParameter, value bool, report bool) (sensors.OTRecord, error) {
	if value {
		return this.NewUint(name, 1, report)
	} else {
		return this.NewUint(name, 0, report)
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *record) String() string {
	name := strings.TrimPrefix(fmt.Sprint(this._Name), "OT_PARAM_")
	typ := strings.TrimPrefix(fmt.Sprint(this._Type), "OT_DATATYPE_")
	req := ""
	if this.report {
		req = " [report]"
	}
	if value, err := this.StringValue(); err == nil {
		return fmt.Sprintf("%v<%v=%v%v>", name, typ, value, req)
	} else {
		return fmt.Sprintf("%v<Type=%v Size=%v Data=%v%v>", name, typ, this._Size, strings.ToUpper(hex.EncodeToString(this._Data)), req)
	}
}

////////////////////////////////////////////////////////////////////////////////
// OTRecord IMPLEMENTATION

func (this *record) Name() sensors.OTParameter {
	return this._Name
}

func (this *record) Type() sensors.OTDataType {
	return this._Type
}

func (this *record) IsReport() bool {
	return this.report
}

func (this *record) IsDuplicate(other sensors.OTRecord) bool {
	if this == other {
		return true
	}
	if other == nil || this.Name() != other.Name() || this.Type() != other.Type() {
		return false
	}
	if other_, ok := other.(*record); ok == false {
		return false
	} else if this._Size != other_._Size {
		return false
	} else if this.report != other_.report {
		return false
	} else {
		if len(this._Data) != len(other_._Data) {
			return false
		}
		for i := range this._Data {
			if this._Data[i] != other_._Data[i] {
				return false
			}
		}
	}
	return true
}

func (this *record) Data() ([]byte, error) {
	// Sanity check the record
	if this._Name == sensors.OT_PARAM_NONE || this._Name > sensors.OT_PARAM_MAX {
		return nil, gopi.ErrBadParameter
	}
	if this._Size > 0x0F {
		return nil, gopi.ErrBadParameter
	}
	if this._Type > sensors.OT_DATATYPE_DEC_24 {
		return nil, gopi.ErrBadParameter
	}

	// Create the encoded data
	encoded := make([]byte, this._Size+2)

	// Zero byte is the request and the name
	encoded[0] = byte(this._Name) & 0x7F
	if this.report {
		encoded[0] |= 0x80
	}

	// First byte is the type and length
	encoded[1] = (byte(this._Type)&0x0F)<<4 | (this._Size & 0x0F)

	// Remaining bytes are the data
	if this._Size > 0 {
		copy(encoded[2:], this._Data[:])
	}

	// Success
	return encoded, nil
}

////////////////////////////////////////////////////////////////////////////////
// OTRecord DECODE VALUES

func (this *record) unsignedDecimalValue() (uint64, error) {
	// Only deal with UDEC values here
	switch this._Type {
	case sensors.OT_DATATYPE_UDEC_0:
		fallthrough
	case sensors.OT_DATATYPE_UDEC_4:
		fallthrough
	case sensors.OT_DATATYPE_UDEC_8:
		fallthrough
	case sensors.OT_DATATYPE_UDEC_12:
		fallthrough
	case sensors.OT_DATATYPE_UDEC_16:
		fallthrough
	case sensors.OT_DATATYPE_UDEC_20:
		fallthrough
	case sensors.OT_DATATYPE_UDEC_24:
		// Check size parameter
		if this._Size == 0 {
			return 0, nil
		}
		if this._Size > 8 {
			return 0, gopi.ErrBadParameter
		}
		// Return unsigned decimal value
		udec := uint64(0)
		for _, v := range this._Data {
			udec = udec<<8 | uint64(v)
		}
		return udec, nil
	default:
		return 0, gopi.ErrNotImplemented
	}
}

func (this *record) signedDecimalValue() (int64, error) {
	switch this._Type {
	case sensors.OT_DATATYPE_DEC_0:
		fallthrough
	case sensors.OT_DATATYPE_DEC_8:
		fallthrough
	case sensors.OT_DATATYPE_DEC_16:
		fallthrough
	case sensors.OT_DATATYPE_DEC_24:
		// Check size parameter
		if this._Size == 0 {
			return 0, nil
		}
		if this._Size > 8 {
			return 0, gopi.ErrBadParameter
		}
		// Create the decimal value
		dec := int64(0)
		sign := false
		for i, v := range this._Data {
			if i == 0 {
				v = v & byte(0x7F)
				sign = v&byte(0x80) == 0x00
			}
			dec = dec<<8 | int64(v)
		}
		if sign {
			return dec, nil
		} else {
			return -dec, nil
		}
	default:
		return 0, gopi.ErrNotImplemented
	}
}

func (this *record) BoolValue() (bool, error) {
	switch this._Type {
	case sensors.OT_DATATYPE_UDEC_0:
		if value, err := this.unsignedDecimalValue(); err != nil {
			return false, err
		} else {
			return value != 0, nil
		}
	case sensors.OT_DATATYPE_DEC_0:
		if value, err := this.signedDecimalValue(); err != nil {
			return false, err
		} else {
			return value != 0, nil
		}
	default:
		return false, gopi.ErrNotImplemented
	}
}

func (this *record) IntValue() (int64, error) {
	switch this._Type {
	case sensors.OT_DATATYPE_DEC_0:
		if value, err := this.signedDecimalValue(); err != nil {
			return 0, err
		} else {
			return value, nil
		}
	default:
		return 0, gopi.ErrNotImplemented
	}
}

func (this *record) UintValue() (uint64, error) {
	switch this._Type {
	case sensors.OT_DATATYPE_UDEC_0:
		if value, err := this.unsignedDecimalValue(); err != nil {
			return 0, err
		} else {
			return value, nil
		}
	default:
		return 0, gopi.ErrNotImplemented
	}
}

func (this *record) FloatValue() (float64, error) {
	switch this._Type {
	case sensors.OT_DATATYPE_UDEC_0:
		if value, err := this.unsignedDecimalValue(); err != nil {
			return 0, err
		} else {
			return float64(value), nil
		}
	case sensors.OT_DATATYPE_UDEC_4:
		if value, err := this.unsignedDecimalValue(); err != nil {
			return 0, err
		} else {
			return float64(value) / float64(1<<4), nil
		}
	case sensors.OT_DATATYPE_UDEC_8:
		if value, err := this.unsignedDecimalValue(); err != nil {
			return 0, err
		} else {
			return float64(value) / float64(1<<8), nil
		}
	case sensors.OT_DATATYPE_UDEC_12:
		if value, err := this.unsignedDecimalValue(); err != nil {
			return 0, err
		} else {
			return float64(value) / float64(1<<12), nil
		}
	case sensors.OT_DATATYPE_UDEC_16:
		if value, err := this.unsignedDecimalValue(); err != nil {
			return 0, err
		} else {
			return float64(value) / float64(1<<16), nil
		}
	case sensors.OT_DATATYPE_UDEC_20:
		if value, err := this.unsignedDecimalValue(); err != nil {
			return 0, err
		} else {
			return float64(value) / float64(1<<20), nil
		}
	case sensors.OT_DATATYPE_UDEC_24:
		if value, err := this.unsignedDecimalValue(); err != nil {
			return 0, err
		} else {
			return float64(value) / float64(1<<24), nil
		}
	case sensors.OT_DATATYPE_DEC_0:
		if value, err := this.signedDecimalValue(); err != nil {
			return 0, err
		} else {
			return float64(value), nil
		}
	case sensors.OT_DATATYPE_DEC_8:
		if value, err := this.signedDecimalValue(); err != nil {
			return 0, err
		} else {
			return float64(value) / float64(1<<8), nil
		}
	case sensors.OT_DATATYPE_DEC_16:
		if value, err := this.signedDecimalValue(); err != nil {
			return 0, err
		} else {
			return float64(value) / float64(1<<16), nil
		}
	case sensors.OT_DATATYPE_DEC_24:
		if value, err := this.signedDecimalValue(); err != nil {
			return 0, err
		} else {
			return float64(value) / float64(1<<24), nil
		}
	default:
		return 0, gopi.ErrNotImplemented
	}

}

func (this *record) StringValue() (string, error) {
	switch this._Type {
	case sensors.OT_DATATYPE_UDEC_0:
		if v, err := this.UintValue(); err != nil {
			return "", err
		} else {
			return fmt.Sprint(v), nil
		}
	case sensors.OT_DATATYPE_UDEC_4:
		fallthrough
	case sensors.OT_DATATYPE_UDEC_8:
		fallthrough
	case sensors.OT_DATATYPE_UDEC_12:
		fallthrough
	case sensors.OT_DATATYPE_UDEC_16:
		fallthrough
	case sensors.OT_DATATYPE_UDEC_20:
		fallthrough
	case sensors.OT_DATATYPE_UDEC_24:
		if v, err := this.FloatValue(); err != nil {
			return "", err
		} else {
			return fmt.Sprint(v), nil
		}
	case sensors.OT_DATATYPE_DEC_0:
		if v, err := this.IntValue(); err != nil {
			return "", err
		} else {
			return fmt.Sprint(v), nil
		}
	case sensors.OT_DATATYPE_DEC_8:
		fallthrough
	case sensors.OT_DATATYPE_DEC_16:
		fallthrough
	case sensors.OT_DATATYPE_DEC_24:
		if v, err := this.FloatValue(); err != nil {
			return "", err
		} else {
			return fmt.Sprint(v), nil
		}
	case sensors.OT_DATATYPE_STRING:
		return string(this._Data), nil // Assume string is in UTF-8
	default:
		return "", gopi.ErrNotImplemented
	}
}
