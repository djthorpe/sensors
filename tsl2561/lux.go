/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package tsl2561

import "github.com/djthorpe/sensors"

const (
	LUX_SCALE     = 14     // scale by 2^14
	RATIO_SCALE   = 9      // scale ratio by 2^9
	CH_SCALE      = 10     // scale channel values by 2^10
	CHSCALE_TINT0 = 0x7517 // 322/11 * 2^CH_SCALE
	CHSCALE_TINT1 = 0x0FE7 // 322/81 * 2^CH_SCALE
)

const (
	K1T = 0x0040 // 0.125 * 2^RATIO_SCALE
	B1T = 0x01f2 // 0.0304 * 2^LUX_SCALE
	M1T = 0x01be // 0.0272 * 2^LUX_SCALE
	K2T = 0x0080 // 0.250 * 2^RATIO_SCA
	B2T = 0x0214 // 0.0325 * 2^LUX_SCALE
	M2T = 0x02d1 // 0.0440 * 2^LUX_SCALE
	K3T = 0x00c0 // 0.375 * 2^RATIO_SCALE
	B3T = 0x023f // 0.0351 * 2^LUX_SCALE
	M3T = 0x037b // 0.0544 * 2^LUX_SCALE
	K4T = 0x0100 // 0.50 * 2^RATIO_SCALE
	B4T = 0x0270 // 0.0381 * 2^LUX_SCALE
	M4T = 0x03fe // 0.0624 * 2^LUX_SCALE
	K5T = 0x0138 // 0.61 * 2^RATIO_SCALE
	B5T = 0x016f // 0.0224 * 2^LUX_SCALE
	M5T = 0x01fc // 0.0310 * 2^LUX_SCALE
	K6T = 0x019a // 0.80 * 2^RATIO_SCALE
	B6T = 0x00d2 // 0.0128 * 2^LUX_SCALE
	M6T = 0x00fb // 0.0153 * 2^LUX_SCALE
	K7T = 0x029a // 1.3 * 2^RATIO_SCALE
	B7T = 0x0018 // 0.00146 * 2^LUX_SCALE
	M7T = 0x0012 // 0.00112 * 2^LUX_SCALE
	K8T = 0x029a // 1.3 * 2^RATIO_SCALE
	B8T = 0x0000 // 0.000 * 2^LUX_SCALE
	M8T = 0x0000 // 0.000 * 2^LUX_SCALE

	K1C = 0x0043 // 0.130 * 2^RATIO_SCALE
	B1C = 0x0204 // 0.0315 * 2^LUX_SCALE
	M1C = 0x01ad // 0.0262 * 2^LUX_SCALE
	K2C = 0x0085 // 0.260 * 2^RATIO_SCALE
	B2C = 0x0228 // 0.0337 * 2^LUX_SCALE
	M2C = 0x02c1 // 0.0430 * 2^LUX_SCALE
	K3C = 0x00c8 // 0.390 * 2^RATIO_SCALE
	B3C = 0x0253 // 0.0363 * 2^LUX_SCALE
	M3C = 0x0363 // 0.0529 * 2^LUX_SCALE
	K4C = 0x010a // 0.520 * 2^RATIO_SCALE
	B4C = 0x0282 // 0.0392 * 2^LUX_SCALE
	M4C = 0x03df // 0.0605 * 2^LUX_SCALE
	K5C = 0x014d // 0.65 * 2^RATIO_SCALE
	B5C = 0x0177 // 0.0229 * 2^LUX_SCALE
	M5C = 0x01dd // 0.0291 * 2^LUX_SCALE
	K6C = 0x019a // 0.80 * 2^RATIO_SCALE
	B6C = 0x0101 // 0.0157 * 2^LUX_SCALE
	M6C = 0x0127 // 0.0180 * 2^LUX_SCALE
	K7C = 0x029a // 1.3 * 2^RATIO_SCALE
	B7C = 0x0037 // 0.00338 * 2^LUX_SCALE
	M7C = 0x002b // 0.00260 * 2^LUX_SCALE
	K8C = 0x029a // 1.3 * 2^RATIO_SCALE
	B8C = 0x0000 // 0.000 * 2^LUX_SCALE
	M8C = 0x0000 // 0.000 * 2^LUX_SCALE
)

func calculate_illuminance_lux(value0, value1 uint16, gain sensors.TSL2561Gain, integrate_time sensors.TSL2561IntegrateTime, package_type string) float64 {
	var ch_scale uint64

	switch integrate_time {
	case sensors.TSL2561_INTEGRATETIME_13P7MS:
		ch_scale = CHSCALE_TINT0
	case sensors.TSL2561_INTEGRATETIME_101MS:
		ch_scale = CHSCALE_TINT1
	default:
		ch_scale = (1 << CH_SCALE)
	}
	switch gain {
	case sensors.TSL2561_GAIN_16:
		ch_scale = ch_scale << 4
	}

	// Perform the scaling
	channel0 := (uint64(value0) * ch_scale) >> CH_SCALE
	channel1 := (uint64(value1) * ch_scale) >> CH_SCALE

	// Find the ratio of the channel values (value1/value0)
	// protect agains divide by zero and round ratio value
	var ratio uint64
	if value0 != 0 {
		ratio = (channel1 << (RATIO_SCALE + 1)) / channel0
	}
	ratio = (ratio + 1) >> 1

	// Depending on package type
	var b, m uint64
	switch package_type {
	case "T": // T package
		if (ratio >= 0) && (ratio <= K1T) {
			b = B1T
			m = M1T
		} else if ratio <= K2T {
			b = B2T
			m = M2T
		} else if ratio <= K3T {
			b = B3T
			m = M3T
		} else if ratio <= K4T {
			b = B4T
			m = M4T
		} else if ratio <= K5T {
			b = B5T
			m = M5T
		} else if ratio <= K6T {
			b = B6T
			m = M6T
		} else if ratio <= K7T {
			b = B7T
			m = M7T
		} else if ratio > K8T {
			b = B8T
			m = M8T
		}
	case "CS": // CS package
		if (ratio >= 0) && (ratio <= K1C) {
			b = B1C
			m = M1C
		} else if ratio <= K2C {
			b = B2C
			m = M2C
		} else if ratio <= K3C {
			b = B3C
			m = M3C
		} else if ratio <= K4C {
			b = B4C
			m = M4C
		} else if ratio <= K5C {
			b = B5C
			m = M5C
		} else if ratio <= K6C {
			b = B6C
			m = M6C
		} else if ratio <= K7C {
			b = B7C
			m = M7C
		}
	default:
		panic("Invalid package_type parameter")
	}

	lux := float64(channel0)*float64(b) - float64(channel1)*float64(m)
	// Don't allow lux value below 0
	if lux < 0 {
		lux = 0
	}

	// round LSB
	//temp += (1 << (LUX_SCALE - 1))
	// strip off fractional portion
	//lux := temp >> LUX_SCALE

	return lux
}
