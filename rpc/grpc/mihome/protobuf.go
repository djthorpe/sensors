/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package mihome

import (

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	sensors "github.com/djthorpe/sensors"

	// Protocol buffers
	pb "github.com/djthorpe/sensors/rpc/protobuf/mihome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Message struct {
	Namespace, Key string
}

////////////////////////////////////////////////////////////////////////////////
// NULL

func toProtobufNullEvent() *pb.Message {
	return &pb.Message{}
}

func (this *Message) IsNullEvent() bool {
	// Null events have empty namespace or key
	return this.Namespace == "" || this.Key == ""
}

////////////////////////////////////////////////////////////////////////////////
// EVENT

func toProtobufEvent(evt gopi.Event) *pb.Message {
	if message_, ok := evt.(sensors.Message); message_ != nil && ok {
		namespace, key := message_.Sender()
		return &pb.Message{
			Namespace: namespace,
			Key:       key,
		}
	} else {
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// MESSAGE

func toProtobufMessage(message *Message) *pb.Message {
	if message == nil {
		return nil
	} else {
		return &pb.Message{
			Namespace: message.Namespace,
			Key:       message.Key,
		}
	}
}

func fromProtobufMessage(pb *pb.Message) *Message {
	if pb == nil {
		return nil
	} else {
		return &Message{
			Namespace: pb.Namespace,
			Key:       pb.Key,
		}
	}
}
