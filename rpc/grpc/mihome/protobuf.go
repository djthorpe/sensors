/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package mihome

import (
	"fmt"
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
	ts time.Time
}

////////////////////////////////////////////////////////////////////////////////
// NULL

func toProtobufNullEvent() *pb.Message {
	return &pb.Message{}
}

/*
func (this *Message) IsNullEvent() bool {
	// Null events have empty namespace or key
	return this.Namespace == "" || this.Key == ""
}
*/

////////////////////////////////////////////////////////////////////////////////
// EVENT

func toProtobufEvent(evt gopi.Event) *pb.Message {
	if message := evt.(sensors.Message); message == nil {
		return nil
	} else if timestamp, err := ptypes.TimestampProto(message.Timestamp()); err != nil {
		return nil
	} else if message_ook, ok := message.(sensors.OOKMessage); ok {
		if product := socketToProduct(message_ook.Socket()); product == sensors.MIHOME_PRODUCT_NONE {
			return nil
		} else if sender := toProtobufSensorKey(message_ook.Name(), sensors.OT_MANUFACTURER_NONE, product, message_ook.Addr()); sender == nil {
			return nil
		} else {
			return &pb.Message{
				Timestamp: timestamp,
				Sender:    sender,
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
// MESSAGE

func fromProtobufMessage(pb *pb.Message) sensors.ProtoMessage {
	if pb == nil {
		return nil
	} else if ts, err := ptypes.Timestamp(pb.Timestamp); err != nil {
		return nil
	} else {
		return &proto_message{
			ts: ts,
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// PROTO MESSAGE IMPLEMENTATION

func (this *proto_message) Name() string {
	return "protobuf"
}

func (this *proto_message) Timestamp() time.Time {
	return this.ts
}

func (this *proto_message) Manufacturer() uint8 {
	return 0
}

func (this *proto_message) Protocol() string {
	return ""
}

func (this *proto_message) Source() gopi.Driver {
	return nil
}

func (this *proto_message) Product() uint8 {
	return 0
}

func (this *proto_message) Sensor() uint32 {
	return 0
}

func (this *proto_message) IsDuplicate(other sensors.Message) bool {
	if other_, ok := other.(*proto_message); ok == false {
		return false
	} else {
		fmt.Println(other_)
		return true
	}
}
