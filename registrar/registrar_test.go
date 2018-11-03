package registrar

import (
	"reflect"
	"testing"
	"time"

	"github.com/jtaylorcpp/gerl/core"
	"github.com/jtaylorcpp/gerl/genserver"
)

func TestRegistrar(t *testing.T) {
	reg := NewRegistrar(core.GlobalScope)

	t.Log("new registrar: ", reg)

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

func TestRegistrarRefresh(t *testing.T) {

	REFRESH_TIMER = 100 * time.Millisecond

	reg1 := NewRegistrar(core.GlobalScope)
	t.Log("new registrar: ", reg1)
	reg2 := NewRegistrar(core.GlobalScope)
	t.Log("new registrar: ", reg2)

	JoinRegistrar(genserver.FromAddr(reg1.Pid.GetAddr()),
		genserver.PidAddr(reg2.Pid.GetAddr()))

	time.Sleep(1 * time.Second)

	reg1.Terminate()
}
