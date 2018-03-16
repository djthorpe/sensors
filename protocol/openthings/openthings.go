/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package openthings

import (
	"github.com/djthorpe/gopi"

	// Protocol Buffer Definition
	"github.com/djthorpe/sensors/protobuf/message_pb"
)

//go:generate protoc protobuf/message_pb/message_pb.proto --go_out=plugins=grpc:.

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

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	OT_ENCRYPTION_ID = 0x01 // Default encryption ID
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
	this.log.Debug("<protocol.openthings.Close>{ }")

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// DECRYPT

func (this *OpenThings) Decode(payload []byte) *message_pb.Payload {
	this.log.Debug("<protocol.openthings.Decode>{ }")

	message := new(message_pb.Payload)

	// Success
	return message
}
