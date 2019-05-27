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
	"io"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	grpc "github.com/djthorpe/gopi-rpc/sys/grpc"
	event "github.com/djthorpe/gopi/util/event"
	sensors "github.com/djthorpe/sensors"

	// Protocol buffers
	pb "github.com/djthorpe/sensors/rpc/protobuf/mihome"
	empty "github.com/golang/protobuf/ptypes/empty"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Client struct {
	pb.MiHomeClient
	conn gopi.RPCClientConn
	event.Publisher
}

////////////////////////////////////////////////////////////////////////////////
// NEW

func NewMiHomeClient(conn gopi.RPCClientConn) gopi.RPCClient {
	return &Client{pb.NewMiHomeClient(conn.(grpc.GRPCClientConn).GRPCConn()), conn, event.Publisher{}}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Client) String() string {
	return fmt.Sprintf("<grpc.service.mihome.Client>{ conn=%v }", this.conn)
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
// FUNCTION STUBS

func (this *Client) Ping() error {
	this.conn.Lock()
	defer this.conn.Unlock()

	if _, err := this.MiHomeClient.Ping(this.NewContext(), &empty.Empty{}); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Client) Reset() error {
	this.conn.Lock()
	defer this.conn.Unlock()

	if _, err := this.MiHomeClient.Reset(this.NewContext(), &empty.Empty{}); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Client) On(product sensors.MiHomeProduct, sensor uint32) error {
	this.conn.Lock()
	defer this.conn.Unlock()

	if _, err := this.MiHomeClient.On(this.NewContext(), toProtoSensorKey(sensors.OT_MANUFACTURER_ENERGENIE, product, sensor)); err != nil {
		return err
	} else {
		return nil
	}
}
func (this *Client) Off(product sensors.MiHomeProduct, sensor uint32) error {
	this.conn.Lock()
	defer this.conn.Unlock()

	if _, err := this.MiHomeClient.Off(this.NewContext(), toProtoSensorKey(sensors.OT_MANUFACTURER_ENERGENIE, product, sensor)); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Client) SendJoin(product sensors.MiHomeProduct, sensor uint32) error {
	this.conn.Lock()
	defer this.conn.Unlock()

	if _, err := this.MiHomeClient.SendJoin(this.NewContext(), toProtoSensorKey(sensors.OT_MANUFACTURER_ENERGENIE, product, sensor)); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Client) RequestDiagnostics(product sensors.MiHomeProduct, sensor uint32) error {
	this.conn.Lock()
	defer this.conn.Unlock()

	if _, err := this.MiHomeClient.RequestDiagnostics(this.NewContext(), toProtoSensorRequest(true, sensors.OT_MANUFACTURER_ENERGENIE, product, sensor)); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Client) RequestIdentify(product sensors.MiHomeProduct, sensor uint32) error {
	this.conn.Lock()
	defer this.conn.Unlock()

	if _, err := this.MiHomeClient.RequestIdentify(this.NewContext(), toProtoSensorRequest(true, sensors.OT_MANUFACTURER_ENERGENIE, product, sensor)); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Client) RequestExercise(product sensors.MiHomeProduct, sensor uint32) error {
	this.conn.Lock()
	defer this.conn.Unlock()

	if _, err := this.MiHomeClient.RequestExercise(this.NewContext(), toProtoSensorRequest(true, sensors.OT_MANUFACTURER_ENERGENIE, product, sensor)); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Client) RequestBatteryLevel(product sensors.MiHomeProduct, sensor uint32) error {
	this.conn.Lock()
	defer this.conn.Unlock()

	if _, err := this.MiHomeClient.RequestBatteryLevel(this.NewContext(), toProtoSensorRequest(true, sensors.OT_MANUFACTURER_ENERGENIE, product, sensor)); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Client) SendTargetTemperature(product sensors.MiHomeProduct, sensor uint32, temperature float64) error {
	this.conn.Lock()
	defer this.conn.Unlock()

	if _, err := this.MiHomeClient.SendTargetTemperature(this.NewContext(), toProtoSensorRequestTemperature(true, sensors.OT_MANUFACTURER_ENERGENIE, product, sensor, temperature)); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Client) SendReportInterval(product sensors.MiHomeProduct, sensor uint32, interval time.Duration) error {
	this.conn.Lock()
	defer this.conn.Unlock()

	if _, err := this.MiHomeClient.SendReportInterval(this.NewContext(), toProtoSensorRequestInterval(true, sensors.OT_MANUFACTURER_ENERGENIE, product, sensor, interval)); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Client) SendValveState(product sensors.MiHomeProduct, sensor uint32, state sensors.MiHomeValveState) error {
	this.conn.Lock()
	defer this.conn.Unlock()

	if _, err := this.MiHomeClient.SendValveState(this.NewContext(), toProtoSensorRequestValveState(true, sensors.OT_MANUFACTURER_ENERGENIE, product, sensor, state)); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Client) SendPowerMode(product sensors.MiHomeProduct, sensor uint32, mode sensors.MiHomePowerMode) error {
	this.conn.Lock()
	defer this.conn.Unlock()

	if _, err := this.MiHomeClient.SendPowerMode(this.NewContext(), toProtoSensorRequestPowerMode(true, sensors.OT_MANUFACTURER_ENERGENIE, product, sensor, mode)); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Client) StreamMessages(ctx context.Context) error {
	this.conn.Lock()
	defer this.conn.Unlock()

	stream, err := this.MiHomeClient.StreamMessages(ctx, &empty.Empty{})
	if err != nil {
		return err
	}

	// Errors channel receives errors from recv
	errors := make(chan error)

	// Receive messages in the background
	go func() {
	FOR_LOOP:
		for {
			if message_, err := stream.Recv(); err == io.EOF {
				break FOR_LOOP
			} else if err != nil {
				errors <- err
				break FOR_LOOP
			} else if message_.Sender == nil {
				// Empty message, do nothing
			} else if evt := fromProtoMessage(message_, this.conn); evt != nil {
				this.Emit(evt)
			}
		}
	}()

	// Continue until error or io.EOF is returned
	for {
		select {
		case err := <-errors:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
