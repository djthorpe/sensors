/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package openthings

import (
	"errors"

	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register protocol/openthings module
	gopi.RegisterModule(gopi.Module{
		Name: "protocol/openthings",
		Type: gopi.MODULE_TYPE_OTHER,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagUint("ot.encryption_id", 0, "OpenThings Encryption ID")
			config.AppFlags.FlagBool("ot.ignore_crc", false, "Ignore CRC checking")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			ignore_crc, _ := app.AppFlags.GetBool("ot.ignore_crc")
			encryption_id, _ := app.AppFlags.GetUint("ot.encryption_id")
			if encryption_id > 0xFF {
				return nil, errors.New("Invalid -ot.encryption_id flag")
			}
			return gopi.Open(Config{
				EncryptionID: uint8(encryption_id),
				IgnoreCRC:    ignore_crc,
			}, app.Logger)
		},
	})
}
