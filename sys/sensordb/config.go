/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2019
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package sensordb

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	event "github.com/djthorpe/gopi/util/event"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type config_ struct {
	// Public Members
	Sensors []*sensor `json:"sensors"`
}

type config struct {
	// Database
	config_

	// Private Members
	log      gopi.Logger
	path     string
	modified bool

	sync.Mutex
	event.Tasks
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	FILENAME_DEFAULT = "sensors.json"
	WRITE_DELTA      = 30 * time.Second
)

////////////////////////////////////////////////////////////////////////////////
// INIT / DESTROY

func (this *config) Init(config SensorDB, logger gopi.Logger) error {
	logger.Debug("<sensordb.config>Init{ config=%+v }", config)

	this.log = logger
	this.Sensors = make([]*sensor, 0)

	// Read or create file
	if config.Path != "" {
		if err := this.ReadPath(config.Path); err != nil {
			return fmt.Errorf("ReadPath: %v: %v", config.Path, err)
		}
	}

	// Start process to write occasionally to disk
	this.Tasks.Start(this.WriteConfigTask)

	// Success
	return nil
}

func (this *config) Destroy() error {
	this.log.Debug("<sensordb.config>Destroy{ path=%v }", strconv.Quote(this.path))

	// Stop all tasks
	if err := this.Tasks.Close(); err != nil {
		return err
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *config) String() string {
	return fmt.Sprintf("<sensordb.config>{ path=%v num_sensors=%v }", strconv.Quote(this.path), len(this.Sensors))
}

////////////////////////////////////////////////////////////////////////////////
// READ AND WRITE CONFIG

// SetModified sets the modified flag to true
func (this *config) SetModified() {
	this.Lock()
	defer this.Unlock()
	this.modified = true
}

// ReadPath creates regular file if it doesn't exist, or else reads from the path
func (this *config) ReadPath(path string) error {
	this.log.Debug2("<sensordb.config>ReadPath{ path=%v }", strconv.Quote(path))

	// Append home directory if relative path
	if filepath.IsAbs(path) == false {
		if homedir, err := os.UserHomeDir(); err != nil {
			return err
		} else {
			path = filepath.Join(homedir, path)
		}
	}

	// Set path
	this.path = path

	// Append filename
	if stat, err := os.Stat(this.path); err == nil && stat.IsDir() {
		// append default filename
		this.path = filepath.Join(this.path, FILENAME_DEFAULT)
	}

	// Read file
	if stat, err := os.Stat(this.path); err == nil && stat.Mode().IsRegular() {
		if err := this.ReadPath_(this.path); err != nil {
			return err
		} else {
			return nil
		}
	} else if os.IsNotExist(err) {
		// Create file
		if fh, err := os.Create(this.path); err != nil {
			return err
		} else if err := fh.Close(); err != nil {
			return err
		} else {
			this.SetModified()
			return nil
		}
	} else {
		return err
	}
}

// WritePath writes the configuration file to disk
func (this *config) WritePath(path string, indent bool) error {
	this.log.Debug2("<sensordb.config>WritePath{ path=%v indent=%v }", strconv.Quote(path), indent)
	this.Lock()
	defer this.Unlock()
	if fh, err := os.Create(path); err != nil {
		return err
	} else {
		defer fh.Close()
		if err := this.Writer(fh, this.Sensors, indent); err != nil {
			return err
		} else {
			this.modified = false
		}
	}

	// Success
	return nil
}

func (this *config) ReadPath_(path string) error {
	this.Lock()
	defer this.Unlock()

	if fh, err := os.Open(path); err != nil {
		return err
	} else {
		defer fh.Close()
		if err := this.Reader(fh); err != nil {
			return err
		} else {
			this.modified = false
		}
	}

	// Success
	return nil
}

// Reader reads the configuration from an io.Reader object
func (this *config) Reader(fh io.Reader) error {
	dec := json.NewDecoder(fh)
	if err := dec.Decode(&this.config_); err != nil {
		return err
	} else {
		// TODO
		/*
			// Re-create the services and groups
			for i, service := range this.config_.Services {
				this.config_.Services[i] = CopyService(service)
			}
			for i, group := range this.config_.ServiceGroups {
				this.config_.ServiceGroups[i] = CopyGroup(group)
			}
		*/
	}

	// Success
	return nil
}

// Writer writes an array of service records to a io.Writer object
func (this *config) Writer(fh io.Writer, records []*sensor, indent bool) error {
	enc := json.NewEncoder(fh)
	if indent {
		enc.SetIndent("", "  ")
	}
	if err := enc.Encode(this.config_); err != nil {
		return err
	}
	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// FIND SENSOR

func (this *config) GetSensorByName(ns, key string) *sensor {
	this.log.Debug2("<sensordb.config>GetSensorByName{ ns=%v key=%v }", strconv.Quote(ns), strconv.Quote(key))

	this.Lock()
	defer this.Unlock()

	for _, sensor := range this.Sensors {
		if sensor.Key_ == key && sensor.Namespace_ == ns {
			return sensor
		}
	}

	// Not found
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// ADD & REMOVE SENSORS

func (this *config) AddSensor(sensor *sensor) error {
	this.log.Debug2("<sensordb.config>AddSensor{ sensor=%v }", sensor)

	if sensor == nil {
		return gopi.ErrBadParameter
	} else if sensor_ := this.GetSensorByName(sensor.Namespace(), sensor.Key()); sensor_ != nil {
		return fmt.Errorf("Duplicate sensor: %v", sensor_)
	} else {
		this.Lock()
		defer this.Unlock()
		this.Sensors = append(this.Sensors, sensor)
		this.modified = true
	}

	// Success
	return nil
}

func (this *config) PingSensor(sensor *sensor) error {
	this.log.Debug2("<sensordb.config>PingSensor{ sensor=%v }", sensor)
	if sensor == nil {
		return gopi.ErrBadParameter
	} else {
		this.Lock()
		defer this.Unlock()
		sensor.TimeSeen_ = time.Now()
		this.modified = true
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// BACKGROUND TASKS

func (this *config) WriteConfigTask(start chan<- event.Signal, stop <-chan event.Signal) error {
	start <- gopi.DONE
	ticker := time.NewTimer(100 * time.Millisecond)
FOR_LOOP:
	for {
		select {
		case <-ticker.C:
			if this.modified {
				if this.path == "" {
					// Do nothing
				} else if err := this.WritePath(this.path, true); err != nil {
					this.log.Warn("Write: %v: %v", this.path, err)
				}
			}
			ticker.Reset(WRITE_DELTA)
		case <-stop:
			break FOR_LOOP
		}
	}

	// Stop the ticker
	ticker.Stop()

	// Try and write
	if this.modified {
		if this.path == "" {
			// Do nothing
		} else if err := this.WritePath(this.path, true); err != nil {
			this.log.Warn("Write: %v: %v", this.path, err)
		}
	}

	// Success
	return nil
}
