/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package mihome

import (
	"fmt"
	"sync"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi/util/event"
	"github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type queue struct {
	log    gopi.Logger
	mihome sensors.MiHome
	queue  []*message

	// Lock queue
	sync.Mutex

	// Receive messages in the background
	event.Tasks
}

type message struct {
	product     sensors.MiHomeProduct
	sensor      uint32
	parameter   sensors.OTParameter
	temperature float64
	interval    time.Duration
	valve_state sensors.MiHomeValveState
	low_power   bool
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open the server
func (this *queue) Init(log gopi.Logger, config Service) error {
	log.Debug("<grpc.service.mihome.Queue>Init{}")

	if log == nil || config.MiHome == nil {
		return gopi.ErrBadParameter
	}
	this.log = log
	this.mihome = config.MiHome
	this.queue = make([]*message, 0)

	// Start background task which reports on all events (for debugging, device collection)
	this.Tasks.Start(this.EventTask)

	// Success
	return nil
}

func (this *queue) Destroy() error {
	this.log.Debug("<grpc.service.mihome.Queue>Destroy{}")

	// Stop background tasks
	if err := this.Tasks.Close(); err != nil {
		return err
	}

	// Release resources
	this.log = nil
	this.mihome = nil
	this.queue = nil

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *message) String() string {
	if this.parameter == sensors.OT_PARAM_TEMPERATURE {
		return fmt.Sprintf("<queue.message>{ product=%v sensor=0x%06X parameter=%v temperature=%v }", this.product, this.sensor, this.parameter, this.temperature)
	} else if this.parameter == sensors.OT_PARAM_REPORT_PERIOD {
		return fmt.Sprintf("<queue.message>{ product=%v sensor=0x%06X parameter=%v interval=%v }", this.product, this.sensor, this.parameter, this.interval)
	} else if this.parameter == sensors.OT_PARAM_VALVE_STATE {
		return fmt.Sprintf("<queue.message>{ product=%v sensor=0x%06X parameter=%v state=%v }", this.product, this.sensor, this.parameter, this.valve_state)
	} else if this.parameter == sensors.OT_PARAM_LOW_POWER {
		return fmt.Sprintf("<queue.message>{ product=%v sensor=0x%06X parameter=%v low_power=%v }", this.product, this.sensor, this.parameter, this.low_power)
	} else {
		return fmt.Sprintf("<queue.message>{ product=%v sensor=0x%06X parameter=%v }", this.product, this.sensor, this.parameter)
	}
}

///////////////////////////////////////////////////////////////////////////////
// QUEUE MESSAGES

func (this *queue) QueueDiagnostics(product sensors.MiHomeProduct, sensor uint32) error {
	this.log.Debug("<grpc.service.mihome.Queue>QueueDiagnostics{ product=%v sensor=0x%08X }", product, sensor)

	if this.Match(product, sensor, sensors.OT_PARAM_DIAGNOSTICS, false) != nil {
		// Ignore if there is an existing message in the queue
		return gopi.ErrNotModified
	} else {
		this.Append(&message{product, sensor, sensors.OT_PARAM_DIAGNOSTICS, 0, 0, 0, false})
	}

	// Return sucess
	return nil
}

func (this *queue) QueueIdentify(product sensors.MiHomeProduct, sensor uint32) error {
	this.log.Debug("<grpc.service.mihome.Queue>QueueIdentify{ product=%v sensor=0x%08X }", product, sensor)

	if this.Match(product, sensor, sensors.OT_PARAM_IDENTIFY, false) != nil {
		// Ignore if there is an existing message in the queue
		return gopi.ErrNotModified
	} else {
		this.Append(&message{product, sensor, sensors.OT_PARAM_IDENTIFY, 0, 0, 0, false})
	}

	// Return sucess
	return nil
}

func (this *queue) QueueExercise(product sensors.MiHomeProduct, sensor uint32) error {
	this.log.Debug("<grpc.service.mihome.Queue>QueueExercise{ product=%v sensor=0x%08X }", product, sensor)

	if this.Match(product, sensor, sensors.OT_PARAM_EXERCISE, false) != nil {
		// Ignore if there is an existing message in the queue
		return gopi.ErrNotModified
	} else {
		this.Append(&message{product, sensor, sensors.OT_PARAM_EXERCISE, 0, 0, 0, false})
	}

	// Return sucess
	return nil
}

func (this *queue) QueueBatteryLevel(product sensors.MiHomeProduct, sensor uint32) error {
	this.log.Debug("<grpc.service.mihome.Queue>QueueBatteryLevel{ product=%v sensor=0x%08X }", product, sensor)

	if this.Match(product, sensor, sensors.OT_PARAM_BATTERY_LEVEL, false) != nil {
		// Ignore if there is an existing message in the queue
		return gopi.ErrNotModified
	} else {
		this.Append(&message{product, sensor, sensors.OT_PARAM_BATTERY_LEVEL, 0, 0, 0, false})
	}

	// Return sucess
	return nil
}

func (this *queue) QueueTargetTemperature(product sensors.MiHomeProduct, sensor uint32, temperature float64) error {
	this.log.Debug("<grpc.service.mihome.Queue>QueueTargetTemperature{ product=%v sensor=0x%08X temperature=%v }", product, sensor, temperature)

	// Set the target temperature
	if temperature < 0.0 || temperature > 30.0 {
		return gopi.ErrBadParameter
	} else if reply := this.Match(product, sensor, sensors.OT_PARAM_TEMPERATURE, false); reply != nil {
		// Check to see if record needs updated
		if temperature == reply.temperature {
			return gopi.ErrNotModified
		}
		// Update existing queued record
		reply.temperature = temperature
	} else {
		this.Append(&message{product, sensor, sensors.OT_PARAM_TEMPERATURE, temperature, 0, 0, false})
	}

	// Return sucess
	return nil
}

func (this *queue) QueueReportInterval(product sensors.MiHomeProduct, sensor uint32, interval time.Duration) error {
	this.log.Debug("<grpc.service.mihome.Queue>QueueReportInterval{ product=%v sensor=0x%08X interval=%v }", product, sensor, interval)

	// Set the reporting interval
	interval = interval.Truncate(time.Second)
	if interval < 1*time.Second || interval > 3600*time.Second {
		return gopi.ErrBadParameter
	} else if reply := this.Match(product, sensor, sensors.OT_PARAM_REPORT_PERIOD, false); reply != nil {
		// Check to see if record needs updated
		if reply.interval == interval {
			return gopi.ErrNotModified
		}
		// Update existing queued record
		reply.interval = interval
	} else {
		this.Append(&message{product, sensor, sensors.OT_PARAM_REPORT_PERIOD, 0, interval, 0, false})
	}

	// Return sucess
	return nil
}

func (this *queue) QueueValveState(product sensors.MiHomeProduct, sensor uint32, state sensors.MiHomeValveState) error {
	this.log.Debug("<grpc.service.mihome.Queue>QueueValveState{ product=%v sensor=0x%08X state=%v }", product, sensor, state)

	// Set the valve_state
	if reply := this.Match(product, sensor, sensors.OT_PARAM_VALVE_STATE, false); reply != nil {
		// Check to see if record needs updated
		if reply.valve_state == state {
			return gopi.ErrNotModified
		}
		// Update existing queued record
		reply.valve_state = state
	} else {
		this.Append(&message{product, sensor, sensors.OT_PARAM_VALVE_STATE, 0, 0, state, false})
	}

	// Return sucess
	return nil
}

func (this *queue) QueueLowPowerMode(product sensors.MiHomeProduct, sensor uint32, low_power bool) error {
	this.log.Debug("<grpc.service.mihome.Queue>QueueLowPowerMode{ product=%v sensor=0x%08X low_power=%v }", product, sensor, low_power)

	// Set the low power mode
	if reply := this.Match(product, sensor, sensors.OT_PARAM_LOW_POWER, false); reply != nil {
		// Check to see if record needs updated
		if reply.low_power == low_power {
			return gopi.ErrNotModified
		}
		// Update existing queued record
		reply.low_power = low_power
	} else {
		this.Append(&message{product, sensor, sensors.OT_PARAM_VALVE_STATE, 0, 0, 0, low_power})
	}

	// Return sucess
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// BACKGROUND TASKS

func (this *queue) EventTask(start chan<- event.Signal, stop <-chan event.Signal) error {
	start <- gopi.DONE
	events := this.mihome.Subscribe()
FOR_LOOP:
	for {
		select {
		case evt := <-events:
			if evt_, ok := evt.(sensors.Message); ok == false {
				this.log.Warn("Ignoring: %v", evt)
			} else if err := this.HandleEvent(evt_); err != nil {
				this.log.Error("%v", err)
			}
		case <-stop:
			break FOR_LOOP
		}
	}
	this.mihome.Unsubscribe(events)

	// Success
	return nil
}

func (this *queue) HandleEvent(message sensors.Message) error {
	if message == nil {
		return gopi.ErrBadParameter
	} else if message_, ok := message.(sensors.OTMessage); ok == false {
		this.log.Warn("Ignoring: %v", message)
		return nil
	} else if reply := this.ReplyFor(message_); reply != nil {
		if err := this.SendReply(reply); err != nil {
			return err
		}
	}

	// Return success
	return nil
}

func (this *queue) ReplyFor(message sensors.OTMessage) *message {
	return this.Match(sensors.MiHomeProduct(message.Product()), message.Sensor(), sensors.OT_PARAM_NONE, true)
}

// Match will return a queued message for product, sensor. Where parameter is OT_PARAM_NONE, the
// first message is returned with any parameter. If a message is returned, it is removed from
// the queue
func (this *queue) Match(product sensors.MiHomeProduct, sensor uint32, parameter sensors.OTParameter, remove bool) *message {
	this.Lock()
	defer this.Unlock()

	// Find a message in the queue which matches
	for pos, reply := range this.queue {
		if reply.product != product {
			continue
		}
		if reply.sensor != sensor {
			continue
		}
		if parameter == sensors.OT_PARAM_NONE || reply.parameter == parameter {
			if remove {
				this.queue = append(this.queue[:pos], this.queue[pos+1:]...)
			}
			return reply
		}
	}
	// Not found, so return nil
	return nil
}

func (this *queue) Append(message *message) {
	this.Lock()
	defer this.Unlock()
	this.queue = append(this.queue, message)
}

func (this *queue) SendReply(message *message) error {
	this.log.Debug("<grpc.service.mihome.Queue>SendReply{ message=%v }", message)
	switch message.parameter {
	case sensors.OT_PARAM_DIAGNOSTICS:
		return this.mihome.RequestDiagnostics(message.product, message.sensor)
	case sensors.OT_PARAM_IDENTIFY:
		return this.mihome.RequestIdentify(message.product, message.sensor)
	case sensors.OT_PARAM_EXERCISE:
		return this.mihome.RequestExercise(message.product, message.sensor)
	case sensors.OT_PARAM_BATTERY_LEVEL:
		return this.mihome.RequestBatteryLevel(message.product, message.sensor)
	case sensors.OT_PARAM_TEMPERATURE:
		return this.mihome.RequestTargetTemperature(message.product, message.sensor, message.temperature)
	case sensors.OT_PARAM_REPORT_PERIOD:
		return this.mihome.RequestReportInterval(message.product, message.sensor, message.interval)
	case sensors.OT_PARAM_VALVE_STATE:
		return this.mihome.RequestValveState(message.product, message.sensor, message.valve_state)
	case sensors.OT_PARAM_LOW_POWER:
		return this.mihome.RequestLowPowerMode(message.product, message.sensor, message.low_power)
	default:
		return fmt.Errorf("Reply unsent: %v", message)
	}
}
