/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"context"
	"fmt"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	event "github.com/djthorpe/gopi/util/event"
	sensors "github.com/djthorpe/sensors"
)

////////////////////////////////////////////////////////////////////////////////

type Runner struct {
	app     *gopi.AppInstance
	stubs   []sensors.MiHomeClient
	db      sensors.Database
	cancels []context.CancelFunc
	errors  chan error

	event.Tasks
	event.Merger
	sync.WaitGroup
}

func NewRunner(app *gopi.AppInstance) *Runner {
	this := new(Runner)
	if app == nil {
		return nil
	}

	this.app = app
	this.stubs = make([]sensors.MiHomeClient, 0)
	this.cancels = make([]context.CancelFunc, 0)
	this.errors = make(chan error)
	if db, ok := this.app.ModuleInstance("sensordb").(sensors.Database); ok {
		this.db = db
	} else {
		this.app.Logger.Fatal("Missing or invalid sensor database")
		return nil
	}

	// Task to receive messages
	this.Tasks.Start(this.EventTask)

	// Success
	return this
}

func (this *Runner) AddStub(stub sensors.MiHomeClient) error {
	this.stubs = append(this.stubs, stub)

	// Create background task to stream messages
	ctx, cancel := context.WithCancel(context.Background())
	this.cancels = append(this.cancels, cancel)
	go func() {
		this.WaitGroup.Add(1)
		this.Merger.Merge(stub)
		if err := stub.StreamMessages(ctx); err != nil && err != context.Canceled {
			this.errors <- err
		}
		this.Merger.Unmerge(stub)
		this.WaitGroup.Done()
	}()

	// return success
	return nil
}

func (this *Runner) Close() error {
	// Call cancels
	for _, cancel := range this.cancels {
		cancel()
	}

	// Wait until all streams are completed
	this.WaitGroup.Wait()

	// Stop tasks
	if err := this.Tasks.Close(); err != nil {
		return err
	}

	// Release resources
	close(this.errors)
	this.cancels = nil
	this.stubs = nil

	// return success
	return nil
}

func (this *Runner) EventTask(start chan<- event.Signal, stop <-chan event.Signal) error {
	start <- gopi.DONE
	events := this.Merger.Subscribe()
FOR_LOOP:
	for {
		select {
		case evt := <-events:
			if message, ok := evt.(sensors.Message); ok == false {
				fmt.Println("Unhandled message:", evt)
			} else if sensor, err := this.db.Register(message); err != nil {
				fmt.Println("Register Error:", err)
			} else if err := this.db.Write(sensor, message); err != nil {
				fmt.Println("Write Error:", err)
			}
		case err := <-this.errors:
			fmt.Println("Error:", err)
		case <-stop:
			break FOR_LOOP
		}
	}

	// Unsubscribe, return success
	this.Merger.Unsubscribe(events)
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func Run(app *gopi.AppInstance, client sensors.MiHomeClient) error {
	fmt.Println("Connected to:", client)
	if runner := NewRunner(app); runner == nil {
		return gopi.ErrAppError
	} else {
		// Add a stub to receive messages
		runner.AddStub(client)
		// Wait for CTRL+C signal, then stop
		app.WaitForSignal()
		return runner.Close()
	}
}
