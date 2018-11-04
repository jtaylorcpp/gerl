package genserver

import (
	"log"
	"testing"
	"time"

	"github.com/jtaylorcpp/gerl/core"
)

func TestGenServer(t *testing.T) {
	genserver := NewGenServer("test state", core.LocalScope, CallTest, CastTest)

	msg1 := core.Message{
		Type:        core.Message_SIMPLE,
		Description: "test call",
	}

	returnMsg1 := Call(PidAddr(genserver.Pid.Addr), FromAddr("localhost"), msg1)
	t.Log(returnMsg1)

	Cast(PidAddr(genserver.Pid.GetAddr()), PidAddr("localhost"), msg1)

	time.Sleep(25 * time.Millisecond)

	t.Log("final state: ", genserver.State)

	genserver.Terminate()

}

func CallTest(_ core.Pid, msg core.Message, _ FromAddr, s State) (core.Message, State) {
	log.Println("call test func called")
	return msg, State(s.(string) + " call")
}

func CastTest(_ core.Pid, msg core.Message, _ FromAddr, s State) State {
	log.Println("cast test func called")
	return State(s.(string) + " cast")
}
