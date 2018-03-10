/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package rfm69

import (
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

	version           uint8
	mode              sensors.RFMMode
	sequencer_off     bool
	listen_on         bool
	data_mode         sensors.RFMDataMode
	modulation        sensors.RFMModulation
	node_address      uint8
	broadcast_address uint8
	afc               int16
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	RFM_SPI_MODE      = gopi.SPI_MODE_0
	RFM_SPI_SPEEDHZ   = 4000000 // 4MHz
	RFM_VERSION_VALUE = 0x24
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

	// Write mode and read back again
	if err := this.setOpMode(mode, this.listen_on, false, this.sequencer_off); err != nil {
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
	} else if sequencer_off_read != this.sequencer_off {
		this.log.Debug2("SetMode expecting sequencer_off=%v, got=%v", this.sequencer_off, sequencer_off_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.mode = mode
	}

	// If RX mode then read AFC value
	if this.mode == sensors.RFM_MODE_RX {
		if afc, err := this.getAfc(); err != nil {
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
// NODE AND BROADCAST ADDRESS

func (this *rfm69) NodeAddress() uint8 {
	return this.node_address
}

func (this *rfm69) BroadcastAddress() uint8 {
	return this.broadcast_address
}

func (this *rfm69) SetNodeAddress(value uint8) error {
	this.log.Debug("<sensors.RFM69.SetNodeAddress{ value=%02X }", value)

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
	this.log.Debug("<sensors.RFM69.SetBroadcastAddress{ value=%02X }", value)

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
		this.node_address = value
	}
	return nil
}
