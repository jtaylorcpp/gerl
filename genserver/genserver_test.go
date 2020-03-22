package genserver

import (
	"log"
	"testing"

	//"time"
	"encoding/json"

	"gerl/core"
)

func TestGenServerCallHandlerParsing(t *testing.T) {
	handler, err := newCustomCallHandler(CallTest)
	if err != nil {
		t.Fatal(err.Error())
	}

	testMsg := TestMessage{
		Body: "hello world",
	}

	testState := TestState{
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

	_, err = newCustomCallHandler(func() {})
	if err == nil {
		t.Fatal("this func is not right and should have errored out")
	}

	t.Log(err.Error())

	_, err = newCustomCallHandler(func(a, b, c, d int) (int, int) { return 0, 0 })
	if err == nil {
		t.Fatal("this func is not right and should have errored out")
	}
	t.Log(err.Error())

	_, err = newCustomCallHandler(func(_ core.Pid, b, c, d int) (int, int) { return 0, 0 })
	if err == nil {
		t.Fatal("this func is not right and should have errored out")
	}
	t.Log(err.Error())

	_, err = newCustomCallHandler(func(_ core.Pid, _ struct{}, c, d int) (int, int) { return 0, 0 })
	if err == nil {
		t.Fatal("this func is not right and should have errored out")
	}
	t.Log(err.Error())

	_, err = newCustomCallHandler(func(_ core.Pid, _ struct{}, _ string, d int) (int, int) { return 0, 0 })
	if err == nil {
		t.Fatal("this func is not right and should have errored out")
	}
	t.Log(err.Error())

	_, err = newCustomCallHandler(func(_ core.Pid, _ struct{}, _ string, _ struct{}) (int, int) { return 0, 0 })
	if err == nil {
		t.Fatal("this func is not right and should have errored out")
	}
	t.Log(err.Error())

	_, err = newCustomCallHandler(func(_ core.Pid, _ struct{}, _ string, _ struct{}) (struct{}, int) { return struct{}{}, 0 })
	if err == nil {
		t.Fatal("this func is not right and should have errored out")
	}
	t.Log(err.Error())
}

func TestGenServerCastHandlerParsing(t *testing.T) {
	handler, err := newCustomCastHandler(CastTest)
	if err != nil {
		t.Fatal(err.Error())
	}

	testMsg := TestMessage{
		Body: "hello world",
	}

	testState := TestState{
		Some: "details",
	}

	rawMsg, err := json.Marshal(&testMsg)
	if err != nil {
		t.Fatal(err.Error())
	}

	rState := handler(core.Pid{}, rawMsg, "", testState)

	stateAssert := rState.(TestState)

	if stateAssert.Some != testState.Some {
		t.Fatal("state returned is not equal")
	}

	_, err = newCustomCallHandler(func() {})
	if err == nil {
		t.Fatal("this func is not right and should have errored out")
	}

	t.Log(err.Error())

	_, err = newCustomCallHandler(func(a, b, c, d int) (int, int) { return 0, 0 })
	if err == nil {
		t.Fatal("this func is not right and should have errored out")
	}
	t.Log(err.Error())

	_, err = newCustomCallHandler(func(_ core.Pid, b, c, d int) (int, int) { return 0, 0 })
	if err == nil {
		t.Fatal("this func is not right and should have errored out")
	}
	t.Log(err.Error())

	_, err = newCustomCallHandler(func(_ core.Pid, _ struct{}, c, d int) (int, int) { return 0, 0 })
	if err == nil {
		t.Fatal("this func is not right and should have errored out")
	}
	t.Log(err.Error())

	_, err = newCustomCallHandler(func(_ core.Pid, _ struct{}, _ string, d int) (int, int) { return 0, 0 })
	if err == nil {
		t.Fatal("this func is not right and should have errored out")
	}
	t.Log(err.Error())

	_, err = newCustomCallHandler(func(_ core.Pid, _ struct{}, _ string, _ struct{}) (int, int) { return 0, 0 })
	if err == nil {
		t.Fatal("this func is not right and should have errored out")
	}
	t.Log(err.Error())

	_, err = newCustomCallHandler(func(_ core.Pid, _ struct{}, _ string, _ struct{}) (struct{}, int) { return struct{}{}, 0 })
	if err == nil {
		t.Fatal("this func is not right and should have errored out")
	}
	t.Log(err.Error())
}

func TestGenServer(t *testing.T) {
	genserver, err := NewGenServer(TestState{"test state"}, core.LocalScope, CallTest, CastTest)
	if err != nil {
		t.Fatal(err.Error())
	}

	genserverStarted := make(chan bool, 1)
	genserverStopped := make(chan bool, 1)
	go func() {
		if err := genserver.Start(genserverStarted); err != nil {
			genserverStopped <- false
			t.Fatal(err.Error())
		}

		genserverStopped <- true
	}()

	<-genserverStarted

	//time.Sleep(25 * time.Microsecond)
	t.Log("waiting for genserver to start")

	msg1 := TestMessage{
		Body: "test1",
	}

	returnMsg1, err := Call(genserver.Pid.Addr, "localhost", msg1)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(returnMsg1)

	/*Cast(PidAddr(genserver.Pid.GetAddr()), PidAddr("localhost"), msg1)

	time.Sleep(25 * time.Millisecond)

	t.Log("final state: ", genserver.State)*/

	//time.Sleep(500 * time.Microsecond)
	t.Log("waiting for genserver to clear inbox before terminate")

	genserver.Terminate()

	<-genserverStopped

	t.Log("test genserver has nil proc")

	if genserver.GetPid() != nil {
		t.Fatalf("genserver proc should be nil but is %v\n", genserver.GetPid())

	}
}

func BenchmarkGenServerStart(b *testing.B) {
	for i := 0; i < b.N; i++ {
		genserver, err := NewGenServer(TestState{"test state"}, core.LocalScope, CallTest, CastTest)
		if err != nil {
			b.Fatal(err.Error())
		}

		genserverStarted := make(chan bool, 1)
		genserverStopped := make(chan bool, 1)
		go func() {
			if err := genserver.Start(genserverStarted); err != nil {
				genserverStopped <- false
				b.Fatal(err.Error())
			}

			genserverStopped <- true
		}()
		<-genserverStarted
		genserver.Terminate()
		<-genserverStopped
	}
}

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

func CastTest(_ core.Pid, msg TestMessage, _ FromAddr, s TestState) TestState {
	log.Println("cast test func called")
	return s
}
