
# Energenie ENER314

You can purchase two boards from Energenie which fit into the Raspberry Pi. The two separate boards
are:

  1. The Pi-Mote Remote Control Board [link](https://energenie4u.co.uk/catalogue/product/ENER314)
     which is a simple OOK transmitter;
  2. The two-way Pi-Mote [link]( https://energenie4u.co.uk/catalogue/product/ENER314-RT) which
     is actually a 433Mhz RFM69 board.

Both of these are implemented. The former communicates simply through GPIO whereas the latter
communicates through SPI and uses GPIO to set the receive and transmit LED on the board and
resets the RFM69 module.

In order to create a driver for each, you can import the following modules
anonymously into your application:

| Import                                    | Module Name        | Interface  |
| ----------------------------------------- | ------------------ | ---------- |
| github.com/djthorpe/sensors/sys/ener314   | sensors/ener314    | GPIO       |
| github.com/djthorpe/sensors/sys/ener314rt | sensors/ener314rt  | SPI & GPIO |

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
  bash% go (run|install) cmd/ener314.go
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

Will switch on sockets 1 and 2, with 3 and 4 switched off.

## The Two-Way Pi-Mote (ENER314-RT)

(To be done)


