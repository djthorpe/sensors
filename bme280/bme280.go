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
	"time"

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
	BME280_I2CSLAVE_DEFAULT   uint8  = 0x77
	BME280_SPI_MAXSPEEDHZ     uint32 = 5000
	BME280_CHIPID_DEFAULT     uint8  = 0x60
	BME280_SOFTRESET_VALUE    uint8  = 0xB6
	BME280_SKIPTEMP_VALUE     int32  = 0x80000
	BME280_SKIPPRESSURE_VALUE int32  = 0x80000
	BME280_SKIPHUMID_VALUE    int32  = 0x8000
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

// Return duty cycle
func (this *bme280) DutyCycle() time.Duration {
	return toMeasurementTime(this.osrs_t, this.osrs_p, this.osrs_h) + toStandbyTime(this.t_sb)
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
	// TODO: TIMEOUT
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

func (this *bme280) SetMode(mode sensors.BME280Mode) error {
	this.log.Debug2("<sensors.BME280.SetMode>{ mode=%v }", mode)
	ctrl_meas := uint8(this.osrs_t)<<5 | uint8(this.osrs_p)<<2 | uint8(mode)
	if err := this.WriteRegister_Uint8(BME280_REG_CONTROL, ctrl_meas); err != nil {
		return err
	} else if _, _, _, mode_read, err := this.readControl(); err != nil {
		return err
	} else if mode != mode_read {
		return fmt.Errorf("SetMode: Expected %v but read %v", mode, mode_read)
	} else {
		this.mode = mode
		return nil
	}
}

func (this *bme280) SetOversample(osrs_t, osrs_p, osrs_h sensors.BME280Oversample) error {
	this.log.Debug2("<sensors.BME280.SetOversample>{ osrs_t=%v osrs_p=%v osrs_h=%v }", osrs_t, osrs_p, osrs_h)

	// Write humidity value first
	if err := this.WriteRegister_Uint8(BME280_REG_CONTROLHUMID, uint8(osrs_h&sensors.BME280_OVERSAMPLE_MAX)); err != nil {
		return err
	}

	// Write pressure and temperature second
	ctrl_meas := uint8(osrs_t&sensors.BME280_OVERSAMPLE_MAX)<<5 | uint8(osrs_p&sensors.BME280_OVERSAMPLE_MAX)<<2 | uint8(this.mode&sensors.BME280_MODE_MAX)
	if err := this.WriteRegister_Uint8(BME280_REG_CONTROL, ctrl_meas); err != nil {
		return err
	}

	// Wait for no measuring or updating
	// TODO: TIMEOUT
	for {
		if measuring, updating, err := this.Status(); err != nil {
			return err
		} else if measuring == false && updating == false {
			break
		}
	}

	// Read values back
	if osrs_t_read, osrs_p_read, osrs_h_read, _, err := this.readControl(); err != nil {
		return err
	} else if osrs_t_read != osrs_t {
		return fmt.Errorf("SetOversample: Expected osrs_t=%v but read %v", osrs_t, osrs_t_read)
	} else if osrs_p_read != osrs_p {
		return fmt.Errorf("SetOversample: Expected osrs_p=%v but read %v", osrs_p, osrs_p_read)
	} else if osrs_h_read != osrs_h {
		return fmt.Errorf("SetOversample: Expected osrs_h=%v but read %v", osrs_h, osrs_h_read)
	} else {
		this.osrs_t = osrs_t_read
		this.osrs_p = osrs_p_read
		this.osrs_h = osrs_h_read
		return nil
	}
}

func (this *bme280) SetFilter(filter sensors.BME280Filter) error {
	this.log.Debug2("<sensors.BME280.SetFilter>{ filter=%v }", filter)
	config := uint8(this.t_sb)<<5 | uint8(filter)<<2 | to_uint8(this.spi3w_en)
	if err := this.WriteRegister_Uint8(BME280_REG_CONFIG, config); err != nil {
		return err
	}

	// Read values back
	if _, filter_read, _, err := this.readConfig(); err != nil {
		return err
	} else if filter != filter_read {
		return fmt.Errorf("SetFilter: Expected filter=%v but read %v", filter, filter_read)
	} else {
		this.filter = filter_read
		return nil
	}
}

func (this *bme280) SetStandby(t_sb sensors.BME280Standby) error {
	this.log.Debug2("<sensors.BME280.SetStandby>{ t_sb=%v }", t_sb)
	config := uint8(t_sb)<<5 | uint8(this.filter)<<2 | to_uint8(this.spi3w_en)
	if err := this.WriteRegister_Uint8(BME280_REG_CONFIG, config); err != nil {
		return err
	}

	// Read values back
	if t_sb_read, _, _, err := this.readConfig(); err != nil {
		return err
	} else if t_sb != t_sb_read {
		return fmt.Errorf("SetStandby: Expected t_sb=%v but read %v", t_sb, t_sb_read)
	} else {
		this.t_sb = t_sb
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// GET SAMPLE DATA

// Return raw sample data for temperature, pressure and humidity
func (this *bme280) ReadSample() (float64, float64, float64, error) {
	this.log.Debug2("<sensors.BME280.ReadSample>{}")

	// Wait for no measuring or updating
	// TODO: Timeout
	for {
		if measuring, updating, err := this.Status(); err != nil {
			return 0, 0, 0, err
		} else if measuring == false && updating == false {
			break
		}
	}

	// Set mode of operation if we're in FORCED or SLEEP mode, and wait until we
	// can read the measurement for the correct amount of time
	if this.mode == sensors.BME280_MODE_FORCED || this.mode == sensors.BME280_MODE_SLEEP {
		if err := this.SetMode(sensors.BME280_MODE_FORCED); err != nil {
			return 0, 0, 0, err
		}
		// Wait until we can measure
		this.log.Debug2("In forced mode, measurement time = %v", toMeasurementTime(this.osrs_t, this.osrs_p, this.osrs_h))
		time.Sleep(toMeasurementTime(this.osrs_t, this.osrs_p, this.osrs_h))
	}

	// Read temperature, return error if temperature reading is skipped
	adc_t, err := this.readTemperature()
	if err != nil {
		return 0, 0, 0, err
	}
	if adc_t == BME280_SKIPTEMP_VALUE {
		return 0, 0, 0, sensors.ErrSampleSkipped
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
// sensors.BME280_PRESSURE_SEALEVEL for sealevel
func (this *bme280) AltitudeForPressure(atmospheric, sealevel float64) float64 {
	return 44330.0 * (1.0 - math.Pow(atmospheric/sealevel, (1.0/5.255)))
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Convert bool to uint8
func to_uint8(value bool) uint8 {
	if value {
		return 1
	}
	return 0
}

func toOversampleNumber(value sensors.BME280Oversample) float64 {
	switch value {
	case sensors.BME280_OVERSAMPLE_SKIP:
		return 0
	case sensors.BME280_OVERSAMPLE_1:
		return 1
	case sensors.BME280_OVERSAMPLE_2:
		return 2
	case sensors.BME280_OVERSAMPLE_4:
		return 4
	case sensors.BME280_OVERSAMPLE_8:
		return 8
	case sensors.BME280_OVERSAMPLE_16:
		return 16
	default:
		return 0
	}
}

func toMeasurementTime(osrs_t, osrs_p, osrs_h sensors.BME280Oversample) time.Duration {
	// Measurement Time as per BME280 datasheet section 9.1
	time_ms := 1.25
	if osrs_t != sensors.BME280_OVERSAMPLE_SKIP {
		time_ms += toOversampleNumber(osrs_t) * 2.3
	}
	if osrs_p != sensors.BME280_OVERSAMPLE_SKIP {
		time_ms += toOversampleNumber(osrs_p)*2.3 + 0.575
	}
	if osrs_h != sensors.BME280_OVERSAMPLE_SKIP {
		time_ms += toOversampleNumber(osrs_h)*2.4 + 0.575
	}
	return time.Millisecond * time.Duration(time_ms)
}

func toStandbyTime(value sensors.BME280Standby) time.Duration {
	switch value {
	case sensors.BME280_STANDBY_0P5MS:
		return time.Microsecond * 500
	case sensors.BME280_STANDBY_62P5MS:
		return time.Microsecond * 62500
	case sensors.BME280_STANDBY_125MS:
		return time.Millisecond * 125
	case sensors.BME280_STANDBY_250MS:
		return time.Millisecond * 250
	case sensors.BME280_STANDBY_500MS:
		return time.Millisecond * 500
	case sensors.BME280_STANDBY_1000MS:
		return time.Millisecond * 1000
	case sensors.BME280_STANDBY_10MS:
		return time.Millisecond * 10
	case sensors.BME280_STANDBY_20MS:
		return time.Millisecond * 20
	default:
		return 0
	}
}

////////////////////////////////////////////////////////////////////////////////
// CONVERT ADC SAMPLES TO FLOATS

// Return compensated temperature in Celcius, and the t_fine value
func (this *bme280) toCelcius(adc int32) (float64, float64) {
	var1 := (float64(adc)/16384.0 - float64(this.calibration.T1)/1024.0) * float64(this.calibration.T2)
	var2 := ((float64(adc)/131072.0 - float64(this.calibration.T1)/8192.0) * (float64(adc)/131072.0 - float64(this.calibration.T1)/8192.0)) * float64(this.calibration.T3)
	t_fine := var1 + var2
	return t_fine / 5120.0, t_fine
}

// Return compensated pressure in Pascals
func (this *bme280) toPascals(adc int32, t_fine float64) float64 {
	// Skip and return 0 if sample value is not valid
	if adc == 0 {
		return 0
	}

	var1 := t_fine/2.0 - 64000.0
	var2 := var1 * var1 * float64(this.calibration.P6) / 32768.0
	var2 = var2 + var1*float64(this.calibration.P5)*2.0
	var2 = var2/4.0 + float64(this.calibration.P4)*65536.0
	var1 = (float64(this.calibration.P3)*var1*var1/524288.0 + float64(this.calibration.P2)*var1) / 524288.0
	var1 = (1.0 + var1/32768.0) * float64(this.calibration.P1)
	if var1 == 0 {
		return 0 // avoid exception caused by division by zero
	}
	// Calculate value
	p := 1048576.0 - float64(adc)
	p = ((p - var2/4096.0) * 6250.0) / var1
	var1 = float64(this.calibration.P9) * p * p / 2147483648.0
	var2 = p * float64(this.calibration.P8) / 32768.0
	p = p + (var1+var2+float64(this.calibration.P7))/16.0
	return p / 100.0
}

// Return compensated humidity in %RH
func (this *bme280) toRelativeHumidity(adc int32, t_fine float64) float64 {
	// Skip and return 0 if sample value is not valid
	if adc == 0 {
		return 0
	}
	// Calculate value
	h := t_fine - 76800.0
	h = (float64(adc) - (float64(this.calibration.H4)*64.0 + float64(this.calibration.H5)/16384.8*h)) * (float64(this.calibration.H2) / 65536.0 * (1.0 + float64(this.calibration.H6)/67108864.0*h*(1.0+float64(this.calibration.H3)/67108864.0*h)))
	h = h * (1.0 - float64(this.calibration.H1)*h/524288.0)
	// Trim value between 0-100%
	switch {
	case h > 100.0:
		return 100.0
	case h < 0.0:
		return 0.0
	default:
		return h
	}
}
