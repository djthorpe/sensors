/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package mihome

import (
	// Protocol buffers
	pb "github.com/djthorpe/sensors/rpc/protobuf/mihome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Message struct {
}

////////////////////////////////////////////////////////////////////////////////
// NULL

func toProtobufNullEvent() *pb.Message {
	return &pb.Message{}
}

////////////////////////////////////////////////////////////////////////////////
// MESSAGE

func toProtobufMessage(message *Message) *pb.Message {
	if message == nil {
		return nil
	} else {
		return &pb.Message{}
	}
}

func fromProtobufMessage(pb *pb.Message) *Message {
	if pb == nil {
		return nil
	} else {
		return &Message{}
	}
}
