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
// OTRecord INTERFACE

func (this *ot_record) Name() sensors.OTParameter {
	return this.name
}

func (this *ot_record) String() string {
	return fmt.Sprintf("%v<req=%v t=%v value=%v>", this.name, this.request, this.datatype, this.GetFloat())
}

func (this *ot_record) GetFloat() string {
	switch this.datatype {
	case sensors.OT_DATATYPE_UDEC_0:
		if v, err := this.get_uint(); err != nil {
			return fmt.Sprint(err)
		} else {
			return fmt.Sprintf("%v [unsigned 0]", v)
		}
	case sensors.OT_DATATYPE_UDEC_4:
		if v, err := this.get_uint(); err != nil {
			return fmt.Sprint(err)
		} else {
			return fmt.Sprintf("%v [unsigned 4]", v)
		}
	case sensors.OT_DATATYPE_UDEC_8:
		if v, err := this.get_uint(); err != nil {
			return fmt.Sprint(err)
		} else {
			return fmt.Sprintf("%v [unsigned 8]", v)
		}
	case sensors.OT_DATATYPE_UDEC_12:
		if v, err := this.get_uint(); err != nil {
			return fmt.Sprint(err)
		} else {
			return fmt.Sprintf("%v [unsigned 12]", v)
		}
	case sensors.OT_DATATYPE_UDEC_16:
		if v, err := this.get_uint(); err != nil {
			return fmt.Sprint(err)
		} else {
			return fmt.Sprintf("%v [unsigned 16]", v)
		}
	case sensors.OT_DATATYPE_UDEC_20:
		if v, err := this.get_uint(); err != nil {
			return fmt.Sprint(err)
		} else {
			return fmt.Sprintf("%v [unsigned 20]", v)
		}
	case sensors.OT_DATATYPE_UDEC_24:
		if v, err := this.get_uint(); err != nil {
			return fmt.Sprint(err)
		} else {
			return fmt.Sprintf("%v [unsigned 24]", v)
		}
	default:
		return "[?? Invalid float]"
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *ot_record) get_uint() (uint64, error) {
	switch this.datatype {
	case sensors.OT_DATATYPE_UDEC_0:
	case sensors.OT_DATATYPE_UDEC_4:
	case sensors.OT_DATATYPE_UDEC_8:
	case sensors.OT_DATATYPE_UDEC_12:
	case sensors.OT_DATATYPE_UDEC_16:
	case sensors.OT_DATATYPE_UDEC_20:
	case sensors.OT_DATATYPE_UDEC_24:
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
			return 0, gopi.ErrOutOfOrder
		}
	}
	return 0, gopi.ErrOutOfOrder
}

func (this *ot_record) get_int() (int64, error) {
	switch this.datatype {
	case sensors.OT_DATATYPE_DEC_0:
	case sensors.OT_DATATYPE_DEC_8:
	case sensors.OT_DATATYPE_DEC_16:
	case sensors.OT_DATATYPE_DEC_24:
		switch len(this.data) {
		case 1:
			return int64(this.data[0]), nil
		case 2:
			return int64(binary.BigEndian.Uint16(this.data)), nil
		case 4:
			return int64(binary.BigEndian.Uint32(this.data)), nil
		case 8:
			return int64(binary.BigEndian.Uint64(this.data)), nil
		}
	}
	return 0, gopi.ErrOutOfOrder
}
