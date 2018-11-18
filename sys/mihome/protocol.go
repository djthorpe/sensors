/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package mihome

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi"
	sensors "github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Protocols struct {
	proto_map  map[string]sensors.Proto
	proto_mode map[sensors.MiHomeMode][]sensors.Proto
}

////////////////////////////////////////////////////////////////////////////////
// RELEASE RESOURCES

func (this *Protocols) Close() {
	this.proto_map = nil
	this.proto_mode = nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Protocols) AddProto(proto sensors.Proto) error {
	// Create data structures as necessary
	if this.proto_map == nil {
		this.proto_map = make(map[string]sensors.Proto, 1)
	}
	if this.proto_mode == nil {
		this.proto_mode = make(map[sensors.MiHomeMode][]sensors.Proto, 1)
	}

	// Check to see if protocol is alreay added
	if _, exists := this.proto_map[proto.Name()]; exists {
		return gopi.ErrBadParameter
	} else {
		this.proto_map[proto.Name()] = proto
	}

	// Create an array for holding the protocol
	arr, exists := this.proto_mode[proto.Mode()]
	if exists == false {
		arr = make([]sensors.Proto, 0, 1)
	}

	// Append the protocol
	this.proto_mode[proto.Mode()] = append(arr, proto)

	// Return success
	return nil
}

// Protos returns registered protocols
func (this *Protocols) Protos() []sensors.Proto {
	protos := make([]sensors.Proto, 0, len(this.proto_map))
	for _, proto := range this.proto_map {
		protos = append(protos, proto)
	}
	return protos
}

// ProtoByName returns a single protocol
func (this *Protocols) ProtoByName(name string) sensors.Proto {
	if proto, exists := this.proto_map[name]; exists == false {
		return nil
	} else {
		return proto
	}
}

// ProtosByMode returns zero or more protocols by mode
func (this *Protocols) ProtosByMode(mode sensors.MiHomeMode) []sensors.Proto {
	if protos, exists := this.proto_mode[mode]; exists == false {
		return nil
	} else {
		return protos
	}
}
