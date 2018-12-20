# Sensors

This repository contains sensor interfaces for hardware sensors which
are interfaced through SPI and/or I2C. There are also protocols for
communicating between sensors and RPC microservices for accessing sensor
data remotely.

In order to use these interfaces, the GOPI application 
framework (http://github.com/djthorpe/gopi) is used, and the associated
set of modules for interfacing hardware and remote procedure calls.

The interfaces and definitions for the sensors are in the package
root: `sensors.go`, `rfm69.go`, `ads1x15.go`, `bme680.go`, `energenie.go`
and `protocol.go`. To create a sensor driver you need to create it using the 
`gopi.Open` method on a concrete driver. You can check the examples
in the `cmd` directory for more information.

For more information on using the modules, the documentation is in the `doc` folder:

  * For Bosch BME280 and BME680 temperature, humidity, pressure and air quality
    sensors please see [`doc/BMEx80.md`](https://github.com/djthorpe/sensors/blob/master/doc/BMEx80.md);
  * For the TAOS TSL2561 luminosity sensor, please see  [`doc/TSL2561.md`](https://github.com/djthorpe/sensors/blob/master/doc/TSL2561.md);
  * For the HopeRF RFM69 radio transceiver series, please see [`doc/RFM69.md`](https://github.com/djthorpe/sensors/blob/master/doc/RFM69.md);
  * For the Texas Instruments ADS1015 and ADS1115 analog-to-digital converters,
    please see [`doc/ADS1x15.md`](https://github.com/djthorpe/sensors/blob/master/doc/ADS1x15.md);
  * For the ENER314 OOK transmitter and OOK/FSK transciever boards,
    please see [`doc/ENER314.md`](https://github.com/djthorpe/sensors/blob/master/doc/ENER314.md);
  * For the implementation of the wire protocol for Energenie MiHome series,
    please see [`doc/mihome.md`](https://github.com/djthorpe/sensors/blob/master/doc/mihome.md).

## Building the examples

## ENER314

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

## RFM69

The [Hope RFM69HW module](http://www.hoperf.com/rf_transceiver/modules/RFM69HW.html)
is an RF transceiver capable of receiving and transmitting signals using FSK
and OOK modulation. Ultimately it's a popular general purpose device for general
data transmission and integration with many cheap electronics which use either the
433MHz or 900MHz bands.

The module is mounted on a custom-made PCB which you'll have to construct. You can
manufacture a PCB from Aisler [here](https://aisler.net/p/MINCPAPN). Total cost
including components would cost about £12 / €14 / $14 per unit. This is the
bill of materials:

| Name  | Part               | Description |
| ----- | ------------------ | ---- |
| X1    |  SMA Connector     | Straight 50Ω PCB Mount Bulkhead Fitting SMA Connector, Solder Termination |
| U1    |  Hope RFM69HCW     | HopeRF RF Transceiver RFM69W-433-S2 433 MHz, 1.8 → 3.6V |
| R1    |  330Ω ±5% 0.25W    | Carbon Resistor, 0.25W, 5%, 330R |
| C1,C2 |  0.1uF             | Ceramic Decoupling Capacitors, 0.1uF |
| LED1  |  Generic LED       | LED, 5mm (T-1 3/4) Through Hole package |
| J1    |  26 Way PCB Header | 2.54mm Pitch 13x2 Rows Straight PCB Socket |

Here's the schematic:

![RFM69 Schematic](https://raw.githubusercontent.com/djthorpe/sensors/master/doc/rfm69-schematic-v1.png)

In addition you'll want to purchase a 433Mhz antenna with an SMA connector which
will cost an extra couple of Euros, or you can simply solder a 50ohm 173mm wire if you
want to omit the expensive SMA connector. 

The pinouts for the module then correspond to the following physical and GPIO pins:

| RFM69 Pin | Physical | GPIO Pin | Description             |
| --------- | -------- | -------- | ----------------------- |
| MISO      |  21      | SPI_MISO | SPI Master In Slave Out |
| MOSI      |  19      | SPI_MOSI | SPI Master Out Slave In |
| SCK       |  23      | SPI_CLK  | SPI Clock               |
| NSS       |  26      | SPI_CE1  | SPI Chip Select CE1     |
| RESET     |  22      | GPIO25   | Module Reset            |
| DIO0      |   7      | GPIO4    | Data In Out 0           |
| DIO1      |  11      | GPIO17   | Data In Out 1           |
| DIO2      |  12      | GPIO18   | Data In Out 2           |
| DIO3      |  15      | GPIO22   | Data In Out 3           |
| DIO4      |  18      | GPIO24   | Data In Out 4           |
| DIO5      |  16      | GPIO23   | Data In Out 5           |

In addition there's a connection for an LED on physical pin 13,
which can be programmed to be GPIO17 or GPIO27 depending on the model
of Raspberry Pi you have.

Once you've constructed the PCB you can plug it in and communicate with
it on the SPI bus. In order to install and/or run the command-line tool, use the
following:

```
  bash% cd $GOPATH/src/github.com/djthorpe/sensors
  bash% go (run|install) cmd/rfm69.go
```

(more information on the module here shortly)


# License

Copyright 2016-2018 David Thorpe All Rights Reserved

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

