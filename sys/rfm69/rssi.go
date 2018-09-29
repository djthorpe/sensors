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
)

////////////////////////////////////////////////////////////////////////////////
// MEASURE RSSI

func (this *rfm69) MeasureRSSI() (float32, error) {
	this.log.Debug("<sensors.RFM69.MeasureRSSI>{ }")

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Trigger RSSI measurement sensing
	if err := this.setRegRSSIStart(); err != nil {
		return 0, err
	}

	// Wait for done
	if err := wait_for_condition(this.getRegRSSIDone, true, time.Millisecond*100); err != nil {
		return 0, err
	}

	// Get RSSI value
	if value, err := this.getRegRSSIValue(); err != nil {
		return 0, err
	} else {
		return -float32(value) / 2.0, nil
	}
}
