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

func read_records(payload []byte) ([]sensors.OTRecord, error) {
	records := make([]sensors.OTRecord, 0)
	state := ot_state_start
	record := &ot_record{}
	for _, v := range payload {
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
	return records, nil
}

////////////////////////////////////////////////////////////////////////////////
// OTRecord INTERFACE

func (this *ot_record) Name() sensors.OTParameter {
	return this.name
}

func (this *ot_record) String() string {
	return fmt.Sprintf("%v<req=%v t=%v sz=%v data=%v>", this.name, this.request, this.datatype, this.datasize, strings.ToUpper(hex.EncodeToString(this.data)))
}
