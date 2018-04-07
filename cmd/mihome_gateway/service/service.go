/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi/util/event"
	"github.com/djthorpe/sensors"

	// Modules
	_ "github.com/djthorpe/gopi/sys/hw/linux"
	_ "github.com/djthorpe/sensors/hw/energenie"
	_ "github.com/djthorpe/sensors/hw/rfm69"
	_ "github.com/djthorpe/sensors/protocol/openthings"

	// Protocol Buffer definitions
	pb "github.com/djthorpe/sensors/protobuf/mihome"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register service/mihome:grpc
	gopi.RegisterModule(gopi.Module{
		Name:     "service/mihome:grpc",
		Type:     gopi.MODULE_TYPE_SERVICE,
		Requires: []string{"rpc/server", "sensors/mihome"},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return gopi.Open(Service{
				Server: app.ModuleInstance("rpc/server").(gopi.RPCServer),
				MiHome: app.ModuleInstance("sensors/mihome").(sensors.MiHome),
			}, app.Logger)
		},
	})
}

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Service struct {
	Server gopi.RPCServer
	MiHome sensors.MiHome
}

type service struct {
	log gopi.Logger

	// The MiHome sensor
	mihome sensors.MiHome

	// Pubsub channel for events emitted for Receive
	pubsub *event.PubSub
	events <-chan gopi.Event

	// Cancel function which stops the receive
	cancel context.CancelFunc

	// Stop relaying messages signal
	relay_done   chan struct{}
	capture_done chan struct{}
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open the server
func (config Service) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<grpc.service.mihome>Open{ server=%v mihome=%v }", config.Server, config.MiHome)

	this := new(service)
	this.log = log
	this.mihome = config.MiHome
	this.pubsub = event.NewPubSub(1)

	// Register service with server
	config.Server.Register(this)

	// Reset the radio
	if err := this.mihome.ResetRadio(); err != nil {
		return nil, err
	}

	// Start goroutine for capturing events from mihome
	this.startCapture()

	// Success
	return this, nil
}

func (this *service) Close() error {
	this.log.Debug("<grpc.service.mihome>Close{}")

	// Stop capturing of events
	var err error
	if err = this.stopCapture(); err != nil {
		this.log.Error("Close: %v", err)
	}

	// Release resources
	this.pubsub.Close()
	this.pubsub = nil
	this.mihome = nil

	// Success
	return err
}

////////////////////////////////////////////////////////////////////////////////
// RPC SERVICE INTERFACE

func (this *service) GRPCHook() reflect.Value {
	return reflect.ValueOf(pb.RegisterMiHomeServer)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *service) String() string {
	return fmt.Sprintf("grpc.service.mihome{}")
}

////////////////////////////////////////////////////////////////////////////////
// CAPTURE EVENTS

func (this *service) startCapture() error {
	// Check to make sure we're not calling wrongly
	if this.events != nil || this.mihome == nil {
		return gopi.ErrOutOfOrder
	}

	// Subscribe to events
	if this.events = this.mihome.Subscribe(); this.events == nil {
		return gopi.ErrOutOfOrder
	}

	// Emit events captured in goroutine
	this.relay_done = make(chan struct{})
	this.capture_done = make(chan struct{})
	go this.relayCapturedEvents()

	// Start receiving from MiHome
	if err := this.startReceive(); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *service) stopCapture() error {

	// Stop receiving from MiHome
	if err := this.stopReceive(); err != nil {
		return err
	}

	// Unsubscribe from mihome events
	if this.events != nil {
		this.mihome.Unsubscribe(this.events)
	}

	// Stop the relay goroutine
	this.relay_done <- gopi.DONE
	close(this.relay_done)

	// Release resources
	this.relay_done = nil
	this.events = nil

	// Return success
	return nil
}

func (this *service) relayCapturedEvents() error {
FOR_LOOP:
	for {
		select {
		case <-this.relay_done:
			break FOR_LOOP
		case evt := <-this.events:
			if evt != nil {
				this.log.Info("Relaying: %v", evt.Name())
				this.pubsub.Emit(evt)
			}
		}
	}
	return nil
}

func (this *service) startReceive() error {

	// Check to make sure there's no existing cancel
	if this.cancel != nil {
		return gopi.ErrOutOfOrder
	}

	// Create the context
	ctx, ctx_cancel := context.WithCancel(context.Background())

	// Run in background
	go func() {
		this.cancel = ctx_cancel
		if err := this.mihome.Receive(ctx, sensors.MIHOME_MODE_MONITOR); err != nil {
			this.log.Error("StartReceive: %v", err)
		}
		this.capture_done <- gopi.DONE
	}()

	// Success
	return nil
}

func (this *service) stopReceive() error {
	if this.cancel != nil {
		this.cancel()
		this.cancel = nil
	}
	// Wait until capture done
	<-this.capture_done
	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RPC METHODS

func (this *service) ResetRadio(ctx context.Context, request *pb.ResetRequest) (*pb.ResetResponse, error) {
	if err := this.mihome.ResetRadio(); err != nil {
		return nil, err
	} else {
		return &pb.ResetResponse{}, nil
	}
}

func (this *service) MeasureTemperature(ctx context.Context, request *pb.MeasureRequest) (*pb.MeasureResponse, error) {
	if temp, err := this.mihome.MeasureTemperature(); err != nil {
		return nil, err
	} else {
		return &pb.MeasureResponse{Celcius: temp}, nil
	}
}

func (this *service) On(ctx context.Context, request *pb.SwitchRequest) (*pb.SwitchResponse, error) {
	var switches []uint

	// Obtain the switch to turn on
	switch request.GetSwitch() {
	case pb.SwitchRequest_SWITCH_ALL:
		switches = []uint{}
	case pb.SwitchRequest_SWITCH_1:
		switches = []uint{1}
	case pb.SwitchRequest_SWITCH_2:
		switches = []uint{2}
	case pb.SwitchRequest_SWITCH_3:
		switches = []uint{3}
	case pb.SwitchRequest_SWITCH_4:
		switches = []uint{4}
	default:
		return nil, gopi.ErrBadParameter
	}

	if err := this.mihome.On(switches...); err != nil {
		return nil, err
	} else {
		return &pb.SwitchResponse{}, nil
	}
}

func (this *service) Off(ctx context.Context, request *pb.SwitchRequest) (*pb.SwitchResponse, error) {
	var switches []uint

	// Obtain the switch to turn on
	switch request.GetSwitch() {
	case pb.SwitchRequest_SWITCH_ALL:
		switches = []uint{}
	case pb.SwitchRequest_SWITCH_1:
		switches = []uint{1}
	case pb.SwitchRequest_SWITCH_2:
		switches = []uint{2}
	case pb.SwitchRequest_SWITCH_3:
		switches = []uint{3}
	case pb.SwitchRequest_SWITCH_4:
		switches = []uint{4}
	default:
		return nil, gopi.ErrBadParameter
	}

	if err := this.mihome.Off(switches...); err != nil {
		return nil, err
	} else {
		return &pb.SwitchResponse{}, nil
	}
}

func (this *service) Receive(request *pb.ReceiveRequest, stream pb.MiHome_ReceiveServer) error {
	// Subscribe to events
	this.log.Debug("Receive: Subscribe")
	events := this.pubsub.Subscribe()

	// Send until loop is broken
FOR_LOOP:
	for {
		select {
		case evt := <-events:
			if evt == nil {
				this.log.Warn("Receive: channel closed: closing request")
				break FOR_LOOP
			} else if reply, err := toReceiveReply(evt); err != nil {
				this.log.Warn("Receive: error sending: %v, contiuing", err)
			} else if err := stream.Send(reply); err != nil {
				this.log.Warn("Receive: error sending: %v: closing request", err)
				break FOR_LOOP
			}
		}
	}

	// Unsubscribe from events
	this.log.Debug("Receive: Unsubscribe")
	this.pubsub.Unsubscribe(events)

	// Return success
	return nil
}

func toReceiveReply(evt gopi.Event) (*pb.ReceiveReply, error) {
	if otevent, ok := evt.(sensors.OTEvent); otevent == nil || ok == false {
		return nil, errors.New("Event emitted is not an OTEvent")
	} else if otevent.Reason() != nil {
		return nil, otevent.Reason()
	} else {
		var parameters []*pb.Parameter
		if len(otevent.Message().Records()) > 0 {
			for _, record := range otevent.Message().Records() {
				parameters = append(parameters, &pb.Parameter{
					Name: pb.Parameter_Name(record.Name()),
				})
			}
		}
		return &pb.ReceiveReply{
			Timestamp:    otevent.Timestamp().Format(time.RFC3339),
			Manufacturer: pb.ReceiveReply_Manufacturer(otevent.Message().Manufacturer()),
			Product:      uint32(otevent.Message().ProductID()),
			Sensor:       otevent.Message().SensorID(),
			Payload:      otevent.Message().Payload(),
			Parameters:   parameters,
		}, nil
	}
}
