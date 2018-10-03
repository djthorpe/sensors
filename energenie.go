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
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	MiHomeMode     uint
	OTManufacturer uint8
	OTParameter    uint8
	OTDataType     uint8
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
// PROTOCOLS  - OOK

type ProtoOOK interface {
	gopi.Driver

	// Create a new message
	New(addr uint32, socket uint, state bool) (OOKMessage, error)

	// Encode a message into a payload
	Encode(OOKMessage) []byte

	// Decode a payload into a message
	Decode(payload []byte) (OOKMessage, error)
}

type OOKMessage interface {
	gopi.Event

	Addr() uint32         // 20-bit address
	Socket() uint         // 0 = all or 1-4
	State() bool          // false = off or true = on
	Timestamp() time.Time // Timestamp for received messages
}

////////////////////////////////////////////////////////////////////////////////
// PROTOCOLS  - OPENTHINGS

type ProtoOT interface {
	gopi.Driver

	// Create a new message
	New(manufacturer OTManufacturer, product uint8, sensor uint32) OTMessage

	// Encode a message into a payload
	Encode(OTMessage) []byte

	// Decode a message
	Decode(payload []byte) (OTMessage, error)
}

type OTMessage interface {
	gopi.Event

	Size() uint8
	Manufacturer() OTManufacturer
	ProductID() uint8
	SensorID() uint32
	CRC() uint16
	Records() []OTRecord
	Timestamp() time.Time // Timestamp for received messages
}

type OTRecord interface {
	Name() OTParameter
	Type() OTDataType
	StringValue() (string, error)
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
	// OTManufacturer - see http://www.o-things.com/
	OT_MANUFACTURER_NONE        OTManufacturer = 0x00
	OT_MANUFACTURER_SENTEC      OTManufacturer = 0x01
	OT_MANUFACTURER_HILDERBRAND OTManufacturer = 0x02
	OT_MANUFACTURER_ENERGENIE   OTManufacturer = 0x04
	OT_MANUFACTURER_MAX                        = OT_MANUFACTURER_ENERGENIE
)

const (
	// OTParameter
	OT_PARAM_NONE              OTParameter = 0x00
	OT_PARAM_ALARM             OTParameter = 0x21
	OT_PARAM_DEBUG_OUTPUT      OTParameter = 0x2D
	OT_PARAM_IDENTIFY          OTParameter = 0x3F
	OT_PARAM_SOURCE_SELECTOR   OTParameter = 0x40
	OT_PARAM_WATER_DETECTOR    OTParameter = 0x41
	OT_PARAM_GLASS_BREAKAGE    OTParameter = 0x42
	OT_PARAM_CLOSURES          OTParameter = 0x43
	OT_PARAM_DOOR_BELL         OTParameter = 0x44
	OT_PARAM_ENERGY            OTParameter = 0x45
	OT_PARAM_FALL_SENSOR       OTParameter = 0x46
	OT_PARAM_GAS_VOLUME        OTParameter = 0x47
	OT_PARAM_AIR_PRESSURE      OTParameter = 0x48
	OT_PARAM_ILLUMINANCE       OTParameter = 0x49
	OT_PARAM_LEVEL             OTParameter = 0x4C
	OT_PARAM_RAINFALL          OTParameter = 0x4D
	OT_PARAM_APPARENT_POWER    OTParameter = 0x50
	OT_PARAM_POWER_FACTOR      OTParameter = 0x51
	OT_PARAM_REPORT_PERIOD     OTParameter = 0x52
	OT_PARAM_SMOKE_DETECTOR    OTParameter = 0x53
	OT_PARAM_TIME_AND_DATE     OTParameter = 0x54
	OT_PARAM_VIBRATION         OTParameter = 0x56
	OT_PARAM_WATER_VOLUME      OTParameter = 0x57
	OT_PARAM_WIND_SPEED        OTParameter = 0x58
	OT_PARAM_GAS_PRESSURE      OTParameter = 0x61
	OT_PARAM_BATTERY_LEVEL     OTParameter = 0x62
	OT_PARAM_CO_DETECTOR       OTParameter = 0x63
	OT_PARAM_DOOR_SENSOR       OTParameter = 0x64
	OT_PARAM_EMERGENCY         OTParameter = 0x65
	OT_PARAM_FREQUENCY         OTParameter = 0x66
	OT_PARAM_GAS_FLOW_RATE     OTParameter = 0x67
	OT_PARAM_RELATIVE_HUMIDITY OTParameter = 0x68
	OT_PARAM_CURRENT           OTParameter = 0x69
	OT_PARAM_JOIN              OTParameter = 0x6A
	OT_PARAM_RF_QUALITY        OTParameter = 0x6B
	OT_PARAM_LIGHT_LEVEL       OTParameter = 0x6C
	OT_PARAM_MOTION_DETECTOR   OTParameter = 0x6D
	OT_PARAM_OCCUPANCY         OTParameter = 0x6F
	OT_PARAM_REAL_POWER        OTParameter = 0x70
	OT_PARAM_REACTIVE_POWER    OTParameter = 0x71
	OT_PARAM_ROTATION_SPEED    OTParameter = 0x72
	OT_PARAM_SWITCH_STATE      OTParameter = 0x73
	OT_PARAM_TEMPERATURE       OTParameter = 0x74
	OT_PARAM_VOLTAGE           OTParameter = 0x76
	OT_PARAM_WATER_FLOW_RATE   OTParameter = 0x77
	OT_PARAM_WATER_PRESSURE    OTParameter = 0x78
	OT_PARAM_3PHASE_POWER1     OTParameter = 0x79
	OT_PARAM_3PHASE_POWER2     OTParameter = 0x7A
	OT_PARAM_3PHASE_POWER3     OTParameter = 0x7B
	OT_PARAM_3PHASE_POWER      OTParameter = 0x7C
)

const (
	// OTDataType
	OT_DATATYPE_UDEC_0  OTDataType = 0x00
	OT_DATATYPE_UDEC_4  OTDataType = 0x01
	OT_DATATYPE_UDEC_8  OTDataType = 0x02
	OT_DATATYPE_UDEC_12 OTDataType = 0x03
	OT_DATATYPE_UDEC_16 OTDataType = 0x04
	OT_DATATYPE_UDEC_20 OTDataType = 0x05
	OT_DATATYPE_UDEC_24 OTDataType = 0x06
	OT_DATATYPE_STRING  OTDataType = 0x07
	OT_DATATYPE_DEC_0   OTDataType = 0x08
	OT_DATATYPE_DEC_8   OTDataType = 0x09
	OT_DATATYPE_DEC_16  OTDataType = 0x0A
	OT_DATATYPE_DEC_24  OTDataType = 0x0B
	OT_DATATYPE_ENUM    OTDataType = 0x0C
	OT_DATATYPE_FLOAT   OTDataType = 0x0F
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (m MiHomeMode) String() string {
	switch m {
	case MIHOME_MODE_MONITOR:
		return "MIHOME_MODE_MONITOR"
	case MIHOME_MODE_CONTROL:
		return "MIHOME_MODE_CONTROL"
	default:
		return "[?? Invalid MiHomeMode value]"
	}
}

func (m OTManufacturer) String() string {
	switch m {
	case OT_MANUFACTURER_SENTEC:
		return "OT_MANUFACTURER_SENTEC"
	case OT_MANUFACTURER_HILDERBRAND:
		return "OT_MANUFACTURER_HILDERBRAND"
	case OT_MANUFACTURER_ENERGENIE:
		return "OT_MANUFACTURER_ENERGENIE"
	default:
		return "[?? Invalid OTManufacturer value]"
	}
}

func (p OTParameter) String() string {
	switch p {
	case OT_PARAM_ALARM:
		return "OT_PARAM_ALARM"
	case OT_PARAM_DEBUG_OUTPUT:
		return "OT_PARAM_DEBUG_OUTPUT"
	case OT_PARAM_IDENTIFY:
		return "OT_PARAM_IDENTIFY"
	case OT_PARAM_SOURCE_SELECTOR:
		return "OT_PARAM_SOURCE_SELECTOR"
	case OT_PARAM_WATER_DETECTOR:
		return "OT_PARAM_WATER_DETECTOR"
	case OT_PARAM_GLASS_BREAKAGE:
		return "OT_PARAM_GLASS_BREAKAGE"
	case OT_PARAM_CLOSURES:
		return "OT_PARAM_CLOSURES"
	case OT_PARAM_DOOR_BELL:
		return "OT_PARAM_DOOR_BELL"
	case OT_PARAM_ENERGY:
		return "OT_PARAM_ENERGY"
	case OT_PARAM_FALL_SENSOR:
		return "OT_PARAM_FALL_SENSOR"
	case OT_PARAM_GAS_VOLUME:
		return "OT_PARAM_GAS_VOLUME"
	case OT_PARAM_AIR_PRESSURE:
		return "OT_PARAM_AIR_PRESSURE"
	case OT_PARAM_ILLUMINANCE:
		return "OT_PARAM_ILLUMINANCE"
	case OT_PARAM_LEVEL:
		return "OT_PARAM_LEVEL"
	case OT_PARAM_RAINFALL:
		return "OT_PARAM_RAINFALL"
	case OT_PARAM_APPARENT_POWER:
		return "OT_PARAM_APPARENT_POWER"
	case OT_PARAM_POWER_FACTOR:
		return "OT_PARAM_POWER_FACTOR"
	case OT_PARAM_REPORT_PERIOD:
		return "OT_PARAM_REPORT_PERIOD"
	case OT_PARAM_SMOKE_DETECTOR:
		return "OT_PARAM_SMOKE_DETECTOR"
	case OT_PARAM_TIME_AND_DATE:
		return "OT_PARAM_TIME_AND_DATE"
	case OT_PARAM_VIBRATION:
		return "OT_PARAM_VIBRATION"
	case OT_PARAM_WATER_VOLUME:
		return "OT_PARAM_WATER_VOLUME"
	case OT_PARAM_WIND_SPEED:
		return "OT_PARAM_WIND_SPEED"
	case OT_PARAM_GAS_PRESSURE:
		return "OT_PARAM_GAS_PRESSURE"
	case OT_PARAM_BATTERY_LEVEL:
		return "OT_PARAM_BATTERY_LEVEL"
	case OT_PARAM_CO_DETECTOR:
		return "OT_PARAM_CO_DETECTOR"
	case OT_PARAM_DOOR_SENSOR:
		return "OT_PARAM_DOOR_SENSOR"
	case OT_PARAM_EMERGENCY:
		return "OT_PARAM_EMERGENCY"
	case OT_PARAM_FREQUENCY:
		return "OT_PARAM_FREQUENCY"
	case OT_PARAM_GAS_FLOW_RATE:
		return "OT_PARAM_GAS_FLOW_RATE"
	case OT_PARAM_RELATIVE_HUMIDITY:
		return "OT_PARAM_RELATIVE_HUMIDITY"
	case OT_PARAM_CURRENT:
		return "OT_PARAM_CURRENT"
	case OT_PARAM_JOIN:
		return "OT_PARAM_JOIN"
	case OT_PARAM_RF_QUALITY:
		return "OT_PARAM_RF_QUALITY"
	case OT_PARAM_LIGHT_LEVEL:
		return "OT_PARAM_LIGHT_LEVEL"
	case OT_PARAM_MOTION_DETECTOR:
		return "OT_PARAM_MOTION_DETECTOR"
	case OT_PARAM_OCCUPANCY:
		return "OT_PARAM_OCCUPANCY"
	case OT_PARAM_REAL_POWER:
		return "OT_PARAM_REAL_POWER"
	case OT_PARAM_REACTIVE_POWER:
		return "OT_PARAM_REACTIVE_POWER"
	case OT_PARAM_ROTATION_SPEED:
		return "OT_PARAM_ROTATION_SPEED"
	case OT_PARAM_SWITCH_STATE:
		return "OT_PARAM_SWITCH_STATE"
	case OT_PARAM_TEMPERATURE:
		return "OT_PARAM_TEMPERATURE"
	case OT_PARAM_VOLTAGE:
		return "OT_PARAM_VOLTAGE"
	case OT_PARAM_WATER_FLOW_RATE:
		return "OT_PARAM_WATER_FLOW_RATE"
	case OT_PARAM_WATER_PRESSURE:
		return "OT_PARAM_WATER_PRESSURE"
	case OT_PARAM_3PHASE_POWER1:
		return "OT_PARAM_3PHASE_POWER1"
	case OT_PARAM_3PHASE_POWER2:
		return "OT_PARAM_3PHASE_POWER2"
	case OT_PARAM_3PHASE_POWER3:
		return "OT_PARAM_3PHASE_POWER3"
	case OT_PARAM_3PHASE_POWER:
		return "OT_PARAM_3PHASE_POWER"
	default:
		return "[?? Invalid OTParameter value]"
	}
}

func (t OTDataType) String() string {
	switch t {
	case OT_DATATYPE_UDEC_0:
		return "OT_DATATYPE_UDEC_0"
	case OT_DATATYPE_UDEC_4:
		return "OT_DATATYPE_UDEC_4"
	case OT_DATATYPE_UDEC_8:
		return "OT_DATATYPE_UDEC_8"
	case OT_DATATYPE_UDEC_12:
		return "OT_DATATYPE_UDEC_12"
	case OT_DATATYPE_UDEC_16:
		return "OT_DATATYPE_UDEC_16"
	case OT_DATATYPE_UDEC_20:
		return "OT_DATATYPE_UDEC_20"
	case OT_DATATYPE_UDEC_24:
		return "OT_DATATYPE_UDEC_24"
	case OT_DATATYPE_STRING:
		return "OT_DATATYPE_STRING"
	case OT_DATATYPE_DEC_0:
		return "OT_DATATYPE_DEC_0"
	case OT_DATATYPE_DEC_8:
		return "OT_DATATYPE_DEC_8"
	case OT_DATATYPE_DEC_16:
		return "OT_DATATYPE_DEC_16"
	case OT_DATATYPE_DEC_24:
		return "OT_DATATYPE_DEC_24"
	case OT_DATATYPE_ENUM:
		return "OT_DATATYPE_ENUM"
	case OT_DATATYPE_FLOAT:
		return "OT_DATATYPE_FLOAT"
	default:
		return "[?? Invalid OTDataType value]"
	}
}
