package supervisor

import (
	"gerl/core"
	"gerl/genserver"
	"log"
	"testing"
	"time"
)

func TestEmptySupervisor(t *testing.T) {
	s := NewSupervisor([]Child{}, ONE_FOR_ALL)

	if len(s.Children) != 0 {
		t.Fatal("too many children")
	}

	if s.Strategy != ONE_FOR_ALL {
		t.Fatal("correct strategy not saved")
	}
}

func TestChildRestartNever(t *testing.T) {
	gserver, err := genserver.NewGenServer(TestState{"test state"}, core.LocalScope, CallTest, CastTest)
	if err != nil {
		t.Fatal(err.Error())
	}

	child := Child{
		Name:            "test",
		Process:         gserver,
		RestartStrategy: RESTART_NEVER,
	}

	childErr := make(chan ChildFailure, 1)

	go func() {
		childErr <- child.Start()
	}()

	if len(childErr) != 0 {
		t.Fatal("child exited prematurely")
	}

	child.Process.Terminate()

	// give time to propogate terminate to child
	time.Sleep(50 * time.Microsecond)

	if len(childErr) != 1 {
		t.Fatal("process did not exit properly")
	}

	exitMsg := <-childErr
	if exitMsg.Name != "test" {
		t.Fatal("child error does not have correct name")
	}

	if exitMsg.Error != nil {
		log.Fatal("child process should have exited cleanly")
	}
}

func TestChildRestartOnce(t *testing.T) {
	gserver, err := genserver.NewGenServer(TestState{"test state"}, core.LocalScope, CallTest, CastTest)
	if err != nil {
		t.Fatal(err.Error())
	}

	childErr := make(chan ChildFailure, 1)
	child := NewChild("test2", gserver, RESTART_ONCE)
	go func() {
		childErr <- child.Start()
	}()

	if len(childErr) != 0 {
		t.Fatal("child closed out prematurely")
	}

	child.Process.Terminate()
	time.Sleep(50 * time.Microsecond)

	// child should restart the process

	if len(childErr) != 0 {
		t.Fatal("child should not close after single process termination")
	}

	child.Process.Terminate()
	time.Sleep(50 * time.Microsecond)

	// child should not restart the process

	if len(childErr) != 1 {
		t.Fatal("child should close after second process termiantion")
	}

	/*
		t.Log("waiting for genserver to start")

		msg1 := TestMessage{
			Body: "test1",
		}

		returnMsg1, err := Call(genserver.Pid.Addr, "localhost", msg1)
		if err != nil {
			t.Fatal(err.Error())
		}
		t.Log(returnMsg1)

		//time.Sleep(500 * time.Microsecond)
		t.Log("waiting for genserver to clear inbox before terminate")

		genserver.Terminate()

		<-genserverStopped
	*/
}

type TestMessage struct {
	Body string
}

type TestState struct {
	Some string
}

func CallTest(_ core.Pid, msg TestMessage, _ genserver.FromAddr, s TestState) (TestMessage, TestState) {
	log.Println("call test func called")
	return msg, s
}

func CastTest(_ core.Pid, msg TestMessage, _ genserver.FromAddr, s TestState) TestState {
	log.Println("cast test func called")
	return s
}
