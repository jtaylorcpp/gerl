package genserver

import (
	"log"
	"testing"
	"time"

	"github.com/jtaylorcpp/gerl/core"
)

func TestGenServer(t *testing.T) {
	genserver := NewGenServer("test state", CallTest, CastTest)
	go func() {
		t.Log("genserver start error ", genserver.Start())
	}()

	for !genserver.Pid.Running {
		time.Sleep(25 * time.Microsecond)
		t.Log("waiting for genserver to start")
	}

	genserver.Terminate()

}

func CallTest(msg core.Message, pa PidAddr, s State) (core.Message, State) {
	log.Println("call test func called")
	return msg, s
}

func CastTest(msg core.Message, pa PidAddr, s State) State {
	log.Println("cast test func called")
	return s
}
