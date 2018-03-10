/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package sensors

import (
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// RFM69 TYPES

type (
	RFMMode       uint8
	RFMDataMode   uint8
	RFMModulation uint8
)

////////////////////////////////////////////////////////////////////////////////
// RFM69 INTERFACE

type RFM69 interface {
	gopi.Driver

	// Mode, Data Mode and Modulation
	Mode() RFMMode
	DataMode() RFMDataMode
	SetMode(device_mode RFMMode) error
	SetDataMode(data_mode RFMDataMode) error
	Modulation() RFMModulation
	SetModulation(modulation RFMModulation) error

	// Addresses
	NodeAddress() uint8
	BroadcastAddress() uint8
	SetNodeAddress(value uint8) error
	SetBroadcastAddress(value uint8) error

	/*
		// Keys
		AESKey() RFMAESKey
		SyncKey() RFMSyncKey
		SetAESKey(key RFMAESKey) error
		SetSyncKey(key RFMSyncKey) error

		// AFC
		AFCMode() RFMAFCMode
		AFCRoutine() RFMAFCRoutine
		AFCHertz() float64
		SetAFCMode(afc_mode RFMAFCMode) error
		SetAFCRoutine(afc_routine RFMAFCRoutine) error

		// Packets
		PacketFormat() RFMPacketFormat
		PacketCoding() RFMPacketCoding
		PacketFilter() RFMPacketFilter
		SetPacketFormat(packet_format RFMPacketFormat) error
		SetPacketCoding(packet_coding RFMPacketCoding) error
		SetPacketFilter(packet_filter RFMPacketFilter) error

		// CRC
		CRCEnabled() bool
		CRCAutoClearOff() bool
		SetCRC(crc_enabled bool, crc_auto_clear_off bool) error

		// Payload & Preamble
		PreambleSize() uint16
		PayloadSize() uint8
		SetPreambleSize(preamble_size uint16) error
		SetPayloadSize(payload_size uint8) error

		// LNA
		LNAImpedance() RFMLNAImpedance
		LNAGain() RFMLNAGain
		LNACurrentGain() (RFMLNAGain, error)
		SetLNA(impedance RFMLNAImpedance, gain RFMLNAGain) error

		// OOK Parameters
		SetOOK(ook_threshold_type RFMOOKThresholdType, ook_threshold_step RFMOOKThresholdStep, ook_threshold_dec RFMOOKThresholdDecrement) error

		// FSK parameters
		SetBitrate(bits_per_second uint) error
		SetBitrateUint16(value uint16) error
		SetFreqDeviation(hertz uint) error
		SetFreqDeviationUint16(value uint16) error
		SetFreqCarrier(hertz uint) error
		SetFreqCarrierUint24(value uint32) error
		Bitrate() uint
		FreqDeviation() uint
		FreqCarrier() uint

		// FIFO
		FIFOFillCondition() bool
		FIFOThreshold() uint8
		SetFIFOFillCondition(fifo_fill_condition bool) error
		SetFIFOThreshold(fifo_threshold uint8) error

		// Other
		SetSequencer(value bool) error
		SequencerEnabled() bool
		SetListen(value bool) error
		ListenEnabled() bool
		SetSyncTolerance(sync_tol uint) error
		SyncTolerance() uint
		SetTXStart(tx_start RFMTXStart) error
		TXStart() RFMTXStart
		SetRXBW(mantissa RFMRXBWMantissa, exponent RFMRXBWExponent, cutoff RFMRXBWCutoff) error

		// Methods
		ReadFEIHertz() (float64, error)
		TriggerAFC() error
		ClearAFC() error
		ClearFIFO() error
		CalibrateRCOsc() error
		MeasureTemperature(calibration float32) (float32, error)
		ReadFIFO() ([]byte, error)
		ReceivePacket(timeout time.Duration) ([]byte, error)
	*/
}

////////////////////////////////////////////////////////////////////////////////
// RFM69 CONSTS

const (
	// RFM69 Mode
	RFM_MODE_SLEEP RFMMode = 0x00
	RFM_MODE_STDBY RFMMode = 0x01
	RFM_MODE_FS    RFMMode = 0x02
	RFM_MODE_TX    RFMMode = 0x03
	RFM_MODE_RX    RFMMode = 0x04
	RFM_MODE_MAX   RFMMode = 0x07
)

const (
	// RFM69 Data Mode
	RFM_DATAMODE_PACKET            RFMDataMode = 0x00
	RFM_DATAMODE_CONTINUOUS_NOSYNC RFMDataMode = 0x02
	RFM_DATAMODE_CONTINUOUS_SYNC   RFMDataMode = 0x03
	RFM_DATAMODE_MAX               RFMDataMode = 0x03
)

const (
	// RFM69 Modulation
	RFM_MODULATION_FSK        RFMModulation = 0x00 // 00000 FSK no shaping
	RFM_MODULATION_FSK_BT_1P0 RFMModulation = 0x08 // 01000 FSK Guassian filter, BT=1.0
	RFM_MODULATION_FSK_BT_0P5 RFMModulation = 0x10 // 10000 FSK Gaussian filter, BT=0.5
	RFM_MODULATION_FSK_BT_0P3 RFMModulation = 0x18 // 11000 FSK Gaussian filter, BT=0.3
	RFM_MODULATION_OOK        RFMModulation = 0x01 // 00001 OOK no shaping
	RFM_MODULATION_OOK_BR     RFMModulation = 0x09 // 01001 OOK Filtering with f(cutoff) = BR
	RFM_MODULATION_OOK_2BR    RFMModulation = 0x0A // 01010 OOK Filtering with f(cutoff) = 2BR
	RFM_MODULATION_MAX        RFMModulation = 0x1F
)

////////////////////////////////////////////////////////////////////////////////
// RFM69 STRINGIFY

func (m RFMMode) String() string {
	switch m {
	case RFM_MODE_SLEEP:
		return "RFM_MODE_SLEEP"
	case RFM_MODE_STDBY:
		return "RFM_MODE_STDBY"
	case RFM_MODE_FS:
		return "RFM_MODE_FS"
	case RFM_MODE_TX:
		return "RFM_MODE_TX"
	case RFM_MODE_RX:
		return "RFM_MODE_RX"
	default:
		return "[?? Invalid RFMMode value]"
	}
}

func (m RFMDataMode) String() string {
	switch m {
	case RFM_DATAMODE_PACKET:
		return "RFM_DATAMODE_PACKET"
	case RFM_DATAMODE_CONTINUOUS_NOSYNC:
		return "RFM_DATAMODE_CONTINUOUS_NOSYNC"
	case RFM_DATAMODE_CONTINUOUS_SYNC:
		return "RFM_DATAMODE_CONTINUOUS_SYNC"
	default:
		return "[?? Invalid RFMDataMode value]"
	}
}

func (m RFMModulation) String() string {
	switch m {
	case RFM_MODULATION_FSK:
		return "RFM_MODULATION_FSK"
	case RFM_MODULATION_FSK_BT_1P0:
		return "RFM_MODULATION_FSK_BT_1P0"
	case RFM_MODULATION_FSK_BT_0P5:
		return "RFM_MODULATION_FSK_BT_0P5"
	case RFM_MODULATION_FSK_BT_0P3:
		return "RFM_MODULATION_FSK_BT_0P3"
	case RFM_MODULATION_OOK:
		return "RFM_MODULATION_OOK"
	case RFM_MODULATION_OOK_BR:
		return "RFM_MODULATION_OOK_BR"
	case RFM_MODULATION_OOK_2BR:
		return "RFM_MODULATION_OOK_2BR"
	default:
		return "[?? Invalid RFMModulation value]"
	}
}
