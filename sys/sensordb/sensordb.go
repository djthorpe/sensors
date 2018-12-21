/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package sensordb

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	sensors "github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// SensorDB is the configuration for the sensor file
type SensorDB struct {
	Path     string
	Filename string
}

type sensordb struct {
	log     gopi.Logger
	path    string
	sensors map[string]map[string]*sensor
}

type root struct {
	Sensors []*sensor `xml:"sensors"`
}

type sensor struct {
	Namespace_   string    `xml:"ns,attr"`
	Key_         string    `xml:"key,attr"`
	Description_ string    `xml:"description"`
	TimeCreated  time.Time `xml:"created,omitempty"`
	TimeSeen     time.Time `xml:"seen,omitempty"`
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	SENSORDB_FILENAME = "sensors.json"
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config SensorDB) path() string {
	if config.Path == "" {
		if root, err := os.Getwd(); err != nil {
			return ""
		} else {
			return root
		}
	} else if stat, err := os.Stat(config.Path); os.IsNotExist(err) || stat.IsDir() == false {
		return ""
	} else {
		return config.Path
	}
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config SensorDB) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sensors.db>{ path=\"%v\" }", config.path())

	this := new(sensordb)
	this.log = log
	this.path = config.path()

	// Return if path is nil
	if this.path == "" {
		return nil, gopi.ErrBadParameter
	}

	// Create sensors map
	this.sensors = make(map[string]map[string]*sensor, 1)

	// Attempt to load the file of sensors - ignore if file doesn't exist
	if err := this.load(); os.IsNotExist(err) {
		// Do nothing
	} else if err != nil {
		return nil, err
	}

	// Return success
	return this, nil
}

func (this *sensordb) Close() error {
	this.log.Debug2("<sensors.db>Close{ path=\"%v\" }", this.path)
	// TODO: Only save if modified
	return this.save()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *sensordb) String() string {
	return fmt.Sprintf("<sensors.db>{ path=\"%v\" }", this.path)
}

////////////////////////////////////////////////////////////////////////////////
// DATABASE LOAD AND SAVE

func (this *sensordb) load() error {
	this.log.Debug2("<sensors.db>Load{ path=\"%v\"}", this.path)

	// Check for regular file
	if stat, err := os.Stat(this.filepath()); os.IsNotExist(err) {
		return err
	} else if stat.IsDir() || stat.Mode().IsRegular() == false {
		return gopi.ErrBadParameter
	}

	// Open the file
	fh, err := os.Open(this.filepath())
	if err != nil {
		return err
	}

	// Create the array of sensors from XML decoding
	var sensors root
	defer fh.Close()
	if err := xml.NewDecoder(fh).Decode(&sensors); err != nil {
		if err == io.EOF {
			return nil
		} else {
			return err
		}
	}

	this.log.Info("<sensors.db>Load: Loading %v sensors", len(sensors.Sensors))

	// Create sensors
	for _, sensor := range sensors.Sensors {
		if err := this.Insert(sensor); err != nil {
			this.log.Error("<sensors.db>Load: %v", err)
		}
	}

	// Return success
	return nil
}

func (this *sensordb) save() error {
	this.log.Debug2("<sensors.db>Save{ path=\"%v\"}", this.path)

	// Compile the array of sensors
	var sensors root
	sensors.Sensors = make([]*sensor, 0)
	if this.sensors != nil {
		for _, sensormap := range this.sensors {
			if sensormap != nil {
				for _, sensor := range sensormap {
					sensors.Sensors = append(sensors.Sensors, sensor)
				}
			}
		}
	}

	this.log.Info("<sensors.db>Save: Saving %v sensors", len(sensors.Sensors))

	// Save the array
	if fh, err := os.Create(this.filepath()); err != nil {
		return err
	} else {
		defer fh.Close()

		// Encode XML
		enc := xml.NewEncoder(fh)
		enc.Indent("", "  ")
		if err := enc.Encode(sensors); err != nil {
			return err
		}

		// Output return
		fh.WriteString("\n\n")
	}

	// Success
	return nil
}

func (this *sensordb) filepath() string {
	return path.Join(this.path, SENSORDB_FILENAME)
}

////////////////////////////////////////////////////////////////////////////////
// REGISTER

// Register checks the source of the message and will create a new
// sensor record if it's not been discovered yet. Returns a true
// value if the sensor is newly created
func (this *sensordb) Register(message sensors.Message) (sensors.Sensor, error) {
	// Return ns and key
	if ns, key, description, err := decode_sensor(message); err != nil {
		return nil, err
	} else if sensor := this.Lookup(ns, key); sensor == nil {
		// Create a new sensor record
		if sensor, err := this.New(ns, key, description); err != nil {
			this.log.Error("<sensors.db>Register{ ns=%v key=%v }: %v", ns, key, err)
			return nil, err
		} else {
			this.log.Info("<sensors.db>New{ sensor=%v }", sensor)
			return sensor, nil
		}
	} else {
		// Bump sensor seen time
		this.Ping(sensor)
		return sensor, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// LOOKUP & NEW

func (this *sensordb) Lookup(ns, key string) sensors.Sensor {
	if this.sensors == nil {
		return nil
	} else if _, exists := this.sensors[ns]; exists == false {
		return nil
	} else if sensor, exists := this.sensors[ns][key]; exists == false {
		return nil
	} else {
		return sensor
	}
}

func (this *sensordb) New(ns, key, description string) (*sensor, error) {
	record := &sensor{ns, key, description, time.Now(), time.Time{}}
	if err := this.Insert(record); err != nil {
		return nil, err
	} else {
		return record, nil
	}
}

func (this *sensordb) Insert(insert *sensor) error {
	if _, exists := this.sensors[insert.Namespace_]; exists == false {
		this.sensors[insert.Namespace_] = make(map[string]*sensor, 1)
	}
	if _, exists := this.sensors[insert.Namespace_][insert.Key_]; exists == true {
		return gopi.ErrBadParameter
	}
	this.sensors[insert.Namespace_][insert.Key_] = insert
	return nil
}

// Sensors returns a list of all sensors in no particular order
func (this *sensordb) Sensors() []sensors.Sensor {
	root := make([]sensors.Sensor, 0)
	if this.sensors != nil {
		for _, sensormap := range this.sensors {
			if sensormap != nil {
				for _, sensor := range sensormap {
					root = append(root, sensor)
				}
			}
		}
	}
	return root
}

////////////////////////////////////////////////////////////////////////////////
// UPDATE SENSOR DETAILS

func (this *sensordb) Ping(s sensors.Sensor) {
	if sensor, ok := s.(*sensor); ok {
		sensor.TimeSeen = time.Now()
	}
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
