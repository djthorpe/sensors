/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package rfm69

import (
	"time"

	"github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	register uint8
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// RFM69 Registers
	RFM_REG_FIFO          register = 0x00 /* FIFO Read/Write Access */
	RFM_REG_OPMODE        register = 0x01 /* Operating modes of the transceiver */
	RFM_REG_DATAMODUL     register = 0x02 /* Data operation mode and modulation settings */
	RFM_REG_BITRATEMSB    register = 0x03 /* Bit Rate setting, most significant bits */
	RFM_REG_BITRATELSB    register = 0x04 /* Bit Rate setting, least significant bits */
	RFM_REG_FDEVMSB       register = 0x05 /* Frequency deviation setting, most significant bits */
	RFM_REG_FDEVLSB       register = 0x06 /* Frequency deviation setting, least significant bits */
	RFM_REG_FRFMSB        register = 0x07 /* RF Carrier Frequency, most significant bits */
	RFM_REG_FRFMID        register = 0x08 /* RF Carrier Frequency, intermediate bits */
	RFM_REG_FRFLSB        register = 0x09 /* RF Carrier Frequency, least significant bits */
	RFM_REG_OSC1          register = 0x0A /* RC Oscillators Settings */
	RFM_REG_AFCCTRL       register = 0x0B /* AFC Control in low modulation index situations */
	RFM_REG_LISTEN1       register = 0x0D /* Listen mode settings */
	RFM_REG_LISTEN2       register = 0x0E /* Listen mode idle duration */
	RFM_REG_LISTEN3       register = 0x0F /* Listen mode Rx duration */
	RFM_REG_VERSION       register = 0x10 /* Module version */
	RFM_REG_PALEVEL       register = 0x11 /* PA selection and output power control */
	RFM_REG_PARAMP        register = 0x12 /* Control of the PA ramp time in FSK mode */
	RFM_REG_OCP           register = 0x13 /* Over Current Protection control */
	RFM_REG_LNA           register = 0x18 /* LNA Settings */
	RFM_REG_RXBW          register = 0x19 /* Channel Filter BW Control */
	RFM_REG_AFCBW         register = 0x1A // Channel Filter BW control during the AFC routine
	RFM_REG_OOKPEAK       register = 0x1B // OOK demodulator selection and control in peak mode
	RFM_REG_OOKAVG        register = 0x1C // Average threshold control of the OOK demodulator
	RFM_REG_OOKFIX        register = 0x1D // Fixed threshold control of the OOK demodulator
	RFM_REG_AFCFEI        register = 0x1E // AFC and FEI control and status
	RFM_REG_AFCMSB        register = 0x1F // MSB of the frequency correction of the AFC
	RFM_REG_AFCLSB        register = 0x20 // LSB of the frequency correction of the AFC
	RFM_REG_FEIMSB        register = 0x21 // MSB of the calculated frequency error
	RFM_REG_FEILSB        register = 0x22 // LSB of the calculated frequency error
	RFM_REG_RSSICONFIG    register = 0x23 // RSSI-related settings
	RFM_REG_RSSIVALUE     register = 0x24 // RSSI value in dBm
	RFM_REG_DIOMAPPING1   register = 0x25 // Mapping of pins DIO0 to DIO3
	RFM_REG_DIOMAPPING2   register = 0x26 // Mapping of pins DIO4 and DIO5, ClkOut frequency
	RFM_REG_IRQFLAGS1     register = 0x27 // Status register: PLL Lock state, Timeout, RSSI > Threshold...
	RFM_REG_IRQFLAGS2     register = 0x28 // Status register: FIFO handling flags...
	RFM_REG_RSSITHRESH    register = 0x29 // RSSI Threshold control
	RFM_REG_RXTIMEOUT1    register = 0x2A // Timeout duration between Rx request and RSSI detection
	RFM_REG_RXTIMEOUT2    register = 0x2B // Timeout duration between RSSI detection and PayloadReady
	RFM_REG_PREAMBLEMSB   register = 0x2C // Preamble length, MSB
	RFM_REG_PREAMBLELSB   register = 0x2D // Preamble length, LSB
	RFM_REG_SYNCCONFIG    register = 0x2E // Sync Word Recognition control
	RFM_REG_SYNCVALUE1    register = 0x2F // Sync Word bytes, 1 through 8
	RFM_REG_SYNCVALUE2    register = 0x30
	RFM_REG_SYNCVALUE3    register = 0x31
	RFM_REG_SYNCVALUE4    register = 0x32
	RFM_REG_SYNCVALUE5    register = 0x33
	RFM_REG_SYNCVALUE6    register = 0x34
	RFM_REG_SYNCVALUE7    register = 0x35
	RFM_REG_SYNCVALUE8    register = 0x36
	RFM_REG_PACKETCONFIG1 register = 0x37 // Packet mode settings
	RFM_REG_PAYLOADLENGTH register = 0x38 // Payload length setting
	RFM_REG_NODEADRS      register = 0x39 // Node address
	RFM_REG_BROADCASTADRS register = 0x3A // Broadcast address
	RFM_REG_AUTOMODES     register = 0x3B // Auto modes settings
	RFM_REG_FIFOTHRESH    register = 0x3C // Fifo threshold, Tx start condition
	RFM_REG_PACKETCONFIG2 register = 0x3D // Packet mode settings
	RFM_REG_AESKEY1       register = 0x3E // 16 bytes of the cypher key
	RFM_REG_AESKEY2       register = 0x3F
	RFM_REG_AESKEY3       register = 0x40
	RFM_REG_AESKEY4       register = 0x41
	RFM_REG_AESKEY5       register = 0x42
	RFM_REG_AESKEY6       register = 0x43
	RFM_REG_AESKEY7       register = 0x44
	RFM_REG_AESKEY8       register = 0x45
	RFM_REG_AESKEY9       register = 0x46
	RFM_REG_AESKEY10      register = 0x47
	RFM_REG_AESKEY11      register = 0x48
	RFM_REG_AESKEY12      register = 0x49
	RFM_REG_AESKEY13      register = 0x4A
	RFM_REG_AESKEY14      register = 0x4B
	RFM_REG_AESKEY15      register = 0x4C
	RFM_REG_AESKEY16      register = 0x4D
	RFM_REG_TEMP1         register = 0x4E // Temperature Sensor control
	RFM_REG_TEMP2         register = 0x4F // Temperature readout
	RFM_REG_TEST          register = 0x50 // Internal test registers
	RFM_REG_TESTLNA       register = 0x58 // Sensitivity boost
	RFM_REG_TESTPA1       register = 0x5A // High Power PA settings
	RFM_REG_TESTPA2       register = 0x5C // High Power PA settings
	RFM_REG_TESTDAGC      register = 0x6F // Fading Margin Improvement
	RFM_REG_TESTAFC       register = 0x71 // AFC offset for low modulation index AFC
	RFM_REG_MAX           register = 0x7F // Last possible register value
	RFM_REG_WRITE         register = 0x80 // Write bit
)

const (
	// RFM69 IRQ Flags
	RFM_IRQFLAGS1_MODEREADY        uint8 = 0x80 // Mode has changed
	RFM_IRQFLAGS1_RXREADY          uint8 = 0x40
	RFM_IRQFLAGS1_TXREADY          uint8 = 0x20
	RFM_IRQFLAGS1_PLLLOCK          uint8 = 0x10
	RFM_IRQFLAGS1_RSSI             uint8 = 0x08
	RFM_IRQFLAGS1_TIMEOUT          uint8 = 0x04
	RFM_IRQFLAGS1_AUTOMODE         uint8 = 0x02
	RFM_IRQFLAGS1_SYNCADDRESSMATCH uint8 = 0x01

	RFM_IRQFLAGS2_CRCOK        uint8 = 0x02
	RFM_IRQFLAGS2_PAYLOADREADY uint8 = 0x04
	RFM_IRQFLAGS2_PACKETSENT   uint8 = 0x08
	RFM_IRQFLAGS2_FIFOOVERRUN  uint8 = 0x10
	RFM_IRQFLAGS2_FIFOLEVEL    uint8 = 0x20
	RFM_IRQFLAGS2_FIFONOTEMPTY uint8 = 0x40
	RFM_IRQFLAGS2_FIFOFULL     uint8 = 0x80
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
	this.log.Debug2("<sensors.RFM69>readreg_uint8_array{ reg=%v length=%v recv=%v }", reg, length, recv[1:])
	if err != nil {
		return nil, err
	}
	return recv[1:], nil
}

func (this *rfm69) readreg_uint16(reg register) (uint16, error) {
	recv, err := this.spi.Transfer([]byte{byte(reg & RFM_REG_MAX), 0, 0})
	this.log.Debug2("<sensors.RFM69>readreg_uint16{ reg=%v recv=%v }", reg, recv[1:])
	if err != nil {
		return 0, err
	}
	return uint16(recv[1])<<8 | uint16(recv[2]), nil
}

func (this *rfm69) readreg_int16(reg register) (int16, error) {
	recv, err := this.spi.Transfer([]byte{byte(reg & RFM_REG_MAX), 0, 0})
	this.log.Debug2("<sensors.RFM69>readreg_uint16{ reg=%v recv=%v }", reg, recv[1:])
	if err != nil {
		return 0, err
	}
	return int16(uint16(recv[1])<<8 | uint16(recv[2])), nil
}

func (this *rfm69) readreg_uint24(reg register) (uint32, error) {
	recv, err := this.spi.Transfer([]byte{byte(reg & RFM_REG_MAX), 0, 0, 0})
	this.log.Debug2("<sensors.RFM69>readreg_uint24{ reg=%v recv=%v }", reg, recv[1:])
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

////////////////////////////////////////////////////////////////////////////////
// IRQ FLAGS

func wait_for_condition(callback func() (bool, error), condition bool, timeout time.Duration) error {
	timeout_chan := time.After(timeout)
	for {
		select {
		case <-timeout_chan:
			return sensors.ErrDeviceTimeout
		default:
			r, err := callback()
			if err != nil {
				return err
			}
			if r == condition {
				return nil
			}
			time.Sleep(time.Millisecond * 100)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// RFM_REG_OPMODE

// Read device mode, listen_on, sequencer_off
func (this *rfm69) getOpMode() (sensors.RFMMode, bool, bool, error) {
	data, err := this.readreg_uint8(RFM_REG_OPMODE)
	if err != nil {
		return 0, false, false, err
	}
	mode := sensors.RFMMode(data>>2) & sensors.RFM_MODE_MAX
	listen_on := to_uint8_bool((data >> 6) & 0x01)
	sequencer_off := to_uint8_bool((data >> 7) & 0x01)
	return mode, listen_on, sequencer_off, nil
}

// Write device_mode, listen_on, listen_abort and sequencer_off values
func (this *rfm69) setOpMode(device_mode sensors.RFMMode, listen_on bool, listen_abort bool, sequencer_off bool) error {
	value :=
		uint8(device_mode&sensors.RFM_MODE_MAX)<<2 |
			to_bool_uint8(listen_on)<<6 |
			to_bool_uint8(listen_abort)<<5 |
			to_bool_uint8(sequencer_off)<<7
	return this.writereg_uint8(RFM_REG_OPMODE, value)
}

////////////////////////////////////////////////////////////////////////////////
// RFM_REG_DATAMODUL

// Read data mode and modulation
func (this *rfm69) getDataModul() (sensors.RFMDataMode, sensors.RFMModulation, error) {
	data, err := this.readreg_uint8(RFM_REG_DATAMODUL)
	if err != nil {
		return 0, 0, err
	}
	data_mode := sensors.RFMDataMode(data>>5) & sensors.RFM_DATAMODE_MAX
	modulation := sensors.RFMModulation(data) & sensors.RFM_MODULATION_MAX
	return data_mode, modulation, nil
}

// Write data mode and modulation
func (this *rfm69) setDataModul(data_mode sensors.RFMDataMode, modulation sensors.RFMModulation) error {
	value :=
		uint8(data_mode&sensors.RFM_DATAMODE_MAX)<<5 |
			uint8(modulation&sensors.RFM_MODULATION_MAX)
	return this.writereg_uint8(RFM_REG_DATAMODUL, value)
}

////////////////////////////////////////////////////////////////////////////////
// RFM_REG_VERSION

// Read version
func (this *rfm69) getVersion() (uint8, error) {
	return this.readreg_uint8(RFM_REG_VERSION)
}

////////////////////////////////////////////////////////////////////////////////
// RFM_REG_AFCMSB, RFM_REG_AFCLSB

// Read Auto Frequency Correction value
func (this *rfm69) getAfc() (int16, error) {
	// TODO: Check LSB is also read?
	return this.readreg_int16(RFM_REG_AFCMSB)
}

////////////////////////////////////////////////////////////////////////////////
// RFM_REG_NODEADRS, RFM_REG_BROADCASTADRS

// Read node address
func (this *rfm69) getNodeAddress() (uint8, error) {
	if value, err := this.readreg_uint8(RFM_REG_NODEADRS); err != nil {
		return 0, err
	} else {
		return value, nil
	}
}

// Read broadcast address
func (this *rfm69) getBroadcastAddress() (uint8, error) {
	if value, err := this.readreg_uint8(RFM_REG_BROADCASTADRS); err != nil {
		return 0, err
	} else {
		return value, nil
	}
}

// Write node address
func (this *rfm69) setNodeAddress(value uint8) error {
	return this.writereg_uint8(RFM_REG_NODEADRS, value)
}

// Write broadcast address
func (this *rfm69) setBroadcastAddress(value uint8) error {
	return this.writereg_uint8(RFM_REG_BROADCASTADRS, value)
}

////////////////////////////////////////////////////////////////////////////////
// RFM_REG_IRQXFLAGS

func (this *rfm69) getIRQFlags1(mask uint8) (uint8, error) {
	value, err := this.readreg_uint8(RFM_REG_IRQFLAGS1)
	return value & mask, err
}
