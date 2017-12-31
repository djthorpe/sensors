# Sensors

This repository contains sensor interfaces for hardware sensors which
are interfaced through SPI and/or I2C. In order to use these interfaces,
the GOPI application framework (http://github.com/djthorpe/gopi) is used.

## BME280

The Bosch BME280 measures temperature, humidity and pressure. You can
interface this device either through I2C or SPI. The datasheet is
provided in the "doc" folder.

In order to connect it to the Raspberry Pi, here are the two 
configurations of the pins on the BME280 device. The AdaFruit
product is listed here as an example but there are other ways
to connect (more information
is avilable at https://www.adafruit.com/product/2652). The pin numbers
here are provided for connecting the AdaFruit product with a Raspberry PI
and are for the physical board pins:

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
  bash% cd $GOPATH/src/github.com/djthorpe/sensors
  bash% go (run|install) -tags i2c cmd/bme280.go
```

When you want to install or run a version which communicates over the 
SPI bus:

```
  bash% cd $GOPATH/src/github.com/djthorpe/sensors
  bash% go (run|install) -tags spi cmd/bme280.go
```

The command line tool demonstrates everything you need to know about
using the sensor interface. You can run it with one of the following
commands:


  * `bme280 reset` Resets the sensor and displays the sensor status
  * `bme280 status` Displays the current sensor status
  * `bme280 measure` Measures Temperature, Pressure and/or Humidity

There are also various flags you can use in order to set filter,
mode, oversampling or standby time. Here are the main flags you can use
on the command line:

```
  -filter uint
    	Set filter co-efficient (0,2,4,8,16)
  -mode string
    	Set sensor mode (normal,forced,sleep)
  -oversample uint
    	Set oversampling of measurements (0,1,2,4,8,16)
  -standby float
    	Standby time between measurements in normal mode, ms (0.5,10,20,62.5,125,250,500,1000)
```

There are a set of additional flags you can also use, some of which are available for I2C and others for SPI:

```
  -i2c.bus uint
    	I2C Bus (default 1)
  -i2c.slave uint
    	I2C Slave address (default 0x77)
  -spi.speed uint
    	SPI Communication Speed
  -debug
    	Set debugging mode
  -verbose
    	Verbose logging
  -log.append
    	When writing log to file, append output to end of file
  -log.file string
    	File for logging (default: log to stderr)
```

The modes of operation are essentailly "forced" or "normal". Forced
mode can be used for one-off measurement, or otherwise you can read
measurements with a particular duty cycle, which is a function of
standby time and oversampling values. Here's an example of status
output demonstrating this:

```
bash% bme280 status
+------------------------+----------------------+
|        REGISTER        |        VALUE         |
+------------------------+----------------------+
| chip_id                | 0x60                 |
| chip_version           | 0x00                 |
| mode                   | BME280_MODE_NORMAL   |
| filter                 | BME280_FILTER_OFF    |
| standby                | BME280_STANDBY_500MS |
| duty cycle             | 558ms                |
| oversample temperature | BME280_OVERSAMPLE_8  |
| oversample pressure    | BME280_OVERSAMPLE_8  |
| oversample humidity    | BME280_OVERSAMPLE_8  |
| measuring              | false                |
| updating               | false                |
+------------------------+----------------------+

```

In this example, the standby time is 500ms, measurement time is 
58ms and you shouldn't poll for new measurements more than once 
every 558ms (the duty cycle).

In forced mode only one measurement is taken before the sensor
returns to sleep mode:

```
bash% bme280 -mode forced measure
+-------------+------------+
| MEASUREMENT |   VALUE    |
+-------------+------------+
| temperature |   20.03 °C |
|    pressure | 998.61 hPa |
|    altitude |   122.64 m |
|    humidity |  53.52 %RH |
+-------------+------------+

```

## TSL2561

The TSL2561 luminosity sensor is a digital light sensor. You can
interface this device through I2C. The datasheet is provided in 
the "doc" folder.

In order to connect it to the Raspberry Pi, here are the pin 
configurations. The AdaFruit product is listed here as an 
example but there are other ways to connect (more information
is available at https://learn.adafruit.com/tsl2561). The pin numbers
here are provided for connecting the AdaFruit product with a Raspberry PI
and are for the physical board pins:

| TSL2561 Pin  | GPIO Pin | Description            |
| ------------ | -------- | ---------------------- |
| Vin          |  2       | 3-5VDC power in        |
| GND          |  6       | Ground                 |
| 3Vo          |          | 3.3V power out         |
| Addr         |          | I2C Address Change     |
| Int          |          | Light Change Interrupt |
| SDA          |  3       | I2C Data               |
| SCL          |  5       | I2C Clock              |

The I2C slave address defaults to 0x39. By connecting the Addr 
pin to ground, this changes to 0x29 and connecting to 3.3V 
it changes to 0x49.

In order to install and/or run the command-line tool, use the
following:

```
  bash% cd $GOPATH/src/github.com/djthorpe/sensors
  bash% go (run|install) cmd/tsl2561.go
```


