/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package rfm69

import (
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// Channel Filter Settings

func (this *rfm69) RXFilterFrequency() sensors.RFMRXBWFrequency {
	return this.rxbw_frequency
}

func (this *rfm69) RXFilterCutoff() sensors.RFMRXBWCutoff {
	return this.rxbw_cutoff
}

func (this *rfm69) SetRXFilterCutoff(value sensors.RFMRXBWCutoff) error {
	return gopi.ErrNotImplemented
}

func (this *rfm69) SetRXFilter(frequency sensors.RFMRXBWFrequency, cutoff sensors.RFMRXBWCutoff) error {
	this.log.Debug("<sensors.RFM69.SetRXFilterFrequency{ frequency=%v cutoff=%v }", frequency, cutoff)

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Write
	if err := this.setRegRXBW(frequency, cutoff); err != nil {
		return err
	}

	// Read
	if frequency_read, cutoff_read, err := this.getRegRXBW(); err != nil {
		return err
	} else if frequency_read != frequency {
		this.log.Debug2("SetRXFilter expecting frequency=%v, got=%v", frequency, frequency_read)
		return sensors.ErrUnexpectedResponse
	} else if cutoff_read != cutoff {
		this.log.Debug2("SetRXFilter expecting cutoff=%v, got=%v", cutoff, cutoff_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.rxbw_frequency = frequency
		this.rxbw_cutoff = cutoff
	}
	return nil
}
