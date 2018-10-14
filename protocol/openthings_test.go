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

func Test_OT_000_create(t *testing.T) {
	// Create an OOK module
	if app, err := gopi.NewAppInstance(gopi.NewAppConfig("sensors/protocol/openthings")); err != nil {
		t.Fatal(err)
	} else if _, ok := app.ModuleInstance("sensors/protocol/openthings").(sensors.Proto); ok == false {
		t.Fatal("Does not comply to Proto interface")
	} else if _, ok := app.ModuleInstance("sensors/protocol/openthings").(sensors.OTProto); ok == false {
		t.Fatal("Does not comply to OTProto interface")
	}
}

func Test_OT_001_message(t *testing.T) {
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

func Test_OT_002_newmessage(t *testing.T) {
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
func Test_OT_003_encode(t *testing.T) {
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

func Test_OT_004_decode(t *testing.T) {
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

func Test_OT_005_null(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		if _, err := proto.NewNull(sensors.OT_PARAM_NONE, false); err == nil {
			t.Error("Expected bad parameter")
		} else if record, err := proto.NewNull(sensors.OT_PARAM_JOIN, false); err != nil {
			t.Error(err)
		} else if record.Name() != sensors.OT_PARAM_JOIN {
			t.Error("Expected name=OT_PARAM_JOIN")
		} else if record.Type() != sensors.OT_DATATYPE_UDEC_0 {
			t.Error("Expected type=OT_DATATYPE_UDEC_0")
		} else if value, err := record.UintValue(); err != nil {
			t.Error(err)
		} else if value != 0 {
			t.Error("Expected value=0")
		} else {
			t.Log("NULL=", record)
		}
	}
}

func Test_OT_006_string(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		if _, err := proto.NewString(sensors.OT_PARAM_NONE, "", false); err == nil {
			t.Error("Expected bad parameter")
		} else if record, err := proto.NewString(sensors.OT_PARAM_DEBUG_OUTPUT, "string", false); err != nil {
			t.Error(err)
		} else if record.Name() != sensors.OT_PARAM_DEBUG_OUTPUT {
			t.Error("Expected name=OT_PARAM_DEBUG_OUTPUT")
		} else if record.Type() != sensors.OT_DATATYPE_STRING {
			t.Error("Expected type=OT_DATATYPE_STRING")
		} else if value, err := record.StringValue(); err != nil {
			t.Error(err)
		} else if value != "string" {
			t.Error("Expected value=string")
		} else {
			t.Log("STRING=", record)
		}
	}
}

func Test_OT_007_string_length(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		if _, err := proto.NewString(sensors.OT_PARAM_NONE, "", false); err == nil {
			t.Error("Expected bad parameter")
		} else if _, err := proto.NewString(sensors.OT_PARAM_DEBUG_OUTPUT, "0123456789ABCDE", false); err != nil {
			t.Error(err)
		} else if _, err := proto.NewString(sensors.OT_PARAM_DEBUG_OUTPUT, "0123456789ABCDEF", false); err == nil {
			t.Error("Expected error")
		}
	}
}

func Test_OT_008_uint(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		if _, err := proto.NewUint(sensors.OT_PARAM_NONE, 0, false); err == nil {
			t.Error("Expected bad parameter")
		} else if record, err := proto.NewUint(sensors.OT_PARAM_LEVEL, 0, false); err != nil {
			t.Error(err)
		} else if record.Name() != sensors.OT_PARAM_LEVEL {
			t.Error("Expected name=OT_PARAM_LEVEL")
		} else if record.Type() != sensors.OT_DATATYPE_UDEC_0 {
			t.Error("Expected type=OT_DATATYPE_UDEC_0")
		} else if value, err := record.UintValue(); err != nil {
			t.Error(err)
		} else if value != 0 {
			t.Error("Expected value=0")
		} else {
			t.Log(record)
		}
	}
}

func Test_OT_009_int(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		// Zero
		if _, err := proto.NewInt(sensors.OT_PARAM_NONE, 0, false); err == nil {
			t.Error("Expected bad parameter")
		} else if record, err := proto.NewInt(sensors.OT_PARAM_LEVEL, 0, false); err != nil {
			t.Error(err)
		} else if record.Name() != sensors.OT_PARAM_LEVEL {
			t.Error("Expected name=OT_PARAM_LEVEL")
		} else if record.Type() != sensors.OT_DATATYPE_DEC_0 {
			t.Error("Expected type=OT_DATATYPE_DEC_0")
		} else if value, err := record.IntValue(); err != nil {
			t.Error(err)
		} else if value != 0 {
			t.Error("Expected value=0")
		} else {
			t.Log(record)
		}

		// Positive
		if _, err := proto.NewInt(sensors.OT_PARAM_NONE, 0xFFFF, false); err == nil {
			t.Error("Expected bad parameter")
		} else if record, err := proto.NewInt(sensors.OT_PARAM_LEVEL, 0xFFFF, false); err != nil {
			t.Error(err)
		} else if record.Name() != sensors.OT_PARAM_LEVEL {
			t.Error("Expected name=OT_PARAM_LEVEL")
		} else if record.Type() != sensors.OT_DATATYPE_DEC_0 {
			t.Error("Expected type=OT_DATATYPE_DEC_0")
		} else if value, err := record.IntValue(); err != nil {
			t.Error(err)
		} else if value != 0xFFFF {
			t.Error("Expected value=0xFFFF")
		} else {
			t.Log(record)
		}

		// Negative -1
		if _, err := proto.NewInt(sensors.OT_PARAM_NONE, -1, false); err == nil {
			t.Error("Expected bad parameter")
		} else if record, err := proto.NewInt(sensors.OT_PARAM_LEVEL, -1, false); err != nil {
			t.Error(err)
		} else if record.Name() != sensors.OT_PARAM_LEVEL {
			t.Error("Expected name=OT_PARAM_LEVEL")
		} else if record.Type() != sensors.OT_DATATYPE_DEC_0 {
			t.Error("Expected type=OT_DATATYPE_DEC_0")
		} else if value, err := record.IntValue(); err != nil {
			t.Error(err)
		} else if value != -1 {
			t.Error("Expected value=-1, got", value)
		} else {
			t.Log(record)
		}

		// Negative -0x1234
		if _, err := proto.NewInt(sensors.OT_PARAM_NONE, -0x1234, false); err == nil {
			t.Error("Expected bad parameter")
		} else if record, err := proto.NewInt(sensors.OT_PARAM_LEVEL, -0x1234, false); err != nil {
			t.Error(err)
		} else if record.Name() != sensors.OT_PARAM_LEVEL {
			t.Error("Expected name=OT_PARAM_LEVEL")
		} else if record.Type() != sensors.OT_DATATYPE_DEC_0 {
			t.Error("Expected type=OT_DATATYPE_DEC_0")
		} else if value, err := record.IntValue(); err != nil {
			t.Error(err)
		} else if value != -0x1234 {
			t.Error("Expected value=", -0x1234, "got", value)
		} else {
			t.Log(record)
		}
	}
}

func Test_OT_010_bool_true(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		if _, err := proto.NewBool(sensors.OT_PARAM_NONE, true, false); err == nil {
			t.Error("Expected bad parameter")
		} else if record, err := proto.NewBool(sensors.OT_PARAM_DOOR_BELL, true, false); err != nil {
			t.Error(err)
		} else if record.Name() != sensors.OT_PARAM_DOOR_BELL {
			t.Error("Expected name=OT_PARAOT_PARAM_DOOR_BELLM_LEVEL")
		} else if record.Type() != sensors.OT_DATATYPE_UDEC_0 {
			t.Error("Expected type=OT_DATATYPE_UDEC_0")
		} else if bool_value, err := record.BoolValue(); err != nil {
			t.Error(err)
		} else if bool_value != true {
			t.Error("Expected value=true")
		} else if uint_value, err := record.UintValue(); err != nil {
			t.Error(err)
		} else if uint_value != 1 {
			t.Error("Expected value=1")
		} else {
			t.Log(record)
		}
	}
}

func Test_OT_011_bool_false(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		if _, err := proto.NewBool(sensors.OT_PARAM_NONE, false, false); err == nil {
			t.Error("Expected bad parameter")
		} else if record, err := proto.NewBool(sensors.OT_PARAM_DOOR_BELL, false, false); err != nil {
			t.Error(err)
		} else if record.Name() != sensors.OT_PARAM_DOOR_BELL {
			t.Error("Expected name=OT_PARAOT_PARAM_DOOR_BELLM_LEVEL")
		} else if record.Type() != sensors.OT_DATATYPE_UDEC_0 {
			t.Error("Expected type=OT_DATATYPE_UDEC_0")
		} else if bool_value, err := record.BoolValue(); err != nil {
			t.Error(err)
		} else if bool_value != false {
			t.Error("Expected value=false")
		} else if uint_value, err := record.UintValue(); err != nil {
			t.Error(err)
		} else if uint_value != 0 {
			t.Error("Expected value=0")
		} else {
			t.Log(record)
		}
	}
}

func Test_OT_012_uint_length(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		value := uint64(0)
		for i := 0; i < 9; i++ {
			if record, err := proto.NewUint(sensors.OT_PARAM_LEVEL, value, false); err != nil {
				t.Error(err)
			} else if data, err := record.Data(); err != nil {
				t.Error(err)
			} else if value_, err := record.UintValue(); err != nil {
				t.Error(err)
			} else if value != value_ {
				t.Error("Unexpected value returned")
			} else if i == 0 && len(data) != 3 {
				t.Error("Expected data for", record, "to be 3 bytes but got", strings.ToUpper(hex.EncodeToString(data)))
			} else if i > 0 && len(data) != i+2 {
				t.Error("Expected data for", record, "to be", i+2, "bytes but got", strings.ToUpper(hex.EncodeToString(data)))
			} else {
				t.Log(record, "=>", strings.ToUpper(hex.EncodeToString(data)))
			}
			value <<= 8
			value |= 0xFF
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
