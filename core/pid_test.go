package core

import (
	"reflect"
	"testing"
)

func TestIPSelection(t *testing.T) {
	ip := getPublicIP()
	t.Log("selected ip: ", ip)
}

func TestCall(t *testing.T) {
	pid := NewPid("", "", LocalScope)

	testMsg := Message{
		RawMsg: []byte("call test"),
	}

	testGerl := GerlMsg{
		Type:     GerlMsg_CALL,
		Fromaddr: "calladdr",
		Msg:      &testMsg,
	}

	pid.Outbox <- testGerl

	returnMsg := PidCall(pid.GetAddr(), "calladdr", testMsg)

	t.Logf("message sent to call: %#v\n", testMsg)
	t.Logf("message returned by call:%#v\n", returnMsg)

	if reflect.DeepEqual(testGerl, <-pid.Inbox) {
		t.Fatal("message sent to inbox not the same")
	}

	pid.Terminate()
}


func BenchmarkCall(b *testing.B) {
	pid := NewPid("", "", LocalScope)

	testMsg := Message{
		RawMsg: []byte("call test"),
	}

	testGerl := GerlMsg{
		Type:     GerlMsg_CALL,
		Fromaddr: "calladdr",
		Msg:      &testMsg,
	}

	for i := 0; i < b.N; i++ {
		pid.Outbox <- testGerl
		returnMsg := PidCall(pid.GetAddr(), "calladdr", testMsg)

		<-pid.Inbox

		if !reflect.DeepEqual(returnMsg, returnMsg) {
			b.Fatal("messages did not match on ")
		}

	}

	pid.Terminate()
}


func TestCast(t *testing.T) {
	pid := NewPid("", "", LocalScope)

	testMsg := Message{
		RawMsg: []byte("cast test"),
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
	pid := NewPid("", "", LocalScope)
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

func TestHealthPublic(t *testing.T) {
	pid := NewPid("", "", GlobalScope)
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

func BenchmarkHealth(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pid := NewPid("", "", LocalScope)

		for {
			health := PidHealthCheck(pid.GetAddr())
			if health == true {
				break
			}
		}

		pid.Terminate()
	}
}