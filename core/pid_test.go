package core

import (
	"reflect"
	"testing"
)

func TestCall(t *testing.T) {
	pid := NewPid("", "")

	testMsg := Message{
		Type:        Message_SIMPLE,
		Description: "call test",
	}

	testGerl := GerlMsg{
		Type:     GerlMsg_CALL,
		Fromaddr: "calladdr",
		Msg:      &testMsg,
	}

	pid.Outbox <- testGerl

	returnMsg := PidCall(pid.GetAddr(), "calladdr", testMsg)

	t.Log("message sent to call: ", testMsg)
	t.Log("message returned by call: ", returnMsg)

	if reflect.DeepEqual(testGerl, <-pid.Inbox) {
		t.Fatal("message sent to inbox not the same")
	}

	if (testMsg.Description != returnMsg.Description) || (testMsg.Type != returnMsg.Type) {
		t.Fatal("message description and type are not equal")
	}

	pid.Terminate()
}

func TestCast(t *testing.T) {
	pid := NewPid("", "")

	testMsg := Message{
		Type:        Message_SIMPLE,
		Description: "cast test",
	}

	testGerl := GerlMsg{
		Type:     GerlMsg_CALL,
		Fromaddr: "castaddr",
		Msg:      &testMsg,
	}

	pid.Outbox <- testGerl

	PidCast(pid.GetAddr(), "castaddr", testMsg)

	t.Log("message sent to cast: ", testMsg)

	if reflect.DeepEqual(testGerl, <-pid.Inbox) {
		t.Fatal("message sent to inbox not the same")
	}

	pid.Terminate()
}

func TestHealth(t *testing.T) {
	pid := NewPid("", "")
	addr := pid.GetAddr()

	health := PidHealthCheck(addr)

	t.Log("health check recieved: ", health)

	if health != true {
		t.Fatal("pid not healthy")
	}

	pid.Terminate()

	health2 := PidHealthCheck(addr)

	t.Log("health from terminated pid: ", health2)

	if health2 != false {
		t.Fatal("pid still healthy")
	}
}
