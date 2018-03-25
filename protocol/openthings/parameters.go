/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package openthings

const (
	ot_state_start = iota
	ot_state_length
	ot_state_data
)

////////////////////////////////////////////////////////////////////////////////
// READ RECORDS

func read_records(payload []byte) ([]*OTRecord, error) {
	records := make([]*OTRecord, 0)
	state := ot_state_start
	record := &OTRecord{}
	for _, v := range payload {
		switch state {
		case ot_state_start:
			record.parameter = OTParameter(v & 0x7F)
			record.request = to_uint8_bool(v & 0x80)
			state = ot_state_length
		case ot_state_length:
			record.datatype = OTDataType((v >> 4) & 0x0F)
			record.datasize = v & 0x0F
			record.data = make([]byte, 0, record.datasize)
			state = ot_state_data
		case ot_state_data:
			record.data = append(record.data, v)
			if len(record.data) == int(record.datasize) {
				state = ot_state_start
				records = append(records, record)
				record = &OTRecord{}
			}
		}
	}
	return records, nil
}
