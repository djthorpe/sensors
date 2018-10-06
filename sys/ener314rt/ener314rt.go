/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package ener314rt

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi/util/event"
	"github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// STRUCTS

// Configuration
type MiHome struct {
	GPIO       gopi.GPIO     // GPIO interface
	Radio      sensors.RFM69 // Radio interface
	PinReset   gopi.GPIOPin  // Reset pin
	PinLED1    gopi.GPIOPin  // LED1 (Green, Rx) pin
	PinLED2    gopi.GPIOPin  // LED2 (Red, Tx) pin
	CID        string        // OOK device address
	Repeat     uint          // Number of times to repeat messages by default
	TempOffset float32       // Temperature Offset
}

// mihome driver
type mihome struct {
	log        gopi.Logger
	gpio       gopi.GPIO
	radio      sensors.RFM69
	reset      gopi.GPIOPin
	addr       uint32
	repeat     uint
	tempoffset float32
	led1       gopi.GPIOPin
	led2       gopi.GPIOPin
	ledrx      gopi.GPIOPin
	ledtx      gopi.GPIOPin
	mode       sensors.MiHomeMode

	// protocol storage
	proto_map  map[string]sensors.Proto
	proto_mode map[sensors.MiHomeMode][]sensors.Proto

	// event publisher
	event.Publisher

	// Locker
	sync.Mutex
}

type LED uint

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS, GLOBAL VARIABLES

const (
	// Default Control Address
	ADDR_DEFAULT = "06C6C6"
	// Default number of times to repeat command
	REPEAT_DEFAULT = 8
)

const (
	LED_ALL LED = iota
	LED_1
	LED_2
	LED_RX
	LED_TX
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config MiHome) Open(log gopi.Logger) (gopi.Driver, error) {
	// Set the default CID
	if config.CID == "" {
		config.CID = ADDR_DEFAULT
	}
	if config.Repeat == 0 {
		config.Repeat = REPEAT_DEFAULT
	}
	log.Debug("<sensors.energenie.MiHome>Open{ reset=%v led1=%v led2=%v cid=\"%v\" repeat=%v tempoffset=%v }", config.PinReset, config.PinLED1, config.PinLED2, config.CID, config.Repeat, config.TempOffset)

	if config.GPIO == nil || config.Radio == nil {
		// Fail when either GPIO or Radio is nil
		return nil, gopi.ErrBadParameter
	}

	this := new(mihome)
	this.log = log
	this.gpio = config.GPIO
	this.radio = config.Radio
	this.reset = config.PinReset

	// Set LED's
	this.led1 = config.PinLED1
	this.led2 = config.PinLED2
	this.ledrx = config.PinLED1
	this.ledtx = config.PinLED2
	if this.ledtx == gopi.GPIO_PIN_NONE {
		// Where the second LED doesn't exist, make it the first LED
		this.ledtx = this.led1
	} else if this.ledrx == gopi.GPIO_PIN_NONE {
		// Where the first LED doesn't exist, make it the second LED
		this.ledrx = this.led2
	}

	// Set the default Control Address for legacy OOK devices
	if addr, err := hex.DecodeString(config.CID); err != nil {
		return nil, err
	} else if len(addr) != 3 {
		return nil, gopi.ErrBadParameter
	} else {
		this.addr = uint32(addr[2]) | uint32(addr[1])<<8 | uint32(addr[0])<<16
	}

	// Set number of times to repeat TX by default
	this.repeat = config.Repeat

	// Set the temperature calibration offset
	this.tempoffset = config.TempOffset

	// Set mode to undefined
	this.mode = sensors.MIHOME_MODE_NONE

	// Create map of protocols
	this.proto_map = make(map[string]sensors.Proto, 0)
	this.proto_mode = make(map[sensors.MiHomeMode][]sensors.Proto, 0)

	// Return success
	return this, nil
}

func (this *mihome) Close() error {
	this.log.Debug("<sensors.energenie.MiHome>Close{ addr=0x%05X }", this.addr)

	// Lock until finished
	this.Lock()
	defer this.Unlock()

	// Close publisher
	this.Publisher.Close()

	// Free resources
	this.gpio = nil
	this.radio = nil
	this.proto_map = nil
	this.proto_mode = nil

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *mihome) String() string {
	proto_names := make([]string, 0, len(this.proto_map))
	for proto_name := range this.proto_map {
		proto_names = append(proto_names, proto_name)
	}
	return fmt.Sprintf("<sensors.energenie.MiHome>{ protos=%v addr=0x%05X mode=%v reset=%v led1=%v led2=%v ledrx=%v ledtx=%v gpio=%v radio=%v }", strings.Join(proto_names, ","), this.addr, this.mode, this.reset, this.led1, this.led2, this.ledrx, this.ledtx, this.gpio, this.radio)
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *mihome) ResetRadio() error {
	// If reset is not defined, then return not implemented
	if this.reset == gopi.GPIO_PIN_NONE {
		return gopi.ErrNotImplemented
	}

	// Ensure pin is output
	this.gpio.SetPinMode(this.reset, gopi.GPIO_OUTPUT)

	// Turn all LED's on
	if err := this.SetLED(LED_ALL, gopi.GPIO_HIGH); err != nil {
		return err
	}

	// Pull reset high for 100ms and then low for 5ms
	this.gpio.WritePin(this.reset, gopi.GPIO_HIGH)
	time.Sleep(time.Millisecond * 100)
	this.gpio.WritePin(this.reset, gopi.GPIO_LOW)
	time.Sleep(time.Millisecond * 5)

	// Turn all LED's off
	if err := this.SetLED(LED_ALL, gopi.GPIO_LOW); err != nil {
		return err
	}

	// Set undefined mode
	this.mode = sensors.MIHOME_MODE_NONE

	return nil
}

func (this *mihome) SetLED(led LED, state gopi.GPIOState) error {
	switch led {
	case LED_ALL:
		if this.led1 != gopi.GPIO_PIN_NONE {
			if err := this.SetLED(LED_1, state); err != nil {
				return err
			}
		}
		if this.led2 != gopi.GPIO_PIN_NONE {
			if err := this.SetLED(LED_2, state); err != nil {
				return err
			}
		}
	case LED_1:
		if this.led1 == gopi.GPIO_PIN_NONE {
			return gopi.ErrNotImplemented
		} else {
			this.gpio.SetPinMode(this.led1, gopi.GPIO_OUTPUT)
			this.gpio.WritePin(this.led1, state)
		}
	case LED_2:
		if this.led2 == gopi.GPIO_PIN_NONE {
			return gopi.ErrNotImplemented
		} else {
			this.gpio.SetPinMode(this.led2, gopi.GPIO_OUTPUT)
			this.gpio.WritePin(this.led2, state)
		}
	case LED_RX:
		if this.ledrx == gopi.GPIO_PIN_NONE {
			// Allow to silently do nothing where device does have RX indicator
			return nil
		} else {
			this.gpio.SetPinMode(this.ledrx, gopi.GPIO_OUTPUT)
			this.gpio.WritePin(this.ledrx, state)
		}
	case LED_TX:
		if this.ledtx == gopi.GPIO_PIN_NONE {
			// Allow to silently do nothing where device does have RX indicator
			return nil
		} else {
			this.gpio.SetPinMode(this.ledtx, gopi.GPIO_OUTPUT)
			this.gpio.WritePin(this.ledtx, state)
		}
	default:
		return gopi.ErrBadParameter
	}
	return nil
}

// Receive OOK and FSK payloads until context is cancelled or timeout
func (this *mihome) Receive(ctx context.Context, mode sensors.MiHomeMode) error {
	this.log.Debug2("<sensors.ener314rt>Receive{ mode=%v }", mode)

	// Lock until finished
	this.Lock()
	defer this.Unlock()

	// Switch into correct mode
	if this.mode != mode {
		this.log.Debug("<sensors.ener314rt>Receive: Switch radio mode=%v", mode)
		if err := this.SetMode(mode); err != nil {
			return err
		} else {
			this.mode = mode
		}
	}

	this.log.Debug("<sensors.ener314rt>Receive: Set RX Mode")

	// Set RX mode
	defer this.radio.SetMode(sensors.RFM_MODE_STDBY)
	if err := this.radio.SetMode(sensors.RFM_MODE_RX); err != nil {
		return err
	} else if err := this.radio.SetSequencer(true); err != nil {
		return err
	}

	// Get protocols which could be used to decode the message
	protocols := this.protocols_for_mode(mode)
	protocols = append(protocols, this.protocols_for_mode(sensors.MIHOME_MODE_NONE)...)
	if len(protocols) == 0 {
		this.log.Error("No protocols registered for mode %v", mode)
		return fmt.Errorf("No protocols registered for mode %v", mode)
	}

	// Repeatedly read until context is done
	this.log.Debug("<sensors.ener314rt>Receive: Start Receive Loop")
FOR_LOOP:
	for {
		select {
		case <-ctx.Done():
			break FOR_LOOP
		default:
			this.log.Debug("<sensors.ener314rt>Receive: ReadPayload")
			if data, _, err := this.radio.ReadPayload(ctx); err != nil {
				return err
			} else if data != nil {
				// RX light on
				this.SetLED(LED_RX, gopi.GPIO_HIGH)
				defer this.SetLED(LED_RX, gopi.GPIO_LOW)

				// Emit payload
				if err := this.emit(data, protocols); err != nil {
					this.log.Warn("<sensors.ener314rt>Receive: %v", err)
					if err := this.radio.ClearFIFO(); err != nil {
						this.log.Error("<sensors.ener314rt>Receive: ClearFIFO: %v", err)
					}
				}
			}
		}
	}

	// Success
	return nil
}

// Send a raw payload
func (this *mihome) Send(payload []byte, repeat uint, mode sensors.MiHomeMode) error {
	this.log.Debug2("<sensors.ener314rt>Send{ mode=%v payload=%v repeat=%v }", mode, strings.ToUpper(hex.EncodeToString(payload)), repeat)

	// Lock until finished
	this.Lock()
	defer this.Unlock()

	// Check parameters
	if len(payload) == 0 || repeat == 0 {
		return gopi.ErrBadParameter
	}

	// Switch into correct mode
	if this.mode != mode {
		if err := this.SetMode(mode); err != nil {
			return err
		} else {
			this.mode = mode
		}
	}

	// Set TX Mode
	defer this.radio.SetMode(sensors.RFM_MODE_STDBY)
	if err := this.radio.SetMode(sensors.RFM_MODE_TX); err != nil {
		return err
	} else if err := this.radio.SetSequencer(true); err != nil {
		return err
	}

	// TX light on
	this.SetLED(LED_TX, gopi.GPIO_HIGH)
	defer this.SetLED(LED_TX, gopi.GPIO_LOW)

	// Write payload
	if err := this.radio.WritePayload(payload, repeat); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *mihome) MeasureTemperature() (float32, error) {
	this.log.Debug2("<sensors.ener314rt>MeasureTemperature{ tempoffset=%v }", this.tempoffset)

	// Lock until finished
	this.Lock()
	defer this.Unlock()

	// Need to put into standby mode to measure the temperature
	old_mode := this.radio.Mode()
	if old_mode != sensors.RFM_MODE_STDBY {
		if err := this.radio.SetMode(sensors.RFM_MODE_STDBY); err != nil {
			return 0, err
		}
	}

	// Perform the measurement
	value, err := this.radio.MeasureTemperature(this.tempoffset)

	// Return to previous mode of operation
	if old_mode != sensors.RFM_MODE_STDBY {
		if err := this.radio.SetMode(old_mode); err != nil {
			return 0, err
		}
	}

	// Return the value and error condition
	return value, err
}

func (this *mihome) SetMode(mode sensors.MiHomeMode) error {
	this.log.Debug2("<sensors.ener314rt>SetMode{ mode=%v }", mode)
	switch mode {
	case sensors.MIHOME_MODE_MONITOR:
		if err := this.setFSKMode(); err != nil {
			return err
		}
	case sensors.MIHOME_MODE_CONTROL:
		if err := this.setOOKMode(); err != nil {
			return err
		}
	default:
		return gopi.ErrBadParameter
	}

	// Success
	this.mode = mode
	return nil
}

func (this *mihome) AddProto(proto sensors.Proto) error {
	this.log.Debug2("<sensors.ener314rt>AddProto{ proto=%v mode=%v }", proto, proto.Mode())

	// Check to see if protocol is alreay added
	if _, exists := this.proto_map[proto.Name()]; exists {
		return gopi.ErrBadParameter
	} else {
		this.proto_map[proto.Name()] = proto
	}

	// Create an array for holding the protocol
	arr, exists := this.proto_mode[proto.Mode()]
	if exists == false {
		arr = make([]sensors.Proto, 0, 1)
	}

	// Append the protocol
	this.proto_mode[proto.Mode()] = append(arr, proto)

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - ENER314

// Satisfies the ENER314 interface to switch sockets on
func (this *mihome) On(sockets ...uint) error {
	this.log.Debug2("<sensors.ener314rt>On{ sockets=%v }", sockets)
	return this.onoff(true, sockets)
}

// Satisfies the ENER314 interface to switch sockets off
func (this *mihome) Off(sockets ...uint) error {
	this.log.Debug2("<sensors.ener314rt>Off{ sockets=%v }", sockets)
	return this.onoff(false, sockets)
}

func (this *mihome) onoff(state bool, sockets []uint) error {
	messages := make(map[uint]sensors.OOKMessage, len(sockets))

	// Append 'all' where no socket arguments
	if len(sockets) == 0 {
		sockets = append(sockets, 0)
	}

	// Get OOK protocol
	if proto := this.protocol_for_name("ook").(sensors.OOKProto); proto == nil {
		return fmt.Errorf("OOK Protocol is not available")
	} else {
		// Create messages
		for _, socket := range sockets {
			if message, err := proto.New(this.addr, socket, state); err != nil {
				return err
			} else {
				messages[socket] = message
			}
		}

		// Send messages
		for _, message := range messages {
			if err := this.Send(proto.Encode(message), this.repeat, sensors.MIHOME_MODE_CONTROL); err != nil {
				return err
			}
		}
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *mihome) setFSKMode() error {

	if err := this.radio.SetMode(sensors.RFM_MODE_STDBY); err != nil {
		return err
	} else if err := this.radio.SetModulation(sensors.RFM_MODULATION_FSK); err != nil {
		return err
	} else if err := this.radio.SetSequencer(true); err != nil {
		return err
	} else if err := this.radio.SetBitrate(4800); err != nil {
		return err
	} else if err := this.radio.SetFreqCarrier(434300000); err != nil {
		return err
	} else if err := this.radio.SetFreqDeviation(30000); err != nil {
		return err
	} else if err := this.radio.SetAFCMode(sensors.RFM_AFCMODE_OFF); err != nil {
		return err
	} else if err := this.radio.SetAFCRoutine(sensors.RFM_AFCROUTINE_STANDARD); err != nil {
		return err
	} else if err := this.radio.SetLNA(sensors.RFM_LNA_IMPEDANCE_50, sensors.RFM_LNA_GAIN_AUTO); err != nil {
		return err
	} else if err := this.radio.SetRXFilter(sensors.RFM_RXBW_FREQUENCY_FSK_62P5, sensors.RFM_RXBW_CUTOFF_4); err != nil {
		return err
	} else if err := this.radio.SetDataMode(sensors.RFM_DATAMODE_PACKET); err != nil {
		return err
	} else if err := this.radio.SetPacketFormat(sensors.RFM_PACKET_FORMAT_VARIABLE); err != nil {
		return err
	} else if err := this.radio.SetPacketCoding(sensors.RFM_PACKET_CODING_MANCHESTER); err != nil {
		return err
	} else if err := this.radio.SetPacketFilter(sensors.RFM_PACKET_FILTER_NONE); err != nil {
		return err
	} else if err := this.radio.SetPacketCRC(sensors.RFM_PACKET_CRC_OFF); err != nil {
		return err
	} else if err := this.radio.SetPreambleSize(3); err != nil {
		return err
	} else if err := this.radio.SetPayloadSize(0x40); err != nil {
		return err
	} else if err := this.radio.SetSyncWord([]byte{0x2D, 0xD4}); err != nil {
		return err
	} else if err := this.radio.SetSyncTolerance(0); err != nil {
		return err
	} else if err := this.radio.SetNodeAddress(0x04); err != nil {
		return err
	} else if err := this.radio.SetBroadcastAddress(0xFF); err != nil {
		return err
	} else if err := this.radio.SetAESKey(nil); err != nil {
		return err
	} else if err := this.radio.SetFIFOThreshold(1); err != nil {
		return err
	}

	// Success
	return nil
}

func (this *mihome) setOOKMode() error {
	this.log.Debug2("<sensors.ener314rt>SetOOKMode{}")

	if err := this.radio.SetMode(sensors.RFM_MODE_STDBY); err != nil {
		return err
	} else if err := this.radio.SetModulation(sensors.RFM_MODULATION_OOK); err != nil {
		return err
	} else if err := this.radio.SetSequencer(true); err != nil {
		return err
	} else if err := this.radio.SetBitrate(4800); err != nil {
		return err
	} else if err := this.radio.SetFreqCarrier(433920000); err != nil {
		return err
	} else if err := this.radio.SetFreqDeviation(0); err != nil {
		return err
	} else if err := this.radio.SetAFCMode(sensors.RFM_AFCMODE_OFF); err != nil {
		return err
	} else if err := this.radio.SetDataMode(sensors.RFM_DATAMODE_PACKET); err != nil {
		return err
	} else if err := this.radio.SetPacketFormat(sensors.RFM_PACKET_FORMAT_VARIABLE); err != nil {
		return err
	} else if err := this.radio.SetPacketCoding(sensors.RFM_PACKET_CODING_NONE); err != nil {
		return err
	} else if err := this.radio.SetPacketFilter(sensors.RFM_PACKET_FILTER_NONE); err != nil {
		return err
	} else if err := this.radio.SetPacketCRC(sensors.RFM_PACKET_CRC_OFF); err != nil {
		return err
	} else if err := this.radio.SetPreambleSize(0); err != nil {
		return err
	} else if err := this.radio.SetPayloadSize(0); err != nil {
		return err
	} else if err := this.radio.SetSyncWord(nil); err != nil {
		return err
	} else if err := this.radio.SetAESKey(nil); err != nil {
		return err
	} else if err := this.radio.SetFIFOThreshold(1); err != nil {
		return err
	}

	// Success
	return nil
}

func (this *mihome) emit(payload []byte, protocols []sensors.Proto) error {
	this.log.Debug("<sensors.ener314rt>Emit{ payload=%v protocols=%v }", strings.ToUpper(hex.EncodeToString(payload)), protocols)

	// Record timestamp
	ts := time.Now()

	if len(payload) == 0 {
		// If payload is empty, then return unexpected response
		return sensors.ErrUnexpectedResponse
	}

	// Cycle through protocols until we get a message. We could use goroutines
	// to look at each protocol in parallel and wait until they all return...
	var last_err error
	for _, proto := range protocols {
		if msg, err := proto.Decode(payload, ts); err == nil {
			// We emit the message
			return this.emit_message(msg)
		} else {
			// Record the error returned
			last_err = err
		}
	}

	// Return the last error
	return last_err
}

func (this *mihome) emit_message(msg sensors.Message) error {
	this.Publisher.Emit(msg)
	return nil
}

func (this *mihome) protocol_for_name(name string) sensors.Proto {
	if proto, exists := this.proto_map[name]; exists == false {
		return nil
	} else {
		return proto
	}
}

func (this *mihome) protocols_for_mode(mode sensors.MiHomeMode) []sensors.Proto {
	if protos, exists := this.proto_mode[mode]; exists == false {
		return nil
	} else {
		return protos
	}
}
