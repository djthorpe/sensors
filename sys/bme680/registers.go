/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package bme680

////////////////////////////////////////////////////////////////////////////////
// TYPES

// BME680 registers and modes
type register uint8

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// Write mask
	BME680_REG_SPI_WRITE register = 0x7F
)

const (
	// https://github.com/BoschSensortec/BME680_driver/blob/master/bme680_defs.h
	BME680_REG_CHIP_ID    register = 0xD0
	BME680_REG_SOFT_RESET register = 0xE0
	BME680_REG_CTRL_GAS0  register = 0x70
	BME680_REG_CTRL_GAS1  register = 0x71
	BME680_REG_CTRL_HUM   register = 0x72
	BME680_REG_CTRL_MEAS  register = 0x74
	BME680_REG_CONFIG     register = 0x75
)

const (
	BME680_SOFTRESET_VALUE uint8 = 0xB6
)

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - READ

func (this *bme680) readChipId() (uint8, error) {
	if chipid, err := this.ReadRegister_Uint8(BME680_REG_CHIP_ID); err != nil {
		return 0, err
	} else {
		return chipid, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - WRITE

func (this *bme680) writeSoftReset() error {
	return this.WriteRegister_Uint8(BME680_REG_SOFT_RESET, BME680_SOFTRESET_VALUE)
}

/*
// Reset the device using the complete power-on-reset procedure
func (this *bme680) writeSoftReset() error {
	if err := this.WriteRegister_Uint8(BME680_REG_SOFT_RESET, BME680_SOFTRESET_VALUE); err != nil {
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
*/
