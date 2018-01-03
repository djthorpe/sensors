/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package rfm69

import "fmt"
import "strings"

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *rfm69) String() string {
	params := []string{
		fmt.Sprintf("version=0x%02X", this.version),
		fmt.Sprintf("mode=%v", this.mode),
		fmt.Sprintf("data_mode=%v", this.data_mode),
		fmt.Sprintf("modulation=%v", this.modulation),
		fmt.Sprintf("sequencer_off=%v", this.sequencer_off),
		fmt.Sprintf("listen_on=%v", this.listen_on),
	}
	return fmt.Sprintf("sensors.RFM69{ spi=%v %v }", this.spi, strings.Join(params, " "))
}

func (r register) String() string {
	switch r {
	case RFM_REG_FIFO:
		return "RFM_REG_FIFO"
	case RFM_REG_OPMODE:
		return "RFM_REG_OPMODE"
	case RFM_REG_DATAMODUL:
		return "RFM_REG_DATAMODUL"
	case RFM_REG_BITRATEMSB:
		return "RFM_REG_BITRATEMSB"
	case RFM_REG_BITRATELSB:
		return "RFM_REG_BITRATELSB"
	case RFM_REG_FDEVMSB:
		return "RFM_REG_FDEVMSB"
	case RFM_REG_FDEVLSB:
		return "RFM_REG_FDEVLSB"
	case RFM_REG_FRFMSB:
		return "RFM_REG_FRFMSB"
	case RFM_REG_FRFMID:
		return "RFM_REG_FRFMID"
	case RFM_REG_FRFLSB:
		return "RFM_REG_FRFLSB"
	case RFM_REG_OSC1:
		return "RFM_REG_OSC1"
	case RFM_REG_AFCCTRL:
		return "RFM_REG_AFCCTRL"
	case RFM_REG_LISTEN1:
		return "RFM_REG_LISTEN1"
	case RFM_REG_LISTEN2:
		return "RFM_REG_LISTEN2"
	case RFM_REG_LISTEN3:
		return "RFM_REG_LISTEN3"
	case RFM_REG_VERSION:
		return "RFM_REG_VERSION"
	case RFM_REG_PALEVEL:
		return "RFM_REG_PALEVEL"
	case RFM_REG_PARAMP:
		return "RFM_REG_PARAMP"
	case RFM_REG_OCP:
		return "RFM_REG_OCP"
	case RFM_REG_LNA:
		return "RFM_REG_LNA"
	case RFM_REG_RXBW:
		return "RFM_REG_RXBW"
	case RFM_REG_AFCBW:
		return "RFM_REG_AFCBW"
	case RFM_REG_OOKPEAK:
		return "RFM_REG_OOKPEAK"
	case RFM_REG_OOKAVG:
		return "RFM_REG_OOKAVG"
	case RFM_REG_OOKFIX:
		return "RFM_REG_OOKFIX"
	case RFM_REG_AFCFEI:
		return "RFM_REG_AFCFEI"
	case RFM_REG_AFCMSB:
		return "RFM_REG_AFCMSB"
	case RFM_REG_AFCLSB:
		return "RFM_REG_AFCLSB"
	case RFM_REG_FEIMSB:
		return "RFM_REG_FEIMSB"
	case RFM_REG_FEILSB:
		return "RFM_REG_FEILSB"
	case RFM_REG_RSSICONFIG:
		return "RFM_REG_RSSICONFIG"
	case RFM_REG_RSSIVALUE:
		return "RFM_REG_RSSIVALUE"
	case RFM_REG_DIOMAPPING1:
		return "RFM_REG_DIOMAPPING1"
	case RFM_REG_DIOMAPPING2:
		return "RFM_REG_DIOMAPPING2"
	case RFM_REG_IRQFLAGS1:
		return "RFM_REG_IRQFLAGS1"
	case RFM_REG_IRQFLAGS2:
		return "RFM_REG_IRQFLAGS2"
	case RFM_REG_RSSITHRESH:
		return "RFM_REG_RSSITHRESH"
	case RFM_REG_RXTIMEOUT1:
		return "RFM_REG_RXTIMEOUT1"
	case RFM_REG_RXTIMEOUT2:
		return "RFM_REG_RXTIMEOUT2"
	case RFM_REG_PREAMBLEMSB:
		return "RFM_REG_PREAMBLEMSB"
	case RFM_REG_PREAMBLELSB:
		return "RFM_REG_PREAMBLELSB"
	case RFM_REG_SYNCCONFIG:
		return "RFM_REG_SYNCCONFIG"
	case RFM_REG_SYNCVALUE1:
		return "RFM_REG_SYNCVALUE1"
	case RFM_REG_SYNCVALUE2:
		return "RFM_REG_SYNCVALUE2"
	case RFM_REG_SYNCVALUE3:
		return "RFM_REG_SYNCVALUE3"
	case RFM_REG_SYNCVALUE4:
		return "RFM_REG_SYNCVALUE4"
	case RFM_REG_SYNCVALUE5:
		return "RFM_REG_SYNCVALUE5"
	case RFM_REG_SYNCVALUE6:
		return "RFM_REG_SYNCVALUE6"
	case RFM_REG_SYNCVALUE7:
		return "RFM_REG_SYNCVALUE7"
	case RFM_REG_SYNCVALUE8:
		return "RFM_REG_SYNCVALUE8"
	case RFM_REG_PACKETCONFIG1:
		return "RFM_REG_PACKETCONFIG1"
	case RFM_REG_PAYLOADLENGTH:
		return "RFM_REG_PAYLOADLENGTH"
	case RFM_REG_NODEADRS:
		return "RFM_REG_NODEADRS"
	case RFM_REG_BROADCASTADRS:
		return "RFM_REG_BROADCASTADRS"
	case RFM_REG_AUTOMODES:
		return "RFM_REG_AUTOMODES"
	case RFM_REG_FIFOTHRESH:
		return "RFM_REG_FIFOTHRESH"
	case RFM_REG_PACKETCONFIG2:
		return "RFM_REG_PACKETCONFIG2"
	case RFM_REG_AESKEY1:
		return "RFM_REG_AESKEY1"
	case RFM_REG_AESKEY2:
		return "RFM_REG_AESKEY2"
	case RFM_REG_AESKEY3:
		return "RFM_REG_AESKEY3"
	case RFM_REG_AESKEY4:
		return "RFM_REG_AESKEY4"
	case RFM_REG_AESKEY5:
		return "RFM_REG_AESKEY5"
	case RFM_REG_AESKEY6:
		return "RFM_REG_AESKEY6"
	case RFM_REG_AESKEY7:
		return "RFM_REG_AESKEY7"
	case RFM_REG_AESKEY8:
		return "RFM_REG_AESKEY8"
	case RFM_REG_AESKEY9:
		return "RFM_REG_AESKEY9"
	case RFM_REG_AESKEY10:
		return "RFM_REG_AESKEY10"
	case RFM_REG_AESKEY11:
		return "RFM_REG_AESKEY11"
	case RFM_REG_AESKEY12:
		return "RFM_REG_AESKEY12"
	case RFM_REG_AESKEY13:
		return "RFM_REG_AESKEY13"
	case RFM_REG_AESKEY14:
		return "RFM_REG_AESKEY14"
	case RFM_REG_AESKEY15:
		return "RFM_REG_AESKEY15"
	case RFM_REG_AESKEY16:
		return "RFM_REG_AESKEY16"
	case RFM_REG_TEMP1:
		return "RFM_REG_TEMP1"
	case RFM_REG_TEMP2:
		return "RFM_REG_TEMP2"
	case RFM_REG_TEST:
		return "RFM_REG_TEST"
	case RFM_REG_TESTLNA:
		return "RFM_REG_TESTLNA"
	case RFM_REG_TESTPA1:
		return "RFM_REG_TESTPA1"
	case RFM_REG_TESTPA2:
		return "RFM_REG_TESTPA2"
	case RFM_REG_TESTDAGC:
		return "RFM_REG_TESTDAGC"
	case RFM_REG_TESTAFC:
		return "RFM_REG_TESTAFC"
	default:
		return "[?? Invalid register value]"
	}
}
