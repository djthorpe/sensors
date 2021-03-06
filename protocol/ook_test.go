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
	_ "github.com/djthorpe/sensors/protocol/ook"
)

func Test_OOK_000(t *testing.T) {
	// Create an OOK module
	if app, err := gopi.NewAppInstance(gopi.NewAppConfig("sensors/protocol/ook")); err != nil {
		t.Fatal(err)
	} else if _, ok := app.ModuleInstance("sensors/protocol/ook").(sensors.OOKProto); ok == false {
		t.Fatal("OOK does not comply to OOK interface")
	}
}

func Test_OOK_001(t *testing.T) {
	if ook := OOK(); ook == nil {
		t.Fatal("Missing OOK module")
	} else if _, err := ook.New(0x12345, 0, false, nil); err != nil {
		t.Fatal(err)
	}
}

func Test_OOK_002(t *testing.T) {
	if ook := OOK(); ook == nil {
		t.Fatal("Missing OOK module")
	} else if _, err := ook.New(0x112345, 0, false, nil); err == nil {
		t.Fatal("Expected parameter error due to bad address")
	} else if _, err := ook.New(0x12345, 5, false, nil); err == nil {
		t.Fatal("Expected parameter error due to bad socket")
	}
}

func Test_OOK_003(t *testing.T) {
	if ook := OOK(); ook == nil {
		t.Fatal("Missing OOK module")
	} else {
		for addr := uint32(0); addr <= uint32(0xFFFFF); addr += 0x1234 {
			for socket := uint(0); socket < uint(5); socket++ {
				if off, err := ook.New(addr, socket, false, nil); err != nil {
					t.Error(err)
				} else {
					t.Log(off)
					if off.Addr() != addr {
						t.Error("Unexpected addr")
					}
					if off.State() != false {
						t.Error("Unexpected state")
					}
					if off.Socket() != socket {
						t.Error("Unexpected socket")
					}
				}
				if on, err := ook.New(addr, socket, true, nil); err != nil {
					t.Error(err)
				} else {
					t.Log(on)
					if on.Addr() != addr {
						t.Error("Unexpected addr")
					}
					if on.State() != true {
						t.Error("Unexpected state")
					}
					if on.Socket() != socket {
						t.Error("Unexpected socket")
					}
				}
			}
		}
	}
}

func Test_OOK_004(t *testing.T) {
	if ook := OOK(); ook == nil {
		t.Fatal("Missing OOK module")
	} else if msg, err := ook.New(0x789AB, 0, true, nil); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("message=%v", msg)
		t.Logf("  payload=%v", strings.ToUpper(hex.EncodeToString(ook.Encode(msg))))
	}
}

func Test_OOK_005(t *testing.T) {
	if ook := OOK(); ook == nil {
		t.Fatal("Missing OOK module")
	} else if msg, err := ook.New(0x789AB, 0, true, nil); err != nil {
		t.Fatal(err)
	} else if msg_out, err := ook.Decode(ook.Encode(msg), time.Time{}); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("  in=%v", msg)
		t.Logf("  payload=%v", strings.ToUpper(hex.EncodeToString(ook.Encode(msg))))
		t.Logf("  out=%v", msg_out)
	}
}
func Test_OOK_006(t *testing.T) {
	if ook := OOK(); ook == nil {
		t.Fatal("Missing OOK module")
	} else {
		for addr := uint32(0); addr < uint32(0xFFFFF); addr += uint32(0x245) {
			t.Logf("addr=0x%05X", addr)
			if msg_in, err := ook.New(addr, uint(addr%5), uint(addr%2) == 0, nil); err != nil {
				t.Fatal(err)
			} else if msg_out, err := ook.Decode(ook.Encode(msg_in), time.Time{}); err != nil {
				t.Fatal(err)
			} else if Equals(msg_in, msg_out) == false {
				t.Errorf("Messages don't match: %v and %v", msg_in, msg_out)
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// OOK

func OOK() sensors.OOKProto {
	if app, err := gopi.NewAppInstance(gopi.NewAppConfig("sensors/protocol/ook")); err != nil {
		return nil
	} else if ook, ok := app.ModuleInstance("sensors/protocol/ook").(sensors.OOKProto); ok == false {
		return nil
	} else {
		return ook
	}
}

func Equals(m1, m2 sensors.Message) bool {
	return m1.IsDuplicate(m2)
}
