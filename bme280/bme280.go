/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package bme280

import (
	"fmt"
	"math"

	gopi "github.com/djthorpe/gopi"
	sensors "github.com/djthorpe/sensors"
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
	spi         gopi.SPI
	i2c         gopi.I2C
	slave       uint8
	chipid      uint8
	version     uint8
	calibration *calibation
	mode        sensors.BME280Mode
	filter      sensors.BME280Filter
	t_sb        sensors.BME280Standby
	osrs_t      sensors.BME280Oversample
	osrs_p      sensors.BME280Oversample
	osrs_h      sensors.BME280Oversample
	spi3w_en    bool
	log         gopi.Logger
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	BME280_I2CSLAVE_DEFAULT   uint8   = 0x77
	BME280_SPI_MAXSPEEDHZ     uint32  = 5000
	BME280_CHIPID_DEFAULT     uint8   = 0x60
	BME280_SOFTRESET_VALUE    uint8   = 0xB6
	BME280_SKIPTEMP_VALUE     int32   = 0x80000
	BME280_SKIPPRESSURE_VALUE int32   = 0x80000
	BME280_SKIPHUMID_VALUE    int32   = 0x8000
	BME280_PRESSURE_SEALEVEL  float64 = 1013.25
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACE - GET

// Return ChipID and Version
func (this *bme280) ChipIDVersion() (uint8, uint8) {
	return this.chipid, this.version
}

// Return current sampling mode
func (this *bme280) Mode() sensors.BME280Mode {
	return this.mode
}

// Return IIR filter co-officient
func (this *bme280) Filter() sensors.BME280Filter {
	return this.filter
}

// Return standby time
func (this *bme280) Standby() sensors.BME280Standby {
	return this.t_sb
}

// Return oversampling values osrs_t, osrs_p, osrs_h
func (this *bme280) Oversample() (sensors.BME280Oversample, sensors.BME280Oversample, sensors.BME280Oversample) {
	return this.osrs_t, this.osrs_p, this.osrs_h
}

// Return current measuring and updating value
func (this *bme280) Status() (bool, bool, error) {
	if status, err := this.ReadRegister_Uint8(BME280_REG_STATUS); err != nil {
		return false, false, err
	} else {
		measuring := ((status>>3)&0x01 != 0x00)
		updating := (status&0x01 != 0x00)
		return measuring, updating, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE - SET

// Reset the device using the complete power-on-reset procedure
func (this *bme280) SoftReset() error {
	this.log.Debug2("<sensors.BME280.SoftReset>")
	if err := this.WriteRegister_Uint8(BME280_REG_SOFTRESET, BME280_SOFTRESET_VALUE); err != nil {
		return err
	}

	// Wait for no measuring or updating
	for {
		if measuring, updating, err := this.Status(); err != nil {
			return err
		} else if measuring == false && updating == false {
			break
		}
	}

	// Read registers and return
	return this.read_registers()
}

func (this *bme280) SetMode(mode BME280Mode) error {
	this.log.Debug2("<sensors.BME280.SetMode>{ mode=%v }", mode)
	ctrl_meas := uint8(this.osrs_t)<<5 | uint8(this.osrs_p)<<2 | uint8(mode)
	if err := this.WriteRegister_Uint8(BME280_REG_CONTROL, ctrl_meas); err != nil {
		return err
	} else if _, _, _, mode_read, err = this.readControl(); err != nil {
		return err
	} else if mode != mode_read {
		return fmt.Errorf("SetMode: Expected %v but read %v", mode, mode_read)
	} else {
		this.mode = mode
		return nil
	}
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
