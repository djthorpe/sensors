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
	"github.com/djthorpe/gopi/util/tasks"
	"github.com/djthorpe/sensors"

	// Protocol buffers
	pb "github.com/djthorpe/sensors/rpc/protobuf/mihome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Service struct {
	Server gopi.RPCServer
	MiHome sensors.MiHome
	Mode   sensors.MiHomeMode
}

type service struct {
	log    gopi.Logger
	mihome sensors.MiHome
	mode   sensors.MiHomeMode

	// Emit events
	event.Publisher

	// Background Tasks
	tasks.Tasks
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open the server
func (config Service) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<grpc.service.mihome.Open>{ server=%v mihome=%v mode=%v }", config.Server, config.MiHome, config.Mode)

	// Check for bad input parameters
	if config.Server == nil || config.MiHome == nil {
		return nil, gopi.ErrBadParameter
	}

	this := new(service)
	this.log = log
	this.mihome = config.MiHome
	this.mode = config.Mode

	// Register service with GRPC server
	pb.RegisterMiHomeServer(config.Server.(grpc.GRPCServer).GRPCServer(), this)

	// If RX mode then start receiving
	this.Tasks.Start(this.receive)

	// Success
	return this, nil
}

func (this *service) Close() error {
	this.log.Debug("<grpc.service.mihome.Close>{}")

	// End background tasks
	this.Tasks.Close()

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

func (this *service) Ping(ctx context.Context, _ *pb.EmptyRequest) (*pb.EmptyReply, error) {
	this.log.Debug2("<grpc.service.mihome.Ping>{ }")
	return &pb.EmptyReply{}, nil
}

// Returns the protocols registered
func (this *service) Status(context.Context, *pb.EmptyRequest) (*pb.StatusReply, error) {
	return nil, gopi.ErrNotImplemented
}

/*
	protos := this.mihome.Protos()
	reply := &pb.ProtocolsReply{
		Protocols: make([]*pb.Protocol, len(protos)),
	}
	for i, proto := range protos {
		reply.Protocols[i] = &pb.Protocol{
			Name: proto.Name(),
			Mode: fmt.Sprint(proto.Mode()),
		}
	}
	return reply, nil
}
*/

// Resets the device
func (this *service) ResetRadio(context.Context, *pb.EmptyRequest) (*pb.EmptyReply, error) {
	if err := this.mihome.ResetRadio(); err != nil {
		return &pb.EmptyReply{}, err
	} else {
		return &pb.EmptyReply{}, nil
	}
}

// Stream received messages from the radio
func (this *service) Receive(_ *pb.EmptyRequest, stream pb.MiHome_ReceiveServer) error {
	this.log.Debug2("<grpc.service.mihome.Receive> Started")

	// Subscribe to events
	requests := this.Publisher.Subscribe()
	timer := time.NewTicker(500 * time.Millisecond)

FOR_LOOP:
	// Send until loop is broken - either due to stream error, cancellation request or
	// once per 500ms timer which sends null events on the channel
	for {
		select {
		case evt := <-requests:
			// We should either receive a NullEvent (which terminates the connection)
			// or a sensors.Message event
			if evt == event.NullEvent {
				break FOR_LOOP
			} else if err := stream.Send(toProtobufEvent(evt)); err != nil {
				if grpc.IsErrUnavailable(err) == false {
					// Client not close connection
					this.log.Warn("<grpc.service.mihome.Receive> Warning: %v: closing request", err)
				}
				break FOR_LOOP
			}
		case <-timer.C:
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

	// Indicate end of sending stream
	this.log.Debug2("<grpc.service.mihome.Receive> Ended")

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RECEIVE MESSAGES

func (this *service) receive(start chan<- struct{}, stop <-chan struct{}) error {
	start <- gopi.DONE
	this.log.Debug("Started receive task")

	ctx, cancel := context.WithCancel(context.Background())
	errchan := make(chan error)
	go func() {
		errchan <- this.mihome.Receive(ctx, this.mode)
	}()

	evts := this.mihome.Subscribe()

FOR_LOOP:
	for {
		select {
		case <-stop:
			cancel()
			break FOR_LOOP
		case evt := <-evts:
			this.Publisher.Emit(evt)
		}
	}

	this.log.Debug("Stopped receive task")

	// Unsubscribe from channel, and wait for Receive to end
	// then return the error condition
	this.mihome.Unsubscribe(evts)
	err := <-errchan
	return err
}

/*
// Measure temperature
func (this *service) MeasureTemperature(context.Context, *pb.EmptyRequest) (*pb.MeasureTemperatureReply, error) {
	if celcius, err := this.mihome.MeasureTemperature(); err != nil {
		return &pb.MeasureTemperatureReply{}, err
	} else {
		return &pb.MeasureTemperatureReply{Celcius: celcius}, nil
	}
}
*/
/*
// Send 'On' signal
func (this *service) On(context.Context, *pb.SwitchRequest) (*pb.SwitchResponse, error) {
	return nil, gopi.ErrNotImplemented
}

// Send 'Off' signal
func (this *service) Off(context.Context, *pb.SwitchRequest) (*pb.SwitchResponse, error) {
	return nil, gopi.ErrNotImplemented
}
*/
