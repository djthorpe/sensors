/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package rfm69

import "fmt"

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *rfm69) String() string {
	return fmt.Sprintf("sensors.RFM69{ spi=%v }", this.spi)
}
