/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package rfm69

import (
	"encoding/hex"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *rfm69) readreg_uint8(reg register) (uint8, error) {
	recv, err := this.spi.Transfer([]byte{byte(reg & RFM_REG_MAX), 0})
	this.log.Debug2("<sensors.RFM69>readreg_uint8{ reg=%v recv=0x%02X }", reg, recv[1:])
	if err != nil {
		return 0, err
	}
	return recv[1], nil
}

func (this *rfm69) readreg_uint8_array(reg register, length uint) ([]byte, error) {
	send := make([]byte, length+1)
	send[0] = byte(reg & RFM_REG_MAX)
	recv, err := this.spi.Transfer(send)
	this.log.Debug2("<sensors.RFM69>readreg_uint8_array{ reg=%v length=%v recv=0x%v }", reg, length, strings.ToUpper(hex.EncodeToString(recv[1:])))
	if err != nil {
		return nil, err
	}
	return recv[1:], nil
}

func (this *rfm69) readreg_uint16(reg register) (uint16, error) {
	recv, err := this.spi.Transfer([]byte{byte(reg & RFM_REG_MAX), 0, 0})
	this.log.Debug2("<sensors.RFM69>readreg_uint16{ reg=%v recv=0x%v }", reg, strings.ToUpper(hex.EncodeToString(recv[1:])))
	if err != nil {
		return 0, err
	}
	return uint16(recv[1])<<8 | uint16(recv[2]), nil
}

func (this *rfm69) readreg_int16(reg register) (int16, error) {
	recv, err := this.spi.Transfer([]byte{byte(reg & RFM_REG_MAX), 0, 0})
	this.log.Debug2("<sensors.RFM69>readreg_uint16{ reg=%v recv=0x%v }", reg, strings.ToUpper(hex.EncodeToString(recv[1:])))
	if err != nil {
		return 0, err
	}
	return int16(uint16(recv[1])<<8 | uint16(recv[2])), nil
}

func (this *rfm69) readreg_uint24(reg register) (uint32, error) {
	recv, err := this.spi.Transfer([]byte{byte(reg & RFM_REG_MAX), 0, 0, 0})
	this.log.Debug2("<sensors.RFM69>readreg_uint24{ reg=%v recv=0x%v }", reg, strings.ToUpper(hex.EncodeToString(recv[1:])))
	if err != nil {
		return 0, err
	}
	return uint32(recv[1])<<16 | uint32(recv[2])<<8 | uint32(recv[3]), nil
}

func (this *rfm69) writereg_uint8(reg register, data uint8) error {
	this.log.Debug2("<sensors.RFM69>writereg_uint8{ reg=%v data=0x%02X }", reg, data)
	return this.spi.Write([]byte{byte((reg & RFM_REG_MAX) | RFM_REG_WRITE), data})
}

func (this *rfm69) writereg_uint16(reg register, data uint16) error {
	this.log.Debug2("<sensors.RFM69>writereg_uint16{ reg=%v data=0x%04X }", reg, data)
	return this.spi.Write([]byte{
		byte((reg & RFM_REG_MAX) | RFM_REG_WRITE),
		uint8(data & 0xFF00 >> 8),
		uint8(data & 0xFF),
	})
}

func (this *rfm69) writereg_uint24(reg register, data uint32) error {
	this.log.Debug2("<sensors.RFM69>writereg_uint24{ reg=%v data=0x%06X }", reg, data)
	return this.spi.Write([]byte{
		byte((reg & RFM_REG_MAX) | RFM_REG_WRITE),
		uint8(data & 0xFF0000 >> 16),
		uint8(data & 0xFF00 >> 8),
		uint8(data & 0xFF),
	})
}

func (this *rfm69) writereg_uint8_array(reg register, data []byte) error {
	this.log.Debug2("<sensors.RFM69>writereg_uint8_array{ reg=%v data=%v }", reg, strings.ToUpper(hex.EncodeToString(data)))
	buf := append([]byte(nil), byte((reg&RFM_REG_MAX)|RFM_REG_WRITE))
	return this.spi.Write(append(buf, data...))
}

////////////////////////////////////////////////////////////////////////////////
// DATA CONVERSIONS

func to_uint8_bool(value uint8) bool {
	return (value != 0x00)
}

func to_bool_uint8(value bool) uint8 {
	if value {
		return 0x01
	} else {
		return 0x00
	}
}

func matches_byte_array(a1, a2 []byte) bool {
	if a1 == nil && a2 == nil {
		return true
	}
	if len(a1) != len(a2) {
		return false
	}
	for i := range a1 {
		if a1[i] != a2[i] {
			return false
		}
	}
	return true
}

func matches_byte(a []byte, b byte) bool {
	if a == nil {
		return false
	}
	for i := range a {
		if a[i] == b {
			return true
		}
	}
	return false
}
