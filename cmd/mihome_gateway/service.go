/*
   Go Language Raspberry Pi Interface
   (c) Copyright David Thorpe 2016-2018
   All Rights Reserved
   Documentation http://djthorpe.github.io/gopi/
   For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"context"
	"reflect"

	// Frameworks
	"github.com/djthorpe/gopi"

	// Protocol Buffer
	"github.com/djthorpe/sensors/protobuf/mihome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Service struct {
}

////////////////////////////////////////////////////////////////////////////////
// CREATE THE SERVICE

func NewService() *Service {
	return new(Service)
}

func (this *Service) Register(server gopi.RPCServer) error {
	// Check to make sure we satisfy the interface
	var _ mihome.MiHomeServer = (*Service)(nil)

	// Fudge to perform the registration
	return server.Fudge(reflect.ValueOf(mihome.RegisterMiHomeServer), this)
}

////////////////////////////////////////////////////////////////////////////////
// RESET THE DEVICE

func (this *Service) Reset(ctx context.Context, req *mihome.ResetRequest) (*mihome.ResetReply, error) {
	return nil, gopi.ErrNotImplemented
}
