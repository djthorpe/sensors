/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package rfm69

import (
	"context"
	"time"

	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// READ & CLEAR FIFO

func (this *rfm69) ClearFIFO() error {
	this.log.Debug("<sensors.RFM69.ClearFIFO>{ }")

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Set IRQ2Flags
	if err := this.setIRQFlags2(); err != nil {
		return err
	} else if fifo_not_empty, err := this.getIRQFlags2(RFM_IRQFLAGS2_FIFONOTEMPTY); err != nil {
		return err
	} else if to_uint8_bool(fifo_not_empty) {
		return sensors.ErrUnexpectedResponse
	}

	return nil
}

func (this *rfm69) ReadFIFO(ctx context.Context) ([]byte, error) {
	this.log.Debug("<sensors.RFM69.ReadFIFO>{ }")

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Check FIFO every 100 milliseconds
	interval := time.NewTicker(100 * time.Millisecond)
	defer interval.Stop()
	for {
		select {
		case <-ctx.Done():
			// Context finished without FIFO
			return nil, nil
		case <-interval.C:
			// Check FIFO
			if fifo_empty, err := this.recvFIFOEmpty(); err != nil {
				return nil, err
			} else if fifo_empty {
				continue
			} else if data, err := this.recvFIFO(); err != nil {
				return nil, err
			} else {
				return data, nil
			}
		}
	}
}

func (this *rfm69) ReadPayload(ctx context.Context) ([]byte, bool, error) {
	this.log.Debug("<sensors.RFM69.ReadPayload>{ }")

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Check FIFO every 100 milliseconds
	interval := time.NewTicker(100 * time.Millisecond)
	defer interval.Stop()
	for {
		select {
		case <-ctx.Done():
			// Context finished without FIFO
			return nil, false, nil
		case <-interval.C:
			// Check Payload
			if payload_ready, err := this.recvPayloadReady(); err != nil {
				return nil, false, err
			} else if payload_ready == false {
				continue
			} else if data, err := this.recvFIFO(); err != nil {
				return nil, false, err
			} else if crc_ok, err := this.recvCRCOk(); err != nil {
				return nil, false, err
			} else {
				return data, crc_ok, nil
			}
		}
	}
}

func (this *rfm69) WriteFIFO(data []byte) error {
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// FIFO THRESHOLD

func (this *rfm69) FIFOThreshold() uint8 {
	return this.fifo_threshold
}

func (this *rfm69) SetFIFOThreshold(fifo_threshold uint8) error {
	this.log.Debug("<sensors.RFM69.SetFIFOThreshold>{ fifo_threshold=%v }", fifo_threshold)

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Write
	if err := this.setFIFOThreshold(this.tx_start, fifo_threshold); err != nil {
		return err
	}

	// Read
	if tx_start_read, fifo_threshold_read, err := this.getFIFOThreshold(); err != nil {
		return err
	} else if tx_start_read != this.tx_start {
		this.log.Debug2("SetFIFOThreshold expecting tx_start=%v, got=%v", this.tx_start, tx_start_read)
		return sensors.ErrUnexpectedResponse
	} else if fifo_threshold_read != fifo_threshold {
		this.log.Debug2("SetFIFOThreshold expecting fifo_threshold=%v, got=%v", fifo_threshold, fifo_threshold_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.tx_start = tx_start_read
		this.fifo_threshold = fifo_threshold_read
	}

	// Success
	return nil
}
