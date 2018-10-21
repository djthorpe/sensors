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

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi-rpc/sys/grpc"
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

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *service) String() string {
	return fmt.Sprintf("<grpc.service.mihome>{ mihome=%v }", this.mihome)
}

////////////////////////////////////////////////////////////////////////////////
// RPC METHODS

// Returns the protocols registered
func (this *service) Protocols(context.Context, *pb.EmptyRequest) (*pb.ProtocolsReply, error) {
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

// Resets the device
func (this *service) ResetRadio(context.Context, *pb.EmptyRequest) (*pb.EmptyReply, error) {
	if err := this.mihome.ResetRadio(); err != nil {
		return &pb.EmptyReply{}, err
	} else {
		return &pb.EmptyReply{}, nil
	}
}

// Measure temperature
func (this *service) MeasureTemperature(context.Context, *pb.EmptyRequest) (*pb.MeasureTemperatureReply, error) {
	if celcius, err := this.mihome.MeasureTemperature(); err != nil {
		return &pb.MeasureTemperatureReply{}, err
	} else {
		return &pb.MeasureTemperatureReply{Celcius: celcius}, nil
	}
}

// Receive data
func (this *service) Receive(*pb.ReceiveRequest, pb.MiHome_ReceiveServer) error {
	return gopi.ErrNotImplemented
}

// Send 'On' signal
func (this *service) On(context.Context, *pb.SwitchRequest) (*pb.SwitchResponse, error) {
	return nil, gopi.ErrNotImplemented
}

// Send 'Off' signal
func (this *service) Off(context.Context, *pb.SwitchRequest) (*pb.SwitchResponse, error) {
	return nil, gopi.ErrNotImplemented
}
