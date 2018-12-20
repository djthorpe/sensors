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
	empty "github.com/golang/protobuf/ptypes/empty"
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
func (this *service) Ping(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	this.log.Debug2("<grpc.service.mihome.Ping>{ }")
	return &empty.Empty{}, nil
}

// Status returns the protocols registered
func (this *service) Status(context.Context, *empty.Empty) (*pb.StatusReply, error) {
	return nil, gopi.ErrNotImplemented
}

// Reset the device
func (this *service) Reset(context.Context, *empty.Empty) (*empty.Empty, error) {
	this.log.Debug2("<grpc.service.mihome.Reset>{}")

	if err := this.mihome.Reset(); err != nil {
		this.log.Error("<grpc.service.mihome>Reset: %v", err)
		return &empty.Empty{}, err
	} else {
		return &empty.Empty{}, nil
	}
}

// Receive streams received messages from the radio
func (this *service) Receive(_ *empty.Empty, stream pb.MiHome_ReceiveServer) error {
	this.log.Debug("<grpc.service.mihome.Receive> Started")

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
			if message, ok := evt.(sensors.Message); ok == false || message == nil {
				this.log.Warn("<grpc.service.mihome.Receive> Warning: Did not receive a message: %v", message)
			} else if protobuf := toProtobufMessage(message); protobuf == nil {
				this.log.Warn("<grpc.service.mihome.Receive> Warning: Cannot create protobuf message: %v", message)
			} else if err := stream.Send(protobuf); err != nil {
				if grpc.IsErrUnavailable(err) == false {
					// Client not close connection
					this.log.Warn("<grpc.service.mihome.Receive> Warning: %v: closing request", err)
				}
				break FOR_LOOP
			} else {
				this.log.Debug2("<grpc.service.mihome.Receive>Send: %v", protobuf)
			}
		case evt := <-requests:
			// We should receive a NullEvent here (which terminates the connection)
			if evt == event.NullEvent {
				break FOR_LOOP
			}
		case <-timer.C:
			// Periodic timer to send a null event: an error is returned if
			// the stream has been broken by the client, so that the stream
			// can be closed
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
	this.log.Debug("<grpc.service.mihome.Receive> Ended")

	// Return success
	return nil
}

func (this *service) On(ctx context.Context, key *pb.SensorKey) (*empty.Empty, error) {
	this.log.Debug2("<grpc.service.mihome>On{ key=%v }", key)

	if protocol, manufacturer, product, sensor, err := fromProtobufSensorKey(key); err != nil {
		return nil, err
	} else if protocol != "ook" || manufacturer != sensors.OT_MANUFACTURER_NONE {
		return nil, gopi.ErrBadParameter
	} else if err := this.mihome.RequestSwitchOn(product, sensor); err != nil {
		this.log.Error("<grpc.service.mihome>On: %v", err)
		return &empty.Empty{}, err
	} else {
		return &empty.Empty{}, nil
	}
}

func (this *service) Off(ctx context.Context, key *pb.SensorKey) (*empty.Empty, error) {
	this.log.Debug2("<grpc.service.mihome>Off{ key=%v }", key)

	if protocol, manufacturer, product, sensor, err := fromProtobufSensorKey(key); err != nil {
		return nil, err
	} else if protocol != "ook" || manufacturer != sensors.OT_MANUFACTURER_NONE {
		return nil, gopi.ErrBadParameter
	} else if err := this.mihome.RequestSwitchOff(product, sensor); err != nil {
		this.log.Error("<grpc.service.mihome>Off: %v", err)
		return &empty.Empty{}, err
	} else {
		return &empty.Empty{}, nil
	}
}

func (this *service) SendJoin(ctx context.Context, key *pb.SensorKey) (*empty.Empty, error) {
	this.log.Debug2("<grpc.service.mihome>SendJoin{ key=%v }", key)

	if _, manufacturer, product, sensor, err := fromProtobufSensorKey(key); err != nil {
		return nil, err
	} else if manufacturer != sensors.OT_MANUFACTURER_ENERGENIE {
		return nil, gopi.ErrBadParameter
	} else if err := this.mihome.SendJoin(product, sensor); err != nil {
		this.log.Error("<grpc.service.mihome>SendJoin: %v", err)
		return &empty.Empty{}, err
	} else {
		return &empty.Empty{}, nil
	}
}
