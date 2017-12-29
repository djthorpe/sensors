/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package bme280

import (
	sensors "github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// BME280 registers and modes
type register uint8

// BME280 calibration
type calibation struct {
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
	BME280_REG_DIG_T1       register = 0x88
	BME280_REG_DIG_T2       register = 0x8A
	BME280_REG_DIG_T3       register = 0x8C
	BME280_REG_DIG_P1       register = 0x8E
	BME280_REG_DIG_P2       register = 0x90
	BME280_REG_DIG_P3       register = 0x92
	BME280_REG_DIG_P4       register = 0x94
	BME280_REG_DIG_P5       register = 0x96
	BME280_REG_DIG_P6       register = 0x98
	BME280_REG_DIG_P7       register = 0x9A
	BME280_REG_DIG_P8       register = 0x9C
	BME280_REG_DIG_P9       register = 0x9E
	BME280_REG_DIG_H1       register = 0xA1
	BME280_REG_DIG_H2       register = 0xE1
	BME280_REG_DIG_H3       register = 0xE3
	BME280_REG_DIG_H4       register = 0xE4
	BME280_REG_DIG_H5       register = 0xE5
	BME280_REG_DIG_H6       register = 0xE7
	BME280_REG_CHIPID       register = 0xD0
	BME280_REG_VERSION      register = 0xD1
	BME280_REG_SOFTRESET    register = 0xE0
	BME280_REG_CAL26        register = 0xE1 // R calibration stored in 0xE1-0xF0
	BME280_REG_CONTROLHUMID register = 0xF2
	BME280_REG_STATUS       register = 0xF3
	BME280_REG_CONTROL      register = 0xF4
	BME280_REG_CONFIG       register = 0xF5
	BME280_REG_PRESSUREDATA register = 0xF7
	BME280_REG_TEMPDATA     register = 0xFA
	BME280_REG_HUMIDDATA    register = 0xFD
)

const (
	// Write mask
	BME280_REG_SPI_WRITE register = 0x7F
)

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *bme280) readCalibration() (*calibation, error) {
	var err error
	calibration := new(calibation)

	// Test data from the Bosch datasheet. Used to check calculations.
	// T1=27504 : T2=26435 : T3=-1000
	// P1=36477 : P2=-10685 : P3=3024
	// P4=2855 : P5=140 : P6=-7
	// P7=15500 : P8=-14600 : P9=6000

	// Test data from SPI
	// calibration=<calibration{ T1=28244 T2=26571 T3=50 P1=37759 P2=-10679 P3=3024 P4=8281 P5=-140 P6=-7 P7=9900 P8=-10230 P9=4285 H1=75 H2=353 H3=0 H4=340 H5=0 H6=30 }

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

func (this *bme280) readChipVersion() (uint8, uint8, error) {
	if chipid, err := this.ReadRegister_Uint8(BME280_REG_CHIPID); err != nil {
		return 0, 0, err
	} else if version, err := this.ReadRegister_Uint8(BME280_REG_VERSION); err != nil {
		return 0, 0, err
	} else {
		return chipid, version, nil
	}
}

// Read values osrs_t, osrs_p, osrs_h, mode
func (this *bme280) readControl() (sensors.BME280Oversample, sensors.BME280Oversample, sensors.BME280Oversample, sensors.BME280Mode, error) {
	if ctrl_meas, err := this.ReadRegister_Uint8(BME280_REG_CONTROL); err != nil {
		return 0, 0, 0, 0, err
	} else if ctrl_hum, err := this.ReadRegister_Uint8(BME280_REG_CONTROLHUMID); err != nil {
		return 0, 0, 0, 0, err
	} else {
		mode := sensors.BME280Mode(ctrl_meas) & sensors.BME280_MODE_MAX
		osrs_t := sensors.BME280Oversample(ctrl_meas>>5) & sensors.BME280_OVERSAMPLE_MAX
		osrs_p := sensors.BME280Oversample(ctrl_meas>>2) & sensors.BME280_OVERSAMPLE_MAX
		osrs_h := sensors.BME280Oversample(ctrl_hum) & sensors.BME280_OVERSAMPLE_MAX
		return osrs_t, osrs_p, osrs_h, mode, nil
	}
}

// Read values t_sb, filter, spi3w_en
func (this *bme280) readConfig() (sensors.BME280Standby, sensors.BME280Filter, bool, error) {
	if config, err := this.ReadRegister_Uint8(BME280_REG_CONFIG); err != nil {
		return 0, 0, false, err
	} else {
		filter := sensors.BME280Filter(config>>2) & sensors.BME280_FILTER_MAX
		t_sb := sensors.BME280Standby(config>>5) & sensors.BME280_STANDBY_MAX
		spi3w_en := bool(config&0x01 != 0x00)
		return t_sb, filter, spi3w_en, nil
	}
}

// Read raw temperature value
func (this *bme280) readTemperature() (int32, error) {
	if msb, err := this.ReadRegister_Uint8(BME280_REG_TEMPDATA); err != nil {
		return int32(0), err
	} else if lsb, err := this.ReadRegister_Uint8(BME280_REG_TEMPDATA + 1); err != nil {
		return int32(0), err
	} else if xlsb, err := this.ReadRegister_Uint8(BME280_REG_TEMPDATA + 2); err != nil {
		return int32(0), err
	} else {
		return ((int32(msb) << 16) | (int32(lsb) << 8) | int32(xlsb)) >> 4, nil
	}
}

// Read raw pressure value, assumes temperature has already been read
func (this *bme280) readPressure() (int32, error) {
	if msb, err := this.ReadRegister_Uint8(BME280_REG_PRESSUREDATA); err != nil {
		return int32(0), err
	} else if lsb, err := this.ReadRegister_Uint8(BME280_REG_PRESSUREDATA + 1); err != nil {
		return int32(0), err
	} else if xlsb, err := this.ReadRegister_Uint8(BME280_REG_PRESSUREDATA + 2); err != nil {
		return int32(0), err
	} else {
		return ((int32(msb) << 16) | (int32(lsb) << 8) | int32(xlsb)) >> 4, nil
	}
}

// Read raw humidity value, assumes temperature has already been read
func (this *bme280) readHumidity() (int32, error) {
	if msb, err := this.ReadRegister_Uint8(BME280_REG_HUMIDDATA); err != nil {
		return int32(0), err
	} else if lsb, err := this.ReadRegister_Uint8(BME280_REG_HUMIDDATA + 1); err != nil {
		return int32(0), err
	} else {
		return (int32(msb) << 8) | int32(lsb), nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// READ AND WRITE REGISTERS

func (this *bme280) ReadRegister_Uint8(reg register) (uint8, error) {
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

func (this *bme280) ReadRegister_Int8(reg register) (int8, error) {
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

func (this *bme280) ReadRegister_Uint16(reg register) (uint16, error) {
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

func (this *bme280) ReadRegister_Int16(reg register) (int16, error) {
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

func (this *bme280) WriteRegister_Uint8(reg register, data uint8) error {
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
