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
	"fmt"
	"time"

	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// CREATE MESSAGE

func (this *ook) New(addr uint32, socket uint, state bool) (sensors.OOKMessage, error) {
	this.log.Debug2("<sensors.protocol.OOK>New{ addr=%05X socket=%v state=%v }", addr, socket, state)

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

	return m, nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *message) String() string {
	return fmt.Sprintf("<OOKMessage>{ addr=0x%05X socket=%v state=%v }", this.addr, this.socket, this.state)
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

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENT gopi.Event INTERFACE

func (this *message) Name() string {
	return "OOKMessage"
}

func (this *message) Source() gopi.Driver {
	return this.source
}
