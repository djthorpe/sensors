/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package bosch /* import "github.com/djthorpe/gopi-hw/device/bosch" */

import (
	"errors"
	"fmt"
	"math"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi/hw"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// I2C Configuration
type BME280_I2C struct {
	// the I2C driver
	I2C gopi.I2C

	// The slave address, usually 0x77 or 0x76
	Slave uint8
}

// SPI Configuration
type BME280_SPI struct {
	// the SPI driver
	SPI gopi.SPI

	// SPI Device speed in Hertz
	Speed uint32
}

// Concrete driver
type bme280 struct {
	spi         hw.SPIDriver
	i2c         hw.I2CDriver
	slave       uint8
	chipid      uint8
	version     uint8
	calibration *bme280Calibation
	mode        BME280Mode
	filter      BME280Filter
	t_sb        BME280Standby
	osrs_t      BME280Oversample
	osrs_p      BME280Oversample
	osrs_h      BME280Oversample
	spi3w_en    bool
	log         gopi.Logger
}

// BME280 registers and modes
type BME280Register uint8
type BME280Mode uint8
type BME280Filter uint8
type BME280Standby uint8
type BME280Oversample uint8

// BME280 calibration
type bme280Calibation struct {
	T1                             uint16
	T2, T3                         int16
	P1                             uint16
	P2, P3, P4, P5, P6, P7, P8, P9 int16
	H1                             uint8
	H2                             int16
	H3                             uint8
	H4, H5                         int16
	H6                             int8
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	BME280_REG_DIG_T1       BME280Register = 0x88
	BME280_REG_DIG_T2       BME280Register = 0x8A
	BME280_REG_DIG_T3       BME280Register = 0x8C
	BME280_REG_DIG_P1       BME280Register = 0x8E
	BME280_REG_DIG_P2       BME280Register = 0x90
	BME280_REG_DIG_P3       BME280Register = 0x92
	BME280_REG_DIG_P4       BME280Register = 0x94
	BME280_REG_DIG_P5       BME280Register = 0x96
	BME280_REG_DIG_P6       BME280Register = 0x98
	BME280_REG_DIG_P7       BME280Register = 0x9A
	BME280_REG_DIG_P8       BME280Register = 0x9C
	BME280_REG_DIG_P9       BME280Register = 0x9E
	BME280_REG_DIG_H1       BME280Register = 0xA1
	BME280_REG_DIG_H2       BME280Register = 0xE1
	BME280_REG_DIG_H3       BME280Register = 0xE3
	BME280_REG_DIG_H4       BME280Register = 0xE4
	BME280_REG_DIG_H5       BME280Register = 0xE5
	BME280_REG_DIG_H6       BME280Register = 0xE7
	BME280_REG_CHIPID       BME280Register = 0xD0
	BME280_REG_VERSION      BME280Register = 0xD1
	BME280_REG_SOFTRESET    BME280Register = 0xE0
	BME280_REG_CAL26        BME280Register = 0xE1 // R calibration stored in 0xE1-0xF0
	BME280_REG_CONTROLHUMID BME280Register = 0xF2
	BME280_REG_STATUS       BME280Register = 0xF3
	BME280_REG_CONTROL      BME280Register = 0xF4
	BME280_REG_CONFIG       BME280Register = 0xF5
	BME280_REG_PRESSUREDATA BME280Register = 0xF7
	BME280_REG_TEMPDATA     BME280Register = 0xFA
	BME280_REG_HUMIDDATA    BME280Register = 0xFD

	// Write mask
	BME280_REG_SPI_WRITE BME280Register = 0x7F
)

const (
	BME280_I2CSLAVE_DEFAULT   uint8   = 0x76
	BME280_SPI_MAXSPEEDHZ     uint32  = 5000
	BME280_CHIPID_DEFAULT     uint8   = 0x60
	BME280_SOFTRESET_VALUE    uint8   = 0xB6
	BME280_SKIPTEMP_VALUE     int32   = 0x80000
	BME280_SKIPPRESSURE_VALUE int32   = 0x80000
	BME280_SKIPHUMID_VALUE    int32   = 0x8000
	BME280_PRESSURE_SEALEVEL  float64 = 1013.25
)

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

////////////////////////////////////////////////////////////////////////////////
// VARIABLES

var (
	ErrNoDevice      = errors.New("Device not found")
	ErrInvalidDevice = errors.New("Unexpected chip_id value")
	ErrSampleSkipped = errors.New("Temperature sampling skipped")
	ErrWriteDevice   = errors.New("Device Write Error")
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config BME280_I2C) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug2("<bosch.BME280>Open{ bus=%v }", config.I2C)

	this := new(BME280Driver)
	this.i2c = config.I2C
	this.log = log
	this.slave = BME280_I2CSLAVE_DEFAULT

	if config.Slave != 0 {
		this.slave = config.Slave
	}

	if this.i2c == nil {
		return nil, ErrNoDevice
	}

	// Detect slave
	detected, err := this.i2c.DetectSlave(this.slave)
	if err != nil {
		return nil, err
	}
	if detected == false {
		return nil, ErrNoDevice
	}

	// Set slave
	if err := this.i2c.SetSlave(this.slave); err != nil {
		return nil, err
	}

	// Now perform additional setup
	if err := this.setup(); err != nil {
		return nil, err
	}

	// Return success
	return this, nil
}

func (config BME280_SPI) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug2("<bosch.BME280>Open{ bus=%v }", config.SPI)

	this := new(BME280Driver)
	this.spi = config.SPI
	this.log = log

	if this.spi == nil {
		return nil, ErrNoDevice
	}

	// Set SPI bus speed
	if config.Speed != 0 {
		if err := this.spi.SetMaxSpeedHz(config.Speed); err != nil {
			return nil, err
		}
	} else {
		if err := this.spi.SetMaxSpeedHz(BME280_SPI_MAXSPEEDHZ); err != nil {
			return nil, err
		}
	}

	// Now perform additional setup
	if err := this.setup(); err != nil {
		return nil, err
	}

	// Return success
	return this, nil
}

func (this *BME280Driver) Close() error {
	this.log.Debug2("<bosch.BME280>Close")

	return nil
}

func (this *BME280Driver) String() string {
	var bus string
	if this.i2c != nil {
		bus = fmt.Sprintf("%v", this.i2c)
	}
	if this.spi != nil {
		bus = fmt.Sprintf("%v", this.spi)
	}
	return fmt.Sprintf("<bosch.BME280>{ chipid=0x%02X version=0x%02X mode=%v filter=%v t_sb=%v spi3w_en=%v osrs_t=%v osrs_p=%v osrs_h=%v bus=%v calibration=%v }", this.chipid, this.version, this.mode, this.filter, this.t_sb, this.spi3w_en, this.osrs_t, this.osrs_p, this.osrs_h, bus, this.calibration)
}

////////////////////////////////////////////////////////////////////////////////
// READ AND WRITE REGISTERS

func (this *BME280Driver) ReadRegister_Uint8(reg BME280Register) (uint8, error) {
	if this.spi != nil {
		recv, err := this.spi.Transfer([]byte{uint8(reg), 0})
		if err != nil {
			return 0, err
		}
		return recv[1], nil
	}
	if this.i2c != nil {
		recv, err := this.i2c.ReadUint8(uint8(reg))
		if err != nil {
			return 0, err
		}
		return recv, nil
	}
	return 0, ErrNoDevice
}

func (this *BME280Driver) ReadRegister_Int8(reg BME280Register) (int8, error) {
	if this.spi != nil {
		recv, err := this.spi.Transfer([]byte{uint8(reg), 0})
		if err != nil {
			return 0, err
		}
		return int8(recv[1]), nil
	}
	if this.i2c != nil {
		recv, err := this.i2c.ReadInt8(uint8(reg))
		if err != nil {
			return 0, err
		}
		return recv, nil
	}
	return 0, ErrNoDevice
}

func (this *BME280Driver) ReadRegister_Uint16(reg BME280Register) (uint16, error) {
	if this.spi != nil {
		recv, err := this.spi.Transfer([]byte{uint8(reg), 0, 0})
		if err != nil {
			return 0, err
		}
		return uint16(recv[2])<<8 | uint16(recv[1]), nil
	}
	if this.i2c != nil {
		recv, err := this.i2c.ReadUint16(uint8(reg))
		if err != nil {
			return 0, err
		}
		return recv, nil
	}
	return 0, ErrNoDevice
}

func (this *BME280Driver) ReadRegister_Int16(reg BME280Register) (int16, error) {
	if this.spi != nil {
		recv, err := this.spi.Transfer([]byte{uint8(reg), 0, 0})
		if err != nil {
			return 0, err
		}
		return int16(uint16(recv[2])<<8 | uint16(recv[1])), nil
	}
	if this.i2c != nil {
		recv, err := this.i2c.ReadInt16(uint8(reg))
		if err != nil {
			return 0, err
		}
		return recv, nil
	}
	return 0, ErrNoDevice
}

func (this *BME280Driver) WriteRegister_Uint8(reg BME280Register, data uint8) error {
	if this.spi != nil {
		send := []byte{uint8(reg & BME280_REG_SPI_WRITE), data}
		_, err := this.spi.Transfer(send)
		return err
	}
	if this.i2c != nil {
		return this.i2c.WriteUint8(uint8(reg), data)
	}
	return ErrNoDevice
}

////////////////////////////////////////////////////////////////////////////////
// GET REGISTERS

// Return ChipID and Version
func (this *BME280Driver) GetChipIDVersion() (uint8, uint8) {
	return this.chipid, this.version
}

// Return current sampling mode
func (this *BME280Driver) GetMode() BME280Mode {
	return this.mode
}

// Return IIR filter co-officient
func (this *BME280Driver) GetFilter() BME280Filter {
	return this.filter
}

// Return standby time
func (this *BME280Driver) GetStandby() BME280Standby {
	return this.t_sb
}

// Return oversampling values osrs_t, osrs_p, osrs_h
func (this *BME280Driver) GetOversample() (BME280Oversample, BME280Oversample, BME280Oversample) {
	return this.osrs_t, this.osrs_p, this.osrs_h
}

// Return current measuring and updating value
func (this *BME280Driver) GetStatus() (bool, bool, error) {
	status, err := this.ReadRegister_Uint8(BME280_REG_STATUS)
	if err != nil {
		return false, false, err
	}
	measuring := ((status>>3)&0x01 != 0x00)
	updating := (status&0x01 != 0x00)
	return measuring, updating, nil
}

////////////////////////////////////////////////////////////////////////////////
// SET REGISTERS

// Reset the device using the complete power-on-reset procedure
func (this *BME280Driver) SoftReset() error {
	this.log.Debug2("<bosch.BME280>SoftReset")
	if err := this.WriteRegister_Uint8(BME280_REG_SOFTRESET, BME280_SOFTRESET_VALUE); err != nil {
		return err
	}

	// Wait for no measuring or updating
	for {
		measuring, updating, err := this.GetStatus()
		if err != nil {
			return err
		}
		if measuring == false && updating == false {
			break
		}
	}

	// Read registers and return
	return this.read_registers()
}

func (this *BME280Driver) SetMode(mode BME280Mode) error {
	this.log.Debug2("<bosch.BME280>SetMode{ mode=%v }", mode)
	ctrl_meas := uint8(this.osrs_t)<<5 | uint8(this.osrs_p)<<2 | uint8(mode)
	err := this.WriteRegister_Uint8(BME280_REG_CONTROL, ctrl_meas)
	if err != nil {
		return err
	}
	_, _, _, this.mode, err = this.readControl()
	if err != nil {
		return err
	}
	if this.mode != mode {
		return fmt.Errorf("Expected %v got %v", mode, this.mode)
	}

	return nil
}

func (this *BME280Driver) SetOversample(osrs_t, osrs_p, osrs_h BME280Oversample) error {
	this.log.Debug2("<bosch.BME280>SetOversample{ osrs_t=%v osrs_p=%v osrs_h=%v }", osrs_t, osrs_p, osrs_h)

	// Write humidity value first
	if err := this.WriteRegister_Uint8(BME280_REG_CONTROLHUMID, uint8(osrs_h&BME280_OVERSAMPLE_MAX)); err != nil {
		return err
	}

	// Write pressure and temperature second
	ctrl_meas := uint8(osrs_t&BME280_OVERSAMPLE_MAX)<<5 | uint8(osrs_p&BME280_OVERSAMPLE_MAX)<<2 | uint8(this.mode&BME280_MODE_MAX)
	if err := this.WriteRegister_Uint8(BME280_REG_CONTROL, ctrl_meas); err != nil {
		return err
	}

	// Wait for no measuring or updating
	for {
		measuring, updating, err := this.GetStatus()
		if err != nil {
			return err
		}
		if measuring == false && updating == false {
			break
		}
	}

	// Read values back
	var err error
	this.osrs_t, this.osrs_p, this.osrs_h, _, err = this.readControl()
	if err != nil {
		return err
	}
	if this.osrs_t != osrs_t || this.osrs_p != osrs_p || this.osrs_h != osrs_h {
		return ErrWriteDevice
	}

	return nil
}

func (this *BME280Driver) SetFilter(filter BME280Filter) error {
	this.log.Debug2("<bosch.BME280>SetFilter{ filter=%v }", filter)
	config := uint8(this.t_sb)<<5 | uint8(filter)<<2 | to_uint8(this.spi3w_en)
	if err := this.WriteRegister_Uint8(BME280_REG_CONFIG, config); err != nil {
		return err
	}

	// Read values back
	var err error
	_, this.filter, _, err = this.readConfig()
	if err != nil {
		return err
	}
	if this.filter != filter {
		return ErrWriteDevice
	}

	return nil
}

func (this *BME280Driver) SetStandby(t_sb BME280Standby) error {
	this.log.Debug2("<bosch.BME280>SetStandby{ t_sb=%v }", t_sb)
	config := uint8(t_sb)<<5 | uint8(this.filter)<<2 | to_uint8(this.spi3w_en)
	if err := this.WriteRegister_Uint8(BME280_REG_CONFIG, config); err != nil {
		return err
	}

	// Read values back
	var err error
	this.t_sb, _, _, err = this.readConfig()
	if err != nil {
		return err
	}
	if this.t_sb != t_sb {
		return ErrWriteDevice
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// GET SAMPLE DATA

// Return raw sample data for temperature, pressure and humidity
func (this *BME280Driver) ReadSample() (float64, float64, float64, error) {
	this.log.Debug2("<bosch.BME280>ReadSample")

	// Wait for no measuring or updating
	for {
		measuring, updating, err := this.GetStatus()
		if err != nil {
			return 0, 0, 0, err
		}
		if measuring == false && updating == false {
			break
		}
	}

	// Obtain the current mode of operation if we're in FORCED mode and
	// return ErrSampleSkipped if the current mode isn't forced
	if this.mode == BME280_MODE_FORCED {
		var err error
		if _, _, _, this.mode, err = this.readControl(); err != nil {
			return 0, 0, 0, err
		}
		if this.mode != BME280_MODE_FORCED {
			return 0, 0, 0, ErrSampleSkipped
		}
	}

	// Read temperature, return error if temperature reading is skipped
	adc_t, err := this.readTemperature()
	if err != nil {
		return 0, 0, 0, err
	}
	if adc_t == BME280_SKIPTEMP_VALUE {
		return 0, 0, 0, ErrSampleSkipped
	}
	t_celcius, t_fine := this.toCelcius(adc_t)

	// Read pressure. Set ADC value to zero if skipped
	adc_p, err := this.readPressure()
	if err != nil {
		return 0, 0, 0, err
	}
	if adc_p == BME280_SKIPPRESSURE_VALUE {
		adc_p = 0
	}
	t_pressure := this.toPascals(adc_p, t_fine)

	// Read humidity. Set ADC value to zero if skipped
	adc_h, err := this.readHumidity()
	if err != nil {
		return 0, 0, 0, err
	}
	if adc_h == BME280_SKIPHUMID_VALUE {
		adc_h = 0
	}
	t_humidity := this.toRelativeHumidity(adc_h, t_fine)

	// Return success
	return t_celcius, t_pressure, t_humidity, nil
}

// Return altitude in metres based on pressure reading in Pascals, given
// the sealevel pressure in Pascals. You can use a standard value of
// BME280_PRESSURE_SEALEVEL for sealevel
func (this *BME280Driver) AltitudeForPressure(atmospheric, sealevel float64) float64 {
	return 44330.0 * (1.0 - math.Pow(atmospheric/sealevel, (1.0/5.255)))
}
