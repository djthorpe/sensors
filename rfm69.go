/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package sensors

import (
	"time"

	"github.com/djthorpe/gopi"
)

type RFM69 interface {
	gopi.Driver

	// Mode, Data Mode and Modulation
	Mode() RFMDeviceMode
	DataMode() RFMDataMode
	SetMode(device_mode RFMDeviceMode) error
	SetDataMode(data_mode RFMDataMode) error
	Modulation() RFMModulation
	SetModulation(modulation RFMModulation) error

	// Addresses
	NodeAddress() uint8
	BroadcastAddress() uint8
	SetNodeAddress(value uint8) error
	SetBroadcastAddress(value uint8) error

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
}
