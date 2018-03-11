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
	"math"
	"strings"
	"sync"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// driver
type rfm69 struct {
	spi  gopi.SPI
	log  gopi.Logger
	lock sync.Mutex

	version               uint8
	mode                  sensors.RFMMode
	sequencer_off         bool
	listen_on             bool
	data_mode             sensors.RFMDataMode
	modulation            sensors.RFMModulation
	bitrate               uint16
	frf                   uint32
	fdev                  uint16
	aes_key               []byte
	aes_on                bool
	sync_word             []byte
	sync_on               bool
	sync_size             uint8
	sync_tol              uint8
	rx_inter_packet_delay uint8
	rx_auto_restart       bool
	tx_start              sensors.RFMTXStart
	fifo_threshold        uint8
	fifo_fill_condition   bool
	node_address          uint8
	broadcast_address     uint8
	preamble_size         uint16
	payload_size          uint8
	packet_format         sensors.RFMPacketFormat
	packet_coding         sensors.RFMPacketCoding
	packet_filter         sensors.RFMPacketFilter
	crc_enabled           bool
	crc_auto_clear_off    bool
	afc                   int16
	afc_mode              sensors.RFMAFCMode
	afc_routine           sensors.RFMAFCRoutine
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	RFM_SPI_MODE       = gopi.SPI_MODE_0
	RFM_SPI_SPEEDHZ    = 4000000 // 4MHz
	RFM_VERSION_VALUE  = 0x24
	RFM_AESKEY_BYTES   = 16
	RFM_SYNCWORD_BYTES = 8
	RFM_FXOSC_MHZ      = 32         // Crystal oscillator frequency MHz
	RFM_FSTEP_HZ       = 61         // Frequency synthesizer step
	RFM_BITRATE_MIN    = 500        // bits per second
	RFM_BITRATE_MAX    = 300 * 1024 // bits per second
	RFM_FDEV_MAX       = 0x3FFF     // Maximum value of FDEV
	RFM_FRF_MAX        = 0xFFFFFF   // Maximum value of FRF
	RFM_FIFO_SIZE      = 66         // Bytes
)

////////////////////////////////////////////////////////////////////////////////
// MODE, DATA MODE AND MODULATION

// Return device mode
func (this *rfm69) Mode() sensors.RFMMode {
	return this.mode
}

// Return data mode
func (this *rfm69) DataMode() sensors.RFMDataMode {
	return this.data_mode
}

// Return modulation
func (this *rfm69) Modulation() sensors.RFMModulation {
	return this.modulation
}

// Set device mode
func (this *rfm69) SetMode(mode sensors.RFMMode) error {
	this.log.Debug("<sensors.RFM69.SetMode>{ mode=%v }", mode)

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Switch listen_on off if not going to sleep
	if mode != sensors.RFM_MODE_SLEEP && this.listen_on {
		if err := this.setOpMode(this.mode, false, true, this.sequencer_off); err != nil {
			return err
		} else {
			this.listen_on = false
		}
	}

	// Switch off sequencer if going into standby
	if mode == sensors.RFM_MODE_STDBY && this.sequencer_off == false {
		if err := this.setOpMode(this.mode, this.listen_on, false, true); err != nil {
			return err
		} else {
			this.sequencer_off = true
		}
	}

	// Write mode and read back again
	if err := this.setOpMode(mode, false, false, this.sequencer_off); err != nil {
		return err
	}

	// Wait for device ready bit
	if err := wait_for_condition(func() (bool, error) {
		value, err := this.getIRQFlags1(RFM_IRQFLAGS1_MODEREADY)
		return to_uint8_bool(value), err
	}, true, time.Millisecond*1000); err != nil {
		return err
	}

	// Read back register
	if mode_read, listen_on_read, sequencer_off_read, err := this.getOpMode(); err != nil {
		return err
	} else if mode_read != mode {
		this.log.Debug2("SetMode expecting mode=%v, got=%v", mode, mode_read)
		return sensors.ErrUnexpectedResponse
	} else if listen_on_read != this.listen_on {
		this.log.Debug2("SetMode expecting listen_on=%v, got=%v", this.listen_on, listen_on_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.mode = mode_read
		this.listen_on = listen_on_read
		this.sequencer_off = sequencer_off_read
	}

	// If RX mode then read AFC value
	if this.mode == sensors.RFM_MODE_RX {
		if afc, err := this.getAFC(); err != nil {
			return err
		} else {
			this.afc = afc
		}
	}

	return nil
}

// Set data mode
func (this *rfm69) SetDataMode(data_mode sensors.RFMDataMode) error {
	this.log.Debug("<sensors.RFM69.SetDataMode>{ data_mode=%v }", data_mode)

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Write
	if err := this.setDataModul(data_mode, this.modulation); err != nil {
		return err
	}

	// Read
	if data_mode_read, modulation_read, err := this.getDataModul(); err != nil {
		return err
	} else if data_mode != data_mode_read {
		this.log.Debug2("SetDataMode expecting date_mode=%v, got=%v", data_mode, data_mode_read)
		return sensors.ErrUnexpectedResponse
	} else if modulation_read != this.modulation {
		this.log.Debug2("SetDataMode expecting modulation=%v, got=%v", this.modulation, modulation_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.data_mode = data_mode_read
	}

	return nil
}

// Set modulation
func (this *rfm69) SetModulation(modulation sensors.RFMModulation) error {
	this.log.Debug("<sensors.RFM69.SetModulation{ modulation=%v }", modulation)

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Write
	if err := this.setDataModul(this.data_mode, modulation); err != nil {
		return err
	}

	// Read
	if data_mode_read, modulation_read, err := this.getDataModul(); err != nil {
		return err
	} else if modulation_read != modulation {
		this.log.Debug2("SetModulation expecting modulation=%v, got=%v", modulation, modulation_read)
		return sensors.ErrUnexpectedResponse
	} else if data_mode_read != this.data_mode {
		this.log.Debug2("SetModulation expecting data_mode=%v, got=%v", this.data_mode, data_mode_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.modulation = modulation
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// BITRATE & FREQUENCY

// Return bitrate in bits per second
func (this *rfm69) Bitrate() uint {
	return uint(float64(RFM_FXOSC_MHZ*1E6) / float64(this.bitrate))
}

// Return frequency carrier in Hz
func (this *rfm69) GetFreqCarrier() uint {
	return uint(RFM_FSTEP_HZ) * uint(this.frf)
}

func (this *rfm69) SetBitrate(bits_per_second uint) error {
	this.log.Debug("<sensors.RFM69.SetBitrate>{ bits_per_second=%v }", bits_per_second)

	if bits_per_second < RFM_BITRATE_MIN || bits_per_second > RFM_BITRATE_MAX {
		return gopi.ErrBadParameter
	}

	msb_lsb := uint16(float64(RFM_FXOSC_MHZ*1E6) / float64(bits_per_second))
	return this.SetBitrateUint16(msb_lsb)
}

// Set bitrate as register value
func (this *rfm69) SetBitrateUint16(value uint16) error {
	this.log.Debug("<sensors.RFM69.SetBitrateUint16>{ value=0x%04X }", value)

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Write
	if err := this.setBitrate(value); err != nil {
		return err
	}

	// Read
	if value_read, err := this.getBitrate(); err != nil {
		return err
	} else if value_read != value {
		this.log.Debug2("SetBitrateUint16 expecting value=0x%04X, got=0x%04X", value, value_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.bitrate = value_read
	}

	// Success
	return nil
}

// Return frequency carrier in Hz
func (this *rfm69) FreqCarrier() uint {
	return uint(RFM_FSTEP_HZ) * uint(this.frf)
}

// Return frequency deviation in Hz
func (this *rfm69) FreqDeviation() uint {
	return uint(RFM_FSTEP_HZ) * uint(this.fdev)
}

func (this *rfm69) SetFreqCarrier(hertz uint) error {
	this.log.Debug("<sensors.RFM69.SetFreqCarrier>{ hertz=%v }", hertz)

	msb_mid_lsb := uint32(math.Ceil(float64(hertz) / float64(RFM_FSTEP_HZ)))
	if msb_mid_lsb > RFM_FRF_MAX {
		return gopi.ErrBadParameter
	}

	return this.SetFreqCarrierUint24(msb_mid_lsb)
}

func (this *rfm69) SetFreqDeviation(hertz uint) error {
	this.log.Debug("<sensors.RFM69.SetFreqDeviation>{ hertz=%v }", hertz)

	msb_lsb := uint16(math.Ceil(float64(hertz) / float64(RFM_FSTEP_HZ)))
	if msb_lsb > RFM_FDEV_MAX {
		return gopi.ErrBadParameter
	}

	return this.SetFreqDeviationUint16(msb_lsb)
}

// Set frequency deviation register value
func (this *rfm69) SetFreqDeviationUint16(value uint16) error {
	this.log.Debug("<sensors.RFM69.SetFreqDeviationUint16>{ value=0x%04X }", value)

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Write
	if err := this.setFreqDeviation(value); err != nil {
		return err
	}

	// Read
	if value_read, err := this.getFreqDeviation(); err != nil {
		return err
	} else if value != value_read {
		this.log.Debug2("SetFreqDeviationUint16 expecting value=0x%04X, got=0x%04X", value, value_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.fdev = value
	}

	// Success
	return nil
}

// Set carrier frequency register
func (this *rfm69) SetFreqCarrierUint24(value uint32) error {
	this.log.Debug("<sensors.RFM69.SetFreqCarrierUint24>{ value=0x%06X }", value)

	if value&0xFF000000 != 0x00000000 {
		return gopi.ErrBadParameter
	}

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Write
	if err := this.setFreqCarrier(value); err != nil {
		return err
	}

	// Read
	if value_read, err := this.getFreqCarrier(); err != nil {
		return err
	} else if value != value_read {
		this.log.Debug2("SetFreqCarrierUint24 expecting value=0x%06X, got=0x%06X", value, value_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.frf = value
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// LISTEN MODE AND SEQUENCER

func (this *rfm69) SequencerEnabled() bool {
	return !this.sequencer_off
}

func (this *rfm69) ListenOn() bool {
	return this.listen_on
}

func (this *rfm69) SetListenOn(value bool) error {
	this.log.Debug("<sensors.RFM69.SetListenOn>{ value=%v }", value)

	// if listen mode is to be switched on, then unit needs to go into sleep
	if value && this.mode != sensors.RFM_MODE_SLEEP {
		if err := this.SetMode(sensors.RFM_MODE_SLEEP); err != nil {
			return err
		}
	}

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// listen_abort=true if listen_on=false
	if err := this.setOpMode(this.mode, value, !value, this.sequencer_off); err != nil {
		return err
	}

	if mode_read, listen_on_read, sequencer_off_read, err := this.getOpMode(); err != nil {
		return err
	} else if mode_read != this.mode {
		this.log.Debug2("SetListenOn expecting mode=%v, got=%v", this.mode, mode_read)
		return sensors.ErrUnexpectedResponse
	} else if listen_on_read != value {
		this.log.Debug2("SetListenOn expecting mode=%v, got=%v", this.mode, mode_read)
		return sensors.ErrUnexpectedResponse
	} else if sequencer_off_read != this.sequencer_off {
		this.log.Debug2("SetListenOn expecting mode=%v, got=%v", this.mode, mode_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.listen_on = value
	}
	return nil
}

func (this *rfm69) SetSequencer(enabled bool) error {
	this.log.Debug("<sensors.RFM69.SetSequencer>{ enabled=%v }", enabled)

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	if err := this.setOpMode(this.mode, this.listen_on, false, !enabled); err != nil {
		return err
	}

	if _, _, sequencer_off_read, err := this.getOpMode(); err != nil {
		return err
	} else if sequencer_off_read == enabled {
		this.log.Debug2("SetSequencer expecting sequencer_off=%v, got=%v", !enabled, sequencer_off_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.sequencer_off = !enabled
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// NODE AND BROADCAST ADDRESS

func (this *rfm69) NodeAddress() uint8 {
	return this.node_address
}

func (this *rfm69) BroadcastAddress() uint8 {
	return this.broadcast_address
}

func (this *rfm69) SetNodeAddress(value uint8) error {
	this.log.Debug("<sensors.RFM69.SetNodeAddress>{ value=%02X }", value)

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Write
	if err := this.setNodeAddress(value); err != nil {
		return err
	}

	// Read
	if value_read, err := this.getNodeAddress(); err != nil {
		return err
	} else if value_read != value {
		this.log.Debug2("SetNodeAddress expecting value=%02X, got=%02X", value, value_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.node_address = value
	}
	return nil
}

func (this *rfm69) SetBroadcastAddress(value uint8) error {
	this.log.Debug("<sensors.RFM69.SetBroadcastAddress>{ value=%02X }", value)

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Write
	if err := this.setBroadcastAddress(value); err != nil {
		return err
	}

	// Read
	if value_read, err := this.getBroadcastAddress(); err != nil {
		return err
	} else if value_read != value {
		this.log.Debug2("SetBroadcastAddress expecting value=%02X, got=%02X", value, value_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.broadcast_address = value
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// AES KEY

func (this *rfm69) SetAESEnabled(enabled bool) error {
	this.log.Debug("<sensors.RFM69.SetAESEnabled>{ enabled=%v }", enabled)

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	if err := this.setPacketConfig2(this.rx_inter_packet_delay, this.rx_auto_restart, enabled); err != nil {
		return err
	} else if rx_inter_packet_delay_read, rx_auto_restart_read, enabled_read, err := this.getPacketConfig2(); err != nil {
		return err
	} else if this.rx_inter_packet_delay != rx_inter_packet_delay_read {
		this.log.Debug2("SetAESEnabled expecting rx_inter_packet_delay=%v, got=%v", this.rx_inter_packet_delay, rx_inter_packet_delay_read)
		return sensors.ErrUnexpectedResponse
	} else if this.rx_auto_restart != rx_auto_restart_read {
		this.log.Debug2("SetAESEnabled expecting rx_auto_restart=%v, got=%v", this.rx_auto_restart, rx_auto_restart_read)
		return sensors.ErrUnexpectedResponse
	} else if enabled_read != enabled {
		this.log.Debug2("SetAESEnabled expecting aes_on=%v, got=%v", enabled, enabled_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.aes_on = enabled
	}

	// Success
	return nil
}

func (this *rfm69) AESKey() []byte {
	if this.aes_on {
		return this.aes_key
	} else {
		return nil
	}
}

func (this *rfm69) SetAESKeyEx(key []byte) error {
	this.log.Debug("<sensors.RFM69.SetAESKeyEx>{ key=%v }", strings.ToUpper(hex.EncodeToString(key)))

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	if len(key) != RFM_AESKEY_BYTES {
		this.log.Debug2("Error: SetAESKeyEx invalid key length")
		return gopi.ErrBadParameter
	} else if err := this.setAESKey(key); err != nil {
		return err
	} else if key_read, err := this.getAESKey(); err != nil {
		return err
	} else if matches_byte_array(key_read, key) == false {
		this.log.Debug2("SetAESKeyEx expecting value=%v, got=%v", key, key_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.aes_key = key
		return nil
	}
}

func (this *rfm69) SetAESKey(key []byte) error {
	this.log.Debug("<sensors.RFM69.SetAESKey>{ key=%v }", strings.ToUpper(hex.EncodeToString(key)))

	if key == nil {
		return this.SetAESEnabled(false)
	} else if err := this.SetAESKeyEx(key); err != nil {
		return err
	} else {
		return this.SetAESEnabled(true)
	}
}

////////////////////////////////////////////////////////////////////////////////
// SYNC WORD

func (this *rfm69) SyncWord() []byte {
	if this.sync_on == false {
		return nil
	} else if this.sync_size >= 0 && this.sync_size < RFM_SYNCWORD_BYTES {
		return this.sync_word[0:(this.sync_size + 1)]
	} else {
		return this.sync_word
	}
}

func (this *rfm69) SyncTolerance() uint8 {
	return this.sync_tol
}

func (this *rfm69) SetSyncWord(word []byte) error {
	this.log.Debug("<sensors.RFM69.SetSyncWord>{ word=%v }", strings.ToUpper(hex.EncodeToString(word)))

	// Write sync word size
	if word == nil {
		// Disable
		return this.SetSyncWordSize(0)
	} else if len(word) > RFM_SYNCWORD_BYTES {
		return gopi.ErrBadParameter
	} else if matches_byte(word, 0x00) == true {
		return gopi.ErrBadParameter
	} else if err := this.SetSyncWordSize(uint8(len(word))); err != nil {
		return err
	}

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Write
	if err := this.setSyncWord(word); err != nil {
		return err
	}

	// Read
	if sync_word_read, err := this.getSyncWord(); err != nil {
		return err
	} else if matches_byte_array(sync_word_read[0:len(word)], word) == false {
		this.log.Debug2("SetSyncWord expecting word=%v, got=%v", word, sync_word_read[0:len(word)])
		return sensors.ErrUnexpectedResponse
	} else {
		this.sync_word = sync_word_read
	}

	// success
	return nil
}

func (this *rfm69) SetSyncTolerance(bits uint8) error {
	this.log.Debug("<sensors.RFM69.SetSyncTolerance>{ bits=%v }", bits)

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Write
	if bits > 7 {
		return gopi.ErrBadParameter
	} else if err := this.setSyncConfig(this.sync_on, this.fifo_fill_condition, this.sync_size, bits); err != nil {
		return err
	}

	// Read
	if sync_on_read, fifo_fill_condition_read, sync_size_read, sync_tol_read, err := this.getSyncConfig(); err != nil {
		return err
	} else if sync_on_read != this.sync_on {
		this.log.Debug2("SetSyncTolerance expecting sync_on=%v, got=%v", this.sync_on, sync_on_read)
		return sensors.ErrUnexpectedResponse
	} else if fifo_fill_condition_read != this.fifo_fill_condition {
		this.log.Debug2("SetSyncTolerance expecting fifo_fill_condition=%v, got=%v", this.fifo_fill_condition, fifo_fill_condition_read)
		return sensors.ErrUnexpectedResponse
	} else if sync_size_read != this.sync_size {
		this.log.Debug2("SetSyncTolerance expecting sync_size=%v, got=%v", this.sync_size, sync_size_read)
		return sensors.ErrUnexpectedResponse
	} else if sync_tol_read != bits {
		this.log.Debug2("SetSyncTolerance expecting sync_tol=%v, got=%v", bits, sync_tol_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.sync_tol = bits
	}
	return nil
}

// 0 = disabled, 1-8 is number of bytes
func (this *rfm69) SetSyncWordSize(bytes uint8) error {
	this.log.Debug("<sensors.RFM69.SetSyncWordSize>{ bytes=%v }", bytes)

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Write
	if bytes == 0 {
		if err := this.setSyncConfig(false, this.fifo_fill_condition, this.sync_size, this.sync_tol); err != nil {
			return err
		}
	} else if bytes >= 1 && bytes <= RFM_SYNCWORD_BYTES {
		if err := this.setSyncConfig(true, this.fifo_fill_condition, bytes-1, this.sync_tol); err != nil {
			return err
		}
	} else {
		return gopi.ErrBadParameter
	}

	// Read
	if sync_on_read, fifo_fill_condition_read, sync_size_read, sync_tol_read, err := this.getSyncConfig(); err != nil {
		return err
	} else if bytes > 0 && sync_on_read != true {
		this.log.Debug2("SetSyncWordSize expecting sync_on=%v, got=%v", true, sync_on_read, bytes)
		return sensors.ErrUnexpectedResponse
	} else if fifo_fill_condition_read != this.fifo_fill_condition {
		this.log.Debug2("SetSyncWordSize expecting fifo_fill_condition=%v, got=%v", this.fifo_fill_condition, fifo_fill_condition_read)
		return sensors.ErrUnexpectedResponse
	} else if bytes > 0 && sync_size_read != (bytes-1) {
		this.log.Debug2("SetSyncWordSize expecting sync_size=%v, got=%v", (bytes - 1), sync_size_read)
		return sensors.ErrUnexpectedResponse
	} else if sync_tol_read != this.sync_tol {
		this.log.Debug2("SetSyncWordSize expecting sync_tol=%v, got=%v", this.sync_tol, sync_tol_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.sync_on = sync_on_read
		this.sync_size = sync_size_read
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PREAMBLE AND PAYLOAD SIZE

// Return Preamble size in bytes
func (this *rfm69) PreambleSize() uint16 {
	return this.preamble_size
}

// Return Payload size in bytes
func (this *rfm69) PayloadSize() uint8 {
	return this.payload_size
}

// Set Preamble size in bytes
func (this *rfm69) SetPreambleSize(preamble_size uint16) error {
	this.log.Debug("<sensors.RFM69.SetPreambleSize>{ preamble_size=%v }", preamble_size)

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Get Preamble Size register value
	if err := this.setPreambleSize(preamble_size); err != nil {
		return err
	} else if preamble_size_read, err := this.getPreambleSize(); err != nil {
		return err
	} else if preamble_size != preamble_size_read {
		this.log.Debug2("SetPreambleSize expecting value=%v, got=%v", preamble_size, preamble_size_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.preamble_size = preamble_size_read
	}

	return nil
}

// Set Payload size in bytes
func (this *rfm69) SetPayloadSize(payload_size uint8) error {
	this.log.Debug("<sensors.RFM69.SetPayloadSize>{ payload_size=%v }", payload_size)

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Get Preamble Size register value
	if err := this.setPayloadSize(payload_size); err != nil {
		return err
	} else if payload_size_read, err := this.getPayloadSize(); err != nil {
		return err
	} else if payload_size != payload_size_read {
		this.log.Debug2("SetPayloadSize expecting value=%v, got=%v", payload_size, payload_size_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.payload_size = payload_size
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PACKETS

func (this *rfm69) PacketFormat() sensors.RFMPacketFormat {
	return this.packet_format
}

func (this *rfm69) PacketCoding() sensors.RFMPacketCoding {
	return this.packet_coding
}

func (this *rfm69) PacketFilter() sensors.RFMPacketFilter {
	return this.packet_filter
}

func (this *rfm69) PacketCRC() sensors.RFMPacketCRC {
	if this.crc_enabled == false {
		return sensors.RFM_PACKET_CRC_OFF
	} else if this.crc_auto_clear_off {
		return sensors.RFM_PACKET_CRC_AUTOCLEAR_OFF
	} else {
		return sensors.RFM_PACKET_CRC_AUTOCLEAR_ON
	}
}

func (this *rfm69) SetPacketFormat(packet_format sensors.RFMPacketFormat) error {
	this.log.Debug("<sensors.RFM69.SetPacketFormat>{ packet_format=%v }", packet_format)

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Write
	if err := this.setPacketConfig1(packet_format, this.packet_coding, this.packet_filter, this.crc_enabled, this.crc_auto_clear_off); err != nil {
		return err
	}

	// Read
	if packet_format_read, _, _, _, _, err := this.getPacketConfig1(); err != nil {
		return err
	} else if packet_format_read != packet_format {
		this.log.Debug2("SetPacketFormat expecting packet_format=%v, got=%v", packet_format, packet_format_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.packet_format = packet_format_read
	}

	// Success
	return nil
}

func (this *rfm69) SetPacketCoding(packet_coding sensors.RFMPacketCoding) error {
	this.log.Debug("<sensors.RFM69.SetPacketCoding>{ packet_coding=%v }", packet_coding)

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Write
	if err := this.setPacketConfig1(this.packet_format, packet_coding, this.packet_filter, this.crc_enabled, this.crc_auto_clear_off); err != nil {
		return err
	}

	// Read
	if _, packet_coding_read, _, _, _, err := this.getPacketConfig1(); err != nil {
		return err
	} else if packet_coding_read != packet_coding {
		this.log.Debug2("SetPacketCoding expecting packet_coding=%v, got=%v", packet_coding, packet_coding_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.packet_coding = packet_coding_read
	}

	// Success
	return nil
}

func (this *rfm69) SetPacketFilter(packet_filter sensors.RFMPacketFilter) error {
	this.log.Debug("<sensors.RFM69.SetPacketFilter>{ packet_filter=%v }", packet_filter)

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Write
	if err := this.setPacketConfig1(this.packet_format, this.packet_coding, packet_filter, this.crc_enabled, this.crc_auto_clear_off); err != nil {
		return err
	}

	// Read
	if _, _, packet_filter_read, _, _, err := this.getPacketConfig1(); err != nil {
		return err
	} else if packet_filter_read != packet_filter {
		this.log.Debug2("SetPacketFilter expecting packet_filter=%v, got=%v", packet_filter, packet_filter_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.packet_filter = packet_filter_read
	}

	// Success
	return nil
}

func (this *rfm69) SetPacketCRC(packet_crc sensors.RFMPacketCRC) error {
	this.log.Debug("<sensors.RFM69.SetPacketCRC>{ packet_crc=%v }", packet_crc)

	var crc_enabled, crc_autoclear_off bool
	switch packet_crc {
	case sensors.RFM_PACKET_CRC_OFF:
		crc_enabled = false
		crc_autoclear_off = false
	case sensors.RFM_PACKET_CRC_AUTOCLEAR_OFF:
		crc_enabled = true
		crc_autoclear_off = true
	case sensors.RFM_PACKET_CRC_AUTOCLEAR_ON:
		crc_enabled = true
		crc_autoclear_off = false
	default:
		return gopi.ErrBadParameter
	}

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Write
	if err := this.setPacketConfig1(this.packet_format, this.packet_coding, this.packet_filter, crc_enabled, crc_autoclear_off); err != nil {
		return err
	}

	// Read
	if _, _, _, crc_enabled_read, crc_autoclear_off_read, err := this.getPacketConfig1(); err != nil {
		return err
	} else if crc_enabled_read != crc_enabled {
		this.log.Debug2("SetPacketCRC expecting crc_enabled=%v, got=%v", crc_enabled, crc_enabled_read)
		return sensors.ErrUnexpectedResponse
	} else if crc_autoclear_off_read != crc_autoclear_off {
		this.log.Debug2("SetPacketCRC expecting crc_autoclear_off=%v, got=%v", crc_autoclear_off, crc_autoclear_off_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.crc_enabled = crc_enabled
		this.crc_auto_clear_off = crc_autoclear_off
	}

	// Success
	return nil
}
