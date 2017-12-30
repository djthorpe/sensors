# Sensors

This repository contains sensor interfaces for hardware sensors which
are interfaced through SPI and/or I2C. In order to use these interfaces,
the GOPI application framework (http://github.com/djthorpe/gopi) is used.

## BME280

The Bosch BME280 measures temperature, humidity and pressure. You can
interface this device either through I2C or SPI. The datasheet is
provided in the "doc" folder.

In order to connect it to the Raspberry Pi, here are the two 
configurations of the pins on the AdaFruit BME280 device (more information
is avilable at https://www.adafruit.com/product/2652). The pin numbers
provided are for the physical board pins:

| BME280 Pin   | I2C | SPI(0) | Description                    |
| ------------ | --- | ------ | ------------------------------ |
| Vin          |  2  |  2     | 3-5VDC power in                |
| 3Vo          |     |        | 3.3V power out                 |
| GND          |  6  |  6     | Ground                         |
| SCK          |  5  | 23     | SPI/I2C Clock Pin              |
| SDO/MISO     |     | 21     | Serial Data Out                |
| SDI/SDA/MOSI |  3  | 19     | Serial Data In / I2C Data Pin  |
| CS           |     | 24     | Chip Select                    |

Using the I2C connection bus, the slave address defaults to 0x77.

In order to install and/or run the command-line tool, you
need to deploy the correct flavour using tags. Here is how
you can can install and/or run when you want to use the I2C
flavour:

```
  cd $GOPATH/src/github.com/djthorpe/sensors
  go (run|install) -tags i2c cmd/bme280.go
```

When you want to install or run a version which communicates over the 
SPI bus:

```
  cd $GOPATH/src/github.com/djthorpe/sensors
  go (run|install) -tags spi cmd/bme280.go
```


