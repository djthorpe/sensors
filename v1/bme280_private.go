/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package bosch /* import "github.com/djthorpe/gopi-hw/device/bosch" */

import (
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *BME280Driver) setup() error {
	var err error

	// Read Chip ID and Version
	this.chipid, this.version, err = this.readChipVersion()
	if err != nil {
		return err
	}

	// If Chip ID is unexpected, then return an error
	if this.chipid != BME280_CHIPID_DEFAULT {
		return ErrInvalidDevice
	}

	return this.read_registers()
}

func (this *BME280Driver) read_registers() error {
	var err error

	// Read calibration values
	this.calibration, err = this.readCalibration()
	if err != nil {
		return err
	}

	// Read control registers
	this.osrs_t, this.osrs_p, this.osrs_h, this.mode, err = this.readControl()
	if err != nil {
		return err
	}

	// Read config registers
	this.t_sb, this.filter, this.spi3w_en, err = this.readConfig()
	if err != nil {
		return err
	}

	return nil
}

func (this *BME280Driver) readCalibration() (*bme280Calibation, error) {
	var err error

	calibration := new(bme280Calibation)

	// Test data from the Bosch datasheet. Used to check calculations.
	// T1=27504 : T2=26435 : T3=-1000
	// P1=36477 : P2=-10685 : P3=3024
	// P4=2855 : P5=140 : P6=-7
	// P7=15500 : P8=-14600 : P9=6000

	// Test data from SPI
	// calibration=<bosch.BME280Calibation{ T1=28244 T2=26571 T3=50 P1=37759 P2=-10679 P3=3024 P4=8281 P5=-140 P6=-7 P7=9900 P8=-10230 P9=4285 H1=75 H2=353 H3=0 H4=340 H5=0 H6=30 }

	// Read temperature calibration values
	if calibration.T1, err = this.ReadRegister_Uint16(BME280_REG_DIG_T1); err != nil {
		return nil, err
	}
	if calibration.T2, err = this.ReadRegister_Int16(BME280_REG_DIG_T2); err != nil {
		return nil, err
	}
	if calibration.T3, err = this.ReadRegister_Int16(BME280_REG_DIG_T3); err != nil {
		return nil, err
	}

	// Read pressure calibration values
	if calibration.P1, err = this.ReadRegister_Uint16(BME280_REG_DIG_P1); err != nil {
		return nil, err
	}
	if calibration.P2, err = this.ReadRegister_Int16(BME280_REG_DIG_P2); err != nil {
		return nil, err
	}
	if calibration.P3, err = this.ReadRegister_Int16(BME280_REG_DIG_P3); err != nil {
		return nil, err
	}
	if calibration.P4, err = this.ReadRegister_Int16(BME280_REG_DIG_P4); err != nil {
		return nil, err
	}
	if calibration.P5, err = this.ReadRegister_Int16(BME280_REG_DIG_P5); err != nil {
		return nil, err
	}
	if calibration.P6, err = this.ReadRegister_Int16(BME280_REG_DIG_P6); err != nil {
		return nil, err
	}
	if calibration.P7, err = this.ReadRegister_Int16(BME280_REG_DIG_P7); err != nil {
		return nil, err
	}
	if calibration.P8, err = this.ReadRegister_Int16(BME280_REG_DIG_P8); err != nil {
		return nil, err
	}
	if calibration.P9, err = this.ReadRegister_Int16(BME280_REG_DIG_P9); err != nil {
		return nil, err
	}

	// Read humidity calibration values
	if calibration.H1, err = this.ReadRegister_Uint8(BME280_REG_DIG_H1); err != nil {
		return nil, err
	}
	if calibration.H2, err = this.ReadRegister_Int16(BME280_REG_DIG_H2); err != nil {
		return nil, err
	}
	if calibration.H3, err = this.ReadRegister_Uint8(BME280_REG_DIG_H3); err != nil {
		return nil, err
	}
	h41, err := this.ReadRegister_Uint8(BME280_REG_DIG_H4)
	if err != nil {
		return nil, err
	}
	h42, err := this.ReadRegister_Uint8(BME280_REG_DIG_H4 + 1)
	if err != nil {
		return nil, err
	}
	h51, err := this.ReadRegister_Uint8(BME280_REG_DIG_H5)
	if err != nil {
		return nil, err
	}
	h52, err := this.ReadRegister_Uint8(BME280_REG_DIG_H5 + 1)
	if err != nil {
		return nil, err
	}

	calibration.H4 = (int16(h41) << 4) | (int16(h42) & 0x0F)
	calibration.H5 = ((int16(h51) & 0xF0) >> 4) | int16(h52<<4)

	if calibration.H6, err = this.ReadRegister_Int8(BME280_REG_DIG_H6); err != nil {
		return nil, err
	}

	// Return calibration values
	return calibration, nil
}

func (this *BME280Driver) readChipVersion() (uint8, uint8, error) {
	chipid, err := this.ReadRegister_Uint8(BME280_REG_CHIPID)
	if err != nil {
		return 0, 0, err
	}
	version, err2 := this.ReadRegister_Uint8(BME280_REG_VERSION)
	if err2 != nil {
		return 0, 0, err2
	}
	return chipid, version, nil
}

// Read values osrs_t, osrs_p, osrs_h, mode
func (this *BME280Driver) readControl() (BME280Oversample, BME280Oversample, BME280Oversample, BME280Mode, error) {
	ctrl_meas, err := this.ReadRegister_Uint8(BME280_REG_CONTROL)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	ctrl_hum, err := this.ReadRegister_Uint8(BME280_REG_CONTROLHUMID)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	mode := BME280Mode(ctrl_meas) & BME280_MODE_MAX
	osrs_t := BME280Oversample(ctrl_meas>>5) & BME280_OVERSAMPLE_MAX
	osrs_p := BME280Oversample(ctrl_meas>>2) & BME280_OVERSAMPLE_MAX
	osrs_h := BME280Oversample(ctrl_hum) & BME280_OVERSAMPLE_MAX
	return osrs_t, osrs_p, osrs_h, mode, nil
}

// Read values t_sb, filter, spi3w_en
func (this *BME280Driver) readConfig() (BME280Standby, BME280Filter, bool, error) {
	config, err := this.ReadRegister_Uint8(BME280_REG_CONFIG)
	if err != nil {
		return 0, 0, false, err
	}
	filter := BME280Filter(config>>2) & BME280_FILTER_MAX
	t_sb := BME280Standby(config>>5) & BME280_STANDBY_MAX
	spi3w_en := bool(config&0x01 != 0x00)
	return t_sb, filter, spi3w_en, nil
}

// Read raw temperature value
func (this *BME280Driver) readTemperature() (int32, error) {
	msb, err := this.ReadRegister_Uint8(BME280_REG_TEMPDATA)
	if err != nil {
		return int32(0), err
	}
	lsb, err := this.ReadRegister_Uint8(BME280_REG_TEMPDATA + 1)
	if err != nil {
		return int32(0), err
	}
	xlsb, err := this.ReadRegister_Uint8(BME280_REG_TEMPDATA + 2)
	if err != nil {
		return int32(0), err
	}
	return ((int32(msb) << 16) | (int32(lsb) << 8) | int32(xlsb)) >> 4, nil
}

// Read raw pressure value, assumes temperature has already been read
func (this *BME280Driver) readPressure() (int32, error) {
	msb, err := this.ReadRegister_Uint8(BME280_REG_PRESSUREDATA)
	if err != nil {
		return int32(0), err
	}
	lsb, err := this.ReadRegister_Uint8(BME280_REG_PRESSUREDATA + 1)
	if err != nil {
		return int32(0), err
	}
	xlsb, err := this.ReadRegister_Uint8(BME280_REG_PRESSUREDATA + 2)
	if err != nil {
		return int32(0), err
	}
	return ((int32(msb) << 16) | (int32(lsb) << 8) | int32(xlsb)) >> 4, nil
}

// Read raw humidity value, assumes temperature has already been read
func (this *BME280Driver) readHumidity() (int32, error) {
	msb, err := this.ReadRegister_Uint8(BME280_REG_HUMIDDATA)
	if err != nil {
		return int32(0), err
	}
	lsb, err := this.ReadRegister_Uint8(BME280_REG_HUMIDDATA + 1)
	if err != nil {
		return int32(0), err
	}
	return (int32(msb) << 8) | int32(lsb), nil
}

////////////////////////////////////////////////////////////////////////////////
// CONVERT ADC SAMPLES TO FLOATS

// Return compensated temperature in Celcius, and the t_fine value
func (this *BME280Driver) toCelcius(adc int32) (float64, float64) {
	var1 := (float64(adc)/16384.0 - float64(this.calibration.T1)/1024.0) * float64(this.calibration.T2)
	var2 := ((float64(adc)/131072.0 - float64(this.calibration.T1)/8192.0) * (float64(adc)/131072.0 - float64(this.calibration.T1)/8192.0)) * float64(this.calibration.T3)
	t_fine := var1 + var2
	return t_fine / 5120.0, t_fine
}

// Return compensated pressure in Pascals
func (this *BME280Driver) toPascals(adc int32, t_fine float64) float64 {
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
func (this *BME280Driver) toRelativeHumidity(adc int32, t_fine float64) float64 {
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

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Convert bool to uint8
func to_uint8(value bool) uint8 {
	if value {
		return 1
	}
	return 0
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *bme280Calibation) String() string {
	return fmt.Sprintf("<bosch.BME280Calibation{ T1=%v T2=%v T3=%v P1=%v P2=%v P3=%v P4=%v P5=%v P6=%v P7=%v P8=%v P9=%v H1=%v H2=%v H3=%v H4=%v H5=%v H6=%v }", this.T1, this.T2, this.T3, this.P1, this.P2, this.P3, this.P4, this.P5, this.P6, this.P7, this.P8, this.P9, this.H1, this.H2, this.H3, this.H4, this.H5, this.H6)
}

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
