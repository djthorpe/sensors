/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package ook

import (
	// Frameworks

	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type OOK struct{}

type ook struct {
	log gopi.Logger
}

type message struct {
	addr   uint32
	state  bool
	socket uint
	source sensors.Proto
	data   []byte
	ts     time.Time
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS, GLOBAL VARIABLES

var (
	OOK_PREAMBLE         = []byte{0x80, 0x00, 0x00, 0x00} // OOK Preamble sent before each command
	OOK_ADDR_MASK uint32 = 0xFFFFF                        // Length of the address is 20 bits
)

const (
	// Definition of a bit
	OOK_ZERO byte = 0x08
	OOK_ONE  byte = 0x0E
)

const (
	OOK_NONE    byte = 0x00
	OOK_ON_ALL  byte = 0x0D
	OOK_OFF_ALL byte = 0x0C
	OOK_ON_1    byte = 0x0F
	OOK_OFF_1   byte = 0x0E
	OOK_ON_2    byte = 0x07
	OOK_OFF_2   byte = 0x06
	OOK_ON_3    byte = 0x0B
	OOK_OFF_3   byte = 0x0A
	OOK_ON_4    byte = 0x03
	OOK_OFF_4   byte = 0x02
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config OOK) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sensors.protocol.OOK>Open{ }")

	this := new(ook)
	this.log = log

	// Return success
	return this, nil
}

func (this *ook) Close() error {
	this.log.Debug("<sensors.protocol.OOK>Close{ }")

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// NAME AND MODE

func (this *ook) String() string {
	return fmt.Sprintf("<sensors.protocol>{ name='%v' mode=%v }", this.Name(), this.Mode())
}

func (this *ook) Name() string {
	return "ook"
}

func (this *ook) Mode() sensors.MiHomeMode {
	return sensors.MIHOME_MODE_CONTROL
}

////////////////////////////////////////////////////////////////////////////////
// ENCODE AND DECODE

/*
 The payload is 16 bytes (preamble 4 bytes, address 10 bytes, command 2 bytes)
*/

func (this *ook) Encode(msg sensors.Message) []byte {
	this.log.Debug2("<sensors.protocol.OOK>Encode{ msg=%v }", msg)

	payload := make([]byte, 0, 16)

	// Ensure message is of type OOKMessage
	if msg_, ok := msg.(sensors.OOKMessage); ok == false {
		return nil
	} else {

		// Four bytes for the payload
		payload = append(payload, OOK_PREAMBLE...)

		// 20 bits for the address
		payload = append(payload, encodeByte(byte(msg_.Addr() >> 16))[2:]...)
		payload = append(payload, encodeByte(byte(msg_.Addr()>>8))...)
		payload = append(payload, encodeByte(byte(msg_.Addr()>>0))...)

		// 16 bits for the on/off and socket
		payload = append(payload, encodeCommand(msg_)[2:]...)

		// Return the payload
		return payload
	}
}

func (this *ook) Decode(payload []byte, ts time.Time) (sensors.Message, error) {
	this.log.Debug2("<sensors.protocol.OOK>Decode{ payload=%v ts=%v }", strings.ToUpper(hex.EncodeToString(payload)), ts)

	// Check for a 16 byte payload
	if len(payload) != 16 {
		return nil, sensors.ErrMessageCorruption
	}
	// Check payload and construct
	n := 0
	v := uint32(0)
	addr := uint32(0)
	socket := uint(0)
	state := bool(false)
	for i, by := range payload {
		// Preamble
		if i < len(OOK_PREAMBLE) {
			if by != OOK_PREAMBLE[i] {
				return nil, sensors.ErrMessageCorruption
			}
			continue
		}
		// Nibbles
		v <<= 1
		switch (by & 0xF0) >> 4 {
		case OOK_ZERO:
			break
		case OOK_ONE:
			v |= 1
		default:
			return nil, sensors.ErrMessageCorruption
		}
		v <<= 1
		switch by & 0x0F {
		case OOK_ZERO:
			break
		case OOK_ONE:
			v |= 1
		default:
			return nil, sensors.ErrMessageCorruption
		}
		n += 1
		if n == 10 {
			addr = v
			v = 0
		}
		if n == 12 {
			if a, b, err := decodeCommand(byte(v)); err != nil {
				return nil, err
			} else {
				socket = a
				state = b
			}
		}
	}
	return this.NewWithTimestamp(addr, socket, state, payload, ts)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func encodeCommand(msg sensors.OOKMessage) []byte {
	socket := msg.Socket()
	state := msg.State()
	switch {
	case socket == 0 && state == true:
		return encodeByte(OOK_ON_ALL)
	case socket == 1 && state == true:
		return encodeByte(OOK_ON_1)
	case socket == 2 && state == true:
		return encodeByte(OOK_ON_2)
	case socket == 3 && state == true:
		return encodeByte(OOK_ON_3)
	case socket == 4 && state == true:
		return encodeByte(OOK_ON_4)
	case socket == 0 && state == false:
		return encodeByte(OOK_OFF_ALL)
	case socket == 1 && state == false:
		return encodeByte(OOK_OFF_1)
	case socket == 2 && state == false:
		return encodeByte(OOK_OFF_2)
	case socket == 3 && state == false:
		return encodeByte(OOK_OFF_3)
	case socket == 4 && state == false:
		return encodeByte(OOK_OFF_4)
	}
	return nil
}

func decodeCommand(value byte) (uint, bool, error) {
	switch value {
	case OOK_ON_ALL:
		return 0, true, nil
	case OOK_ON_1:
		return 1, true, nil
	case OOK_ON_2:
		return 2, true, nil
	case OOK_ON_3:
		return 3, true, nil
	case OOK_ON_4:
		return 4, true, nil
	case OOK_OFF_ALL:
		return 0, false, nil
	case OOK_OFF_1:
		return 1, false, nil
	case OOK_OFF_2:
		return 2, false, nil
	case OOK_OFF_3:
		return 3, false, nil
	case OOK_OFF_4:
		return 4, false, nil
	default:
		return 0, false, sensors.ErrMessageCorruption
	}
}

func encodeByte(value byte) []byte {
	// A byte is encoded as 4 bytes (each bit is converted to an 8 or an E - or 4 bits)
	encoded := make([]byte, 4)
	for i := 0; i < 4; i++ {
		by := byte(0)
		for j := 0; j < 2; j++ {
			by <<= 4
			if (value & 0x80) == 0 {
				by |= OOK_ZERO
			} else {
				by |= OOK_ONE
			}
			value <<= 1
		}
		encoded[i] = by
	}
	return encoded
}
