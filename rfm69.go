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
	"time"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// RFM69 TYPES

type (
	RFMMode          uint8
	RFMDataMode      uint8
	RFMModulation    uint8
	RFMPacketFormat  uint8
	RFMPacketCoding  uint8
	RFMPacketFilter  uint8
	RFMPacketCRC     uint8
	RFMAFCMode       uint8
	RFMAFCRoutine    uint8
	RFMTXStart       uint8
	RFMLNAImpedance  uint8
	RFMLNAGain       uint8
	RFMRXBWFrequency uint8
	RFMRXBWCutoff    uint8
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

	// Low Noise Amplifier Settings
	LNAImpedance() RFMLNAImpedance
	LNAGain() RFMLNAGain
	LNACurrentGain() (RFMLNAGain, error)
	SetLNA(impedance RFMLNAImpedance, gain RFMLNAGain) error

	// Channel Filter Settings
	RXFilterFrequency() RFMRXBWFrequency
	RXFilterCutoff() RFMRXBWCutoff
	SetRXFilter(RFMRXBWFrequency, RFMRXBWCutoff) error

	// FIFO
	FIFOThreshold() uint8
	SetFIFOThreshold(fifo_threshold uint8) error
	ReadFIFO(ctx context.Context) ([]byte, error)
	WriteFIFO(data []byte) error
	ClearFIFO() error

	// ReadPayload
	ReadPayload(ctx context.Context) ([]byte, bool, error)

	// WritePayload writes a packet a number of times, with a delay between each
	// when the repeat is greater than zero
	WritePayload(data []byte, repeat uint, delay time.Duration) error

	// MeasureTemperature and return after calibration
	MeasureTemperature(calibration float32) (float32, error)

	/*
		// OOK Parameters
		SetOOK(ook_threshold_type RFMOOKThresholdType, ook_threshold_step RFMOOKThresholdStep, ook_threshold_dec RFMOOKThresholdDecrement) error

		// FIFO
		FIFOFillCondition() bool
		SetFIFOFillCondition(fifo_fill_condition bool) error

		// Other
		SetTXStart(tx_start RFMTXStart) error
		TXStart() RFMTXStart

		// Methods
		ReadFEIHertz() (float64, error)
		ClearAFC() error
		CalibrateRCOsc() error
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
	RFM_MODULATION_FSK_BT_1P0 RFMModulation = 0x01 // 01000 FSK Guassian filter, BT=1.0
	RFM_MODULATION_FSK_BT_0P5 RFMModulation = 0x02 // 10000 FSK Gaussian filter, BT=0.5
	RFM_MODULATION_FSK_BT_0P3 RFMModulation = 0x03 // 11000 FSK Gaussian filter, BT=0.3
	RFM_MODULATION_OOK        RFMModulation = 0x08 // 00001 OOK no shaping
	RFM_MODULATION_OOK_BR     RFMModulation = 0x09 // 01001 OOK Filtering with f(cutoff) = BR
	RFM_MODULATION_OOK_2BR    RFMModulation = 0x0A // 01010 OOK Filtering with f(cutoff) = 2BR
	RFM_MODULATION_MAX        RFMModulation = 0x0A
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

const (
	// Low Noise Amplifier Impedance
	RFM_LNA_IMPEDANCE_50  RFMLNAImpedance = 0x00 // 50 Ohms
	RFM_LNA_IMPEDANCE_100 RFMLNAImpedance = 0x01 // 100 Ohms
	RFM_LNA_IMPEDANCE_MAX RFMLNAImpedance = 0x01
)

const (
	// Low Noise Amplifier Gain
	RFM_LNA_GAIN_AUTO RFMLNAGain = 0x00 // Gain set by internal AGC Loop
	RFM_LNA_GAIN_G1   RFMLNAGain = 0x01 // Highest gain
	RFM_LNA_GAIN_G2   RFMLNAGain = 0x02 // Highest gain minus 6dB
	RFM_LNA_GAIN_G3   RFMLNAGain = 0x03 // Highest gain minus 12dB
	RFM_LNA_GAIN_G4   RFMLNAGain = 0x04 // Highest gain minus 24dB
	RFM_LNA_GAIN_G5   RFMLNAGain = 0x05 // Highest gain minus 36dB
	RFM_LNA_GAIN_G6   RFMLNAGain = 0x06 // Highest gain minus 48dB
	RFM_LNA_GAIN_MAX  RFMLNAGain = 0x07 // Gain mask
)

const (
	RFM_RXBW_CUTOFF_16    RFMRXBWCutoff = 0x00
	RFM_RXBW_CUTOFF_8     RFMRXBWCutoff = 0x01
	RFM_RXBW_CUTOFF_4     RFMRXBWCutoff = 0x02
	RFM_RXBW_CUTOFF_2     RFMRXBWCutoff = 0x03
	RFM_RXBW_CUTOFF_1     RFMRXBWCutoff = 0x04
	RFM_RXBW_CUTOFF_0P5   RFMRXBWCutoff = 0x05
	RFM_RXBW_CUTOFF_0P25  RFMRXBWCutoff = 0x06
	RFM_RXBW_CUTOFF_0P125 RFMRXBWCutoff = 0x07
	RFM_RXBW_CUTOFF_MAX   RFMRXBWCutoff = RFM_RXBW_CUTOFF_0P125
)

const (
	RFM_RXBW_FREQUENCY_FSK_2P6   RFMRXBWFrequency = 2<<3 | 7
	RFM_RXBW_FREQUENCY_FSK_3P1   RFMRXBWFrequency = 1<<3 | 7
	RFM_RXBW_FREQUENCY_FSK_3P9   RFMRXBWFrequency = 0<<3 | 7
	RFM_RXBW_FREQUENCY_FSK_5P2   RFMRXBWFrequency = 2<<3 | 6
	RFM_RXBW_FREQUENCY_FSK_6P3   RFMRXBWFrequency = 1<<3 | 6
	RFM_RXBW_FREQUENCY_FSK_7P8   RFMRXBWFrequency = 0<<3 | 6
	RFM_RXBW_FREQUENCY_FSK_10P4  RFMRXBWFrequency = 2<<3 | 5
	RFM_RXBW_FREQUENCY_FSK_12P5  RFMRXBWFrequency = 1<<3 | 5
	RFM_RXBW_FREQUENCY_FSK_15P6  RFMRXBWFrequency = 0<<3 | 5
	RFM_RXBW_FREQUENCY_FSK_20P8  RFMRXBWFrequency = 2<<3 | 4
	RFM_RXBW_FREQUENCY_FSK_25P0  RFMRXBWFrequency = 1<<3 | 4
	RFM_RXBW_FREQUENCY_FSK_31P3  RFMRXBWFrequency = 0<<3 | 4
	RFM_RXBW_FREQUENCY_FSK_41P7  RFMRXBWFrequency = 2<<3 | 3
	RFM_RXBW_FREQUENCY_FSK_50P0  RFMRXBWFrequency = 1<<3 | 3
	RFM_RXBW_FREQUENCY_FSK_62P5  RFMRXBWFrequency = 0<<3 | 3
	RFM_RXBW_FREQUENCY_FSK_83P3  RFMRXBWFrequency = 2<<3 | 2
	RFM_RXBW_FREQUENCY_FSK_100P0 RFMRXBWFrequency = 1<<3 | 2
	RFM_RXBW_FREQUENCY_FSK_125P0 RFMRXBWFrequency = 0<<3 | 2
	RFM_RXBW_FREQUENCY_FSK_166P7 RFMRXBWFrequency = 2<<3 | 1
	RFM_RXBW_FREQUENCY_FSK_200P0 RFMRXBWFrequency = 1<<3 | 1
	RFM_RXBW_FREQUENCY_FSK_250P0 RFMRXBWFrequency = 0<<3 | 1
	RFM_RXBW_FREQUENCY_FSK_333P3 RFMRXBWFrequency = 2<<3 | 0
	RFM_RXBW_FREQUENCY_FSK_400P0 RFMRXBWFrequency = 1<<3 | 0
	RFM_RXBW_FREQUENCY_FSK_500P0 RFMRXBWFrequency = 0<<3 | 0
	RFM_RXBW_FREQUENCY_MAX       RFMRXBWFrequency = 0x1F
)

const (
	RFM_RXBW_FREQUENCY_OOK_1P3   = RFM_RXBW_FREQUENCY_FSK_2P6
	RFM_RXBW_FREQUENCY_OOK_1P6   = RFM_RXBW_FREQUENCY_FSK_3P1
	RFM_RXBW_FREQUENCY_OOK_2P0   = RFM_RXBW_FREQUENCY_FSK_3P9
	RFM_RXBW_FREQUENCY_OOK_2P6   = RFM_RXBW_FREQUENCY_FSK_5P2
	RFM_RXBW_FREQUENCY_OOK_3P1   = RFM_RXBW_FREQUENCY_FSK_6P3
	RFM_RXBW_FREQUENCY_OOK_3P9   = RFM_RXBW_FREQUENCY_FSK_7P8
	RFM_RXBW_FREQUENCY_OOK_5P2   = RFM_RXBW_FREQUENCY_FSK_10P4
	RFM_RXBW_FREQUENCY_OOK_6P3   = RFM_RXBW_FREQUENCY_FSK_12P5
	RFM_RXBW_FREQUENCY_OOK_7P8   = RFM_RXBW_FREQUENCY_FSK_15P6
	RFM_RXBW_FREQUENCY_OOK_10P4  = RFM_RXBW_FREQUENCY_FSK_20P8
	RFM_RXBW_FREQUENCY_OOK_12P5  = RFM_RXBW_FREQUENCY_FSK_25P0
	RFM_RXBW_FREQUENCY_OOK_15P6  = RFM_RXBW_FREQUENCY_FSK_31P3
	RFM_RXBW_FREQUENCY_OOK_20P8  = RFM_RXBW_FREQUENCY_FSK_41P7
	RFM_RXBW_FREQUENCY_OOK_25P0  = RFM_RXBW_FREQUENCY_FSK_50P0
	RFM_RXBW_FREQUENCY_OOK_31P3  = RFM_RXBW_FREQUENCY_FSK_62P5
	RFM_RXBW_FREQUENCY_OOK_41P7  = RFM_RXBW_FREQUENCY_FSK_83P3
	RFM_RXBW_FREQUENCY_OOK_50P0  = RFM_RXBW_FREQUENCY_FSK_100P0
	RFM_RXBW_FREQUENCY_OOK_62P5  = RFM_RXBW_FREQUENCY_FSK_125P0
	RFM_RXBW_FREQUENCY_OOK_83P3  = RFM_RXBW_FREQUENCY_FSK_166P7
	RFM_RXBW_FREQUENCY_OOK_100P0 = RFM_RXBW_FREQUENCY_FSK_200P0
	RFM_RXBW_FREQUENCY_OOK_125P0 = RFM_RXBW_FREQUENCY_FSK_250P0
	RFM_RXBW_FREQUENCY_OOK_166P7 = RFM_RXBW_FREQUENCY_FSK_333P3
	RFM_RXBW_FREQUENCY_OOK_200P0 = RFM_RXBW_FREQUENCY_FSK_400P0
	RFM_RXBW_FREQUENCY_OOK_250P0 = RFM_RXBW_FREQUENCY_FSK_500P0
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

func (v RFMLNAImpedance) String() string {
	switch v {
	case RFM_LNA_IMPEDANCE_50:
		return "RFM_LNA_IMPEDANCE_50"
	case RFM_LNA_IMPEDANCE_100:
		return "RFM_LNA_IMPEDANCE_100"
	default:
		return "[?? Invalid RFMLNAImpedance value]"
	}
}

func (v RFMLNAGain) String() string {
	switch v {
	case RFM_LNA_GAIN_AUTO:
		return "RFM_LNA_GAIN_AUTO"
	case RFM_LNA_GAIN_G1:
		return "RFM_LNA_GAIN_G1"
	case RFM_LNA_GAIN_G2:
		return "RFM_LNA_GAIN_G2"
	case RFM_LNA_GAIN_G3:
		return "RFM_LNA_GAIN_G3"
	case RFM_LNA_GAIN_G4:
		return "RFM_LNA_GAIN_G4"
	case RFM_LNA_GAIN_G5:
		return "RFM_LNA_GAIN_G5"
	case RFM_LNA_GAIN_G6:
		return "RFM_LNA_GAIN_G6"
	default:
		return "[?? Invalid RFMLNAGain value]"
	}
}

func (v RFMRXBWCutoff) String() string {
	switch v {
	case RFM_RXBW_CUTOFF_16:
		return "RFM_RXBW_CUTOFF_16"
	case RFM_RXBW_CUTOFF_8:
		return "RFM_RXBW_CUTOFF_8"
	case RFM_RXBW_CUTOFF_4:
		return "RFM_RXBW_CUTOFF_4"
	case RFM_RXBW_CUTOFF_2:
		return "RFM_RXBW_CUTOFF_2"
	case RFM_RXBW_CUTOFF_1:
		return "RFM_RXBW_CUTOFF_1"
	case RFM_RXBW_CUTOFF_0P5:
		return "RFM_RXBW_CUTOFF_0P5"
	case RFM_RXBW_CUTOFF_0P25:
		return "RFM_RXBW_CUTOFF_0P25"
	case RFM_RXBW_CUTOFF_0P125:
		return "RFM_RXBW_CUTOFF_0P125"
	default:
		return "[?? Invalid RFMRXBWCutoff value]"
	}
}

func (v RFMRXBWFrequency) String() string {
	switch v {
	case RFM_RXBW_FREQUENCY_FSK_2P6:
		return "RFM_RXBW_FREQUENCY_FSK_2P6,RFM_RXBW_FREQUENCY_OOK_1P3"
	case RFM_RXBW_FREQUENCY_FSK_3P1:
		return "RFM_RXBW_FREQUENCY_FSK_3P1,RFM_RXBW_FREQUENCY_OOK_1P6"
	case RFM_RXBW_FREQUENCY_FSK_3P9:
		return "RFM_RXBW_FREQUENCY_FSK_3P9,RFM_RXBW_FREQUENCY_OOK_2P0"
	case RFM_RXBW_FREQUENCY_FSK_5P2:
		return "RFM_RXBW_FREQUENCY_FSK_5P2,RFM_RXBW_FREQUENCY_OOK_2P6"
	case RFM_RXBW_FREQUENCY_FSK_6P3:
		return "RFM_RXBW_FREQUENCY_FSK_6P3,RFM_RXBW_FREQUENCY_OOK_3P1"
	case RFM_RXBW_FREQUENCY_FSK_7P8:
		return "RFM_RXBW_FREQUENCY_FSK_7P8,RFM_RXBW_FREQUENCY_OOK_3P9"
	case RFM_RXBW_FREQUENCY_FSK_10P4:
		return "RFM_RXBW_FREQUENCY_FSK_10P4,RFM_RXBW_FREQUENCY_OOK_5P2"
	case RFM_RXBW_FREQUENCY_FSK_12P5:
		return "RFM_RXBW_FREQUENCY_FSK_12P5,RFM_RXBW_FREQUENCY_OOK_6P3"
	case RFM_RXBW_FREQUENCY_FSK_15P6:
		return "RFM_RXBW_FREQUENCY_FSK_15P6,RFM_RXBW_FREQUENCY_OOK_7P8"
	case RFM_RXBW_FREQUENCY_FSK_20P8:
		return "RFM_RXBW_FREQUENCY_FSK_20P8,RFM_RXBW_FREQUENCY_OOK_10P4"
	case RFM_RXBW_FREQUENCY_FSK_25P0:
		return "RFM_RXBW_FREQUENCY_FSK_25P0,RFM_RXBW_FREQUENCY_OOK_12P5"
	case RFM_RXBW_FREQUENCY_FSK_31P3:
		return "RFM_RXBW_FREQUENCY_FSK_31P3,RFM_RXBW_FREQUENCY_OOK_15P6"
	case RFM_RXBW_FREQUENCY_FSK_41P7:
		return "RFM_RXBW_FREQUENCY_FSK_41P7,RFM_RXBW_FREQUENCY_OOK_20P8"
	case RFM_RXBW_FREQUENCY_FSK_50P0:
		return "RFM_RXBW_FREQUENCY_FSK_50P0,RFM_RXBW_FREQUENCY_OOK_25P0"
	case RFM_RXBW_FREQUENCY_FSK_62P5:
		return "RFM_RXBW_FREQUENCY_FSK_62P5,RFM_RXBW_FREQUENCY_OOK_31P3"
	case RFM_RXBW_FREQUENCY_FSK_83P3:
		return "RFM_RXBW_FREQUENCY_FSK_83P3,RFM_RXBW_FREQUENCY_OOK_41P7"
	case RFM_RXBW_FREQUENCY_FSK_100P0:
		return "RFM_RXBW_FREQUENCY_FSK_100P0,RFM_RXBW_FREQUENCY_OOK_50P0"
	case RFM_RXBW_FREQUENCY_FSK_125P0:
		return "RFM_RXBW_FREQUENCY_FSK_125P0,RFM_RXBW_FREQUENCY_OOK_62P5"
	case RFM_RXBW_FREQUENCY_FSK_166P7:
		return "RFM_RXBW_FREQUENCY_FSK_166P7,RFM_RXBW_FREQUENCY_OOK_83P3"
	case RFM_RXBW_FREQUENCY_FSK_200P0:
		return "RFM_RXBW_FREQUENCY_FSK_200P0,RFM_RXBW_FREQUENCY_OOK_100P0"
	case RFM_RXBW_FREQUENCY_FSK_250P0:
		return "RFM_RXBW_FREQUENCY_FSK_250P0,RFM_RXBW_FREQUENCY_OOK_125P0"
	case RFM_RXBW_FREQUENCY_FSK_333P3:
		return "RFM_RXBW_FREQUENCY_FSK_333P3,RFM_RXBW_FREQUENCY_OOK_166P7"
	case RFM_RXBW_FREQUENCY_FSK_400P0:
		return "RFM_RXBW_FREQUENCY_FSK_400P0,RFM_RXBW_FREQUENCY_OOK_200P0"
	case RFM_RXBW_FREQUENCY_FSK_500P0:
		return "RFM_RXBW_FREQUENCY_FSK_500P0,RFM_RXBW_FREQUENCY_OOK_250P0"
	default:
		return "[?? Invalid RFMRXBFrequency value]"

	}
}
