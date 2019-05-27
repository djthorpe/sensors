/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package mihome

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	gopi "github.com/djthorpe/gopi"
	sensors "github.com/djthorpe/sensors"
	pb "github.com/djthorpe/sensors/rpc/protobuf/mihome"
	ptypes "github.com/golang/protobuf/ptypes"
	duration "github.com/golang/protobuf/ptypes/duration"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type pb_message struct {
	pb   *pb.Message
	conn gopi.RPCClientConn
}

type pb_record struct {
	pb *pb.Parameter
}

////////////////////////////////////////////////////////////////////////////////
// MISC

func toProtoProtocols(protos []sensors.Proto) []string {
	protostr := make([]string, len(protos))
	if protos == nil {
		return nil
	}
	for i, proto := range protos {
		protostr[i] = proto.Name()
	}
	return protostr
}

func fromProtoDuration(proto *duration.Duration) time.Duration {
	if duration, err := ptypes.Duration(proto); err != nil {
		return 0
	} else {
		return duration
	}
}

func fromProtoPowerMode(proto pb.SensorRequestPowerMode_PowerMode) sensors.MiHomePowerMode {
	switch proto {
	case pb.SensorRequestPowerMode_LOW:
		return sensors.MIHOME_POWER_LOW
	case pb.SensorRequestPowerMode_NONE:
		return sensors.MIHOME_POWER_NORMAL
	default:
		return sensors.MIHOME_POWER_NONE
	}
}

func fromProtoValveState(proto pb.SensorRequestValveState_ValveState) sensors.MiHomeValveState {
	switch proto {
	case pb.SensorRequestValveState_CLOSED:
		return sensors.MIHOME_VALVE_STATE_CLOSED
	case pb.SensorRequestValveState_OPEN:
		return sensors.MIHOME_VALVE_STATE_OPEN
	default:
		return sensors.MIHOME_VALVE_STATE_NORMAL
	}
}

func toProtoPowerMode(power_mode sensors.MiHomePowerMode) pb.SensorRequestPowerMode_PowerMode {
	switch power_mode {
	case sensors.MIHOME_POWER_LOW:
		return pb.SensorRequestPowerMode_LOW
	default:
		return pb.SensorRequestPowerMode_NONE
	}
}

////////////////////////////////////////////////////////////////////////////////
// MESSAGES

func toProtoMessage(msg sensors.Message) *pb.Message {
	if msg == nil {
		return nil
	} else if ts, err := ptypes.TimestampProto(msg.Timestamp()); err != nil {
		return nil
	} else if msg_, ok := msg.(sensors.OTMessage); ok {
		return &pb.Message{
			Sender: toProtoSensorKey(msg_.Manufacturer(), sensors.MiHomeProduct(msg_.Product()), msg_.Sensor()),
			Ts:     ts,
			Data:   msg_.Data(),
			Params: toProtoParameterArray(msg_.Records()),
		}
	} else if msg_, ok := msg.(sensors.OOKMessage); ok {
		return &pb.Message{
			Sender: toProtoSensorKeyOOK(msg_.Addr(), msg_.Socket()),
			Ts:     ts,
			Data:   msg_.Data(),
		}
	} else {
		return nil
	}
}

func fromProtoMessage(message *pb.Message, conn gopi.RPCClientConn) sensors.OTMessage {
	return &pb_message{message, conn}
}

////////////////////////////////////////////////////////////////////////////////
// SENSOR KEY

func fromProtobufSensorKey(key *pb.SensorKey) (sensors.OTManufacturer, sensors.MiHomeProduct, uint32, error) {
	if key == nil {
		return 0, 0, 0, gopi.ErrBadParameter
	} else {
		return sensors.OTManufacturer(key.Manufacturer), sensors.MiHomeProduct(key.Product), key.Sensor, nil
	}
}

func toProtoSensorKey(manufacturer sensors.OTManufacturer, product sensors.MiHomeProduct, sensor uint32) *pb.SensorKey {
	return &pb.SensorKey{
		Manufacturer: uint32(manufacturer),
		Product:      uint32(product),
		Sensor:       sensor,
	}
}

func toProtoSensorKeyOOK(addr uint32, socket uint) *pb.SensorKey {
	if product := sensors.SocketProduct(socket); product == sensors.MIHOME_PRODUCT_NONE {
		return nil
	} else {
		return &pb.SensorKey{
			Manufacturer: uint32(sensors.OT_MANUFACTURER_ENERGENIE),
			Product:      uint32(product),
			Sensor:       addr,
		}
	}
}

func toProtoSensorRequest(queue_request bool, manufacturer sensors.OTManufacturer, product sensors.MiHomeProduct, sensor uint32) *pb.SensorRequest {
	return &pb.SensorRequest{
		QueueRequest: queue_request,
		Sensor:       toProtoSensorKey(manufacturer, product, sensor),
	}
}

func toProtoSensorRequestTemperature(queue_request bool, manufacturer sensors.OTManufacturer, product sensors.MiHomeProduct, sensor uint32, temperature float64) *pb.SensorRequestTemperature {
	return &pb.SensorRequestTemperature{
		QueueRequest: queue_request,
		Sensor:       toProtoSensorKey(manufacturer, product, sensor),
		Temperature:  temperature,
	}
}

func toProtoSensorRequestInterval(queue_request bool, manufacturer sensors.OTManufacturer, product sensors.MiHomeProduct, sensor uint32, interval time.Duration) *pb.SensorRequestInterval {
	return &pb.SensorRequestInterval{
		QueueRequest: queue_request,
		Sensor:       toProtoSensorKey(manufacturer, product, sensor),
		Interval:     ptypes.DurationProto(interval),
	}
}

func toProtoSensorRequestValveState(queue_request bool, manufacturer sensors.OTManufacturer, product sensors.MiHomeProduct, sensor uint32, state sensors.MiHomeValveState) *pb.SensorRequestValveState {
	return &pb.SensorRequestValveState{
		QueueRequest: queue_request,
		Sensor:       toProtoSensorKey(manufacturer, product, sensor),
		ValveState:   pb.SensorRequestValveState_ValveState(state),
	}
}

func toProtoSensorRequestPowerMode(queue_request bool, manufacturer sensors.OTManufacturer, product sensors.MiHomeProduct, sensor uint32, mode sensors.MiHomePowerMode) *pb.SensorRequestPowerMode {
	return &pb.SensorRequestPowerMode{
		QueueRequest: queue_request,
		Sensor:       toProtoSensorKey(manufacturer, product, sensor),
		PowerMode:    toProtoPowerMode(mode),
	}
}

////////////////////////////////////////////////////////////////////////////////
// PARAMETERS

func toProtoParameterArray(records []sensors.OTRecord) []*pb.Parameter {
	if records == nil {
		return nil
	}
	params := make([]*pb.Parameter, len(records))
	for i, record := range records {
		params[i] = toProtoParameter(record)
	}
	return params
}

func toProtoParameter(record sensors.OTRecord) *pb.Parameter {
	if record == nil {
		return nil
	} else if data, err := record.Data(); err != nil {
		return nil
	} else {
		param := &pb.Parameter{
			Name:   pb.Parameter_Name(record.Name()),
			Report: record.IsReport(),
			Data:   data,
		}

		switch record.Type() {
		case sensors.OT_DATATYPE_UDEC_0:
			if udec, err := record.UintValue(); err != nil {
				return nil
			} else {
				param.Value = &pb.Parameter_UintValue{
					UintValue: udec,
				}
			}
		case sensors.OT_DATATYPE_UDEC_4, sensors.OT_DATATYPE_UDEC_8, sensors.OT_DATATYPE_UDEC_12, sensors.OT_DATATYPE_UDEC_16, sensors.OT_DATATYPE_UDEC_20, sensors.OT_DATATYPE_UDEC_24:
			if udec, err := record.FloatValue(); err != nil {
				return nil
			} else {
				param.Value = &pb.Parameter_FloatValue{
					FloatValue: udec,
				}
			}
		case sensors.OT_DATATYPE_STRING:
			if str, err := record.StringValue(); err != nil {
				return nil
			} else {
				param.Value = &pb.Parameter_StringValue{
					StringValue: str,
				}
			}
		case sensors.OT_DATATYPE_DEC_0:
			if dec, err := record.IntValue(); err != nil {
				return nil
			} else {
				param.Value = &pb.Parameter_IntValue{
					IntValue: dec,
				}
			}
		case sensors.OT_DATATYPE_DEC_8, sensors.OT_DATATYPE_DEC_16, sensors.OT_DATATYPE_DEC_24:
			if dec, err := record.FloatValue(); err != nil {
				return nil
			} else {
				param.Value = &pb.Parameter_FloatValue{
					FloatValue: dec,
				}
			}
		default:
			return nil
		}
		return param
	}
}

////////////////////////////////////////////////////////////////////////////////
// OPENTHINGS MESSAGE IMPLEMENTATION

func (this *pb_message) Append(...sensors.OTRecord) sensors.OTMessage {
	// NOT IMPLEMENTED
	return this
}

func (this *pb_message) Manufacturer() sensors.OTManufacturer {
	if this.pb == nil {
		return sensors.OT_MANUFACTURER_NONE
	} else {
		return sensors.OTManufacturer(this.pb.Sender.Manufacturer)
	}
}

func (this *pb_message) Product() uint8 {
	if this.pb == nil {
		return uint8(sensors.OT_MANUFACTURER_NONE)
	} else {
		return uint8(this.pb.Sender.Product)
	}
}

func (this *pb_message) Sensor() uint32 {
	if this.pb == nil {
		return 0
	} else {
		return uint32(this.pb.Sender.Sensor)
	}
}

func (this *pb_message) Records() []sensors.OTRecord {
	if this.pb == nil {
		return nil
	} else {
		records := make([]sensors.OTRecord, len(this.pb.Params))
		for i, record := range this.pb.Params {
			records[i] = &pb_record{record}
		}
		return records
	}
}

func (this *pb_message) Data() []byte {
	if this.pb == nil {
		return nil
	} else {
		return this.pb.Data
	}
}

func (this *pb_message) Timestamp() time.Time {
	if this.pb == nil {
		return time.Time{}
	} else if ts, err := ptypes.Timestamp(this.pb.Ts); err != nil {
		return time.Time{}
	} else {
		return ts
	}
}

func (this *pb_message) IsDuplicate(other sensors.Message) bool {
	if this.pb == nil || other == nil {
		return false
	} else if other_, ok := other.(sensors.OTMessage); ok == false {
		return false
	} else if this.Manufacturer() != other_.Manufacturer() {
		return false
	} else if this.Product() != other_.Product() {
		return false
	} else if this.Sensor() != other_.Sensor() {
		return false
	} else if len(this.Records()) != len(other_.Records()) {
		return false
	} else {
		other_records := other_.Records()
		for i, record := range this.Records() {
			if record.IsDuplicate(other_records[i]) == false {
				return false
			}
		}
		return true
	}
}

func (this *pb_message) Name() string {
	return "OTMessage"
}

func (this *pb_message) Source() gopi.Driver {
	return this.conn
}

func (this *pb_message) String() string {
	data := strings.ToUpper(hex.EncodeToString(this.Data()))
	ts := this.Timestamp().Format(time.Kitchen)
	addr := "<nil>"
	if this.conn != nil {
		addr = this.conn.Addr()
	}
	if this.Manufacturer() == sensors.OT_MANUFACTURER_ENERGENIE {
		product := strings.TrimPrefix(fmt.Sprint(sensors.MiHomeProduct(this.Product())), "MIHOME_PRODUCT_")
		return fmt.Sprintf("<sensors.OTMessage>{ manufacturer=\"ENERGENIE\" product=\"%v\" sensor=0x%08X records=%v ts=%v data=%v src=%v }",
			product, this.Sensor(), this.Records(), ts, data, addr)
	} else {
		return fmt.Sprintf("<sensors.OTMessage>{ manufacturer=%v product=0x%02X sensor=0x%08X records=%v ts=%v data=%v src=%v }",
			this.Manufacturer(), this.Product(), this.Sensor(), this.Records(), ts, data, addr)
	}
}

////////////////////////////////////////////////////////////////////////////////
// OPENTHINGS RECORD IMPLEMENTATION

func (this *pb_record) Name() sensors.OTParameter {
	if this.pb == nil {
		return sensors.OT_PARAM_NONE
	} else {
		return sensors.OTParameter(this.pb.Name)
	}
}

func (this *pb_record) Type() sensors.OTDataType {
	return 0
}

func (this *pb_record) IsReport() bool {
	if this.pb == nil {
		return false
	} else {
		return this.pb.Report
	}
}

func (this *pb_record) Data() ([]byte, error) {
	if this.pb == nil {
		return nil, gopi.ErrAppError
	} else {
		return this.pb.Data, nil
	}
}

func (this *pb_record) BoolValue() (bool, error) {
	return false, gopi.ErrAppError
}

func (this *pb_record) StringValue() (string, error) {
	return "", gopi.ErrAppError
}

func (this *pb_record) UintValue() (uint64, error) {
	return 0, gopi.ErrAppError
}

func (this *pb_record) IntValue() (int64, error) {
	return 0, gopi.ErrAppError
}

func (this *pb_record) FloatValue() (float64, error) {
	return 0, gopi.ErrAppError
}

// Compares one record against another and returns true if identical
func (this *pb_record) IsDuplicate(other sensors.OTRecord) bool {
	if this.pb == nil || other == nil {
		return false
	} else if other_data, err := other.Data(); err != nil {
		return false
	} else if this.pb.Data == nil && other_data == nil {
		return true
	} else if len(this.pb.Data) != len(other_data) {
		return false
	} else {
		for i, v := range this.pb.Data {
			if v != other_data[i] {
				return false
			}
		}
		return true
	}
}

func (this *pb_record) String() string {
	if this.pb == nil {
		return "<nil>"
	} else if this.IsReport() {
		return fmt.Sprintf("<%v=%v [report]>", this.Name(), this.Value())
	} else {
		return fmt.Sprintf("<%v=%v>", this.Name(), this.Value())
	}
}

func (this *pb_record) Value() interface{} {
	if this.pb == nil {
		return nil
	}
	switch this.pb.Value.(type) {
	case *pb.Parameter_StringValue:
		return this.pb.GetStringValue()
	case *pb.Parameter_FloatValue:
		return this.pb.GetFloatValue()
	case *pb.Parameter_UintValue:
		return this.pb.GetUintValue()
	case *pb.Parameter_IntValue:
		return this.pb.GetIntValue()
	default:
		return nil
	}
}
