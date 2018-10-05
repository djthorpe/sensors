/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"reflect"

	"github.com/djthorpe/gopi"

	// Protocol Buffer Definition
	protobuf "github.com/djthorpe/sensors/cmd/rpc/sensors"
)

//go:generate protoc sensors/sensors.proto --go_out=plugins=grpc:.

////////////////////////////////////////////////////////////////////////////////
// SensorService implementation

type SensorService struct{}

func (this *SensorService) Register(server gopi.RPCServer) error {
	// Check to make sure we satisfy the interface
	var _ proto.SensorService = (*SensorService)(nil)
	return server.Fudge(reflect.ValueOf(protobuf.RegisterSensorService), this)
}
