/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package rfm69

import (
	gopi "github.com/djthorpe/gopi"
	sensors "github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Configuration
type RFM69 struct {
	// the SPI driver
	SPI gopi.SPI

	// Device speed
	Speed uint32
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config RFM69) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sensors.RFM69.Open>{ spi=%v speed=%v }", config.SPI, config.Speed)

	this := new(rfm69)
	this.spi = config.SPI
	this.log = log

	if this.spi == nil {
		return nil, gopi.ErrBadParameter
	}

	// Set SPI mode
	if err := this.spi.SetMode(RFM_SPI_MODE); err != nil {
		return nil, err
	}

	// Set SPI speed
	if config.Speed > 0 {
		if err := this.spi.SetMaxSpeedHz(config.Speed); err != nil {
			return nil, err
		}
	} else {
		if err := this.spi.SetMaxSpeedHz(RFM_SPI_SPEEDHZ); err != nil {
			return nil, err
		}
	}

	// Get version - and check against expected value
	if version, err := this.getVersion(); err != nil {
		return nil, sensors.ErrNoDevice
	} else if version != RFM_VERSION_VALUE {
		return nil, sensors.ErrNoDevice
	} else {
		this.version = version
	}

	// Reset Registers
	if err := this.ResetRegisters(); err != nil {
		return nil, err
	}

	// Return success
	return this, nil
}

func (this *rfm69) Close() error {
	this.log.Debug("<sensors.RFM69.Close>{ }")

	// Lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Blank out SPI value
	this.spi = nil

	return nil
}

func (this *rfm69) ResetRegisters() error {
	this.log.Debug2("<sensors.RFM69.ResetRegisters>{ }")

	// Lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Get operational mode
	if mode, listen_on, sequencer_off, err := this.getOpMode(); err != nil {
		return err
	} else {
		this.mode = mode
		this.listen_on = listen_on
		this.sequencer_off = sequencer_off
	}

	// Get data mode and modulation
	if data_mode, modulation, err := this.getDataModul(); err != nil {
		return err
	} else {
		this.data_mode = data_mode
		this.modulation = modulation
	}

	// Get bitrate
	if bitrate, err := this.getBitrate(); err != nil {
		return err
	} else {
		this.bitrate = bitrate
	}

	// Get Frequency Carrier & Deviation
	if frf, err := this.getFreqCarrier(); err != nil {
		return err
	} else if fdev, err := this.getFreqDeviation(); err != nil {
		return err
	} else {
		this.frf = frf
		this.fdev = fdev
	}

	// Automatic frequency correction
	if afc, err := this.getAFC(); err != nil {
		return err
	} else if afc_routine, err := this.getAFCRoutine(); err != nil {
		return err
	} else if afc_mode, _, _, err := this.getAFCControl(); err != nil {
		return err
	} else {
		this.afc = afc
		this.afc_routine = afc_routine
		this.afc_mode = afc_mode
	}

	// Low Noise Amplifer values (last value ignored is the current gain setting)
	if impedance, gain, _, err := this.getRegLNA(); err != nil {
		return err
	} else {
		this.lna_impedance = impedance
		this.lna_gain = gain
	}

	// Channel filter settings
	if frequency, cutoff, err := this.getRegRXBW(); err != nil {
		return err
	} else {
		this.rxbw_frequency = frequency
		this.rxbw_cutoff = cutoff
	}

	// Get Node address and Broadcast address
	if node_address, err := this.getNodeAddress(); err != nil {
		return err
	} else if broadcast_address, err := this.getBroadcastAddress(); err != nil {
		return err
	} else {
		this.node_address = node_address
		this.broadcast_address = broadcast_address
	}

	// AES Key
	if aes_key, err := this.getAESKey(); err != nil {
		return err
	} else {
		this.aes_key = aes_key
	}

	// PacketConfig1 values
	if packet_format, packet_coding, packet_filter, crc_enabled, crc_auto_clear_off, err := this.getPacketConfig1(); err != nil {
		return err
	} else {
		this.packet_format = packet_format
		this.packet_coding = packet_coding
		this.packet_filter = packet_filter
		this.crc_enabled = crc_enabled
		this.crc_auto_clear_off = crc_auto_clear_off
	}

	// PacketConfig2 values
	if rx_inter_packet_delay, rx_auto_restart, aes_on, err := this.getPacketConfig2(); err != nil {
		return err
	} else {
		this.aes_on = aes_on
		this.rx_inter_packet_delay = rx_inter_packet_delay
		this.rx_auto_restart = rx_auto_restart
	}

	// Get preamble and payload sizes
	if preamble_size, err := this.getPreambleSize(); err != nil {
		return err
	} else if payload_size, err := this.getPayloadSize(); err != nil {
		return err
	} else {
		this.preamble_size = preamble_size
		this.payload_size = payload_size
	}

	// Sync word
	if sync_word, err := this.getSyncWord(); err != nil {
		return err
	} else if sync_on, fifo_fill_condition, sync_size, sync_tol, err := this.getSyncConfig(); err != nil {
		return err
	} else {
		this.sync_word = sync_word
		this.sync_on = sync_on
		this.fifo_fill_condition = fifo_fill_condition
		this.sync_size = sync_size
		this.sync_tol = sync_tol
	}

	// Get TX and FIFO parameters
	if tx_start, fifo_threshold, err := this.getFIFOThreshold(); err != nil {
		return err
	} else {
		this.tx_start = tx_start
		this.fifo_threshold = fifo_threshold
	}

	// Return success
	return nil
}
