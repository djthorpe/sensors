/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package sensors

import (
	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES - ADS1X15 Analog to Digital Convertors
// Note this driver is still in development

type ADS1X15 interface {
	gopi.Driver

	// Return product
	Product() ADS1X15Product
}

type ADS1015 interface {
	ADS1X15
}

type ADS1115 interface {
	ADS1X15
}

type ADS1X15Product uint
type ADS1X15Rate uint16

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	ADS1X15_PRODUCT_NONE ADS1X15Product = iota
	ADS1X15_PRODUCT_1015                // 12-bit ADC with 4 channels
	ADS1X15_PRODUCT_1115                // 16-bit ADC with 4 channels
	ADS1X15_PRODUCT_MAX  = ADS1X15_PRODUCT_1115
)

const (
	ADS1X15_RATE_NONE ADS1X15Rate = iota
	ADS1X15_RATE_8
	ADS1X15_RATE_16
	ADS1X15_RATE_32
	ADS1X15_RATE_64
	ADS1X15_RATE_128
	ADS1X15_RATE_250
	ADS1X15_RATE_475
	ADS1X15_RATE_490
	ADS1X15_RATE_860
	ADS1X15_RATE_920
	ADS1X15_RATE_1600
	ADS1X15_RATE_2400
	ADS1X15_RATE_3300
	ADS1X15_RATE_MAX = ADS1X15_RATE_3300
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (p ADS1X15Product) String() string {
	switch p {
	case ADS1X15_PRODUCT_NONE:
		return "ADS1X15_PRODUCT_NONE"
	case ADS1X15_PRODUCT_1015:
		return "ADS1X15_PRODUCT_1015"
	case ADS1X15_PRODUCT_1115:
		return "ADS1X15_PRODUCT_1115"
	default:
		return "[?? Invalid ADS1X15Product]"
	}
}
