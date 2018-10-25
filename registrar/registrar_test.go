package registrar

import (
	"reflect"
	"testing"

	"github.com/jtaylorcpp/gerl/core"
)

func TestRegistrar(t *testing.T) {
	reg := New(core.GlobalScope)

	t.Log("new registrar: ", reg)
}

func TestRegister(t *testing.T) {
	rec := record{
		name:    "test",
		address: "1.1.1.1",
		scope:   core.LocalScope,
	}

	reg := newRegister()

	reg = reg.addRecords(rec)

	if _, ok := reg.recordmap[rec.name]; ok {
		if rectest, ok2 := reg.recordmap[rec.name][rec.address]; ok2 {
			if !reflect.DeepEqual(rec, rectest) {
				t.Fatal("records not the same")
			}
		} else {
			t.Fatal("record not found for addr: ", rec.address)
		}
	} else {
		t.Fatal("svc not in register: ", rec.name, reg)
	}
}
