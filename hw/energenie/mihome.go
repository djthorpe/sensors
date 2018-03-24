/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package energenie

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"
	"github.com/olekukonko/tablewriter"
)

////////////////////////////////////////////////////////////////////////////////
// STRUCTS

// Configuration
type MiHome struct {
	GPIO     gopi.GPIO     // GPIO interface
	Radio    sensors.RFM69 // Radio interface
	PinReset gopi.GPIOPin  // Reset pin
	PinLED1  gopi.GPIOPin  // LED1 (Green, Rx) pin
	PinLED2  gopi.GPIOPin  // LED2 (Red, Tx) pin
	CID      string        // OOK device address
	Repeat   uint          // Number of times to repeat messages by default
}

// mihome driver
type mihome struct {
	log    gopi.Logger
	gpio   gopi.GPIO
	radio  sensors.RFM69
	reset  gopi.GPIOPin
	cid    []byte // 10 bytes for the OOK address
	repeat uint
	led1   gopi.GPIOPin
	led2   gopi.GPIOPin
	ledrx  gopi.GPIOPin
	ledtx  gopi.GPIOPin
	mode   sensors.MiHomeMode
}

type LED uint
type Command byte

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS, GLOBAL VARIABLES

const (
	// Default Control ID
	CID_DEFAULT = "6C6C6"
	// Default number of times to repeat command
	REPEAT_DEFAULT = 8
)

var (
	// OOK Preamble sent before each command
	OOK_PREAMBLE = []byte{0x80, 0x00, 0x00, 0x00}
)

const (
	OOK_ZERO byte = 0x08
	OOK_ONE  byte = 0x0E
)

const (
	OOK_NONE    Command = 0x00
	OOK_ON_ALL  Command = 0x0D
	OOK_OFF_ALL Command = 0x0C
	OOK_ON_1    Command = 0x0F
	OOK_OFF_1   Command = 0x0E
	OOK_ON_2    Command = 0x07
	OOK_OFF_2   Command = 0x06
	OOK_ON_3    Command = 0x0B
	OOK_OFF_3   Command = 0x0A
	OOK_ON_4    Command = 0x03
	OOK_OFF_4   Command = 0x02
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
		config.CID = CID_DEFAULT
	}
	if config.Repeat == 0 {
		config.Repeat = REPEAT_DEFAULT
	}
	log.Debug2("<sensors.energenie.MiHome>Open{ reset=%v led1=%v led2=%v cid=\"%v\" repeat=%v }", config.PinReset, config.PinLED1, config.PinLED2, config.CID, config.Repeat)

	if config.GPIO == nil || config.Radio == nil {
		// Fail when either GPIO is nil or Radio is nil
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

	// Set the default Control ID for legacy OOK devices
	if cid, err := decodeHexString(config.CID); err != nil {
		return nil, err
	} else {
		this.cid = cid
	}

	// Set number of times to repeat TX by default
	this.repeat = config.Repeat

	// Set mode to undefined
	this.mode = sensors.MIHOME_MODE_NONE

	// Return success
	return this, nil
}

func (this *mihome) Close() error {
	this.log.Debug2("<sensors.energenie.MiHome>Close{ cid=0x%v }", strings.ToUpper(hex.EncodeToString(this.cid)))

	this.gpio = nil
	this.radio = nil
	this.cid = nil

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *mihome) String() string {
	return fmt.Sprintf("<sensors.energenie.MiHome>{ gpio=%v radio=%v reset=%v led1=%v led2=%v ledrx=%v ledtx=%v cid=0x%v mode=%v }", this.gpio, this.radio, this.reset, this.led1, this.led2, this.ledrx, this.ledtx, strings.ToUpper(hex.EncodeToString(this.cid)), this.mode)
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
	// We only support the MONITOR mode (FSK) for the moment
	if mode != sensors.MIHOME_MODE_MONITOR {
		return gopi.ErrNotImplemented
	}

	// Switch into FSK mode
	if this.radio.Modulation() != sensors.RFM_MODULATION_FSK || this.mode != sensors.MIHOME_MODE_MONITOR {
		if err := this.setFSKMode(); err != nil {
			return err
		} else {
			this.mode = sensors.MIHOME_MODE_MONITOR
		}
	}

	// Switch into RX mode
	if this.radio.Mode() != sensors.RFM_MODE_RX {
		if err := this.radio.SetMode(sensors.RFM_MODE_RX); err != nil {
			return err
		}
	} else if err := this.radio.ClearFIFO(); err != nil {
		return err
	}

	// Repeatedly read until context is done
FOR_LOOP:
	for {
		select {
		case <-ctx.Done():
			break FOR_LOOP
		default:
			if data, crc_ok, err := this.radio.ReadPayload(ctx); err != nil {
				return err
			} else if data == nil {
				continue
			} else {
				// Output register information
				table := tablewriter.NewWriter(os.Stdout)

				table.SetHeader([]string{"Payload", "Value"})
				table.Append([]string{"payload", fmt.Sprintf("%v", strings.ToUpper(hex.EncodeToString(data)))})
				table.Append([]string{"crc_ok", fmt.Sprintf("%v", crc_ok)})

				table.Render()
			}
		}
	}

	// Success
	return nil
}

// Send Command TX in Control Mode (aka Legacy mode, or using OOK
func (this *mihome) SendControl(cid []byte, cmd Command, repeat uint) error {
	this.log.Debug("<sensors.energenie.MiHome.SendControl{ cid=%v cmd=%v repeat=%v }", strings.ToUpper(hex.EncodeToString(cid)), cmd, repeat)

	if repeat == 0 || cid == nil {
		return gopi.ErrBadParameter
	} else if payload, err := encodeCommandPayload(cid, cmd); err != nil {
		return err
	} else if this.radio.Modulation() != sensors.RFM_MODULATION_OOK || this.mode != sensors.MIHOME_MODE_CONTROL {
		if err := this.setOOKMode(); err != nil {
			return err
		} else {
			this.mode = sensors.MIHOME_MODE_CONTROL
		}
	} else if err := this.radio.SetMode(sensors.RFM_MODE_TX); err != nil {
		return err
	} else if err := this.radio.SetSequencer(true); err != nil {
		return err
	} else if err := this.radio.WritePayload(payload, repeat); err != nil {
		return err
	}
	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - ENER314

// Satisfies the ENER314 interface to switch sockets on
func (this *mihome) On(sockets ...uint) error {
	// Switch on TX LED and switch off later
	this.SetLED(LED_TX, gopi.GPIO_HIGH)
	defer this.SetLED(LED_TX, gopi.GPIO_LOW)

	if len(sockets) == 0 {
		// all on
		return this.SendControl(this.cid, OOK_ON_ALL, this.repeat)
	} else {
		for _, socket := range sockets {
			if cmd, err := onCommandForSocket(socket); err != nil {
				return err
			} else if err := this.SendControl(this.cid, cmd, this.repeat); err != nil {
				return err
			}
		}
	}

	// Success
	return nil
}

// Satisfies the ENER314 interface to switch sockets off
func (this *mihome) Off(sockets ...uint) error {
	// Switch on TX LED and switch off later
	this.SetLED(LED_TX, gopi.GPIO_HIGH)
	defer this.SetLED(LED_TX, gopi.GPIO_LOW)

	if len(sockets) == 0 {
		// all off
		return this.SendControl(this.cid, OOK_OFF_ALL, this.repeat)
	} else {
		for _, socket := range sockets {
			if cmd, err := offCommandForSocket(socket); err != nil {
				return err
			} else if err := this.SendControl(this.cid, cmd, this.repeat); err != nil {
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
	} else if err := this.radio.SetAFCMode(sensors.RFM_AFCMODE_ON); err != nil {
		return err
	} else if err := this.radio.SetAFCRoutine(sensors.RFM_AFCROUTINE_STANDARD); err != nil {
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
	} else if err := this.radio.SetPayloadSize(66); err != nil {
		return err
	} else if err := this.radio.SetSyncWord([]byte{0xD4, 0x2D}); err != nil {
		return err
	} else if err := this.radio.SetSyncTolerance(3); err != nil {
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

// Convert hex string into bytes
func decodeHexString(value string) ([]byte, error) {
	// Pad with zeros
	for len(value)%2 != 0 {
		value = "0" + value
	}
	// Return hex
	return hex.DecodeString(value)
}

func onCommandForSocket(socket uint) (Command, error) {
	switch socket {
	case 1:
		return OOK_ON_1, nil
	case 2:
		return OOK_ON_2, nil
	case 3:
		return OOK_ON_3, nil
	case 4:
		return OOK_ON_4, nil
	default:
		return OOK_NONE, gopi.ErrBadParameter
	}
}

func offCommandForSocket(socket uint) (Command, error) {
	switch socket {
	case 1:
		return OOK_OFF_1, nil
	case 2:
		return OOK_OFF_2, nil
	case 3:
		return OOK_OFF_3, nil
	case 4:
		return OOK_OFF_4, nil
	default:
		return OOK_NONE, gopi.ErrBadParameter
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

func encodeByteArray(array []byte) []byte {
	// Return nil in the case of empty array
	if array == nil {
		return nil
	}
	// A byte is encoded as 4 bytes
	encoded := make([]byte, 0, 4*len(array))
	for i := range array {
		b := encodeByte(array[i])
		encoded = append(encoded, b...)
	}
	return encoded
}

func encodeCommandPayload(cid []byte, cmd Command) ([]byte, error) {
	if encoded_cmd := encodeByte(byte(cmd)); len(encoded_cmd) != 4 {
		return nil, gopi.ErrAppError
	} else if encoded_cid := encodeByteArray(cid); len(encoded_cid) != 12 {
		return nil, gopi.ErrAppError
	} else {
		// The payload is 16 bytes (preamble 4 bytes, address 10 bytes, command 2 bytes)
		// we chop off some bytes from both the command and the address
		payload := make([]byte, 0, 16)
		payload = append(payload, OOK_PREAMBLE...)
		payload = append(payload, encoded_cid[2:]...)
		payload = append(payload, encoded_cmd[2:]...)
		return payload, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (c Command) String() string {
	switch c {
	case OOK_ON_ALL:
		return "OOK_ON_ALL"
	case OOK_OFF_ALL:
		return "OOK_OFF_ALL"
	case OOK_ON_1:
		return "OOK_ON_1"
	case OOK_ON_2:
		return "OOK_ON_2"
	case OOK_OFF_2:
		return "OOK_OFF_2"
	case OOK_ON_3:
		return "OOK_ON_3"
	case OOK_OFF_3:
		return "OOK_OFF_3"
	case OOK_ON_4:
		return "OOK_ON_4"
	case OOK_OFF_4:
		return "OOK_OFF_4"
	default:
		return "[?? Invalid Command value]"
	}
}
