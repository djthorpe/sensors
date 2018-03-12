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
// GET PARAMETERS

// Return AFC in Hertz
func (this *rfm69) AFC() uint {
	return uint(this.afc) * RFM_FSTEP_HZ
}

// Return AFC Mode
func (this *rfm69) AFCMode() sensors.RFMAFCMode {
	return this.afc_mode
}

// Return AFC Routine
func (this *rfm69) AFCRoutine() sensors.RFMAFCRoutine {
	return this.afc_routine
}

// Set AFC Routine
func (this *rfm69) SetAFCRoutine(afc_routine sensors.RFMAFCRoutine) error {
	this.log.Debug("<sensors.RFM69.SetAFCRoutine>{ afc_routine=%v }", afc_routine)

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Write
	if err := this.setAFCRoutine(afc_routine); err != nil {
		return err
	}

	// Read
	if afc_routine_read, err := this.getAFCRoutine(); err != nil {
		return err
	} else if afc_routine != afc_routine_read {
		this.log.Debug2("SetAFCRoutine expecting afc_routine=%v, got=%v", afc_routine, afc_routine_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.afc_routine = afc_routine
	}

	// Success
	return nil
}

// Set AFC Mode
func (this *rfm69) SetAFCMode(afc_mode sensors.RFMAFCMode) error {
	this.log.Debug("<sensors.RFM69.SetAFCMode>{ afc_mode=%v }", afc_mode)

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Write
	if err := this.setAFCControl(afc_mode, false, false, false); err != nil {
		return err
	}

	// Read
	if afc_mode_read, _, _, err := this.getAFCControl(); err != nil {
		return err
	} else if afc_mode != afc_mode_read {
		this.log.Debug2("SetAFCMode expecting afc_mode=%v, got=%v", afc_mode, afc_mode_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.afc_mode = afc_mode
	}

	// Clear AFC if AFC is off, and read the value back
	if afc_mode == sensors.RFM_AFCMODE_OFF {
		if err := this.setAFCControl(this.afc_mode, false, true, false); err != nil {
			return err
		} else if afc, err := this.getAFC(); err != nil {
			return err
		} else {
			this.afc = afc
		}
	}

	// Success
	return nil
}

// Trigger AFC
func (this *rfm69) TriggerAFC() error {
	this.log.Debug("<sensors.RFM69.TriggerAFC>{}")

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	if err := this.setAFCControl(this.afc_mode, false, false, true); err != nil {
		return err
	}

	// Wait for afc_done bit
	if err := wait_for_condition(func() (bool, error) {
		_, _, afc_done, err := this.getAFCControl()
		return afc_done, err
	}, true, time.Millisecond*1000); err != nil {
		return err
	} else if afc, err := this.getAFC(); err != nil {
		return err
	} else {
		this.afc = afc
	}

	return nil
}
