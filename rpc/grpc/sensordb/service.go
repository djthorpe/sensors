/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package sensordb

import (
	"context"
	"fmt"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi-rpc/sys/grpc"
	"github.com/djthorpe/gopi/util/event"
	"github.com/djthorpe/sensors"

	// Protocol buffers
	pb "github.com/djthorpe/sensors/rpc/protobuf/sensordb"
	empty "github.com/golang/protobuf/ptypes/empty"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Service struct {
	Server   gopi.RPCServer
	Database sensors.Database
}

type service struct {
	log      gopi.Logger
	database sensors.Database

	// Emit events
	event.Publisher
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open the server
func (config Service) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<grpc.service.sensordb.Open>{ server=%v database=%v }", config.Server, config.Database)

	// Check for bad input parameters
	if config.Server == nil || config.Database == nil {
		return nil, gopi.ErrBadParameter
	}

	this := new(service)
	this.log = log
	this.database = config.Database

	// Register service with GRPC server
	pb.RegisterSensorDBServer(config.Server.(grpc.GRPCServer).GRPCServer(), this)

	// Success
	return this, nil
}

func (this *service) Close() error {
	this.log.Debug("<grpc.service.sensordb.Close>{}")

	// Close publisher
	this.Publisher.Close()

	// Release resources
	this.database = nil

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *service) String() string {
	return fmt.Sprintf("<grpc.service.sensordb>{ database=%v }", this.database)
}

////////////////////////////////////////////////////////////////////////////////
// CANCEL STREAMING REQUESTS

func (this *service) CancelRequests() error {
	this.log.Debug2("<grpc.service.sensordb.CancelRequests>{}")

	// Cancel any streaming requests
	this.Publisher.Emit(event.NullEvent)

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RPC METHODS

// Ping returns an empty response
func (this *service) Ping(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	this.log.Debug2("<grpc.service.sensordb.Ping>{ }")
	return &empty.Empty{}, nil
}
