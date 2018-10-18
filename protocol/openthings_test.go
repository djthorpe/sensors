package protocol_test

import (
	"encoding/hex"
	"math"
	"math/rand"
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

func Test_OT_012_uint8(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		for value := uint64(0); value <= uint64(math.MaxUint8); value += 1 {
			if record, err := proto.NewUint(sensors.OT_PARAM_LEVEL, value, false); err != nil {
				t.Error(err)
			} else if data, err := record.Data(); err != nil {
				t.Error(err)
			} else if value_, err := record.UintValue(); err != nil {
				t.Error(err)
			} else if value != value_ {
				t.Error("Unexpected value returned, expected", value, "but got", value_, "[", strings.ToUpper(hex.EncodeToString(data)), "]")
			} else if len(data) != 3 {
				t.Error("Unexpected data length", len(data))
			} else {
				t.Log(record, "=>", strings.ToUpper(hex.EncodeToString(data)))
			}
		}
	}
}

func Test_OT_013_uint16(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		for value := uint64(math.MaxUint8 + 1); value <= uint64(math.MaxUint16); value += 123 {
			if record, err := proto.NewUint(sensors.OT_PARAM_LEVEL, value, false); err != nil {
				t.Error(err)
			} else if data, err := record.Data(); err != nil {
				t.Error(err)
			} else if value_, err := record.UintValue(); err != nil {
				t.Error(err)
			} else if value != value_ {
				t.Error("Unexpected value returned, expected", value, "but got", value_, "[", strings.ToUpper(hex.EncodeToString(data)), "]")
			} else if len(data) != 4 {
				t.Error("Unexpected data length", len(data))
			} else {
				t.Log(record, "=>", strings.ToUpper(hex.EncodeToString(data)))
			}
		}
	}
}

func Test_OT_014_uint32(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		for value := uint64(math.MaxUint16 + 1); value <= uint64(math.MaxUint32); value += 123456789 {
			if record, err := proto.NewUint(sensors.OT_PARAM_LEVEL, value, false); err != nil {
				t.Error(err)
			} else if data, err := record.Data(); err != nil {
				t.Error(err)
			} else if value_, err := record.UintValue(); err != nil {
				t.Error(err)
			} else if value != value_ {
				t.Error("Unexpected value returned, expected", value, "but got", value_, "[", strings.ToUpper(hex.EncodeToString(data)), "]")
			} else if len(data) != 6 {
				t.Error("Unexpected data length", len(data))
			} else {
				t.Log(record, "=>", strings.ToUpper(hex.EncodeToString(data)))
			}
		}
	}
}

func Test_OT_015_uint64(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		step := uint64(0xF123456789AB)
		for value := uint64(math.MaxUint32 + 1); value <= uint64(math.MaxUint64)-step-1; value += step {
			if record, err := proto.NewUint(sensors.OT_PARAM_LEVEL, value, false); err != nil {
				t.Error(err)
			} else if data, err := record.Data(); err != nil {
				t.Error(err)
			} else if value_, err := record.UintValue(); err != nil {
				t.Error(err)
			} else if value != value_ {
				t.Error("Unexpected value returned, expected", value, "but got", value_, "[", strings.ToUpper(hex.EncodeToString(data)), "]")
			} else if len(data) != 10 {
				t.Error("Unexpected data length", len(data))
			} else {
				t.Log(record, "=>", strings.ToUpper(hex.EncodeToString(data)))
			}
		}
	}
}

func Test_OT_016_int8(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		for value := int64(math.MinInt8); value <= int64(math.MaxInt8); value += 1 {
			if record, err := proto.NewInt(sensors.OT_PARAM_LEVEL, value, false); err != nil {
				t.Error(err)
			} else if data, err := record.Data(); err != nil {
				t.Error(err)
			} else if value_, err := record.IntValue(); err != nil {
				t.Error(err)
			} else if value != value_ {
				t.Error("Unexpected value returned, expected", value, "but got", value_, "[", strings.ToUpper(hex.EncodeToString(data)), "]")
			} else if len(data) != 3 {
				t.Error("Unexpected data length", len(data))
			} else {
				t.Log(record, "=>", strings.ToUpper(hex.EncodeToString(data)))
			}
		}
	}
}

func Test_OT_017_int16(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		for i := 0; i < 100000; i++ {
			// Create a random 16-bit value value
			value := rand.Int63n(math.MaxInt16)
			if rand.Intn(2) == 1 {
				value = -value
			}
			// Create a new integer record
			if record, err := proto.NewInt(sensors.OT_PARAM_LEVEL, value, false); err != nil {
				t.Error(err)
			} else if data, err := record.Data(); err != nil {
				t.Error(err)
			} else if value_, err := record.IntValue(); err != nil {
				t.Error(err)
			} else if value != value_ {
				t.Error("Unexpected value returned, expected", value, "but got", value_, "[", strings.ToUpper(hex.EncodeToString(data)), "]")
			} else {
				t.Log(value, "=>", record)
				switch {
				case value >= math.MinInt8 && value <= math.MaxInt8:
					if len(data) != 3 {
						t.Error("Unexpected data length", len(data), "for 8-bit value", value)
					}
				case value >= math.MinInt16 && value <= math.MaxInt16:
					if len(data) != 4 {
						t.Error("Unexpected data length", len(data), "for 16-bit value", value)
					}
				default:
					t.Error("value length is unknown", value)
				}
			}
		}
	}
}

func Test_OT_018_int32(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		for i := 0; i < 100000; i++ {
			// Create a random 32-bit value value
			value := rand.Int63n(math.MaxInt32)
			if rand.Intn(2) == 1 {
				value = -value
			}
			// Create a new integer record
			if record, err := proto.NewInt(sensors.OT_PARAM_LEVEL, value, false); err != nil {
				t.Error(err)
			} else if data, err := record.Data(); err != nil {
				t.Error(err)
			} else if value_, err := record.IntValue(); err != nil {
				t.Error(err)
			} else if value != value_ {
				t.Error("Unexpected value returned, expected", value, "but got", value_, "[", strings.ToUpper(hex.EncodeToString(data)), "]")
			} else {
				t.Log(value, "=>", record)
				switch {
				case value >= math.MinInt8 && value <= math.MaxInt8:
					if len(data) != 3 {
						t.Error("Unexpected data length", len(data), "for 8-bit value", value)
					}
				case value >= math.MinInt16 && value <= math.MaxInt16:
					if len(data) != 4 {
						t.Error("Unexpected data length", len(data), "for 16-bit value", value)
					}
				case value >= math.MinInt32 && value <= math.MaxInt32:
					if len(data) != 6 {
						t.Error("Unexpected data length", len(data), "for 32-bit value", value)
					}
				default:
					t.Error("value length is unknown", value)
				}
			}
		}
	}
}

func Test_OT_019_int64(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		for i := 0; i < 100000; i++ {
			// Create a random 64-bit value value
			value := rand.Int63n(math.MaxInt64)
			if rand.Intn(2) == 1 {
				value = -value
			}
			// Create a new integer record
			if record, err := proto.NewInt(sensors.OT_PARAM_LEVEL, value, false); err != nil {
				t.Error(err)
			} else if data, err := record.Data(); err != nil {
				t.Error(err)
			} else if value_, err := record.IntValue(); err != nil {
				t.Error(err)
			} else if value != value_ {
				t.Error("Unexpected value returned, expected", value, "but got", value_, "[", strings.ToUpper(hex.EncodeToString(data)), "]")
			} else {
				t.Log(value, "=>", record)
				switch {
				case value >= math.MinInt8 && value <= math.MaxInt8:
					if len(data) != 3 {
						t.Error("Unexpected data length", len(data), "for 8-bit value", value)
					}
				case value >= math.MinInt16 && value <= math.MaxInt16:
					if len(data) != 4 {
						t.Error("Unexpected data length", len(data), "for 16-bit value", value)
					}
				case value >= math.MinInt32 && value <= math.MaxInt32:
					if len(data) != 6 {
						t.Error("Unexpected data length", len(data), "for 32-bit value", value)
					}
				case value >= math.MinInt64 && value <= math.MaxInt64:
					if len(data) != 10 {
						t.Error("Unexpected data length", len(data), "for 64-bit value", value)
					}
				default:
					t.Error("value length is unknown", value)
				}
			}
		}
	}
}

func Test_OT_020_float_dec0(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		for i := 0; i < 1000; i++ {
			// Create a random 64-bit value value
			value := rand.Int63n(math.MaxInt64)
			if rand.Intn(2) == 1 {
				value = -value
			}
			// Create float with the int64 value
			if record, err := proto.NewFloat(sensors.OT_PARAM_FREQUENCY, sensors.OT_DATATYPE_DEC_0, float64(value), false); err != nil {
				t.Error(err)
			} else if record.Type() != sensors.OT_DATATYPE_DEC_0 {
				t.Error("Expected type=OT_DATATYPE_DEC_0")
			} else if value_, err := record.FloatValue(); err != nil {
				t.Error(err)
			} else if float64(value) != value_ {
				t.Error("Unexpected value", value_, "expected", float64(value))
			} else {
				t.Log(record)
			}
		}
	}
}

func Test_OT_021_float_udec0(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		for i := 0; i < 1000; i++ {
			// Create a random 64-bit value value
			value := rand.Uint64()
			// Create float with the int64 value
			if record, err := proto.NewFloat(sensors.OT_PARAM_FREQUENCY, sensors.OT_DATATYPE_UDEC_0, float64(value), false); err != nil {
				t.Error(err)
			} else if record.Type() != sensors.OT_DATATYPE_UDEC_0 {
				t.Error("Expected type=OT_DATATYPE_UDEC_0")
			} else if value_, err := record.FloatValue(); err != nil {
				t.Error(err)
			} else if float64(value) != value_ {
				t.Error("Unexpected value", value_, "expected", float64(value))
			} else {
				t.Log(record)
			}
		}
	}
}

func Test_OT_022_float_udec4(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		value := float64(50.0)
		if record, err := proto.NewFloat(sensors.OT_PARAM_FREQUENCY, sensors.OT_DATATYPE_UDEC_4, value, false); err != nil {
			t.Error(err)
		} else if record.Type() != sensors.OT_DATATYPE_UDEC_4 {
			t.Error("Expected type=OT_DATATYPE_UDEC_4")
		} else if value_, err := record.FloatValue(); err != nil {
			t.Error(err)
		} else if float64(value) != value_ {
			t.Error("Unexpected value", value_, "expected", float64(value))
		} else {
			t.Log(record)
		}
	}
}
func Test_OT_023_float_udec8(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		value := float64(50.0)
		if record, err := proto.NewFloat(sensors.OT_PARAM_FREQUENCY, sensors.OT_DATATYPE_UDEC_8, value, false); err != nil {
			t.Error(err)
		} else if record.Type() != sensors.OT_DATATYPE_UDEC_8 {
			t.Error("Expected type=OT_DATATYPE_UDEC_8")
		} else if value_, err := record.FloatValue(); err != nil {
			t.Error(err)
		} else if float64(value) != value_ {
			t.Error("Unexpected value", value_, "expected", float64(value))
		} else {
			t.Log(record)
		}
	}
}

func Test_OT_024_float_udec12(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		value := float64(50.0)
		if record, err := proto.NewFloat(sensors.OT_PARAM_FREQUENCY, sensors.OT_DATATYPE_UDEC_12, value, false); err != nil {
			t.Error(err)
		} else if record.Type() != sensors.OT_DATATYPE_UDEC_12 {
			t.Error("Expected type=OT_DATATYPE_UDEC_12")
		} else if value_, err := record.FloatValue(); err != nil {
			t.Error(err)
		} else if float64(value) != value_ {
			t.Error("Unexpected value", value_, "expected", float64(value))
		} else {
			t.Log(record)
		}
	}
}

func Test_OT_025_float_udec16(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		value := float64(50.0)
		if record, err := proto.NewFloat(sensors.OT_PARAM_FREQUENCY, sensors.OT_DATATYPE_UDEC_16, value, false); err != nil {
			t.Error(err)
		} else if record.Type() != sensors.OT_DATATYPE_UDEC_16 {
			t.Error("Expected type=OT_DATATYPE_UDEC_16")
		} else if value_, err := record.FloatValue(); err != nil {
			t.Error(err)
		} else if value != value_ {
			t.Error("Unexpected value", value_, "expected", float64(value))
		} else {
			t.Log(record)
		}
	}
}
func Test_OT_025_float_udec20(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		value := float64(50.0)
		if record, err := proto.NewFloat(sensors.OT_PARAM_FREQUENCY, sensors.OT_DATATYPE_UDEC_20, value, false); err != nil {
			t.Error(err)
		} else if record.Type() != sensors.OT_DATATYPE_UDEC_20 {
			t.Error("Expected type=OT_DATATYPE_UDEC_20")
		} else if value_, err := record.FloatValue(); err != nil {
			t.Error(err)
		} else if value != value_ {
			t.Error("Unexpected value", value_, "expected", float64(value))
		} else {
			t.Log(record)
		}
	}
}

func Test_OT_025_float_udec24(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		value := float64(50.0)
		if record, err := proto.NewFloat(sensors.OT_PARAM_FREQUENCY, sensors.OT_DATATYPE_UDEC_24, value, false); err != nil {
			t.Error(err)
		} else if record.Type() != sensors.OT_DATATYPE_UDEC_24 {
			t.Error("Expected type=OT_DATATYPE_UDEC_24")
		} else if value_, err := record.FloatValue(); err != nil {
			t.Error(err)
		} else if value != value_ {
			t.Error("Unexpected value", value_, "expected", float64(value))
		} else {
			t.Log(record)
		}
	}
}

func Test_OT_025_float_dec8(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		for _, value := range []float64{-50.0, 50.0, -1, 1, 0, -9E10, 9E10} {
			if record, err := proto.NewFloat(sensors.OT_PARAM_FREQUENCY, sensors.OT_DATATYPE_DEC_8, value, false); err != nil {
				t.Error(err)
			} else if record.Type() != sensors.OT_DATATYPE_DEC_8 {
				t.Error("Expected type=OT_DATATYPE_DEC_8")
			} else if value_, err := record.FloatValue(); err != nil {
				t.Error(err)
			} else if value != value_ {
				t.Error("Unexpected value", value_, "expected", float64(value))
			} else {
				t.Log(record)
			}
		}
	}
}

func Test_OT_026_float_dec16(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		for _, value := range []float64{-50.0, 50.0, -1, 1, 0, -9E10, 9E10} {
			if record, err := proto.NewFloat(sensors.OT_PARAM_FREQUENCY, sensors.OT_DATATYPE_DEC_16, value, false); err != nil {
				t.Error(err)
			} else if record.Type() != sensors.OT_DATATYPE_DEC_16 {
				t.Error("Expected type=OT_DATATYPE_DEC_16")
			} else if value_, err := record.FloatValue(); err != nil {
				t.Error(err)
			} else if value != value_ {
				t.Error("Unexpected value", value_, "expected", float64(value))
			} else {
				t.Log(record)
			}
		}
	}
}

func Test_OT_027_float_dec24(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		for _, value := range []float64{-50.0, 50.0, -1, 1, 0, -9E10, 9E10} {
			if record, err := proto.NewFloat(sensors.OT_PARAM_FREQUENCY, sensors.OT_DATATYPE_DEC_24, value, false); err != nil {
				t.Error(err)
			} else if record.Type() != sensors.OT_DATATYPE_DEC_24 {
				t.Error("Expected type=OT_DATATYPE_DEC_24")
			} else if value_, err := record.FloatValue(); err != nil {
				t.Error(err)
			} else if value != value_ {
				t.Error("Unexpected value", value_, "expected", float64(value))
			} else {
				t.Log(record)
			}
		}
	}
}

func Test_OT_028_encode_empty(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else if msg, err := proto.New(sensors.OT_MANUFACTURER_ENERGENIE, 0xFF, 0x12345); err != nil {
		t.Fatal(err)
	} else if encoded := proto.Encode(msg); len(encoded) == 0 {
		t.Error("Expected encoded value")
	} else if len(encoded) != 11 {
		t.Error("Expected encoded to be 11 bytes, got", len(encoded))
	} else {
		t.Log(msg, "=>", strings.ToUpper(hex.EncodeToString(encoded)))
	}
}

func Test_OT_029_encode_decode_empty(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else if msg, err := proto.New(sensors.OT_MANUFACTURER_ENERGENIE, 0xFF, 0x12345); err != nil {
		t.Fatal(err)
	} else if encoded := proto.Encode(msg); len(encoded) == 0 {
		t.Error("Expected encoded value")
	} else if decoded, err := proto.Decode(encoded, time.Now()); err != nil {
		t.Error(err)
	} else if msg.IsDuplicate(decoded) == false {
		t.Error("Messages not identical", msg, " and ", decoded)
	} else {
		t.Log(msg, "=>", strings.ToUpper(hex.EncodeToString(encoded)), "=>", decoded)
	}
}

func Test_OT_030_encode_decode_join(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else if msg, err := proto.New(sensors.OT_MANUFACTURER_ENERGENIE, 0xFF, 0x12345); err != nil {
		t.Error(err)
	} else if join, err := proto.NewNull(sensors.OT_PARAM_JOIN, true); err != nil {
		t.Error(err)
	} else if join.IsReport() == false {
		t.Error("Expected report=true")
	} else if msg.Append(join); false {
		//
	} else if encoded := proto.Encode(msg); len(encoded) == 0 {
		t.Error("Expected encoded value")
	} else if decoded, err := proto.Decode(encoded, time.Now()); err != nil {
		t.Error(err)
	} else if msg.IsDuplicate(decoded) == false {
		t.Error("Messages not identical", msg, " and ", decoded)
	} else {
		t.Log(msg, "=>", strings.ToUpper(hex.EncodeToString(encoded)), "=>", decoded)
	}
}

func Test_OT_031_encode_decode_level(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else if msg, err := proto.New(sensors.OT_MANUFACTURER_ENERGENIE, 0xFF, 0x12345); err != nil {
		t.Error(err)
	} else if level, err := proto.NewUint(sensors.OT_PARAM_LEVEL, 100, true); err != nil {
		t.Error(err)
	} else if level.IsReport() == false {
		t.Error("Expected report=true")
	} else if msg.Append(level); false {
		//
	} else if encoded := proto.Encode(msg); len(encoded) == 0 {
		t.Error("Expected encoded value")
	} else if decoded, err := proto.Decode(encoded, time.Now()); err != nil {
		t.Error(err)
	} else if msg.IsDuplicate(decoded) == false {
		t.Error("Messages not identical", msg, " and ", decoded)
	} else {
		t.Log(msg, "=>", strings.ToUpper(hex.EncodeToString(encoded)), "=>", decoded)
	}
}

func Test_OT_032_encode_decode_two_records(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else if msg, err := proto.New(sensors.OT_MANUFACTURER_ENERGENIE, 0xFF, 0x12345); err != nil {
		t.Error(err)
	} else {
		if level1, err := proto.NewUint(sensors.OT_PARAM_LEVEL, 100, false); err != nil {
			t.Error(err)
		} else {
			msg.Append(level1)
		}
		if level2, err := proto.NewUint(sensors.OT_PARAM_LEVEL, 200, false); err != nil {
			t.Error(err)
		} else {
			msg.Append(level2)
		}
		if encoded := proto.Encode(msg); len(encoded) == 0 {
			t.Error("Expected encoded value")
		} else if decoded, err := proto.Decode(encoded, time.Now()); err != nil {
			t.Error(err)
		} else if msg.IsDuplicate(decoded) == false {
			t.Error("Messages not identical", msg, " and ", decoded)
		}
	}
}

func Test_OT_033_decode_encode(t *testing.T) {
	if proto := OTProto(); proto == nil {
		t.Fatal("Missing OTProto module")
	} else {
		for _, payload_str := range received_good {
			if payload, err := hex.DecodeString(payload_str); err != nil {
				t.Error(err)
			} else if msg, err := proto.Decode(payload, time.Now()); err != nil {
				t.Error(err)
			} else if encoded := proto.Encode(msg); encoded == nil {
				t.Error(err)
			} else if encoded_str := strings.ToUpper(hex.EncodeToString(encoded)); encoded_str != payload_str {
				t.Error("Expected encoded payload to match decoded payload")
			} else {
				t.Log(payload_str, "=>", encoded_str)
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
