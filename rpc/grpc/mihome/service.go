/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package mihome

import (
	"context"
	"fmt"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi-rpc/sys/grpc"
	"github.com/djthorpe/gopi/util/event"
	"github.com/djthorpe/sensors"

	// Protocol buffers
	pb "github.com/djthorpe/sensors/rpc/protobuf/mihome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Service struct {
	Server gopi.RPCServer
	MiHome sensors.MiHome
}

type service struct {
	log    gopi.Logger
	mihome sensors.MiHome

	// Emit events
	event.Publisher
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open the server
func (config Service) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<grpc.service.mihome.Open>{ server=%v mihome=%v }", config.Server, config.MiHome)

	// Check for bad input parameters
	if config.Server == nil || config.MiHome == nil {
		return nil, gopi.ErrBadParameter
	}

	this := new(service)
	this.log = log
	this.mihome = config.MiHome

	// Register service with GRPC server
	pb.RegisterMiHomeServer(config.Server.(grpc.GRPCServer).GRPCServer(), this)

	// Success
	return this, nil
}

func (this *service) Close() error {
	this.log.Debug("<grpc.service.mihome.Close>{}")

	// Close publisher
	this.Publisher.Close()

	// Release resources
	this.mihome = nil

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *service) String() string {
	return fmt.Sprintf("<grpc.service.mihome>{ mihome=%v }", this.mihome)
}

////////////////////////////////////////////////////////////////////////////////
// CANCEL STREAMING REQUESTS

func (this *service) CancelRequests() error {
	this.log.Debug2("<grpc.service.mihome.CancelRequests>{}")

	// Cancel any streaming requests
	this.Publisher.Emit(event.NullEvent)

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RPC METHODS

// Ping returns an empty response
func (this *service) Ping(ctx context.Context, _ *pb.EmptyRequest) (*pb.EmptyReply, error) {
	this.log.Debug2("<grpc.service.mihome.Ping>{ }")
	return &pb.EmptyReply{}, nil
}

// Status returns the protocols registered
func (this *service) Status(context.Context, *pb.EmptyRequest) (*pb.StatusReply, error) {
	return nil, gopi.ErrNotImplemented
}

// Reset the device
func (this *service) Reset(context.Context, *pb.EmptyRequest) (*pb.EmptyReply, error) {
	this.log.Debug2("<grpc.service.mihome.Reset>{}")

	if err := this.mihome.Reset(); err != nil {
		this.log.Error("<grpc.service.mihome>Reset: %v", err)
		return &pb.EmptyReply{}, err
	} else {
		return &pb.EmptyReply{}, nil
	}
}

// Receive streams received messages from the radio
func (this *service) Receive(_ *pb.EmptyRequest, stream pb.MiHome_ReceiveServer) error {
	this.log.Debug2("<grpc.service.mihome.Receive> Started")

	// Subscribe to events
	requests := this.Publisher.Subscribe()
	messages := this.mihome.Subscribe()
	timer := time.NewTicker(500 * time.Millisecond)

FOR_LOOP:
	// Send until loop is broken - either due to stream error, cancellation request or
	// once per 500ms timer which sends null events on the channel
	for {
		select {
		case evt := <-messages:
			fmt.Println(evt)
			if err := stream.Send(toProtobufEvent(evt)); err != nil {
				if grpc.IsErrUnavailable(err) == false {
					// Client not close connection
					this.log.Warn("<grpc.service.mihome.Receive> Warning: %v: closing request", err)
				}
				break FOR_LOOP
			}
		case evt := <-requests:
			// We should receive a NullEvent here (which terminates the connection)
			if evt == event.NullEvent {
				break FOR_LOOP
			}
		case <-timer.C:
			// Periodic timer to send a null event
			if err := stream.Send(toProtobufNullEvent()); err != nil {
				if grpc.IsErrUnavailable(err) == false {
					// Client not close connection
					this.log.Warn("<grpc.service.mihome.Receive> Warning: %v: closing request", err)
				}
				break FOR_LOOP
			}
		}
	}

	// Unsubscribe from events
	timer.Stop()
	this.Publisher.Unsubscribe(requests)
	this.mihome.Unsubscribe(messages)

	// Indicate end of sending stream
	this.log.Debug2("<grpc.service.mihome.Receive> Ended")

	// Return success
	return nil
}

func (this *service) On(ctx context.Context, _ *pb.SensorKey) (*pb.EmptyReply, error) {
	this.log.Debug2("<grpc.service.mihome>On{ sensor=%v }", sensor)

	// TODO: Convert sensor to product,sensor

	if err := this.mihome.RequestSwitchOn(product, sensor); err != nil {
		this.log.Error("<grpc.service.mihome>On: %v", err)
		return &pb.EmptyReply{}, err
	} else {
		return &pb.EmptyReply{}, nil
	}
}

func (this *service) Off(ctx context.Context, _ *pb.SensorKey) (*pb.EmptyReply, error) {
	this.log.Debug2("<grpc.service.mihome>Off{ sensor=%v }", sensor)

	// TODO: Convert sensor to product,sensor

	if err := this.mihome.RequestSwitchOff(product, sensor); err != nil {
		this.log.Error("<grpc.service.mihome>Off: %v", err)
		return &pb.EmptyReply{}, err
	} else {
		return &pb.EmptyReply{}, nil
	}
}
