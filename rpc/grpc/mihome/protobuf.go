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

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	sensors "github.com/djthorpe/sensors"

	// Protocol buffers
	pb "github.com/djthorpe/sensors/rpc/protobuf/mihome"
	ptypes "github.com/golang/protobuf/ptypes"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type proto_message struct {
	source       gopi.Driver
	ts           time.Time
	protocol     string
	manufacturer uint8
	product      uint8
	sensor       uint32
	parameters   []*proto_parameter
	data         []byte
}

type proto_parameter struct {
	name   sensors.OTParameter
	report bool
	data   []byte
}

type OTRecord interface {
	// Name is the parameter name
	Name() OTParameter

	// Type is the type of data
	Type() OTDataType

	// IsReport returns the report bit for the record
	IsReport() bool

	// Data returns the record encoded as data
	Data() ([]byte, error)

	// BoolValue returns the boolean value, when type is UDEC_0
	BoolValue() (bool, error)

	// StringValue returns the value for all types except FLOAT and ENUM
	StringValue() (string, error)

	// UintValue returns the value for UDEC_0 types
	UintValue() (uint64, error)

	// IntValue returns the value for DEC_0 types
	IntValue() (int64, error)

	// FloatValue returns the value for all UDEC and DEC types
	FloatValue() (float64, error)

	// Compares one record against another and returns true if identical
	IsDuplicate(OTRecord) bool
}

////////////////////////////////////////////////////////////////////////////////
// NULL MESSAGES

func toProtobufNullEvent() *pb.Message {
	// Return an empty message which has no sender
	return &pb.Message{}
}

func isNullProtobufMessage(pb *pb.Message) bool {
	if pb == nil || pb.Sender == nil {
		return true
	} else {
		return false
	}
}

////////////////////////////////////////////////////////////////////////////////
// MESSAGE

func toProtobufMessage(message sensors.Message) *pb.Message {
	if message == nil {
		return nil
	} else if timestamp, err := ptypes.TimestampProto(message.Timestamp()); err != nil {
		return nil
	} else if message_ook, ok := message.(sensors.OOKMessage); ok {
		if product := socketToProduct(message_ook.Socket()); product == sensors.MIHOME_PRODUCT_NONE {
			return nil
		} else if sender := toProtobufSensorKey(message_ook.Name(), sensors.OT_MANUFACTURER_NONE, product, message_ook.Addr()); sender == nil {
			return nil
		} else if parameter := toProtobufBoolParameter(pb.Parameter_SWITCH_STATE, message_ook.State()); parameter == nil {
			return nil
		} else {
			return &pb.Message{
				Timestamp:  timestamp,
				Sender:     sender,
				Parameters: []*pb.Parameter{parameter},
				Data:       message.Data(),
			}
		}
	} else if message_ot, ok := message.(sensors.OTMessage); ok {
		if sender := toProtobufSensorKey(message.Name(), message_ot.Manufacturer(), sensors.MiHomeProduct(message_ot.Product()), message_ot.Sensor()); sender == nil {
			return nil
		} else if parameters := toProtobufParameters(message_ot.Records()); parameters == nil {
			return nil
		} else {
			return &pb.Message{
				Timestamp:  timestamp,
				Sender:     sender,
				Parameters: parameters,
				Data:       message.Data(),
			}
		}
	} else {
		return nil
	}
}

func socketToProduct(socket uint) sensors.MiHomeProduct {
	switch socket {
	case 0:
		return sensors.MIHOME_PRODUCT_CONTROL_ALL
	case 1:
		return sensors.MIHOME_PRODUCT_CONTROL_ONE
	case 2:
		return sensors.MIHOME_PRODUCT_CONTROL_TWO
	case 3:
		return sensors.MIHOME_PRODUCT_CONTROL_THREE
	case 4:
		return sensors.MIHOME_PRODUCT_CONTROL_FOUR
	default:
		return sensors.MIHOME_PRODUCT_NONE
	}
}

func fromProtobufMessage(source gopi.Driver, pb *pb.Message) sensors.ProtoMessage {
	if pb == nil {
		return nil
	} else if ts, err := ptypes.Timestamp(pb.Timestamp); err != nil {
		return nil
	} else if sender := pb.Sender; sender == nil {
		return nil
	} else if parameters := fromProtobufParameters(ph.Parameters); parameters == nil {
		return nil
	} else {
		return &proto_message{
			source:       source,
			ts:           ts,
			protocol:     sender.Protocol,
			manufacturer: uint8(sender.Manufacturer),
			product:      uint8(sender.Product),
			sensor:       sender.Sensor,
			parameters:   parameters,
			data:         pb.Data,
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// PARAMETERS

func toProtobufParameters(records []sensors.OTRecord) []*pb.Parameter {
	parameters := make([]*pb.Parameter, len(records))
	for i, record := range records {
		if data_value, err := record.Data(); err != nil {
			return nil
		} else {
			parameters[i] = &pb.Parameter{
				Name:   pb.Parameter_Name(record.Name()),
				Report: record.IsReport(),
				Data:   data_value,
			}
		}
	}
	return parameters
}

func toProtobufBoolParameter(name pb.Parameter_Name, value bool) *pb.Parameter {
	return &pb.Parameter{
		Name: name,
	}
}

func fromProtobufParameters(pb []*pb.Parameter) []*Parameter {

}

////////////////////////////////////////////////////////////////////////////////
// SENSORKEY

func toProtobufSensorKey(protocol string, manufacturer sensors.OTManufacturer, product sensors.MiHomeProduct, sensor uint32) *pb.SensorKey {
	return &pb.SensorKey{
		Protocol:     protocol,
		Manufacturer: uint32(manufacturer),
		Product:      uint32(product),
		Sensor:       sensor,
	}
}

func fromProtobufSensorKey(key *pb.SensorKey) (string, sensors.OTManufacturer, sensors.MiHomeProduct, uint32, error) {
	if key == nil {
		return "", 0, 0, 0, gopi.ErrBadParameter
	} else {
		return key.Protocol, sensors.OTManufacturer(key.Manufacturer), sensors.MiHomeProduct(key.Product), key.Sensor, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// gopi.Event IMPLEMENTATION

func (this *proto_message) Name() string {
	return "sensors.ProtoMessage"
}

func (this *proto_message) Source() gopi.Driver {
	return this.source
}

////////////////////////////////////////////////////////////////////////////////
// sensors.ProtoMessage IMPLEMENTATION

func (this *proto_message) Timestamp() time.Time {
	return this.ts
}

func (this *proto_message) Manufacturer() uint8 {
	return this.manufacturer
}

func (this *proto_message) Protocol() string {
	return this.protocol
}

func (this *proto_message) Product() uint8 {
	return this.product
}

func (this *proto_message) Sensor() uint32 {
	return this.sensor
}

func (this *proto_message) Data() []byte {
	return this.data
}

func (this *proto_message) IsDuplicate(other sensors.Message) bool {
	if other_, ok := other.(*proto_message); ok == false {
		return false
	} else if other_.Manufacturer() != this.Manufacturer() {
		return false
	} else if other_.Protocol() != this.Protocol() {
		return false
	} else if other_.Product() != this.Product() {
		return false
	} else if other_.Sensor() != this.Sensor() {
		return false
	} else {
		// TODO: Parameters
		return true
	}
}

func (this *proto_message) String() string {
	return fmt.Sprintf("%v{ protocol='%v' manufacturer=0x%02X product=0x%02X sensor=0x%08X data=%v ts=%v }", this.Name(), this.Protocol(), this.Manufacturer(), this.Product(), this.Sensor(), strings.ToUpper(hex.EncodeToString(this.Data())), this.Timestamp().Format(time.Kitchen))
}
