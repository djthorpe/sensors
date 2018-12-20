
# HopeRF RFM69 Transceiver

The [Hope RFM69HW module](http://www.hoperf.com/rf_transceiver/modules/RFM69HW.html)
is an RF transceiver capable of receiving and transmitting signals using FSK
and OOK modulation. Ultimately it's a popular general purpose device for general
data transmission and integration with many cheap electronics which use either the
433MHz or 900MHz bands.

  * The RFM69HW datasheet [link](RFM69HW-V1.3.pdf)
  * List of HopeRF transceiver modules [link](https://www.hoperf.com/modules/rf_transceiver/index.html)

## Printed Circuit Board 

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

![RFM69 Schematic](rfm69-schematic-v1.png)

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

## The Command Line Tool 

Once you've constructed the PCB you can plug it in and communicate with
it on the SPI bus. In order to install and/or run the command-line tool, use the
following:

```
  bash% cd $GOPATH/src/github.com/djthorpe/sensors
  bash% go (run|install) cmd/rfm69.go
```

## RFM69 Interface

The interface for the RFM69 is as follows:

```

type RFM69 interface {
	gopi.Driver

	// Mode, Data Mode and Modulation
	Mode() RFMMode
	DataMode() RFMDataMode
	SetMode(device_mode RFMMode) error
	SetDataMode(data_mode RFMDataMode) error
	Modulation() RFMModulation
	SetModulation(modulation RFMModulation) error

	// Bitrate & Frequency
	Bitrate() uint
	FreqCarrier() uint
	FreqDeviation() uint
	SetBitrate(bits_per_second uint) error
	SetFreqCarrier(hertz uint) error
	SetFreqDeviation(hertz uint) error

	// Listen Mode and Sequencer
	SetSequencer(enabled bool) error
	SequencerEnabled() bool
	SetListenOn(value bool) error
	ListenOn() bool

	// Packets
	PacketFormat() RFMPacketFormat
	PacketCoding() RFMPacketCoding
	PacketFilter() RFMPacketFilter
	PacketCRC() RFMPacketCRC
	SetPacketFormat(packet_format RFMPacketFormat) error
	SetPacketCoding(packet_coding RFMPacketCoding) error
	SetPacketFilter(packet_filter RFMPacketFilter) error
	SetPacketCRC(packet_crc RFMPacketCRC) error

	// Addresses
	NodeAddress() uint8
	BroadcastAddress() uint8
	SetNodeAddress(value uint8) error
	SetBroadcastAddress(value uint8) error

	// Payload & Preamble
	PreambleSize() uint16
	PayloadSize() uint8
	SetPreambleSize(preamble_size uint16) error
	SetPayloadSize(payload_size uint8) error

	// Encryption Key & Sync Words for Packet mode
	AESKey() []byte
	SetAESKey(key []byte) error
	SyncWord() []byte
	SetSyncWord(word []byte) error
	SyncTolerance() uint8
	SetSyncTolerance(bits uint8) error

	// AFC
	AFC() uint
	AFCMode() RFMAFCMode
	AFCRoutine() RFMAFCRoutine
	SetAFCRoutine(afc_routine RFMAFCRoutine) error
	SetAFCMode(afc_mode RFMAFCMode) error
	TriggerAFC() error

	// Low Noise Amplifier Settings
	LNAImpedance() RFMLNAImpedance
	LNAGain() RFMLNAGain
	LNACurrentGain() (RFMLNAGain, error)
	SetLNA(impedance RFMLNAImpedance, gain RFMLNAGain) error

	// Channel Filter Settings
	RXFilterFrequency() RFMRXBWFrequency
	RXFilterCutoff() RFMRXBWCutoff
	SetRXFilter(RFMRXBWFrequency, RFMRXBWCutoff) error

	// FIFO
	FIFOThreshold() uint8
	SetFIFOThreshold(fifo_threshold uint8) error
	ReadFIFO(ctx context.Context) ([]byte, error)
	WriteFIFO(data []byte) error
	ClearFIFO() error

	// ReadPayload listens for a packet and returns it. If the data is
	// read then it will also return true if the CRC value was
	// correct, or false otherwise
	ReadPayload(ctx context.Context) ([]byte, bool, error)

	// WritePayload writes a packet a number of times, with a delay between each
	// when the repeat is greater than zero
	WritePayload(data []byte, repeat uint, delay time.Duration) error

	// MeasureTemperature and return after calibration
	MeasureTemperature(calibration float32) (float32, error)
}
```

There are some missing methods for setting power levels and measuring signal
noise.

(more information on the module here shortly)

