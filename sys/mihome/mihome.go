/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package mihome

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	event "github.com/djthorpe/gopi/util/event"
	tasks "github.com/djthorpe/gopi/util/tasks"
	sensors "github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type MiHome struct {
	Radio      sensors.ENER314RT
	Mode       sensors.MiHomeMode
	Metrics    gopi.Metrics
	Repeat     uint    // Number of times to repeat messages by default
	TempOffset float32 // Temperature Offset
}

type mihome struct {
	log        gopi.Logger
	radio      sensors.ENER314RT
	repeat     uint
	tempoffset float32
	mode       sensors.MiHomeMode
	metrics    gopi.Metrics
	cancel     context.CancelFunc
	err        chan error
	payload    chan []byte

	Protocols
	event.Publisher
	tasks.Tasks
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// Default number of times to repeat command
	REPEAT_DEFAULT = 3
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open the server
func (config MiHome) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sensors.mihome>Open{ mode=%v radio=%v metrics=%v repeat=%v tempoffset=%v }", config.Mode, config.Radio, config.Metrics, config.Repeat, config.TempOffset)

	// Check for bad input parameters
	if config.Repeat == 0 {
		config.Repeat = REPEAT_DEFAULT
	}
	if config.Radio == nil {
		return nil, gopi.ErrBadParameter
	}

	this := new(mihome)
	this.log = log
	this.radio = config.Radio
	this.metrics = config.Metrics
	this.mode = config.Mode
	this.err = make(chan error)
	this.payload = make(chan []byte)
	this.repeat = config.Repeat
	this.tempoffset = config.TempOffset

	// Start receiving and recording device temperature
	this.Tasks.Start(this.receive)

	// Initiate receiving mode in background
	if err := this.rx_mode(true); err != nil {
		return nil, err
	}

	// Success
	return this, nil
}

func (this *mihome) Close() error {
	this.log.Debug("<sensors.mihome>Close{ mode=%v }", this.mode)

	// Cancel receive mode in foreground
	if err := this.rx_mode(false); err != nil {
		return err
	}

	// Close publisher
	this.Publisher.Close()

	// End background tasks
	if err := this.Tasks.Close(); err != nil {
		return err
	}

	// Close protocol map and publisher
	this.Protocols.Close()

	// Release resources
	this.radio = nil
	this.metrics = nil

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *mihome) String() string {
	return fmt.Sprintf("<sensors.mihome>{ mode=%v }", this.mode)
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - RESET

func (this *mihome) Reset() error {
	this.log.Debug2("<sensors.mihome>Reset{}")

	if err := this.rx_mode(false); err != nil {
		return err
	} else if err := this.radio.ResetRadio(); err != nil {
		return err
	} else if err := this.rx_mode(true); err != nil {
		return err
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - MEASURE TEMPERATURE

func (this *mihome) MeasureTemperature() (float32, error) {
	this.log.Debug2("<sensors.mihome>MeasureTemperature{ offset=%v }", this.tempoffset)

	if err := this.rx_mode(false); err != nil {
		return 0, err
	} else if celcius, err := this.radio.MeasureTemperature(this.tempoffset); err != nil {
		return 0, err
	} else if err := this.rx_mode(true); err != nil {
		return 0, err
	} else {
		return celcius, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - SWTICH

func (this *mihome) RequestSwitchOn(product sensors.MiHomeProduct, sensor uint32) error {
	return this.RequestSwitchState(product, sensor, true)
}

func (this *mihome) RequestSwitchOff(product sensors.MiHomeProduct, sensor uint32) error {
	return this.RequestSwitchState(product, sensor, false)
}

func (this *mihome) RequestSwitchState(product sensors.MiHomeProduct, sensor uint32, state bool) error {
	this.log.Debug2("<sensors.mihome>RequestSwitchState{ product=%v sensor=0x%05X state=%v }", product, sensor, state)

	if mode := product.Mode(); mode == sensors.MIHOME_MODE_NONE {
		// Invalid mode for product
		return gopi.ErrBadParameter
	} else if protos := this.ProtosByMode(mode); len(protos) == 0 {
		// Invalid mode for product
		return gopi.ErrBadParameter
	} else if proto, ok := protos[0].(sensors.OOKProto); ok && proto != nil {
		// OOK Protocol
		if message, err := proto.New(sensor, product.Socket(), state, nil); err != nil {
			return err
		} else if err := this.tx_mode(proto, message); err != nil {
			return err
		} else {
			return nil
		}
	} else if proto, ok := protos[0].(sensors.OTProto); ok && proto != nil {
		// FSK (OpenThings) Protocol
		if message, err := proto.New(sensors.OT_MANUFACTURER_ENERGENIE, uint8(product), sensor); err != nil {
			return err
		} else if state, err := proto.NewBool(sensors.OT_PARAM_SWITCH_STATE, state, true); err != nil {
			return err
		} else if err := this.tx_mode(proto, message.Append(state)); err != nil {
			return err
		} else {
			return nil
		}
	}

	// Unimplemented for this protocol/mode combination
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - ETRV

func (this *mihome) RequestIdentify(product sensors.MiHomeProduct, sensor uint32) error {
	this.log.Debug2("<sensors.mihome>RequestIdentify{ product=%v sensor=0x%08X }", product, sensor)

	// We only support this with the openthings protocol in monitor mode
	if proto_ := this.ProtoByName("openthings"); proto_ == nil {
		return gopi.ErrBadParameter
	} else if proto, ok := proto_.(sensors.OTProto); ok == false || proto == nil {
		return gopi.ErrBadParameter
	} else if msg, err := proto.New(sensors.OT_MANUFACTURER_ENERGENIE, uint8(product), sensor); err != nil {
		return err
	} else if record, err := proto.NewNull(sensors.OT_PARAM_IDENTIFY, true); err != nil {
		return err
	} else if err := this.tx_mode(proto, msg.Append(record)); err != nil {
		return err
	}

	// Success
	return nil
}

func (this *mihome) RequestDiagnostics(product sensors.MiHomeProduct, sensor uint32) error {
	this.log.Debug2("<sensors.mihome>RequestDiagnostics{ product=%v sensor=0x%08X }", product, sensor)

	// We only support this with the openthings protocol in monitor mode
	if proto_ := this.ProtoByName("openthings"); proto_ == nil {
		return gopi.ErrBadParameter
	} else if proto, ok := proto_.(sensors.OTProto); ok == false || proto == nil {
		return gopi.ErrBadParameter
	} else if msg, err := proto.New(sensors.OT_MANUFACTURER_ENERGENIE, uint8(product), sensor); err != nil {
		return err
	} else if record, err := proto.NewNull(sensors.OT_PARAM_DIAGNOSTICS, true); err != nil {
		return err
	} else if err := this.tx_mode(proto, msg.Append(record)); err != nil {
		return err
	}

	// Success
	return nil
}
func (this *mihome) RequestExercise(product sensors.MiHomeProduct, sensor uint32) error {
	this.log.Debug2("<sensors.mihome>RequestExercise{ product=%v sensor=0x%08X }", product, sensor)

	// We only support this with the openthings protocol in monitor mode
	if proto_ := this.ProtoByName("openthings"); proto_ == nil {
		return gopi.ErrBadParameter
	} else if proto, ok := proto_.(sensors.OTProto); ok == false || proto == nil {
		return gopi.ErrBadParameter
	} else if msg, err := proto.New(sensors.OT_MANUFACTURER_ENERGENIE, uint8(product), sensor); err != nil {
		return err
	} else if record, err := proto.NewNull(sensors.OT_PARAM_EXERCISE, true); err != nil {
		return err
	} else if err := this.tx_mode(proto, msg.Append(record)); err != nil {
		return err
	}

	// Success
	return nil
}

func (this *mihome) RequestBatteryLevel(product sensors.MiHomeProduct, sensor uint32) error {
	this.log.Debug2("<sensors.mihome>RequestBatteryLevel{ product=%v sensor=0x%08X }", product, sensor)

	// We only support this with the openthings protocol in monitor mode
	if proto_ := this.ProtoByName("openthings"); proto_ == nil {
		return gopi.ErrBadParameter
	} else if proto, ok := proto_.(sensors.OTProto); ok == false || proto == nil {
		return gopi.ErrBadParameter
	} else if msg, err := proto.New(sensors.OT_MANUFACTURER_ENERGENIE, uint8(product), sensor); err != nil {
		return err
	} else if record, err := proto.NewNull(sensors.OT_PARAM_BATTERY_LEVEL, true); err != nil {
		return err
	} else if err := this.tx_mode(proto, msg.Append(record)); err != nil {
		return err
	}

	// Success
	return nil
}

func (this *mihome) SendJoin(product sensors.MiHomeProduct, sensor uint32) error {
	this.log.Debug2("<sensors.mihome>SendJoin{ product=%v sensor=0x%08X }", product, sensor)

	// We only support this with the openthings protocol in monitor mode
	if proto_ := this.ProtoByName("openthings"); proto_ == nil {
		return gopi.ErrBadParameter
	} else if proto, ok := proto_.(sensors.OTProto); ok == false || proto == nil {
		return gopi.ErrBadParameter
	} else if msg, err := proto.New(sensors.OT_MANUFACTURER_ENERGENIE, uint8(product), sensor); err != nil {
		return err
	} else if record, err := proto.NewNull(sensors.OT_PARAM_JOIN, false); err != nil {
		return err
	} else if err := this.tx_mode(proto, msg.Append(record)); err != nil {
		return err
	}

	// Success
	return nil
}

func (this *mihome) RequestTargetTemperature(product sensors.MiHomeProduct, sensor uint32, celcius float64) error {
	this.log.Debug2("<sensors.mihome>RequestTargetTemperature{ product=%v sensor=0x%08X celcius=%v }", product, sensor, celcius)

	// Return bad parameter if calcius is less than 0 or greater than 200
	if celcius < 0 || celcius > 200 {
		return gopi.ErrBadParameter
	}

	// We only support this with the openthings protocol in monitor mode
	if proto_ := this.ProtoByName("openthings"); proto_ == nil {
		return gopi.ErrBadParameter
	} else if proto, ok := proto_.(sensors.OTProto); ok == false || proto == nil {
		return gopi.ErrBadParameter
	} else if msg, err := proto.New(sensors.OT_MANUFACTURER_ENERGENIE, uint8(product), sensor); err != nil {
		return err
	} else if record, err := proto.NewFloat(sensors.OT_PARAM_TEMPERATURE, sensors.OT_DATATYPE_DEC_8, celcius, true); err != nil {
		return err
	} else if err := this.tx_mode(proto, msg.Append(record)); err != nil {
		return err
	}

	// Success
	return nil
}

func (this *mihome) RequestReportInterval(product sensors.MiHomeProduct, sensor uint32, interval time.Duration) error {
	this.log.Debug2("<sensors.mihome>RequestReportInterval{ product=%v sensor=0x%08X interval=%v }", product, sensor, interval)

	// Return bad parameter if interval is less than 0 seconds or greater than uint16
	seconds := uint64(interval.Seconds())
	if seconds < 0 || seconds > math.MaxUint16 {
		return gopi.ErrBadParameter
	}

	// We only support this with the openthings protocol in monitor mode
	if proto_ := this.ProtoByName("openthings"); proto_ == nil {
		return gopi.ErrBadParameter
	} else if proto, ok := proto_.(sensors.OTProto); ok == false || proto == nil {
		return gopi.ErrBadParameter
	} else if msg, err := proto.New(sensors.OT_MANUFACTURER_ENERGENIE, uint8(product), sensor); err != nil {
		return err
	} else if record, err := proto.NewUint16(sensors.OT_PARAM_REPORT_PERIOD, uint16(seconds), true); err != nil {
		return err
	} else if err := this.tx_mode(proto, msg.Append(record)); err != nil {
		return err
	}

	// Success
	return nil
}

func (this *mihome) RequestValveState(product sensors.MiHomeProduct, sensor uint32, state sensors.MiHomeValveState) error {
	this.log.Debug2("<sensors.mihome>RequestValveState{ product=%v sensor=0x%08X state=%v }", product, sensor, state)

	// We only support this with the openthings protocol in monitor mode
	if proto_ := this.ProtoByName("openthings"); proto_ == nil {
		return gopi.ErrBadParameter
	} else if proto, ok := proto_.(sensors.OTProto); ok == false || proto == nil {
		return gopi.ErrBadParameter
	} else if msg, err := proto.New(sensors.OT_MANUFACTURER_ENERGENIE, uint8(product), sensor); err != nil {
		return err
	} else if record, err := proto.NewUint8(sensors.OT_PARAM_VALVE_STATE, uint8(state), true); err != nil {
		return err
	} else if err := this.tx_mode(proto, msg.Append(record)); err != nil {
		return err
	}

	// Success
	return nil
}

func (this *mihome) RequestLowPowerMode(product sensors.MiHomeProduct, sensor uint32, mode bool) error {
	this.log.Debug2("<sensors.mihome>RequestLowPowerMode{ product=%v sensor=0x%08X mode=%v }", product, sensor, mode)

	// We only support this with the openthings protocol in monitor mode
	if proto_ := this.ProtoByName("openthings"); proto_ == nil {
		return gopi.ErrBadParameter
	} else if proto, ok := proto_.(sensors.OTProto); ok == false || proto == nil {
		return gopi.ErrBadParameter
	} else if msg, err := proto.New(sensors.OT_MANUFACTURER_ENERGENIE, uint8(product), sensor); err != nil {
		return err
	} else if record, err := proto.NewBool(sensors.OT_PARAM_JOIN, mode, true); err != nil {
		return err
	} else if err := this.tx_mode(proto, msg.Append(record)); err != nil {
		return err
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - TX DATA

func (this *mihome) tx_mode(proto sensors.Proto, message sensors.Message) error {
	this.log.Debug("<sensors.mihome>TXMode{ proto=%v messgage=%v }", proto, message)

	// Encode the message, switch off RX mode, send then return to RX mode
	if encoded := proto.Encode(message); len(encoded) == 0 {
		return sensors.ErrMessageCorruption
	} else if err := this.rx_mode(false); err != nil {
		return err
	} else if err := this.radio.Send(encoded, this.repeat, proto.Mode()); err != nil {
		return err
	} else if err := this.rx_mode(true); err != nil {
		return err
	} else {
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - START AND STOP RX MODE

func (this *mihome) rx_cancel() {
	this.log.Debug("<sensors.mihome>RXCancel{}")
	if this.cancel != nil {
		this.cancel()
		this.cancel = nil
	}
}

func (this *mihome) rx_mode(state bool) error {
	this.log.Debug("<sensors.mihome>RXMode{ state=%v }", state)

	// Lock critical section
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if state == false {
		// If state is OFF then cancel in foreground
		if this.cancel != nil {
			this.rx_cancel()
		}
	} else if this.mode == sensors.MIHOME_MODE_NONE {
		// Do nothing with RX mode here
	} else if state && this.cancel == nil {
		// If state is ON, then run it in the background until we receive an error or nil
		ctx, cancel := context.WithCancel(context.Background())
		this.cancel = cancel
		go func(ctx context.Context) {
			err := this.radio.Receive(ctx, this.mode, this.payload)
			this.err <- err
		}(ctx)
	} else {
		// Assume RX is already running
		//this.log.Warn("<sensors.mihome>RXMode: Invalid state, state=%v cancel=%v", state, this.cancel)
		return nil //gopi.ErrAppError
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RECEIVE AND DECODE DATA

func (this *mihome) receive(start chan<- struct{}, stop <-chan struct{}) error {
	this.log.Debug("<sensors.mihome>receive: Started")
	start <- gopi.DONE

	// Obtain the protocols we will use to decode the message
	var protocols []sensors.Proto

FOR_LOOP:
	for {
		select {
		case data := <-this.payload:
			if protocols == nil {
				protocols = this.ProtosByMode(this.mode)
				if this.mode != sensors.MIHOME_MODE_NONE {
					protocols = append(protocols, this.ProtosByMode(sensors.MIHOME_MODE_NONE)...)
				}
			} else if len(protocols) == 0 {
				this.log.Warn("<sensors.mihome>Receive: No protocols found for mode %v", this.mode)
			} else if err := this.decode(data, protocols); err != nil {
				this.log.Warn("<sensors.mihome>Receive: %v", err)
			}
		case err := <-this.err:
			if err != context.Canceled && err != sensors.ErrDeviceTimeout {
				this.log.Warn("<sensors.mihome>Receive: %v", err)
				// Perform a reset after a short interval, if no cancel
				time.Sleep(time.Second)
				this.log.Warn("<sensors.mihome>Resetting device")
				if err := this.Reset(); err != nil {
					this.log.Warn("<sensors.mihome>Receive: %v", err)
				}
			} else {
				this.log.Debug("<sensors.mihome>Receive: %v", err)
			}
		case <-stop:
			this.log.Debug("<sensors.mihome>Receive: Ended")
			break FOR_LOOP
		}
	}
	return nil
}

func (this *mihome) decode(payload []byte, protos []sensors.Proto) error {
	// Check arguments
	if len(payload) == 0 || len(protos) == 0 {
		return gopi.ErrBadParameter
	}

	// Decode through protocols until we find one which decodes the payload
	var last_err error
	for _, proto := range protos {
		if msg, err := proto.Decode(payload, time.Now()); err == nil {
			this.Emit(msg)
			return nil
		} else {
			// Record the error returned
			last_err = err
		}
	}

	// We were not able to decode the message, return the last error
	// recorded
	return last_err
}
