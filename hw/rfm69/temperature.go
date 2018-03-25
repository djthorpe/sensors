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

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// MEASURE TEMPERATURE

func (this *rfm69) MeasureTemperature(calibration float32) (float32, error) {
	this.log.Debug("<sensors.RFM69.MeasureTemperature>{ calibration=%v }", calibration)

	// Mode needs to be in standby or frequency synth
	if mode := this.Mode(); mode != sensors.RFM_MODE_STDBY && mode != sensors.RFM_MODE_FS {
		return 0, gopi.ErrOutOfOrder
	}

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Wait for not running
	if err := wait_for_condition(this.getRegTemp1, false, time.Millisecond*1000); err != nil {
		return 0, err
	}

	// Trigger temperature sensing
	if err := this.setRegTemp1(); err != nil {
		return 0, err
	}

	// Wait for not running
	if err := wait_for_condition(this.getRegTemp1, false, time.Millisecond*5000); err != nil {
		return 0, err
	}

	// Get temperature value
	temp, err := this.getRegTemp2()
	if err != nil {
		return 0, err
	}

	return float32(RFM_TEMP_COEF-int(temp)) + calibration, nil
}
