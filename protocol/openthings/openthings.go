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

type OpenThings struct {
	EncryptionID uint8
	IgnoreCRC    bool
}

type openthings struct {
	log           gopi.Logger
	encryption_id uint8
	ignore_crc    bool
}

type message struct {
	//payload   []byte
	//crc     uint16
	//records []sensors.OTRecord
	manufacturer sensors.OTManufacturer
	product      uint8
	sensor       uint32
	source       sensors.Proto
	ts           time.Time
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	OT_ENCRYPTION_ID       = 0xF2                                            // Default encryption ID
	OT_MESSAGE_HEADER_SIZE = 8                                               // Size of a header in bytes
	OT_MESSAGE_FOOTER_SIZE = 3                                               // Size of a footer in bytes
	OT_PAYLOAD_MINSIZE     = OT_MESSAGE_HEADER_SIZE + OT_MESSAGE_FOOTER_SIZE // Minimum size of a payload
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config OpenThings) Open(log gopi.Logger) (gopi.Driver, error) {
	this := new(openthings)
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

func (this *openthings) Close() error {
	this.log.Debug("<protocol.openthings.Close>{ EncryptionID=0x%02X }", this.encryption_id)

	// No resources to free

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// NAME AND MODE

func (this *openthings) Name() string {
	return "openthings"
}

func (this *openthings) Mode() sensors.MiHomeMode {
	return sensors.MIHOME_MODE_MONITOR
}

func (this *openthings) String() string {
	return fmt.Sprintf("<sensors.protocol>{ name='%v' mode=%v encryption_id=0x%02X ignore_crc=%v }", this.Name(), this.Mode(), this.encryption_id, this.ignore_crc)
}

////////////////////////////////////////////////////////////////////////////////
// CREATE NEW MESSAGE

// Create a new message
func (this *openthings) New(manufacturer sensors.OTManufacturer, product uint8, sensor uint32) (sensors.OTMessage, error) {
	return this.NewWithTimestamp(manufacturer, product, sensor, time.Time{})
}

func (this *openthings) NewWithTimestamp(manufacturer sensors.OTManufacturer, product uint8, sensor uint32, ts time.Time) (sensors.OTMessage, error) {
	this.log.Debug2("<protocol.openthings>NewWithTimestamp{ manufacturer=%v product=%02X sensor=%08X ts=%v }", manufacturer, product, sensor, ts)

	// Check incoming parameters
	if manufacturer == sensors.OT_MANUFACTURER_NONE || manufacturer > sensors.OT_MANUFACTURER_MAX {
		return nil, gopi.ErrBadParameter
	}
	if sensor&0xFFFFFF != sensor {
		return nil, gopi.ErrBadParameter
	}

	// Create message
	message := new(message)
	message.manufacturer = manufacturer
	message.product = product
	message.sensor = sensor
	message.ts = ts
	message.source = this

	// Return message
	return message, nil
}

// Encode a message into a payload
func (this *openthings) Encode(msg sensors.Message) []byte {
	this.log.Debug2("<protocol.openthings>Encode{ msg=%v }", msg)

	// Check for incoming message
	if msg_, ok := msg.(*message); msg_ == nil || ok == false {
		return nil
	} else {
		// Make a message
		payload := msg_.encode_header(uint16(0))
		payload = append(payload, msg_.encode_records()...)
		crc := uint16(0)
		payload = append(payload, msg_.encode_footer(crc)...)
		return payload
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// DECODE

func (this *openthings) Decode(payload []byte, ts time.Time) (sensors.Message, error) {
	this.log.Debug2("<protocol.openthings>Decode>{ payload=%v ts=%v }", strings.ToUpper(hex.EncodeToString(payload)), ts)
	return nil, gopi.ErrNotImplemented
}

/*
func (this *openthings) Decode(payload []byte, ts time.Time) (sensors.Message, error) {
	this.log.Debug2("<protocol.openthings.Decode>{ payload=%v ts=%v }", strings.ToUpper(hex.EncodeToString(payload)), ts)

	message := new(message)
	message.payload = payload
	message.ts = ts

	// Check minimum message size
	if len(message.payload) < OT_PAYLOAD_MINSIZE {
		this.log.Debug2("protocol.openthings.Decode: Payload size too short")
		return message, sensors.ErrMessageCorruption
	}
	// Check size byte vs size of message
	if int(message.payload[0]) != len(payload)-1 {
		this.log.Debug2("protocol.openthings.Decode: Size byte mismatch")
		return message, sensors.ErrMessageCorruption
	}
	// Check manufacturer is known
	if message.Manufacturer() == sensors.OT_MANUFACTURER_NONE {
		this.log.Debug2("protocol.openthings.Decode: Invalid manufacturer code")
		return message, sensors.ErrMessageCorruption
	}

	// Decrypt packet, sanity check to make sure the payload is at least 7 bytes
	decrypted := this.decrypt_message(payload[5:], binary.BigEndian.Uint16(payload[3:]))
	if len(decrypted) < OT_MESSAGE_MINSIZE {
		this.log.Debug2("protocol.openthings.Decode: Message size too short")
		return message, sensors.ErrMessageCorruption
	}

	// Set the sensor ID
	message.sensor_id = binary.BigEndian.Uint32(decrypted[0:]) & 0xFFFFFF00 >> 8

	// Set the CRC value
	message.crc = binary.BigEndian.Uint16(decrypted[len(decrypted)-2:])

	// Check the zero-byte before the CRC value
	if decrypted[len(decrypted)-3] != 0x00 {
		this.log.Debug2("protocol.openthings.Decode: Missing zero byte before CRC")
		return message, sensors.ErrMessageCorruption
	}

	// Check CRC
	if this.ignore_crc == false {
		expected_crc := compute_crc(decrypted[0 : len(decrypted)-2])
		if expected_crc != message.crc {
			this.log.Debug2("protocol.openthings.Decode: CRC mismatch")
			return message, sensors.ErrMessageCRC
		}
	}

	// Read Records
	if records, err := read_records(decrypted[3 : len(decrypted)-2]); err != nil {
		return message, err
	} else {
		message.records = records
	}

	// Success
	return message, nil
}
*/

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Function to decrypt an incoming message
func (this *openthings) decrypt_message(buf []byte, pip uint16) []byte {
	random := seed(this.encryption_id, pip)
	for i := range buf {
		buf[i], random = encrypt_decrypt(buf[i], random)
	}
	return buf
}

// Function to update the seed to match the pip received in the message
func seed(encryption_id uint8, pip uint16) uint16 {
	return (uint16(encryption_id) << 8) ^ pip
}

// Function to encrypt or decrypt the next byte of data in the stream
func encrypt_decrypt(value byte, random uint16) (byte, uint16) {
	for i := 0; i < 5; i++ {
		if random&0x01 > 0x00 {
			random = (random >> 1) ^ uint16(62965)
		} else {
			random = (random >> 1)
		}
	}
	return uint8((random ^ uint16(value) ^ 90)), random
}

// Function to compute the CRC value
func compute_crc(buf []byte) uint16 {
	rem := uint16(0)
	for _, v := range buf {
		rem = rem ^ (uint16(v) << 8)
		for bit := 0; bit < 8; bit++ {
			if rem&(1<<15) != 0 {
				rem = ((rem << 1) ^ 0x1021)
			} else {
				rem = (rem << 1)
			}
		}
	}
	return rem
}

// Return boolean true value when parameter is non-zero
func to_uint8_bool(value uint8) bool {
	if value != 0x00 {
		return true
	} else {
		return false
	}
}
