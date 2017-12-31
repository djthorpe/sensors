/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package tsl2561

import (
	sensors "github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type register uint8

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// Registers
	REG_COMMAND    register = 0x80
	REG_CONTROL    register = 0x00
	REG_TIMING     register = 0x01
	REG_CHAN0_LOW  register = 0x0C
	REG_CHAN0_HIGH register = 0x0D
	REG_CHAN1_LOW  register = 0x0E
	REG_CHAN1_HIGH register = 0x0F
	REG_ID         register = 0x0A

	// Power on and off
	CONTROL_POWERON  uint8 = 0x03
	CONTROL_POWEROFF uint8 = 0x00

	// Word reading
	WORD_BIT register = 0x20
)

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *tsl2561) readChipVersion() (uint8, uint8, error) {
	if value, err := this.ReadRegister_Uint8(REG_ID); err != nil {
		return 0, 0, err
	} else {
		return (value & 0xF0) >> 4, (value & 0x0F), nil
	}
}

func (this *tsl2561) readTiming() (sensors.TSL2561Gain, sensors.TSL2561IntegrateTime, error) {
	if value, err := this.ReadRegister_Uint8(REG_TIMING); err != nil {
		return 0, 0, err
	} else {
		gain := sensors.TSL2561Gain(value>>4) & sensors.TSL2561_GAIN_MAX
		timing := sensors.TSL2561IntegrateTime(value>>0) & sensors.TSL2561_INTEGRATETIME_MAX
		return gain, timing, nil
	}
}

func (this *tsl2561) writeTiming(gain sensors.TSL2561Gain, integrate_time sensors.TSL2561IntegrateTime) error {
	value := uint8(gain&sensors.TSL2561_GAIN_MAX)<<4 | uint8(integrate_time&sensors.TSL2561_INTEGRATETIME_MAX)
	return this.WriteRegister_Uint8(REG_TIMING, value)
}

func (this *tsl2561) poweredOn() (bool, error) {
	if value, err := this.ReadRegister_Uint8(REG_CONTROL); err != nil {
		return false, err
	} else if value == CONTROL_POWERON {
		return true, nil
	} else if value == CONTROL_POWEROFF {
		return false, nil
	} else {
		return false, sensors.ErrUnexpectedResponse
	}
}

func (this *tsl2561) powerOn() error {
	if err := this.WriteRegister_Uint8(REG_CONTROL, CONTROL_POWERON); err != nil {
		return err
	} else if value, err := this.ReadRegister_Uint8(REG_CONTROL); err != nil {
		return err
	} else if value != CONTROL_POWERON {
		return sensors.ErrUnexpectedResponse
	} else {
		return nil
	}
}

func (this *tsl2561) powerOff() error {
	if err := this.WriteRegister_Uint8(REG_CONTROL, CONTROL_POWEROFF); err != nil {
		return err
	} else if value, err := this.ReadRegister_Uint8(REG_CONTROL); err != nil {
		return err
	} else if value != CONTROL_POWEROFF {
		return sensors.ErrUnexpectedResponse
	} else {
		return nil
	}
}

func (this *tsl2561) getADC0Sample() (uint16, error) {
	if lsb, err := this.ReadRegister_Uint8(REG_CHAN0_LOW | WORD_BIT); err != nil {
		return 0, err
	} else if msb, err := this.ReadRegister_Uint8(REG_CHAN0_HIGH | WORD_BIT); err != nil {
		return 0, err
	} else {
		return uint16(lsb) | (uint16(msb) << 8), nil
	}
}

func (this *tsl2561) getADC1Sample() (uint16, error) {
	if lsb, err := this.ReadRegister_Uint8(REG_CHAN1_LOW | WORD_BIT); err != nil {
		return 0, err
	} else if msb, err := this.ReadRegister_Uint8(REG_CHAN1_HIGH | WORD_BIT); err != nil {
		return 0, err
	} else {
		return uint16(lsb) | (uint16(msb) << 8), nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// READ AND WRITE REGISTERS

func (this *tsl2561) ReadRegister_Uint8(reg register) (uint8, error) {
	if this.i2c != nil {
		recv, err := this.i2c.ReadUint8(uint8(reg | REG_COMMAND))
		if err != nil {
			return 0, err
		}
		return recv, nil
	} else {
		return 0, sensors.ErrNoDevice
	}
}

func (this *tsl2561) WriteRegister_Uint8(reg register, data uint8) error {
	if this.i2c != nil {
		return this.i2c.WriteUint8(uint8(reg|REG_COMMAND), data)
	}
	return sensors.ErrNoDevice
}
