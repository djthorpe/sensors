
# Energenie ENER314

You can purchase two boards from Energenie which fit into the Raspberry Pi. The two separate boards
are:

  1. The Pi-Mote Remote Control Board [link](https://energenie4u.co.uk/catalogue/product/ENER314)
     which is a simple OOK transmitter;
  2. The two-way Pi-Mote [link](https://energenie4u.co.uk/catalogue/product/ENER314-RT) which
     is actually a 433Mhz RFM69 board.

Here are the datasheets for these boards:

  * The Pi-Mote ENER314 [link](ENER314-UM.pdf)
  * The Two way Pi-Mote ENER314-RT [link](ENER314-RT.pdf)

Both of these are implemented. The former communicates simply through GPIO whereas the latter
communicates through SPI and uses GPIO to set the receive and transmit LED on the board and
resets the RFM69 module. In order to create a driver for each, you can import the following 
modules anonymously into your application:

| Import                                    | Module Name        | Interface  |
| ----------------------------------------- | ------------------ | ---------- |
| github.com/djthorpe/sensors/sys/ener314   | sensors/ener314    | GPIO       |
| github.com/djthorpe/sensors/sys/ener314rt | sensors/ener314rt  | SPI & GPIO |

Here's an example for both:

```
package main

import (
  "os"
  "fmt"

  // Frameworks
  "github.com/djthorpe/gopi"
  "github.com/djthorpe/sensors"

  // Modules
  _ "github.com/djthorpe/sensors/sys/ener314"
  _ "github.com/djthorpe/sensors/sys/ener314rt"
  _ "github.com/djthorpe/sensors/sys/rfm69"
  _ "github.com/djthorpe/gopi-hw/sys/spi"
  _ "github.com/djthorpe/gopi-hw/sys/gpio"
)

const (
  ENER314 = "sensors/ener314"
  ENER314RT = "sensors/ener314rt"
)

func Main(app *gopi.AppInstance, done chan<- struct{}) error {
  ener314 := app.ModuleInstance(ENER314).(sensors.ENER314)
  ener314rt := app.ModuleInstance(ENER314RT).(sensors.ENER314RT)
  fmt.Println("ENER314=",ener314)
  fmt.Println("ENER314RT=",ener314rt)
  return nil
}

func main() {
  config := gopi.NewAppConfig(ENER314,ENER314RT)
  os.Exit(gopi.CommandLineTool(config, Main))
}
```

## The Pi-Mote Remote Control (ENER314)

The ENER314 is a simple OOK transmitter which communicates with
some of the "legacy" Energenie devices such as switches. It's a transmit-only
device and is paired with a switch by holding down the button on
the switch for an extended time, then transmitting either an on or off
signal.

The board is bought as-is and requires no additional set-up but
here are the pinouts:

| ENER314 Pin  | GPIO Pin | Description            |
| ------------ | -------- | ---------------------- |
| K0           |  17      | Switch address         |
| K1           |  22      | Switch address         |
| K2           |  23      | Switch address         |
| K3           |  27      | Switch address         |
| MODSEL       |  24      | Low OOK High FSK       |
| CE           |  25      | Low off High on        |

In order to install and/or run the command-line tool, use the
following:

```
  bash% cd $GOPATH/src/github.com/djthorpe/sensors
  bash% go (run|install) -tags rpi ./cmd/ener314/...
```

The command line tool demonstrates everything you need to know about
using the interface. You can run it with one of the following
commands:

  * `ener314 -on` Switches on
  * `ener314 -off` Switches off

You can append comma-separated socket numbers in order
to switch one or more sockets on or off. Without this argument,
all switches are controlled simultaneously. For example,

```
bash% ener314 -on 1,2
bash% ener314 -off 3,4
```

Will switch on sockets 1 and 2, with 3 and 4 switched off. Here is
the interface:

```
type ENER314 interface {
	gopi.Driver

	// Send on signal - when no sockets specified then
	// sends to all sockets
	On(sockets ...uint) error

	// Send off signal - when no sockets specified then
	// sends to all sockets
	Off(sockets ...uint) error
}
```

## The Two-Way Pi-Mote (ENER314-RT)

The ENER314-RT implements sending for both the OOK and FSK variants
of the Energenie product line, and receiving for the FSK variant:

  * The OOK variants are usually simple on/off switches and are called
    "Control" devices in the Energenie product catalogue;
  * The FSK variants require more complicated communication or are
    devices which send information, such as sensors. These are called
    "Monitor" devices in the Energenie product catalogue.

The ENER314-RT interface is as follows:

```
type ENER314RT interface {
  gopi.Driver

	// Receive payloads with radio until context deadline exceeded or cancel,
	// this blocks sending
	Receive(ctx context.Context, mode MiHomeMode, payload chan<- []byte) error

	// Send a raw payload with radio
	Send(payload []byte, repeat uint, mode MiHomeMode) error

	// Measure device temperature
	MeasureTemperature(offset float32) (float32, error)

	// ResetRadio device
	ResetRadio() error
}
```

The interface abstracts out sending and receiving of payloads, controlling
the LED's and resetting the radio. In fact, it's not likely you would create
a driver of this type. You are more likely to create a [MiHome](mihome.md)
driver which encodes and decodes with wire protocols, so please take a look
at more information there in order to understand how to use the ENER314-RT
board.





