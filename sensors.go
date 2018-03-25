/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package sensors

import (
	"errors"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type BME280Mode uint8
type BME280Filter uint8
type BME280Standby uint8
type BME280Oversample uint8

type TSL2561Gain uint8
type TSL2561IntegrateTime uint8

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type BME280 interface {
	gopi.Driver

	// Get Version
	ChipIDVersion() (uint8, uint8)

	// Get Mode
	Mode() BME280Mode

	// Return IIR filter co-officient
	Filter() BME280Filter

	// Return standby time
	Standby() BME280Standby

	// Return oversampling values osrs_t, osrs_p, osrs_h
	Oversample() (BME280Oversample, BME280Oversample, BME280Oversample)

	// Return current measuring and updating value
	Status() (bool, bool, error)

	// Return the measurement duty cycle (minimum duration between subsequent readings)
	// in normal mode
	DutyCycle() time.Duration

	// Reset
	SoftReset() error

	// Set BME280 mode
	SetMode(mode BME280Mode) error

	// Set Oversampling
	SetOversample(osrs_t, osrs_p, osrs_h BME280Oversample) error

	// Set Filter
	SetFilter(filter BME280Filter) error

	// Set Standby mode
	SetStandby(t_sb BME280Standby) error

	// Return raw sample data for temperature, pressure and humidity
	// Temperature in Celcius, Pressure in hPa and humidity in
	// %age
	ReadSample() (float64, float64, float64, error)

	// Return altitude in meters for given pressure
	AltitudeForPressure(atmospheric, sealevel float64) float64
}

type TSL2561 interface {
	gopi.Driver

	// Get Version
	ChipIDVersion() (uint8, uint8)

	// Get Gain
	Gain() TSL2561Gain

	// Get Integrate Time
	IntegrateTime() TSL2561IntegrateTime

	// Set Gain
	SetGain(TSL2561Gain) error

	// Set Integrate Time
	SetIntegrateTime(TSL2561IntegrateTime) error

	// Read Luminosity Value in Lux
	ReadSample() (float64, error)
}

////////////////////////////////////////////////////////////////////////////////
// BME280 CONSTANTS

// BME280 Mode
const (
	BME280_MODE_SLEEP   BME280Mode = 0x00
	BME280_MODE_FORCED  BME280Mode = 0x01
	BME280_MODE_FORCED2 BME280Mode = 0x02
	BME280_MODE_NORMAL  BME280Mode = 0x03
	BME280_MODE_MAX     BME280Mode = 0x03
)

// BME280 Filter Co-efficient
const (
	BME280_FILTER_OFF BME280Filter = 0x00
	BME280_FILTER_2   BME280Filter = 0x01
	BME280_FILTER_4   BME280Filter = 0x02
	BME280_FILTER_8   BME280Filter = 0x03
	BME280_FILTER_16  BME280Filter = 0x04
	BME280_FILTER_MAX BME280Filter = 0x07
)

// BME280 Standby time
const (
	BME280_STANDBY_0P5MS  BME280Standby = 0x00
	BME280_STANDBY_62P5MS BME280Standby = 0x01
	BME280_STANDBY_125MS  BME280Standby = 0x02
	BME280_STANDBY_250MS  BME280Standby = 0x03
	BME280_STANDBY_500MS  BME280Standby = 0x04
	BME280_STANDBY_1000MS BME280Standby = 0x05
	BME280_STANDBY_10MS   BME280Standby = 0x06
	BME280_STANDBY_20MS   BME280Standby = 0x07
	BME280_STANDBY_MAX    BME280Standby = 0x07
)

// BME280 Oversampling value
const (
	BME280_OVERSAMPLE_SKIP BME280Oversample = 0x00
	BME280_OVERSAMPLE_1    BME280Oversample = 0x01
	BME280_OVERSAMPLE_2    BME280Oversample = 0x02
	BME280_OVERSAMPLE_4    BME280Oversample = 0x03
	BME280_OVERSAMPLE_8    BME280Oversample = 0x04
	BME280_OVERSAMPLE_16   BME280Oversample = 0x05
	BME280_OVERSAMPLE_MAX  BME280Oversample = 0x07
)

// Sealevel pressure approximation
const (
	BME280_PRESSURE_SEALEVEL float64 = 1013.25
)

////////////////////////////////////////////////////////////////////////////////
// TSL2561 CONSTANTS

const (
	TSL2561_INTEGRATETIME_13P7MS TSL2561IntegrateTime = 0x00
	TSL2561_INTEGRATETIME_101MS  TSL2561IntegrateTime = 0x01
	TSL2561_INTEGRATETIME_402MS  TSL2561IntegrateTime = 0x02
	TSL2561_INTEGRATETIME_MAX    TSL2561IntegrateTime = 0x03
)

const (
	TSL2561_GAIN_1   TSL2561Gain = 0x00
	TSL2561_GAIN_16  TSL2561Gain = 0x01
	TSL2561_GAIN_MAX TSL2561Gain = 0x01
)

////////////////////////////////////////////////////////////////////////////////
// ERRORS

var (
	ErrNoDevice           = errors.New("Missing or invalid hardware device")
	ErrSampleSkipped      = errors.New("Sampling skipped or not enabled")
	ErrUnexpectedResponse = errors.New("Unexpected response from sensor")
	ErrDeviceTimeout      = errors.New("Device timeout")
	ErrMessageCorruption  = errors.New("Message Corruption")
	ErrMessageCRC         = errors.New("Message CRC Error")
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (m BME280Mode) String() string {
	switch m {
	case BME280_MODE_SLEEP:
		return "BME280_MODE_SLEEP"
	case BME280_MODE_FORCED:
		return "BME280_MODE_FORCED"
	case BME280_MODE_FORCED2:
		return "BME280_MODE_FORCED"
	case BME280_MODE_NORMAL:
		return "BME280_MODE_NORMAL"
	default:
		return "[?? Invalid BME280Mode value]"
	}
}

func (f BME280Filter) String() string {
	switch f {
	case BME280_FILTER_OFF:
		return "BME280_FILTER_OFF"
	case BME280_FILTER_2:
		return "BME280_FILTER_2"
	case BME280_FILTER_4:
		return "BME280_FILTER_4"
	case BME280_FILTER_8:
		return "BME280_FILTER_8"
	case BME280_FILTER_16:
		return "BME280_FILTER_16"
	default:
		return "BME280_FILTER_16"
	}
}

func (t BME280Standby) String() string {
	switch t {
	case BME280_STANDBY_0P5MS:
		return "BME280_STANDBY_0P5MS"
	case BME280_STANDBY_62P5MS:
		return "BME280_STANDBY_62P5MS"
	case BME280_STANDBY_125MS:
		return "BME280_STANDBY_125MS"
	case BME280_STANDBY_250MS:
		return "BME280_STANDBY_250MS"
	case BME280_STANDBY_500MS:
		return "BME280_STANDBY_500MS"
	case BME280_STANDBY_1000MS:
		return "BME280_STANDBY_1000MS"
	case BME280_STANDBY_10MS:
		return "BME280_STANDBY_10MS"
	case BME280_STANDBY_20MS:
		return "BME280_STANDBY_20MS"
	default:
		return "[?? Invalid BME280Standby value]"
	}
}

func (o BME280Oversample) String() string {
	switch o {
	case BME280_OVERSAMPLE_SKIP:
		return "BME280_OVERSAMPLE_SKIP"
	case BME280_OVERSAMPLE_1:
		return "BME280_OVERSAMPLE_1"
	case BME280_OVERSAMPLE_2:
		return "BME280_OVERSAMPLE_2"
	case BME280_OVERSAMPLE_4:
		return "BME280_OVERSAMPLE_4"
	case BME280_OVERSAMPLE_8:
		return "BME280_OVERSAMPLE_8"
	case BME280_OVERSAMPLE_16:
		return "BME280_OVERSAMPLE_16"
	default:
		return "[?? Invalid BME280Oversample value]"
	}
}

func (t TSL2561IntegrateTime) String() string {
	switch t {
	case TSL2561_INTEGRATETIME_13P7MS:
		return "TSL2561_INTEGRATETIME_13P7MS"
	case TSL2561_INTEGRATETIME_101MS:
		return "TSL2561_INTEGRATETIME_101MS"
	case TSL2561_INTEGRATETIME_402MS:
		return "TSL2561_INTEGRATETIME_402MS"
	default:
		return "[?? Invalid TSL2561IntegrateTime value]"
	}
}

func (g TSL2561Gain) String() string {
	switch g {
	case TSL2561_GAIN_1:
		return "TSL2561_GAIN_1"
	case TSL2561_GAIN_16:
		return "TSL2561_GAIN_16"
	default:
		return "[?? Invalid TSL2561Gain value]"
	}
}
