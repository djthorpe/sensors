/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package bme680

import (
	// Frameworks
	"github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// READ AND WRITE REGISTERS

func (this *bme680) ReadRegister_Uint8(reg register) (uint8, error) {
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
	return 0, sensors.ErrNoDevice
}

func (this *bme680) ReadRegister_Int8(reg register) (int8, error) {
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
	return 0, sensors.ErrNoDevice
}

func (this *bme680) ReadRegister_Uint16(reg register) (uint16, error) {
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
	return 0, sensors.ErrNoDevice
}

func (this *bme680) ReadRegister_Int16(reg register) (int16, error) {
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
	return 0, sensors.ErrNoDevice
}

func (this *bme680) WriteRegister_Uint8(reg register, data uint8) error {
	if this.spi != nil {
		send := []byte{uint8(reg & BME680_REG_SPI_WRITE), data}
		_, err := this.spi.Transfer(send)
		return err
	}
	if this.i2c != nil {
		return this.i2c.WriteUint8(uint8(reg), data)
	}
	return sensors.ErrNoDevice
}
