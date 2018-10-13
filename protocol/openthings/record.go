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

	"github.com/djthorpe/gopi"

	// Frameworks
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
				r._Request = (v & 0x80) != 0x00
				state = ot_state_length
			}
		case ot_state_length:
			r._Type = sensors.OTDataType((v >> 4) & 0x0F)
			r._Size = v & 0x0F
			state = ot_state_data
			if r._Size > 0 {
				// Sanity check size
				if r._Size > uint8(len(data)-i-2) {
					return nil, sensors.ErrMessageCorruption
				}
				r._Data = make([]byte, 0, r._Size)
				state = ot_state_data
			} else if r._Size == 0 {
				// Zero-sized record
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
// STRINGIFY

func (this *record) String() string {
	name := strings.TrimPrefix(fmt.Sprint(this._Name), "OT_PARAM_")
	typ := strings.TrimPrefix(fmt.Sprint(this._Type), "OT_DATATYPE_")
	req := ""
	if this._Request {
		req = " Req=true"
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
	} else if this._Request != other_._Request {
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

/*
////////////////////////////////////////////////////////////////////////////////
// OTRECORD INTERFACE

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
*/
