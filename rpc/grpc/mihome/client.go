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
	gopi "github.com/djthorpe/gopi"
	grpc "github.com/djthorpe/gopi-rpc/sys/grpc"

	// Protocol buffers
	pb "github.com/djthorpe/sensors/rpc/protobuf/mihome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Client struct {
	pb.MiHomeClient
	conn gopi.RPCClientConn
}

type Protocol struct {
	Name string
	Mode string
}

////////////////////////////////////////////////////////////////////////////////
// NEW

func NewMiHomeClient(conn gopi.RPCClientConn) gopi.RPCClient {
	return &Client{pb.NewMiHomeClient(conn.(grpc.GRPCClientConn).GRPCConn()), conn}
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Client) Conn() gopi.RPCClientConn {
	return this.conn
}

func (this *Client) NewContext() context.Context {
	if this.conn.Timeout() == 0 {
		return context.Background()
	} else {
		ctx, _ := context.WithTimeout(context.Background(), this.conn.Timeout())
		return ctx
	}
}

////////////////////////////////////////////////////////////////////////////////
// CALLS

func (this *Client) MeasureTemperature() (float32, error) {
	if reply, err := this.MiHomeClient.MeasureTemperature(this.NewContext(), &pb.EmptyRequest{}); err != nil {
		return 0, err
	} else {
		return reply.Celcius, nil
	}
}

func (this *Client) Protocols() ([]Protocol, error) {
	if reply, err := this.MiHomeClient.Protocols(this.NewContext(), &pb.EmptyRequest{}); err != nil {
		return nil, err
	} else {
		protocols := make([]Protocol, len(reply.Protocols))
		for i, proto := range reply.Protocols {
			protocols[i].Name = proto.Name
			protocols[i].Mode = proto.Mode
		}
		return protocols, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Client) String() string {
	return fmt.Sprintf("<sensors.MiHome>{ conn=%v }", this.conn)
}
