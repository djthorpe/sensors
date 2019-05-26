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
	"sync"
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

	// Lock
	sync.Mutex

	// Emit events
	event.Publisher

	// Queue for transmitting messages through a queue
	// which is triggered on received message
	queue
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

	// Init queue
	if err := this.queue.Init(log, config); err != nil {
		return nil, err
	}

	// Register service with GRPC server
	pb.RegisterMiHomeServer(config.Server.(grpc.GRPCServer).GRPCServer(), this)

	// Success
	return this, nil
}

func (this *service) Close() error {
	this.log.Debug("<grpc.service.mihome.Close>{}")

	// Close publisher
	this.Publisher.Close()

	// Destroy queue
	if err := this.queue.Destroy(); err != nil {
		return err
	}

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
	this.log.Debug2("<grpc.service.mihome>CancelRequests{}")

	// Cancel any streaming requests
	this.Publisher.Emit(event.NullEvent)

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RPC METHODS

// Ping returns an empty response
func (this *service) Ping(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	this.log.Debug("<grpc.service.mihome.Ping>{ }")

	this.Lock()
	defer this.Unlock()

	return &empty.Empty{}, nil
}

// Reset the device
func (this *service) Reset(context.Context, *empty.Empty) (*empty.Empty, error) {
	this.log.Debug("<grpc.service.mihome>Reset>{}")

	this.Lock()
	defer this.Unlock()

	if err := this.mihome.Reset(); err != nil {
		this.log.Error("Reset: %v", err)
		return nil, err
	} else {
		return &empty.Empty{}, nil
	}
}

// Send an On signal
func (this *service) On(ctx context.Context, key *pb.SensorKey) (*empty.Empty, error) {
	this.log.Debug("<grpc.service.mihome>On{ key=%v }", key)

	this.Lock()
	defer this.Unlock()

	if manufacturer, product, sensor, err := fromProtobufSensorKey(key); err != nil {
		return nil, err
	} else if manufacturer != sensors.OT_MANUFACTURER_ENERGENIE {
		return nil, gopi.ErrBadParameter
	} else if err := this.mihome.RequestSwitchOn(product, sensor); err != nil {
		this.log.Error("On: %v", err)
		return nil, err
	} else {
		return &empty.Empty{}, nil
	}
}

// Send an Off signal
func (this *service) Off(ctx context.Context, key *pb.SensorKey) (*empty.Empty, error) {
	this.log.Debug("<grpc.service.mihome>Off{ key=%v }", key)

	this.Lock()
	defer this.Unlock()

	if manufacturer, product, sensor, err := fromProtobufSensorKey(key); err != nil {
		return nil, err
	} else if manufacturer != sensors.OT_MANUFACTURER_ENERGENIE {
		return nil, gopi.ErrBadParameter
	} else if err := this.mihome.RequestSwitchOff(product, sensor); err != nil {
		this.log.Error("Off: %v", err)
		return nil, err
	} else {
		return &empty.Empty{}, nil
	}
}

// Status returns the protocols registered
func (this *service) Status(context.Context, *empty.Empty) (*pb.StatusReply, error) {
	this.log.Debug("<grpc.service.mihome>Status{}")

	this.Lock()
	defer this.Unlock()

	if celcius, err := this.mihome.MeasureTemperature(); err != nil {
		return nil, err
	} else {
		reply := &pb.StatusReply{
			Protocol:      toProtoProtocols(this.mihome.Protos()),
			DeviceCelcius: celcius,
		}
		return reply, nil
	}
}

// Receive streams received messages from the radio
func (this *service) StreamMessages(_ *empty.Empty, stream pb.MiHome_StreamMessagesServer) error {
	this.log.Debug("<grpc.service.mihome>StreamMessages Started")

	// Subscribe to channel for incoming events, and continue until cancel request is received, send
	// empty events occasionally to ensure the channel is still alive
	events := this.mihome.Subscribe()
	cancel := this.Subscribe()
	ticker := time.NewTicker(time.Second)

FOR_LOOP:
	for {
		select {
		case evt := <-events:
			if evt == nil {
				break FOR_LOOP
			} else if evt_, ok := evt.(sensors.Message); ok {
				if err := stream.Send(toProtoMessage(evt_)); err != nil {
					this.log.Warn("StreamMessages: %v", err)
					break FOR_LOOP
				}
			} else {
				this.log.Warn("StreamMessages: Ignoring event: %v", evt)
			}
		case <-ticker.C:
			if err := stream.Send(&pb.Message{}); err != nil {
				this.log.Warn("StreamMessages: %v", err)
				break FOR_LOOP
			}
		case <-cancel:
			break FOR_LOOP
		}
	}

	// Stop ticker, unsubscribe from events
	ticker.Stop()
	this.mihome.Unsubscribe(events)
	this.Unsubscribe(cancel)

	this.log.Debug2("StreamMessages: Ended")

	// Return success
	return nil
}

func (this *service) SendJoin(ctx context.Context, key *pb.SensorKey) (*empty.Empty, error) {
	this.log.Debug("<grpc.service.mihome>SendJoin{ key=%v }", key)

	this.Lock()
	defer this.Unlock()

	if manufacturer, product, sensor, err := fromProtobufSensorKey(key); err != nil {
		return nil, err
	} else if manufacturer != sensors.OT_MANUFACTURER_ENERGENIE {
		return nil, gopi.ErrBadParameter
	} else if err := this.mihome.SendJoin(product, sensor); err != nil {
		this.log.Error("SendJoin: %v", err)
		return nil, err
	} else {
		return &empty.Empty{}, nil
	}
}

func (this *service) RequestDiagnostics(ctx context.Context, req *pb.SensorRequest) (*empty.Empty, error) {
	this.log.Debug("<grpc.service.mihome>RequestDiagnostics{ req=%v }", req)

	this.Lock()
	defer this.Unlock()

	if manufacturer, product, sensor, err := fromProtobufSensorKey(req.Sensor); err != nil {
		return nil, err
	} else if manufacturer != sensors.OT_MANUFACTURER_ENERGENIE {
		return nil, gopi.ErrBadParameter
	} else if req.QueueRequest {
		if err := this.queue.QueueDiagnostics(product, sensor); err != nil {
			this.log.Error("QueueDiagnostics: %v", err)
			return nil, err
		}
	} else {
		if err := this.mihome.RequestDiagnostics(product, sensor); err != nil {
			this.log.Error("RequestDiagnostics: %v", err)
			return nil, err
		}
	}

	// Success
	return &empty.Empty{}, nil
}

func (this *service) RequestIdentify(ctx context.Context, req *pb.SensorRequest) (*empty.Empty, error) {
	this.log.Debug("<grpc.service.mihome>RequestIdentify{ req=%v }", req)

	this.Lock()
	defer this.Unlock()

	if manufacturer, product, sensor, err := fromProtobufSensorKey(req.Sensor); err != nil {
		return nil, err
	} else if manufacturer != sensors.OT_MANUFACTURER_ENERGENIE {
		return nil, gopi.ErrBadParameter
	} else if req.QueueRequest {
		if err := this.queue.QueueIdentify(product, sensor); err != nil {
			this.log.Error("QueueIdentify: %v", err)
			return nil, err
		}
	} else {
		if err := this.mihome.RequestIdentify(product, sensor); err != nil {
			this.log.Error("RequestIdentify: %v", err)
			return nil, err
		}
	}

	// Success
	return &empty.Empty{}, nil
}

func (this *service) RequestExercise(ctx context.Context, req *pb.SensorRequest) (*empty.Empty, error) {
	this.log.Debug("<grpc.service.mihome>RequestExercise{ req=%v }", req)

	this.Lock()
	defer this.Unlock()

	if manufacturer, product, sensor, err := fromProtobufSensorKey(req.Sensor); err != nil {
		return nil, err
	} else if manufacturer != sensors.OT_MANUFACTURER_ENERGENIE {
		return nil, gopi.ErrBadParameter
	} else if req.QueueRequest {
		if err := this.queue.QueueExercise(product, sensor); err != nil {
			this.log.Error("QueueExercise: %v", err)
			return nil, err
		}
	} else {
		if err := this.mihome.RequestExercise(product, sensor); err != nil {
			this.log.Error("RequestExercise: %v", err)
			return nil, err
		}
	}

	// Success
	return &empty.Empty{}, nil
}

func (this *service) RequestBatteryLevel(ctx context.Context, req *pb.SensorRequest) (*empty.Empty, error) {
	this.log.Debug("<grpc.service.mihome>RequestBatteryLevel{ req=%v }", req)

	this.Lock()
	defer this.Unlock()

	if manufacturer, product, sensor, err := fromProtobufSensorKey(req.Sensor); err != nil {
		return nil, err
	} else if manufacturer != sensors.OT_MANUFACTURER_ENERGENIE {
		return nil, gopi.ErrBadParameter
	} else if req.QueueRequest {
		if err := this.queue.QueueBatteryLevel(product, sensor); err != nil {
			this.log.Error("QueueBatteryLevel: %v", err)
			return nil, err
		}
	} else {
		if err := this.mihome.RequestBatteryLevel(product, sensor); err != nil {
			this.log.Error("RequestBatteryLevel: %v", err)
			return nil, err
		}
	}

	// Success
	return &empty.Empty{}, nil
}

func (this *service) SendTargetTemperature(ctx context.Context, req *pb.SensorRequestTemperature) (*empty.Empty, error) {
	this.log.Debug("<grpc.service.mihome>SendTargetTemperature{ req=%v }", req)

	this.Lock()
	defer this.Unlock()

	if manufacturer, product, sensor, err := fromProtobufSensorKey(req.Sensor); err != nil {
		return nil, err
	} else if manufacturer != sensors.OT_MANUFACTURER_ENERGENIE {
		return nil, gopi.ErrBadParameter
	} else if req.QueueRequest {
		if err := this.queue.QueueTargetTemperature(product, sensor, req.Temperature); err != nil {
			this.log.Error("QueueTargetTemperature: %v", err)
			return nil, err
		}
	} else {
		if err := this.mihome.RequestTargetTemperature(product, sensor, req.Temperature); err != nil {
			this.log.Error("RequestTargetTemperature: %v", err)
			return nil, err
		}
	}

	// Success
	return &empty.Empty{}, nil
}

func (this *service) SendReportInterval(ctx context.Context, req *pb.SensorRequestInterval) (*empty.Empty, error) {
	this.log.Debug("<grpc.service.mihome>SendReportInterval{ req=%v }", req)

	this.Lock()
	defer this.Unlock()

	if manufacturer, product, sensor, err := fromProtobufSensorKey(req.Sensor); err != nil {
		return nil, err
	} else if manufacturer != sensors.OT_MANUFACTURER_ENERGENIE {
		return nil, gopi.ErrBadParameter
	} else if duration := fromProtoDuration(req.Interval); duration == 0 {
		return nil, gopi.ErrBadParameter
	} else if req.QueueRequest {
		if err := this.queue.QueueReportInterval(product, sensor, duration); err != nil {
			this.log.Error("QueueReportInterval: %v", err)
			return nil, err
		}
	} else {
		if err := this.mihome.RequestReportInterval(product, sensor, duration); err != nil {
			this.log.Error("RequestReportInterval: %v", err)
			return nil, err
		}
	}

	// Success
	return &empty.Empty{}, nil
}

func (this *service) SendPowerMode(ctx context.Context, req *pb.SensorRequestPowerMode) (*empty.Empty, error) {
	this.log.Debug("<grpc.service.mihome>SendPowerMode{ req=%v }", req)

	this.Lock()
	defer this.Unlock()

	if manufacturer, product, sensor, err := fromProtobufSensorKey(req.Sensor); err != nil {
		return nil, err
	} else if manufacturer != sensors.OT_MANUFACTURER_ENERGENIE {
		return nil, gopi.ErrBadParameter
	} else if req.QueueRequest {
		if err := this.queue.QueueLowPowerMode(product, sensor, fromProtoPowerMode(req.PowerMode)); err != nil {
			this.log.Error("QueueLowPowerMode: %v", err)
			return nil, err
		}
	} else {
		if err := this.mihome.RequestLowPowerMode(product, sensor, fromProtoPowerMode(req.PowerMode)); err != nil {
			this.log.Error("RequestLowPowerMode: %v", err)
			return nil, err
		}
	}

	// Success
	return &empty.Empty{}, nil
}

func (this *service) SendValueState(ctx context.Context, req *pb.SensorRequestValveState) (*empty.Empty, error) {
	this.log.Debug("<grpc.service.mihome>SendValueState{ req=%v }", req)

	this.Lock()
	defer this.Unlock()

	if manufacturer, product, sensor, err := fromProtobufSensorKey(req.Sensor); err != nil {
		return nil, err
	} else if manufacturer != sensors.OT_MANUFACTURER_ENERGENIE {
		return nil, gopi.ErrBadParameter
	} else if req.QueueRequest {
		if err := this.queue.QueueValveState(product, sensor, fromProtoValueState(req.ValueState)); err != nil {
			this.log.Error("QueueValveState: %v", err)
			return nil, err
		}
	} else {
		if err := this.mihome.RequestValveState(product, sensor, fromProtoValueState(req.ValueState)); err != nil {
			this.log.Error("RequestValveState: %v", err)
			return nil, err
		}
	}

	// Success
	return &empty.Empty{}, nil
}
