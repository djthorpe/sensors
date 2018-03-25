/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package openthings

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	EncryptionID uint8
	IgnoreCRC    bool
}

type OpenThings struct {
	log           gopi.Logger
	encryption_id uint8
	ignore_crc    bool
}

type Message struct {
	timestamp time.Time
	payload   []byte
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	OT_ENCRYPTION_ID   = 0x01 // Default encryption ID
	OT_PAYLOAD_MINSIZE = 11   // Minimum size of a payload
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config Config) Open(log gopi.Logger) (gopi.Driver, error) {
	this := new(OpenThings)
	this.log = log
	this.ignore_crc = config.IgnoreCRC

	if config.EncryptionID != 0 {
		this.encryption_id = config.EncryptionID
	} else {
		this.encryption_id = OT_ENCRYPTION_ID
	}

	log.Debug("<protocol.openthings.Open>{ EncryptionID=0x%02X IgnoreCRC=%v }", this.encryption_id, config.IgnoreCRC)

	// Return success
	return this, nil
}

func (this *OpenThings) Close() error {
	this.log.Debug("<protocol.openthings.Close>{ EncryptionID=0x%02X }", this.encryption_id)

	// No resources to free

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// DECRYPT

func (this *OpenThings) Decode(payload []byte, ts time.Time) (sensors.OTMessage, error) {
	this.log.Debug2("<protocol.openthings.Decode>{ payload=%v ts=%v }", strings.ToUpper(hex.EncodeToString(payload)), ts)

	message := new(Message)
	message.timestamp = ts
	message.payload = payload

	// Check minimum message size
	if len(message.payload) < OT_PAYLOAD_MINSIZE {
		return message, sensors.ErrMessageCorruption
	}
	// Check size byte vs size of message
	if int(message.payload[0]) != len(payload)-1 {
		return message, sensors.ErrMessageCorruption
	}
	// Check manufacturer is known
	if message.Manufacturer() == sensors.OT_MANUFACTURER_NONE {
		return message, sensors.ErrMessageCorruption
	}

	// Success
	return message, nil
}

////////////////////////////////////////////////////////////////////////////////
// MESSAGE IMPLEMENTATION

func (this *Message) Timestamp() time.Time {
	return this.timestamp
}

func (this *Message) Payload() []byte {
	return this.payload
}

func (this *Message) Size() uint8 {
	if len(this.payload) > 0 {
		return this.payload[0]
	} else {
		return 0
	}
}

func (this *Message) Manufacturer() sensors.OTManufacturer {
	if len(this.payload) >= 2 {
		m := sensors.OTManufacturer(this.payload[1])
		if m <= sensors.OT_MANUFACTURER_MAX {
			return m
		}
	}
	return sensors.OT_MANUFACTURER_NONE
}

func (this *Message) ProductID() uint8 {
	if len(this.payload) >= 3 {
		return this.payload[2]
	} else {
		return 0
	}
}

func (this *Message) SensorID() uint32 {
	return 0
}

func (this *Message) CRC() uint16 {
	return 0
}

func (this *Message) Packet() []byte {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Message) String() string {
	var params []string
	if this.timestamp.IsZero() == false {
		params = append(params, fmt.Sprintf("ts=%v", this.timestamp.Format(time.RFC3339Nano)))
	}
	if this.Size() > 0 {
		params = append(params, fmt.Sprintf("payload_size=%v", this.Size()))
	}
	if this.Manufacturer() != sensors.OT_MANUFACTURER_NONE {
		params = append(params, fmt.Sprintf("manufacturer=%v", this.Manufacturer()))
		params = append(params, fmt.Sprintf("product_id=0x%02X", this.ProductID()))
	} else {
		params = append(params, fmt.Sprintf("payload=%v", strings.ToUpper(hex.EncodeToString(this.payload))))
	}
	return fmt.Sprintf("<protocol.openthings.Message>{ %v }", strings.Join(params, " "))
}
