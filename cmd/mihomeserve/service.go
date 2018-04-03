/*
   Go Language Raspberry Pi Interface
   (c) Copyright David Thorpe 2016-2018
   All Rights Reserved
   Documentation http://djthorpe.github.io/gopi/
   For Licensing and Usage information, please see LICENSE.md
*/

package main

//go:generate protoc ../../protobuf/mihome/mihome.proto --go_out=plugins=grpc:.

import (
	// MiHome Service definition
	"reflect"

	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors/protobuf/mihome"
)

type Service struct {
}

func NewService() *Service {
	return new(Service)
}

func (this *Service) Register(server gopi.RPCServer) error {
	// Check to make sure we satisfy the interface
	var _ mihome.Service = (*Service)(nil)
	// Fudge!
	return server.Fudge(reflect.ValueOf(mihome.Service), this)
}
