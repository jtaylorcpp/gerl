package genserver

import (
	"log"
	"testing"
	//"time"
	"encoding/json"

	"gerl/core"
)

func TestGenServerCallHandlerParsing(t *testing.T){
	handler, err := newCustomCallHandler(CallTest)
	if err != nil {
		t.Fatal(err.Error())
	}

	testMsg := TestMessage {
		Body: "hello world",
	}

	testState := TestState {
		Some: "details",
	}

	rawMsg, err := json.Marshal(&testMsg)
	if err != nil {
		t.Fatal(err.Error())
	}

	rMsg, rState := handler(core.Pid{}, rawMsg, "", testState)

	stateAssert := rState.(TestState)

	if stateAssert.Some != testState.Some {
		t.Fatal("state returned is not equal")
	}

	returnedMsg := &TestMessage{}
	err = json.Unmarshal(rMsg, returnedMsg)
	if err != nil {
		t.Fatal(err.Error())
	}

	if returnedMsg.Body != testMsg.Body {
		t.Fatal("returned message is not equal to sent message")
	}
}
/*func TestGenServer(t *testing.T) {
	genserver := NewGenServer("test state", core.LocalScope, CallTest, CastTest)
	go func() {
		t.Log("genserver start error ", genserver.Start())
	}()

	for !core.PidHealthCheck(genserver.Pid.GetAddr()) {
		time.Sleep(25 * time.Microsecond)
		t.Log("waiting for genserver to start")
	}

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
*/

type TestMessage struct {
	Body string
}

type TestState struct {
	Some string
}

func CallTest(_ core.Pid, msg TestMessage, _ FromAddr, s TestState) (TestMessage, TestState) {
	log.Println("call test func called")
	return msg, s
}

func CastTest(_ core.Pid, msg TestMessage, _ FromAddr, s State) State {
	log.Println("cast test func called")
	return State(s.(string) + " cast")
}

