package protocol_test

import (
	"encoding/hex"
	"strings"
	"testing"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/sensors"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/sensors/protocol/openthings"
)

var (
	// An array of messages which are correct (with good CRC)
	received_good = []string{
		"0E040303DD35D6C56D328716BDBFB3",
		"0E04031AC01C395F52DBAECB39A52A",
		"0D040D33F2995BC2153D4AF764D4",
		"0D040D33F2995BC2153D4BF757E5",
		"0D040D33F2995BC2153D4BF757E5",
		"0D040D33F2995BC2153D4AF764D4",
		"0D040D33F2995BC2153D4AF764D4",
		"0D040D33F2995BC2153D4BF757E5",
		"0E04032244AA1AD6307144E48C25E9",
		"0E040303DD35D6C56D328716BDBFB3",
		"0E04035DEA13B49753F22078A248DD",
		"0D040C33F29950491C3D4BF7A91F",
		"0E0403714F6D5639FAA12FDEB80389",
		"0D040C33F29950491C3D4AF79A2E",
		"0D040C33F29950491C3D4BF7A91F",
		"1C0402F2E134FBB048D28C4211CA964949E31C29CE30C400D4DBCF5ED3",
		"1C0402BFEA03453C89756BC796BCA35E6ED819AEA8E86890D5125F5C83",
		"0D0402FE2C4614C2443455CE6012",
		"0C040364CEA883B988691EE331",
		"0C04032020CBE0F2F5D350BDAA",
		"0C04033543304D13839E86259C",
	}
)

func Test_OT_000(t *testing.T) {
	// Create an OOK module
	if app, err := gopi.NewAppInstance(gopi.NewAppConfig("sensors/protocol/openthings")); err != nil {
		t.Fatal(err)
	} else if _, ok := app.ModuleInstance("sensors/protocol/openthings").(sensors.Proto); ok == false {
		t.Fatal("Does not comply to Proto interface")
	} else if _, ok := app.ModuleInstance("sensors/protocol/openthings").(sensors.OTProto); ok == false {
		t.Fatal("Does not comply to OTProto interface")
	}
}

func Test_OT_001(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else if msg_, err := proto.New(sensors.OT_MANUFACTURER_ENERGENIE, 0, 0); err != nil {
		t.Fatal(err)
	} else if msg, ok := msg_.(sensors.OTMessage); ok == false {
		t.Fatal("Not an OTMessage")
	} else if msg.Manufacturer() != sensors.OT_MANUFACTURER_ENERGENIE {
		t.Error("Unexpected manufacturer")
	} else if msg.Product() != 0 {
		t.Error("Unexpected product")
	} else if msg.Sensor() != 0 {
		t.Error("Unexpected sensor")
	}
}

func Test_OT_002(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		for sensor := uint32(0); sensor < 0xFFFFFF; sensor += 0x1234 {
			product := uint8(sensor % 0x7F)
			if msg, err := proto.New(sensors.OT_MANUFACTURER_ENERGENIE, product, sensor); err != nil {
				t.Error(err)
			} else {
				t.Log(msg)
			}
		}
	}
}
func Test_OT_003(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		for sensor := uint32(0); sensor < 0xFFFFFF; sensor += 0x1234 {
			product := uint8(sensor % 0x7F)
			if msg, err := proto.New(sensors.OT_MANUFACTURER_ENERGENIE, product, sensor); err != nil {
				t.Error(err)
			} else if payload := proto.Encode(msg); len(payload) == 0 {
				t.Error("Empty payload")
			} else {
				t.Log(msg)
				t.Logf("...%v", strings.ToUpper(hex.EncodeToString(payload)))
			}
		}
	}
}

func Test_OT_004(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		for i, received_str := range received_good {
			if payload, err := hex.DecodeString(received_str); err != nil {
				t.Fatal("Message", i, err)
			} else if msg, err := proto.Decode(payload, time.Now()); err != nil {
				t.Fatal("Message", i, err)
			} else {
				t.Log("Message", i, msg)
			}
		}
	}

}

////////////////////////////////////////////////////////////////////////////////
// OT

func OTProto() sensors.OTProto {
	config := gopi.NewAppConfig("sensors/protocol/openthings")
	if testing.Verbose() {
		config.Debug = true
		config.Verbose = true
	}
	if app, err := gopi.NewAppInstance(config); err != nil {
		return nil
	} else if proto, ok := app.ModuleInstance("sensors/protocol/openthings").(sensors.OTProto); ok == false {
		return nil
	} else {
		return proto
	}
}
