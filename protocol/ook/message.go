/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package ook

import (
	"encoding/hex"
	"strings"
	// Frameworks
	"fmt"
	"time"

	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// CREATE MESSAGE

func (this *ook) New(addr uint32, socket uint, state bool, data []byte) (sensors.OOKMessage, error) {
	return this.NewWithTimestamp(addr, socket, state, data, time.Time{})
}

func (this *ook) NewWithTimestamp(addr uint32, socket uint, state bool, data []byte, ts time.Time) (sensors.OOKMessage, error) {
	this.log.Debug2("<sensors.protocol.OOK>New{ addr=%05X socket=%v state=%v data=%v ts=%v }", addr, socket, state, strings.ToUpper(hex.EncodeToString(data)), ts)

	// Address is 20-bits
	if addr&OOK_ADDR_MASK != addr {
		return nil, gopi.ErrBadParameter
	}
	// Socket is 0-4
	if socket > 4 {
		return nil, gopi.ErrBadParameter
	}

	// Set up message
	m := new(message)
	m.addr = addr
	m.state = state
	m.socket = socket
	m.source = this
	m.data = data
	m.ts = ts

	return m, nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *message) String() string {
	if this.ts.IsZero() {
		return fmt.Sprintf("<sensors.Message>{ name='%v' addr=0x%05X socket=%v state=%v data=%v }", this.Name(), this.addr, this.socket, this.state, strings.ToUpper(hex.EncodeToString(this.data)))
	} else {
		return fmt.Sprintf("<sensors.Message>{ name='%v' addr=0x%05X socket=%v state=%v data=%v ts=%v }", this.Name(), this.addr, this.socket, this.state, strings.ToUpper(hex.EncodeToString(this.data)), this.ts.Format(time.Kitchen))
	}
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENT OOKMessage INTERFACE

func (this *message) Addr() uint32 {
	return this.addr & OOK_ADDR_MASK
}

func (this *message) State() bool {
	return this.state
}

func (this *message) Socket() uint {
	return this.socket
}

func (this *message) Timestamp() time.Time {
	return this.ts
}

func (this *message) Data() []byte {
	return this.data
}

func (this *message) IsDuplicate(other sensors.Message) bool {
	if this.Name() != other.Name() {
		return false
	}
	if this.Addr() != other.(sensors.OOKMessage).Addr() {
		return false
	}
	if this.State() != other.(sensors.OOKMessage).State() {
		return false
	}
	if this.Socket() != other.(sensors.OOKMessage).Socket() {
		return false
	}
	return true
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENT gopi.Event INTERFACE

func (this *message) Name() string {
	return this.source.Name()
}

func (this *message) Source() gopi.Driver {
	return this.source
}
