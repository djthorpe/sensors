/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package sensors

import (
	"context"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// RFM69 TYPES

type (
	RFMMode         uint8
	RFMDataMode     uint8
	RFMModulation   uint8
	RFMPacketFormat uint8
	RFMPacketCoding uint8
	RFMPacketFilter uint8
	RFMPacketCRC    uint8
	RFMAFCMode      uint8
	RFMAFCRoutine   uint8
	RFMTXStart      uint8
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

	// Bitrate & Frequency
	Bitrate() uint
	FreqCarrier() uint
	FreqDeviation() uint
	SetBitrate(bits_per_second uint) error
	SetFreqCarrier(hertz uint) error
	SetFreqDeviation(hertz uint) error

	// Listen Mode and Sequencer
	SetSequencer(enabled bool) error
	SequencerEnabled() bool
	SetListenOn(value bool) error
	ListenOn() bool

	// Packets
	PacketFormat() RFMPacketFormat
	PacketCoding() RFMPacketCoding
	PacketFilter() RFMPacketFilter
	PacketCRC() RFMPacketCRC
	SetPacketFormat(packet_format RFMPacketFormat) error
	SetPacketCoding(packet_coding RFMPacketCoding) error
	SetPacketFilter(packet_filter RFMPacketFilter) error
	SetPacketCRC(packet_crc RFMPacketCRC) error

	// Addresses
	NodeAddress() uint8
	BroadcastAddress() uint8
	SetNodeAddress(value uint8) error
	SetBroadcastAddress(value uint8) error

	// Payload & Preamble
	PreambleSize() uint16
	PayloadSize() uint8
	SetPreambleSize(preamble_size uint16) error
	SetPayloadSize(payload_size uint8) error

	// Encryption Key & Sync Words for Packet mode
	AESKey() []byte
	SetAESKey(key []byte) error
	SyncWord() []byte
	SetSyncWord(word []byte) error
	SyncTolerance() uint8
	SetSyncTolerance(bits uint8) error

	// AFC
	AFC() uint
	AFCMode() RFMAFCMode
	AFCRoutine() RFMAFCRoutine
	SetAFCRoutine(afc_routine RFMAFCRoutine) error
	SetAFCMode(afc_mode RFMAFCMode) error
	TriggerAFC() error

	// FIFO
	FIFOThreshold() uint8
	SetFIFOThreshold(fifo_threshold uint8) error
	ReadFIFO(ctx context.Context) ([]byte, error)
	ReadPayload(ctx context.Context) ([]byte, bool, error)
	WriteFIFO(data []byte) error
	ClearFIFO() error

	/*
		AFCRoutine() RFMAFCRoutine
		AFCHertz() float64
		SetAFCRoutine(afc_routine RFMAFCRoutine) error

		// LNA
		LNAImpedance() RFMLNAImpedance
		LNAGain() RFMLNAGain
		LNACurrentGain() (RFMLNAGain, error)
		SetLNA(impedance RFMLNAImpedance, gain RFMLNAGain) error

		// OOK Parameters
		SetOOK(ook_threshold_type RFMOOKThresholdType, ook_threshold_step RFMOOKThresholdStep, ook_threshold_dec RFMOOKThresholdDecrement) error

		// FIFO
		FIFOFillCondition() bool
		FIFOThreshold() uint8
		SetFIFOFillCondition(fifo_fill_condition bool) error
		SetFIFOThreshold(fifo_threshold uint8) error

		// Other
		SetTXStart(tx_start RFMTXStart) error
		TXStart() RFMTXStart
		SetRXBW(mantissa RFMRXBWMantissa, exponent RFMRXBWExponent, cutoff RFMRXBWCutoff) error

		// Methods
		ReadFEIHertz() (float64, error)
		TriggerAFC() error
		ClearAFC() error
		CalibrateRCOsc() error
		MeasureTemperature(calibration float32) (float32, error)
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
	// RFM69 Packet format
	RFM_PACKET_FORMAT_FIXED    RFMPacketFormat = 0x00 // Fixed Packet Format
	RFM_PACKET_FORMAT_VARIABLE RFMPacketFormat = 0x01 // Variable Packet Format
)

const (
	// RFM69 Packet coding
	RFM_PACKET_CODING_NONE       RFMPacketCoding = 0x00 // No Packet Coding
	RFM_PACKET_CODING_MANCHESTER RFMPacketCoding = 0x01 // Manchester
	RFM_PACKET_CODING_WHITENING  RFMPacketCoding = 0x02 // Whitening
	RFM_PACKET_CODING_MAX        RFMPacketCoding = 0x03 // Mask
)

const (
	// RFM69 Packet filtering
	RFM_PACKET_FILTER_NONE      RFMPacketFilter = 0x00 // Promiscious mode
	RFM_PACKET_FILTER_NODE      RFMPacketFilter = 0x01 // Matches Node
	RFM_PACKET_FILTER_BROADCAST RFMPacketFilter = 0x02 // Matches Node or Broadcast
	RFM_PACKET_FILTER_MAX       RFMPacketFilter = 0x03 // Mask
)

const (
	// RFM69 Packet CRC
	RFM_PACKET_CRC_OFF           RFMPacketCRC = 0x00 // CRC off
	RFM_PACKET_CRC_AUTOCLEAR_OFF RFMPacketCRC = 0x01 // CRC on
	RFM_PACKET_CRC_AUTOCLEAR_ON  RFMPacketCRC = 0x02 // CRC on
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

const (
	// Automatic Frequency Correction Mode
	RFM_AFCMODE_OFF       RFMAFCMode = 0x00 // AFC is performed only when triggered
	RFM_AFCMODE_ON        RFMAFCMode = 0x01 // AFC register is not cleared before a new AFC phase
	RFM_AFCMODE_AUTOCLEAR RFMAFCMode = 0x03 // AFC register is cleared before a new AFC phase
	RFM_AFCMODE_MASK      RFMAFCMode = 0x03
)

const (
	// Automatic Frequency Correction Routine
	RFM_AFCROUTINE_STANDARD RFMAFCRoutine = 0x00 // Standard AFC Routine
	RFM_AFCROUTINE_IMPROVED RFMAFCRoutine = 0x01 // Improved AFC Routine
	RFM_AFCROUTINE_MASK     RFMAFCRoutine = 0x01
)

const (
	// RFM69 TX Start Condition
	RFM_TXSTART_FIFOLEVEL    RFMTXStart = 0x00 // When FIFO threshold is exceeded
	RFM_TXSTART_FIFONOTEMPTY RFMTXStart = 0x01 // When FIFO is not empty
	RFM_TXSTART_MAX          RFMTXStart = 0x01 // Mask
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

func (f RFMPacketFormat) String() string {
	switch f {
	case RFM_PACKET_FORMAT_FIXED:
		return "RFM_PACKET_FORMAT_FIXED"
	case RFM_PACKET_FORMAT_VARIABLE:
		return "RFM_PACKET_FORMAT_VARIABLE"
	default:
		return "[?? Invalid RFMPacketFormat value]"
	}
}

func (c RFMPacketCoding) String() string {
	switch c {
	case RFM_PACKET_CODING_NONE:
		return "RFM_PACKET_CODING_NONE"
	case RFM_PACKET_CODING_MANCHESTER:
		return "RFM_PACKET_CODING_MANCHESTER"
	case RFM_PACKET_CODING_WHITENING:
		return "RFM_PACKET_CODING_WHITENING"
	default:
		return "[?? Invalid RFMPacketCoding value]"
	}
}

func (f RFMPacketFilter) String() string {
	switch f {
	case RFM_PACKET_FILTER_NONE:
		return "RFM_PACKET_FILTER_NONE"
	case RFM_PACKET_FILTER_NODE:
		return "RFM_PACKET_FILTER_NODE"
	case RFM_PACKET_FILTER_BROADCAST:
		return "RFM_PACKET_FILTER_BROADCAST"
	default:
		return "[?? Invalid RFMPacketFilter value]"
	}
}

func (c RFMPacketCRC) String() string {
	switch c {
	case RFM_PACKET_CRC_OFF:
		return "RFM_PACKET_CRC_OFF"
	case RFM_PACKET_CRC_AUTOCLEAR_OFF:
		return "RFM_PACKET_CRC_AUTOCLEAR_OFF"
	case RFM_PACKET_CRC_AUTOCLEAR_ON:
		return "RFM_PACKET_CRC_AUTOCLEAR_ON"
	default:
		return "[?? Invalid RFMPacketCRC value]"
	}
}

func (m RFMAFCMode) String() string {
	switch m {
	case RFM_AFCMODE_OFF:
		return "RFM_AFCMODE_OFF"
	case RFM_AFCMODE_ON:
		return "RFM_AFCMODE_ON"
	case RFM_AFCMODE_AUTOCLEAR:
		return "RFM_AFCMODE_AUTOCLEAR"
	default:
		return "[?? Invalid RFMAFCMode value]"
	}
}

func (r RFMAFCRoutine) String() string {
	switch r {
	case RFM_AFCROUTINE_STANDARD:
		return "RFM_AFCROUTINE_STANDARD"
	case RFM_AFCROUTINE_IMPROVED:
		return "RFM_AFCROUTINE_IMPROVED"
	default:
		return "[?? Invalid RFMAFCRoutine value]"
	}
}

func (v RFMTXStart) String() string {
	switch v {
	case RFM_TXSTART_FIFOLEVEL:
		return "RFM_TXSTART_FIFOLEVEL"
	case RFM_TXSTART_FIFONOTEMPTY:
		return "RFM_TXSTART_FIFONOTEMPTY"
	default:
		return "[?? Invalid RFMTXStart value]"
	}
}
