/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package openthings

import (
	"encoding/binary"
	"encoding/hex"
	"strings"

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

type Message struct {
	payload   []byte
	sensor_id uint32
	crc       uint16
	records   []sensors.OTRecord
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	OT_ENCRYPTION_ID   = 0xF2 // Default encryption ID
	OT_PAYLOAD_MINSIZE = 11   // Minimum size of a payload
	OT_MESSAGE_MINSIZE = 7    // Minimum size of a decypted message
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
// DECRYPT

func (this *openthings) Decode(payload []byte) (sensors.OTMessage, error) {
	this.log.Debug("<protocol.openthings.Decode>{ payload=%v }", strings.ToUpper(hex.EncodeToString(payload)))

	message := new(Message)
	message.payload = payload

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
