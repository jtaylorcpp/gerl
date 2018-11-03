package registrar

import (
	"reflect"
	"testing"

	"github.com/jtaylorcpp/gerl/core"
)

func TestRegister(t *testing.T) {
	rec := Record{
		Name:    "test",
		Address: "1.1.1.1",
		Scope:   core.LocalScope,
	}

	reg := newRegister()

	reg = reg.addRecords(rec)

	if _, ok := reg.recordmap[rec.Name]; ok {
		if rectest, ok2 := reg.recordmap[rec.Name][rec.Address]; ok2 {
			if !reflect.DeepEqual(rec, rectest) {
				t.Fatal("records not the same")
			}
		} else {
			t.Fatal("record not found for addr: ", rec.Address)
		}
	} else {
		t.Fatal("svc not in register: ", rec.Name, reg)
	}

	getReg := reg.getRecords("test")

	t.Log("got records: ", getReg)

	if len(getReg) != 1 {
		t.Fatal("did not get one record back")
	}

	if !reflect.DeepEqual(rec, getReg[0]) {

	}

}
