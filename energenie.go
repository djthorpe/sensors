/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package sensors

import (
	"context"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	MiHomeMode    uint
	MiHomeProduct uint
)

////////////////////////////////////////////////////////////////////////////////
// ENER314 AND ENER314RT

type ENER314 interface {
	gopi.Driver

	// Send on signal - when no sockets specified then
	// sends to all sockets
	On(sockets ...uint) error

	// Send off signal - when no sockets specified then
	// sends to all sockets
	Off(sockets ...uint) error
}

type MiHome interface {
	gopi.Publisher
	ENER314

	// Add a wire protocol which encodes/decodes messages
	AddProto(Proto) error

	// Set a protocol mode
	SetMode(MiHomeMode) error

	// Reset the radio device
	ResetRadio() error

	// Receive payloads with radio until context deadline exceeded or cancel
	Receive(ctx context.Context, mode MiHomeMode) error

	// Send a raw payload with radio
	Send(payload []byte, repeat uint, mode MiHomeMode) error

	// Measure Temperature
	MeasureTemperature() (float32, error)
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MIHOME_MODE_NONE    MiHomeMode = iota
	MIHOME_MODE_MONITOR            // FSK
	MIHOME_MODE_CONTROL            // OOK
	MIHOME_MODE_MAX     = MIHOME_MODE_CONTROL
)

const (
	MIHOME_PRODUCT_NONE    MiHomeProduct = 0x00
	MIHOME_PRODUCT_MIHO004 MiHomeProduct = 0x01 // Adaptor Monitor
	MIHOME_PRODUCT_MIHO005 MiHomeProduct = 0x02 // Adaptor Plus
	MIHOME_PRODUCT_MIHO013 MiHomeProduct = 0x03 // eTRV
	MIHOME_PRODUCT_MIHO006 MiHomeProduct = 0x05 // House Monitor
	MIHOME_PRODUCT_MIHO032 MiHomeProduct = 0x0C // Motion sensor
	MIHOME_PRODUCT_MIHO033 MiHomeProduct = 0x0D // Door sensor
	MIHOME_PRODUCT_MAX                   = MIHOME_PRODUCT_MIHO033
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (m MiHomeMode) String() string {
	switch m {
	case MIHOME_MODE_NONE:
		return "MIHOME_MODE_NONE"
	case MIHOME_MODE_MONITOR:
		return "MIHOME_MODE_MONITOR"
	case MIHOME_MODE_CONTROL:
		return "MIHOME_MODE_CONTROL"
	default:
		return "[?? Invalid MiHomeMode value]"
	}
}

func (p MiHomeProduct) String() string {
	switch p {
	case MIHOME_PRODUCT_NONE:
		return "MIHOME_PRODUCT_NONE"
	case MIHOME_PRODUCT_MIHO004:
		return "MIHOME_PRODUCT_MIHO004"
	case MIHOME_PRODUCT_MIHO005:
		return "MIHOME_PRODUCT_MIHO005"
	case MIHOME_PRODUCT_MIHO013:
		return "MIHOME_PRODUCT_MIHO013"
	case MIHOME_PRODUCT_MIHO006:
		return "MIHOME_PRODUCT_MIHO006"
	case MIHOME_PRODUCT_MIHO032:
		return "MIHOME_PRODUCT_MIHO032"
	case MIHOME_PRODUCT_MIHO033:
		return "MIHOME_PRODUCT_MIHO033"
	default:
		return "[?? Invalid MiHomeProduct value]"
	}
}
