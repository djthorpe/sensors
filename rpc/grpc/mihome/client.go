/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package mihome

import (

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

////////////////////////////////////////////////////////////////////////////////
// NEW

func NewMiHomeClient(conn gopi.RPCClientConn) gopi.RPCClient {
	return &Client{pb.NewMiHomeClient(conn.(grpc.GRPCClientConn).GRPCConn()), conn}
}

/*

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

func (this *Client) Ping() error {
	this.conn.Lock()
	defer this.conn.Unlock()

	if _, err := this.MiHomeClient.Ping(this.NewContext(), &empty.Empty{}); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Client) Receive(done <-chan struct{}, messages chan<- sensors.Message) error {
	this.conn.Lock()
	defer this.conn.Unlock()

	// Create a context with a cancel function, and wait for the 'done'
	// in background
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-done
		cancel()
	}()

	// Receive a stream of messages, when done is received then
	// context.Cancel() is called to end the loop, which returns nil
	if stream, err := this.MiHomeClient.Receive(ctx, &empty.Empty{}); err != nil {
		close(messages)
		return err
	} else {
		for {
			if message_, err := stream.Recv(); err == io.EOF {
				break
			} else if err != nil {
				close(messages)
				return err
			} else if message := fromProtobufMessage(this.conn, message_); message != nil {
				messages <- message
			}
		}
	}

	// Success
	close(messages)
	return nil
}

func (this *Client) On(product sensors.MiHomeProduct, sensor uint32) error {
	this.conn.Lock()
	defer this.conn.Unlock()

	if _, err := this.MiHomeClient.On(this.NewContext(), toProtobufSensorKey("ook", sensors.OT_MANUFACTURER_NONE, product, sensor)); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Client) Off(product sensors.MiHomeProduct, sensor uint32) error {
	this.conn.Lock()
	defer this.conn.Unlock()

	if _, err := this.MiHomeClient.Off(this.NewContext(), toProtobufSensorKey("ook", sensors.OT_MANUFACTURER_NONE, product, sensor)); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Client) SendJoin(product sensors.MiHomeProduct, sensor uint32) error {
	this.conn.Lock()
	defer this.conn.Unlock()

	if _, err := this.MiHomeClient.SendJoin(this.NewContext(), toProtobufSensorKey("openthings", sensors.OT_MANUFACTURER_ENERGENIE, product, sensor)); err != nil {
		return err
	} else {
		return nil
	}
}

/*
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
*/
