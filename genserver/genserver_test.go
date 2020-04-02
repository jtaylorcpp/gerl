package genserver

import (
	"testing"
	"time"

	//"time"
	"encoding/json"

	"gerl/core"

	log "github.com/sirupsen/logrus"
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
	config := &GenServerV2Config{
		StartState:  TestState{"test state"},
		Scope:       core.LocalScope,
		CallHandler: CallTest,
		CastHandler: CastTest,
	}

	genserver, err := NewGenServerV2(config)
	if err != nil {
		t.Fatal(err.Error())
	}

	go func() {
		log.Errorln("error from running genserver: ", genserver.Start())
	}()

	for !genserver.IsReady() {
		time.Sleep(1 * time.Second)
	}

	msg1 := TestMessage{
		Body: "test1",
	}

	returnMsg1, err := Call(genserver.pid.GetAddr(), "localhost", msg1)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(returnMsg1)

	genserver.Terminate()
	log.Println("genserver terminated")
}

func BenchmarkGenServerStart(b *testing.B) {
	for i := 0; i < b.N; i++ {
		config := &GenServerV2Config{
			StartState:  TestState{"test state"},
			Scope:       core.LocalScope,
			CallHandler: CallTest,
			CastHandler: CastTest,
		}

		genserver, err := NewGenServerV2(config)
		if err != nil {
			b.Fatal(err.Error())
		}

		go func() {
			log.Errorln("error from genserver: ", genserver.Start())
		}()
		for !genserver.IsReady() {
			// do nothing
		}

		genserver.Terminate()
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
