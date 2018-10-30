package registrar

import (
	"reflect"
	"testing"
	"time"

	"github.com/jtaylorcpp/gerl/core"
)

func TestRegistrar(t *testing.T) {
	reg := New(core.GlobalScope)

	t.Log("new registrar: ", reg)

	go func() {
		t.Log("error from registrar server: ", reg.Start())
	}()

	for !core.PidHealthCheck(reg.Pid.GetAddr()) {
		time.Sleep(25 * time.Microsecond)
		t.Log("waiting for registrar to start")
	}

	added := AddRecords(reg.Pid.GetAddr(), "local", Record{
		Name:    "test",
		Address: "test:local",
		Scope:   core.LocalScope,
	})

	t.Log("added record: ", added)

	recs := GetRecords(reg.Pid.GetAddr(), "local", "test")

	if len(recs) != 1 {
		t.Fatal("only 1 record should be returned: ", recs)
	}

	if !reflect.DeepEqual(recs[0], Record{
		Name:    "test",
		Address: "test:local",
		Scope:   core.LocalScope,
	}) {
		t.Fatal("should be same record returned")
	}
}

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
