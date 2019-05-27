/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2019
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package sensordb

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	sensors "github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type SensorDB struct {
	Path           string
	InfluxAddr     string
	InfluxTimeout  time.Duration
	InfluxDatabase string
}

type sensordb struct {
	log gopi.Logger

	// Config and Influxdb
	config
	influxdb
}

type sensor struct {
	Namespace_   string    `json:"ns"`
	Key_         string    `json:"key"`
	Description_ string    `json:"description"`
	TimeCreated_ time.Time `json:"ts_created"`
	TimeSeen_    time.Time `json:"ts_seen"`
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config SensorDB) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sensordb>Open{ config=%+v }", config)

	this := new(sensordb)
	this.log = log

	if err := this.config.Init(config, log); err != nil {
		return nil, err
	}
	if err := this.influxdb.Init(config, log); err != nil {
		return nil, err
	}

	// Return success
	return this, nil
}

func (this *sensordb) Close() error {
	this.log.Debug("<sensordb>Close{ config=%v influxdb=%v }", this.config.String(), this.influxdb.String())

	if err := this.influxdb.Destroy(); err != nil {
		return err
	}
	if err := this.config.Destroy(); err != nil {
		return err
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *sensordb) String() string {
	return fmt.Sprintf("<sensordb>{ config=%v influxdb=%v }", this.config.String(), this.influxdb.String())
}

////////////////////////////////////////////////////////////////////////////////
// DATABASE IMPLEMENTATION

// Return an array of all sensors
func (this *sensordb) Sensors() []sensors.Sensor {
	sensors := make([]sensors.Sensor, len(this.config.Sensors))
	for i, sensor := range this.config.Sensors {
		sensors[i] = sensor
	}
	return sensors
}

// Register a sensor from a message, recording sensor details
// as necessary
func (this *sensordb) Register(message sensors.Message) (sensors.Sensor, error) {
	this.log.Debug2("<sensordb>Register{ message=%v }", message)

	// Return ns and key
	if ns, key, description, err := decode_sensor(message); err != nil {
		return nil, err
	} else if sensor_ := this.config.GetSensorByName(ns, key); sensor_ == nil {
		// Create a new sensor record
		if sensor_ := NewSensor(ns, key, description); sensor_ == nil {
			this.log.Warn("NewSensor: Failed")
			return nil, gopi.ErrAppError
		} else if err := this.config.AddSensor(sensor_); err != nil {
			this.log.Warn("NewSensor: Failed: %v", err)
			return nil, err
		} else {
			return sensor_, nil
		}
	} else if err := this.config.PingSensor(sensor_); err != nil {
		this.log.Warn("PingSensor: Failed: %v", err)
		return nil, err
	} else {
		return sensor_, nil
	}
}

// Lookup an existing sensor based on namespace and key, or nil if not found
func (this *sensordb) Lookup(ns, key string) sensors.Sensor {
	this.log.Debug2("<sensordb>Lookup{ ns=%v key=%v }", strconv.Quote(ns), strconv.Quote(key))
	if ns == "" || key == "" {
		return nil
	}
	return this.config.GetSensorByName(ns, key)
}

// Write out a message to the database
func (this *sensordb) Write(sensor sensors.Sensor, message sensors.Message) error {
	this.log.Debug2("<sensordb>Write{ message=%v }", message)
	return this.influxdb.Write(sensor, message)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// decode_sensor converts a message into a namespace and key
func decode_sensor(message sensors.Message) (string, string, string, error) {
	if message == nil {
		return "", "", "", gopi.ErrBadParameter
	} else if message_, ok := message.(sensors.OOKMessage); ok {
		if product := decode_ook_socket(message_); product == sensors.MIHOME_PRODUCT_NONE {
			return "", "", "", fmt.Errorf("Invalid or unknown product for message: %v", message_)
		} else {
			return message_.Name(), fmt.Sprintf("%02X:%06X", product, message_.Addr()), "Switch", nil
		}
	} else if message_, ok := message.(sensors.OTMessage); ok {
		product := fmt.Sprintf("%v", sensors.MiHomeProduct(message_.Product()))
		if strings.HasPrefix(product, "MIHOME_PRODUCT_") {
			product = strings.TrimPrefix(product, "MIHOME_PRODUCT_")
		}
		return message_.Name(), fmt.Sprintf("%02X:%06X", message_.Product(), message_.Sensor()), product, nil
	} else {
		return "", "", "", sensors.ErrUnexpectedResponse
	}
}

func decode_ook_socket(message sensors.OOKMessage) sensors.MiHomeProduct {
	switch message.Socket() {
	case 0:
		return sensors.MIHOME_PRODUCT_CONTROL_ALL
	case 1:
		return sensors.MIHOME_PRODUCT_CONTROL_ONE
	case 2:
		return sensors.MIHOME_PRODUCT_CONTROL_TWO
	case 3:
		return sensors.MIHOME_PRODUCT_CONTROL_THREE
	case 4:
		return sensors.MIHOME_PRODUCT_CONTROL_FOUR
	default:
		return sensors.MIHOME_PRODUCT_NONE
	}
}
