/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2019
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package sensordb

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	sensors "github.com/djthorpe/sensors"
	influx "github.com/influxdata/influxdb1-client/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type influxdb struct {
	// Private Members
	log    gopi.Logger
	client influx.Client
	db     string
}

////////////////////////////////////////////////////////////////////////////////
// INIT / DESTROY

func (this *influxdb) Init(config SensorDB, logger gopi.Logger) error {
	logger.Debug("<sensordb.influxdb>Init{ config=%+v }", config)

	this.log = logger

	// Set up influxdb
	if config.InfluxAddr == "" {
		// No influx client, return nil
		return nil
	} else if client, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr: config.InfluxAddr,
	}); err != nil {
		return err
	} else if interval, version, err := client.Ping(config.InfluxTimeout); err != nil {
		return err
	} else {
		logger.Info("<sensordb.influxdb>Init{ version=%v interval=%v }", strconv.Quote(version), interval)
		this.client = client
	}

	// Database
	if config.InfluxDatabase == "" {
		return gopi.ErrBadParameter
	} else {
		this.db = config.InfluxDatabase
	}

	// Success
	return nil
}

func (this *influxdb) Destroy() error {
	this.log.Debug("<sensordb.influxdb>Destroy{}")

	// Close client
	if this.client != nil {
		if err := this.client.Close(); err != nil {
			return err
		}
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *influxdb) String() string {
	return fmt.Sprintf("<sensordb.influxdb>{ client=%v db=%v }", this.client, strconv.Quote(this.db))
}

////////////////////////////////////////////////////////////////////////////////
// REGISTER MESSAGE

func (this *influxdb) Write(sensor sensors.Sensor, message sensors.Message) error {
	this.log.Debug2("<sensordb.influxdb>Write{ msg=%v }", message)
	if sensor == nil || message == nil {
		return gopi.ErrBadParameter
	}
	if this.client == nil {
		// Where there is no client, return nil
		return nil
	}
	if batch, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database: this.db,
	}); err != nil {
		return err
	} else if message_, ok := message.(sensors.OTMessage); ok {
		if point, err := this.PointForOTMessage(sensor, message_); err != nil {
			return err
		} else {
			batch.AddPoint(point)
		}
		if err := this.client.Write(batch); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Don't know how to generate data for: %v", message.Name())
	}

	// Success
	return nil
}

func (this *influxdb) PointForOTMessage(sensor sensors.Sensor, message sensors.OTMessage) (*influx.Point, error) {
	// Check parameters
	if sensor == nil || message == nil {
		return nil, gopi.ErrBadParameter
	}

	// Create point
	tags := make(map[string]string)
	fields := make(map[string]interface{})
	tags["data"] = strings.ToUpper(hex.EncodeToString(message.Data()))
	tags["product"] = fmt.Sprintf("0x%02X", message.Product())
	tags["manufacturer"] = fmt.Sprint(message.Manufacturer())
	tags["sensor"] = fmt.Sprintf("0x%06X", message.Sensor())
	tags["ns"] = sensor.Namespace()
	tags["key"] = sensor.Key()
	tags["description"] = sensor.Description()
	if src, ok := message.Source().(gopi.RPCClientConn); ok {
		tags["source"] = src.Addr()
	}

	// Set fields
	for _, record := range message.Records() {
		name := strings.ToLower(strings.TrimPrefix(fmt.Sprint(record.Name()), "OT_PARAM_"))
		// Convert integers and unsigned integers into floats
		v := record.Value()
		switch v.(type) {
		case int:
			fields[name] = float64(v.(int))
		case int8:
			fields[name] = float64(v.(int8))
		case int16:
			fields[name] = float64(v.(int16))
		case int32:
			fields[name] = float64(v.(int32))
		case int64:
			fields[name] = float64(v.(int64))
		case uint:
			fields[name] = float64(v.(uint))
		case uint8:
			fields[name] = float64(v.(uint8))
		case uint16:
			fields[name] = float64(v.(uint16))
		case uint32:
			fields[name] = float64(v.(uint32))
		case uint64:
			fields[name] = float64(v.(uint64))
		default:
			fields[name] = v
		}
	}

	// Return point
	return influx.NewPoint(message.Name(), tags, fields, message.Timestamp())
}
