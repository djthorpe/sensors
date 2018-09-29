/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// Low Noise Amplifier Settings

package rfm69

import "github.com/djthorpe/sensors"

func (this *rfm69) LNAImpedance() sensors.RFMLNAImpedance {
	return this.lna_impedance
}

func (this *rfm69) LNAGain() sensors.RFMLNAGain {
	return this.lna_gain
}

func (this *rfm69) LNACurrentGain() (sensors.RFMLNAGain, error) {
	if _, _, lna_gain, err := this.getRegLNA(); err != nil {
		return sensors.RFM_LNA_GAIN_AUTO, err
	} else {
		return lna_gain, nil
	}
}

func (this *rfm69) SetLNA(impedance sensors.RFMLNAImpedance, gain sensors.RFMLNAGain) error {
	this.log.Debug("<sensors.RFM69.SetLNA{ impedance=%v gain=%v }", impedance, gain)

	// Mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Write
	if err := this.setRegLNA(impedance, gain); err != nil {
		return err
	}

	// Read
	if impedance_read, gain_read, _, err := this.getRegLNA(); err != nil {
		return err
	} else if impedance_read != impedance {
		this.log.Debug2("SetLNA expecting impedance=%v, got=%v", impedance, impedance_read)
		return sensors.ErrUnexpectedResponse
	} else if gain_read != gain {
		this.log.Debug2("SetLNA expecting gain=%v, got=%v", gain, gain_read)
		return sensors.ErrUnexpectedResponse
	} else {
		this.lna_impedance = impedance
		this.lna_gain = gain
	}
	return nil
}
