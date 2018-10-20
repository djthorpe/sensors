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

	// Frameworks
	"github.com/djthorpe/gopi"
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
func (this *rfm69) getAFC() (int16, error) {
	return this.readreg_int16(RFM_REG_AFCMSB)
}

// Read RegAfcCtrl register
func (this *rfm69) getAFCRoutine() (sensors.RFMAFCRoutine, error) {
	if afc_routine, err := this.readreg_uint8(RFM_REG_AFCCTRL); err != nil {
		return 0, err
	} else {
		return sensors.RFMAFCRoutine(afc_routine>>5) & sensors.RFM_AFCROUTINE_MASK, nil
	}
}

func (this *rfm69) setAFCRoutine(afc_routine sensors.RFMAFCRoutine) error {
	value := uint8(afc_routine&sensors.RFM_AFCROUTINE_MASK) << 5
	return this.writereg_uint8(RFM_REG_AFCCTRL, value)
}

// Read RFM_REG_AFCFEI - mode, afc_done, fei_done
func (this *rfm69) getAFCControl() (sensors.RFMAFCMode, bool, bool, error) {
	if value, err := this.readreg_uint8(RFM_REG_AFCFEI); err != nil {
		return 0, false, false, err
	} else {
		fei_done := to_uint8_bool(value & 0x40)
		afc_done := to_uint8_bool(value & 0x10)
		afc_mode := sensors.RFMAFCMode(value>>2) & sensors.RFM_AFCMODE_MASK
		return afc_mode, afc_done, fei_done, nil
	}
}

// Write RFM_REG_AFCFEI register
func (this *rfm69) setAFCControl(afc_mode sensors.RFMAFCMode, fei_start, afc_clear, afc_start bool) error {
	value :=
		to_bool_uint8(fei_start)<<5 |
			uint8(afc_mode&sensors.RFM_AFCMODE_MASK)<<2 |
			to_bool_uint8(afc_clear)<<1 |
			to_bool_uint8(afc_start)<<0
	return this.writereg_uint8(RFM_REG_AFCFEI, value)
}

////////////////////////////////////////////////////////////////////////////////
// RFM_REG_BITRATE

// Read bitrate (two bytes)
func (this *rfm69) getBitrate() (uint16, error) {
	return this.readreg_uint16(RFM_REG_BITRATEMSB)
}

// Write bitrate (two bytes)
func (this *rfm69) setBitrate(bitrate uint16) error {
	return this.writereg_uint16(RFM_REG_BITRATEMSB, bitrate)
}

////////////////////////////////////////////////////////////////////////////////
// RFM_REG_FRF

// Read FRF (three bytes)
func (this *rfm69) getFreqCarrier() (uint32, error) {
	if frf, err := this.readreg_uint24(RFM_REG_FRFMSB); err != nil {
		return 0, err
	} else {
		return frf & RFM_FRF_MAX, nil
	}
}

// Read FDEV (two bytes)
func (this *rfm69) getFreqDeviation() (uint16, error) {
	if fdev, err := this.readreg_uint16(RFM_REG_FDEVMSB); err != nil {
		return 0, err
	} else {
		return fdev & RFM_FDEV_MAX, nil
	}
}

// Write FRF (three bytes)
func (this *rfm69) setFreqCarrier(value uint32) error {
	// write MSB, MIDDLE and LSB in that order
	if err := this.writereg_uint8(RFM_REG_FRFMSB, uint8(value>>16)); err != nil {
		return err
	}
	if err := this.writereg_uint8(RFM_REG_FRFMID, uint8(value>>8)); err != nil {
		return err
	}
	if err := this.writereg_uint8(RFM_REG_FRFLSB, uint8(value)); err != nil {
		return err
	}
	return nil
}

// Write FRF (three bytes)
func (this *rfm69) setFreqDeviation(value uint16) error {
	return this.writereg_uint16(RFM_REG_FDEVMSB, value)
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
// RFM_REG_PREAMBLE

// Read Preamble size
func (this *rfm69) getPreambleSize() (uint16, error) {
	return this.readreg_uint16(RFM_REG_PREAMBLEMSB)
}

// Write Preamble size
func (this *rfm69) setPreambleSize(preamble_size uint16) error {
	return this.writereg_uint16(RFM_REG_PREAMBLEMSB, preamble_size)
}

////////////////////////////////////////////////////////////////////////////////
// RFM_REG_PAYLOADLENGTH

// Read Payload size
func (this *rfm69) getPayloadSize() (uint8, error) {
	return this.readreg_uint8(RFM_REG_PAYLOADLENGTH)
}

// Write Payload size
func (this *rfm69) setPayloadSize(payload_size uint8) error {
	return this.writereg_uint8(RFM_REG_PAYLOADLENGTH, payload_size)
}

////////////////////////////////////////////////////////////////////////////////
// RFM_REG_AESKEY, RFM_REG_SYNCKEY, RFM_REG_SYNCCONFIG

// Read RFM_REG_AESKEY, RFM_REG_SYNCKEY registers
func (this *rfm69) getAESKey() ([]byte, error) {
	if key, err := this.readreg_uint8_array(RFM_REG_AESKEY1, RFM_AESKEY_BYTES); err != nil {
		return nil, err
	} else {
		return key, nil
	}
}

// Read RFM_REG_SYNCVALUE register
func (this *rfm69) getSyncWord() ([]byte, error) {
	if key, err := this.readreg_uint8_array(RFM_REG_SYNCVALUE1, RFM_SYNCWORD_BYTES); err != nil {
		return nil, err
	} else {
		return key, nil
	}
}

// Write RFM_REG_SYNCVALUE register
func (this *rfm69) setSyncWord(word []byte) error {
	if len(word) > RFM_SYNCWORD_BYTES {
		return gopi.ErrBadParameter
	}
	return this.writereg_uint8_array(RFM_REG_SYNCVALUE1, word)
}

func (this *rfm69) setAESKey(aes_key []byte) error {
	if len(aes_key) != RFM_AESKEY_BYTES {
		this.log.Debug2("setAESKey: invalid AES key length (%v bytes, should be %v bytes)", len(aes_key), RFM_AESKEY_BYTES)
		return gopi.ErrBadParameter
	} else {
		return this.writereg_uint8_array(RFM_REG_AESKEY1, aes_key)
	}
}

// Read RFM_REG_SYNCCONFIG registers
// Returns SyncOn, FifoFillCondition, SyncSize, SyncTol
// Note sync_size is one less than the SyncSize
func (this *rfm69) getSyncConfig() (bool, bool, uint8, uint8, error) {
	if value, err := this.readreg_uint8(RFM_REG_SYNCCONFIG); err != nil {
		return false, false, 0, 0, err
	} else {
		return to_uint8_bool(value & 0x80), to_uint8_bool(value & 0x40), (uint8(value) >> 3) & 0x07, uint8(value & 0x07), nil
	}
}

// Write RFM_REG_SYNCCONFIG registers
// Note sync_size is one less than the SyncSize
func (this *rfm69) setSyncConfig(sync_on, fifo_fill_condition bool, sync_size, sync_tol uint8) error {
	value :=
		to_bool_uint8(sync_on)<<7 |
			to_bool_uint8(fifo_fill_condition)<<6 |
			(sync_size&0x07)<<3 |
			(sync_tol & 0x07)
	return this.writereg_uint8(RFM_REG_SYNCCONFIG, value)
}

////////////////////////////////////////////////////////////////////////////////
// RFM_REG_FIFOTHRESH

// Read RFM_REG_FIFOTHRESH register
func (this *rfm69) getFIFOThreshold() (sensors.RFMTXStart, uint8, error) {
	if value, err := this.readreg_uint8(RFM_REG_FIFOTHRESH); err != nil {
		return 0, 0, err
	} else {
		tx_start := sensors.RFMTXStart(value>>7) & sensors.RFM_TXSTART_MAX
		fifo_threshold := value & 0x7F
		return tx_start, fifo_threshold, nil
	}
}

// Write RFM_REG_FIFOTHRESH register
func (this *rfm69) setFIFOThreshold(tx_start sensors.RFMTXStart, fifo_threshold uint8) error {
	value := uint8(tx_start)<<7 | fifo_threshold&0x7F
	return this.writereg_uint8(RFM_REG_FIFOTHRESH, value)
}

////////////////////////////////////////////////////////////////////////////////
// RFM_REG_PACKETCONFIG1

// Read RegPacketConfig1 register - PacketFormat, PacketCoding, AddressFiltering, CRCOn, CRCAutoClearOff
func (this *rfm69) getPacketConfig1() (sensors.RFMPacketFormat, sensors.RFMPacketCoding, sensors.RFMPacketFilter, bool, bool, error) {
	if value, err := this.readreg_uint8(RFM_REG_PACKETCONFIG1); err != nil {
		return 0, 0, 0, false, false, err
	} else {
		packet_format := sensors.RFMPacketFormat((value >> 7) & 0x01)
		packet_coding := sensors.RFMPacketCoding(value>>5) & sensors.RFM_PACKET_CODING_MAX
		packet_filter := sensors.RFMPacketFilter(value) & sensors.RFM_PACKET_FILTER_MAX
		crc_on := to_uint8_bool(value & 0x10)
		crc_auto_clear_off := to_uint8_bool(value & 0x08)
		return packet_format, packet_coding, packet_filter, crc_on, crc_auto_clear_off, nil
	}
}

// Write RegPacketConfig1 register
func (this *rfm69) setPacketConfig1(packet_format sensors.RFMPacketFormat, packet_coding sensors.RFMPacketCoding, packet_filter sensors.RFMPacketFilter, crc_on bool, crc_auto_clear_off bool) error {
	value :=
		(uint8(packet_format)&0x01)<<7 |
			uint8(packet_coding&sensors.RFM_PACKET_CODING_MAX)<<5 |
			uint8(packet_filter&sensors.RFM_PACKET_FILTER_MAX) |
			to_bool_uint8(crc_on)<<4 |
			to_bool_uint8(crc_auto_clear_off)<<3
	return this.writereg_uint8(RFM_REG_PACKETCONFIG1, value)
}

////////////////////////////////////////////////////////////////////////////////
// RFM_REG_PACKETCONFIG2

// Read RegPacketConfig2 register - InterPacketRxDelay, AutoRxRestartOn, AesOn
func (this *rfm69) getPacketConfig2() (uint8, bool, bool, error) {
	if value, err := this.readreg_uint8(RFM_REG_PACKETCONFIG2); err != nil {
		return 0, false, false, err
	} else {
		rx_inter_packet_delay := uint8(value&0xF0) >> 4
		rx_auto_restart := to_uint8_bool(value & 0x02)
		aes_on := to_uint8_bool(value & 0x01)
		return rx_inter_packet_delay, rx_auto_restart, aes_on, nil
	}
}

// Write RegPacketConfig2 register
func (this *rfm69) setPacketConfig2(rx_inter_packet_delay uint8, rx_auto_restart bool, aes_on bool) error {
	value := (rx_inter_packet_delay&0x0F)<<4 | to_bool_uint8(rx_auto_restart)<<1 | to_bool_uint8(aes_on)
	return this.writereg_uint8(RFM_REG_PACKETCONFIG2, value)
}

////////////////////////////////////////////////////////////////////////////////
// RFM_REG_IRQXFLAGS

func (this *rfm69) getIRQFlags1(mask uint8) (uint8, error) {
	value, err := this.readreg_uint8(RFM_REG_IRQFLAGS1)
	return value & mask, err
}

func (this *rfm69) setIRQFlags2() error {
	// Set "FifoOverrun" flag to clear the FIFO buffer
	return this.writereg_uint8(RFM_REG_IRQFLAGS2, 0x10)
}

func (this *rfm69) getIRQFlags2(mask uint8) (uint8, error) {
	value, err := this.readreg_uint8(RFM_REG_IRQFLAGS2)
	if value&RFM_IRQFLAGS2_CRCOK != 0 {
		this.log.Debug2("RFM_IRQFLAGS2_CRCOK")
	}
	if value&RFM_IRQFLAGS2_CRCOK != 0 {
		this.log.Debug2("RFM_IRQFLAGS2_CRCOK")
	}
	if value&RFM_IRQFLAGS2_PAYLOADREADY != 0 {
		this.log.Debug2("RFM_IRQFLAGS2_PAYLOADREADY")
	}
	if value&RFM_IRQFLAGS2_PACKETSENT != 0 {
		this.log.Debug2("RFM_IRQFLAGS2_PACKETSENT")
	}
	if value&RFM_IRQFLAGS2_FIFOOVERRUN != 0 {
		this.log.Debug2("RFM_IRQFLAGS2_FIFOOVERRUN")
	}
	if value&RFM_IRQFLAGS2_FIFOLEVEL != 0 {
		this.log.Debug2("RFM_IRQFLAGS2_FIFOLEVEL")
	}
	if value&RFM_IRQFLAGS2_FIFONOTEMPTY != 0 {
		this.log.Debug2("RFM_IRQFLAGS2_FIFONOTEMPTY")
	}
	if value&RFM_IRQFLAGS2_FIFOFULL != 0 {
		this.log.Debug2("RFM_IRQFLAGS2_FIFOFULL")
	}
	return value & mask, err
}

////////////////////////////////////////////////////////////////////////////////
// RFM_REG_FIFO

func (this *rfm69) recvFIFOEmpty() (bool, error) {
	if fifo_not_empty, err := this.getIRQFlags2(RFM_IRQFLAGS2_FIFONOTEMPTY); err != nil {
		return false, err
	} else {
		return (fifo_not_empty != RFM_IRQFLAGS2_FIFONOTEMPTY), nil
	}
}

func (this *rfm69) irqFIFOLevel() (bool, error) {
	if fifo_level, err := this.getIRQFlags2(RFM_IRQFLAGS2_FIFOLEVEL); err != nil {
		return false, err
	} else {
		return (fifo_level != RFM_IRQFLAGS2_FIFOLEVEL), nil
	}
}

func (this *rfm69) recvCRCOk() (bool, error) {
	if crc_ok, err := this.getIRQFlags2(RFM_IRQFLAGS2_CRCOK); err != nil {
		return false, err
	} else {
		return (crc_ok == RFM_IRQFLAGS2_FIFONOTEMPTY), nil
	}
}

func (this *rfm69) recvPayloadReady() (bool, error) {
	if payload_ready, err := this.getIRQFlags2(RFM_IRQFLAGS2_PAYLOADREADY); err != nil {
		return false, err
	} else {
		return payload_ready == RFM_IRQFLAGS2_PAYLOADREADY, nil
	}
}

func (this *rfm69) recvPacketSent() (bool, error) {
	if packet_sent, err := this.getIRQFlags2(RFM_IRQFLAGS2_PACKETSENT); err != nil {
		return false, err
	} else {
		return packet_sent == RFM_IRQFLAGS2_PACKETSENT, nil
	}
}

func (this *rfm69) recvFIFO() ([]byte, error) {
	buffer := make([]byte, 0, RFM_FIFO_SIZE)
	for i := 0; i < RFM_FIFO_SIZE; i++ {
		if fifo_empty, err := this.recvFIFOEmpty(); err != nil {
			return nil, err
		} else if fifo_empty {
			break
		} else if value, err := this.readreg_uint8(RFM_REG_FIFO); err != nil {
			return nil, err
		} else {
			buffer = append(buffer, value)
		}
	}
	return buffer, nil
}

func (this *rfm69) writeFIFO(data []byte) error {
	return this.writereg_uint8_array(RFM_REG_FIFO, data)
}

////////////////////////////////////////////////////////////////////////////////
// RFM_REG_RXBW

func (this *rfm69) writeRXBW(value byte) error {
	return this.writereg_uint8(RFM_REG_RXBW, value)
}

////////////////////////////////////////////////////////////////////////////////
// RFM_REG_TEMP1, RFM_REG_TEMP2

// Get running bit
func (this *rfm69) getRegTemp1() (bool, error) {
	if value, err := this.readreg_uint8(RFM_REG_TEMP1); err != nil {
		return false, err
	} else {
		running := to_uint8_bool(value & 0x04)
		return running, nil
	}
}

// Set start measurement bit high
func (this *rfm69) setRegTemp1() error {
	return this.writereg_uint8(RFM_REG_TEMP1, 0x08)
}

// Read uncalibrated temperature
func (this *rfm69) getRegTemp2() (uint8, error) {
	return this.readreg_uint8(RFM_REG_TEMP2)
}

////////////////////////////////////////////////////////////////////////////////
// RFM_REG_RSSICONFIG, RFM_REG_RSSIVALUE

// Get RSSI done
func (this *rfm69) getRegRSSIDone() (bool, error) {
	if value, err := this.readreg_uint8(RFM_REG_RSSICONFIG); err != nil {
		return false, err
	} else {
		done := to_uint8_bool(value & 0x02)
		return done, nil
	}
}

// Set RSSI start
func (this *rfm69) setRegRSSIStart() error {
	return this.writereg_uint8(RFM_REG_RSSICONFIG, 0x01)
}

// Return RFM_REG_RSSIVALUE
func (this *rfm69) getRegRSSIValue() (uint8, error) {
	return this.readreg_uint8(RFM_REG_RSSIVALUE)
}

////////////////////////////////////////////////////////////////////////////////
// RFM_REG_LNA

// Read LNA settings - returns the impedance, the gain setting and the current gain value
func (this *rfm69) getRegLNA() (sensors.RFMLNAImpedance, sensors.RFMLNAGain, sensors.RFMLNAGain, error) {
	if value, err := this.readreg_uint8(RFM_REG_LNA); err != nil {
		return 0, 0, 0, err
	} else {
		lna_impedance := sensors.RFMLNAImpedance(value>>7) & sensors.RFM_LNA_IMPEDANCE_MAX
		set_gain := sensors.RFMLNAGain(value) & sensors.RFM_LNA_GAIN_MAX
		current_gain := sensors.RFMLNAGain(value>>3) & sensors.RFM_LNA_GAIN_MAX
		return lna_impedance, set_gain, current_gain, nil
	}
}

// Write LNA settings
func (this *rfm69) setRegLNA(impedance sensors.RFMLNAImpedance, gain sensors.RFMLNAGain) error {
	value :=
		uint8(impedance&sensors.RFM_LNA_IMPEDANCE_MAX)<<7 |
			uint8(gain&sensors.RFM_LNA_GAIN_MAX)
	return this.writereg_uint8(RFM_REG_LNA, value)
}

////////////////////////////////////////////////////////////////////////////////
// RFM_REG_RXBW

func (this *rfm69) getRegRXBW() (sensors.RFMRXBWFrequency, sensors.RFMRXBWCutoff, error) {
	if value, err := this.readreg_uint8(RFM_REG_RXBW); err != nil {
		return 0, 0, err
	} else {
		cutoff := sensors.RFMRXBWCutoff(value>>5) & sensors.RFM_RXBW_CUTOFF_MAX
		frequency := sensors.RFMRXBWFrequency(value) & sensors.RFM_RXBW_FREQUENCY_MAX
		return frequency, cutoff, nil
	}
}

func (this *rfm69) setRegRXBW(frequency sensors.RFMRXBWFrequency, cutoff sensors.RFMRXBWCutoff) error {
	value :=
		uint8(frequency&sensors.RFM_RXBW_FREQUENCY_MAX) | uint8(cutoff&sensors.RFM_RXBW_CUTOFF_MAX)<<5
	return this.writereg_uint8(RFM_REG_RXBW, value)
}
