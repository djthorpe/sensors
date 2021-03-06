# Sensors

This repository contains sensor interfaces for hardware sensors which
are interfaced through GPIO, SPI and/or I2C. There are also protocols 
for communicating between sensors and RPC microservices for accessing sensor
data remotely.

In order to use these interfaces, the GOPI application 
framework (http://github.com/djthorpe/gopi) is used, and the associated
set of modules for interfacing hardware and remote procedure calls.

The interfaces and definitions for the sensors are in the package
root: `sensors.go`, `rfm69.go`, `ads1x15.go`, `bme680.go`, `energenie.go`
and `protocol.go`. You can check the examples in the `cmd` directory for more
information on using the drivers.

For more information on using the drivers, the documentation is in the `doc` folder:

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

There is a makefile which will test and make all the example commands for both MacOS and Raspberry Pi 
target. On MacOS, only the client command-line utility which allows you to connect to a remote
`mihome-service` is built. On the Raspberry Pi, all the examples are built.

On MacOS:

```bash
% make -f Makefile.darwin
```

On Raspberry Pi:

```bash
% make -f Makefile.rpi
```

# License

Copyright 2016-2019 David Thorpe All Rights Reserved

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

