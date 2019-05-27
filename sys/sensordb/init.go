/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2019
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package sensordb

import (
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	gopi.RegisterModule(gopi.Module{
		Name: "sensordb",
		Type: gopi.MODULE_TYPE_OTHER,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagString("sensordb.path", "", "Path to sensor database")
			config.AppFlags.FlagString("sensordb.influxdb.addr", "", "URL to influxdb database")
			config.AppFlags.FlagDuration("sensordb.influxdb.timeout", 5*time.Second, "InfluxDB timeout")
			config.AppFlags.FlagString("sensordb.influxdb.db", "sensordb", "InfluxDB database name")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			path, _ := app.AppFlags.GetString("sensordb.path")
			influxdb_addr, _ := app.AppFlags.GetString("sensordb.influxdb.addr")
			influxdb_timeout, _ := app.AppFlags.GetDuration("sensordb.influxdb.timeout")
			influxdb_db, _ := app.AppFlags.GetString("sensordb.influxdb.db")
			return gopi.Open(SensorDB{
				Path:           path,
				InfluxAddr:     influxdb_addr,
				InfluxTimeout:  influxdb_timeout,
				InfluxDatabase: influxdb_db,
			}, app.Logger)
		},
	})
}
