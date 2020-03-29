package supervisor

import (
	"gerl/core"
	"gerl/genserver"
	"log"
	"testing"
)

func TestEmptySupervisor(t *testing.T) {
	s := NewSupervisor([]*Child{}, ONE_FOR_ALL)

	if len(s.Children) != 0 {
		t.Fatal("too many children")
	}

	if s.Strategy != ONE_FOR_ALL {
		t.Fatal("correct strategy not saved")
	}
}

func TestRestartNever(t *testing.T) {
	gserver, err := genserver.NewGenServer(TestState{"test state"}, core.LocalScope, CallTest, CastTest)
	if err != nil {
		t.Fatal(err.Error())
	}

	// close process directly
	child := Child{
		Name:            "test",
		Process:         gserver,
		RestartStrategy: RESTART_NEVER,
	}

	started := make(chan bool, 1)
	errorChan := make(chan error, 1)
	go func() {
		errorChan <- child.restartNever(started)
	}()
	<-started

	child.Process.Terminate()

	err = <-errorChan
	switch err {
	case nil:
		t.Log("process terminated with no error")
	default:
		t.Fatalf("process terminated with error: %s\n", err.Error())
	}
}

func TestStartRestartNever(t *testing.T) {
	gserver, err := genserver.NewGenServer(TestState{"test state"}, core.LocalScope, CallTest, CastTest)
	if err != nil {
		t.Fatal(err.Error())
	}

	// close process directly
	child := NewChild("test", gserver, RESTART_NEVER)
	started := make(chan bool, 1)
	errorChan := make(chan error, 1)
	go func() {
		errorChan <- child.restartNever(started)
	}()
	<-started

	child.Process.Terminate()

	err = <-errorChan
	switch err {
	case nil:
		t.Log("process terminated with no error")
	default:
		t.Fatalf("process terminated with error: %s\n", err.Error())
	}
}

func TestRestartOnce(t *testing.T) {
	gserver, err := genserver.NewGenServer(TestState{"test state"}, core.LocalScope, CallTest, CastTest)
	if err != nil {
		t.Fatal(err.Error())
	}

	// close process directly
	child := NewChild("test", gserver, RESTART_ONCE)

	started := make(chan bool, 1)
	errorChan := make(chan error, 1)
	go func() {
		errorChan <- child.restartOnce(started)
	}()
	<-started
	t.Log("process started")

	t.Log("terminating process")
	child.Process.Terminate()

	if len(errorChan) != 0 {
		t.Fatal("should be no errors for child process termination")
	}

	t.Log("waiting for child proc to restart")
	<-started
	t.Log("child proc restarted")
	child.Process.Terminate()

	t.Log("error chan len: ", len(errorChan))
	err = <-errorChan
	if err != nil {
		t.Fatal("error should be nil: ", err.Error())
	}
}

func TestRestartAlways(t *testing.T) {
	gserver, err := genserver.NewGenServer(TestState{"test state"}, core.LocalScope, CallTest, CastTest)
	if err != nil {
		t.Fatal(err.Error())
	}

	// close process directly
	child := NewChild("test", gserver, RESTART_ALWAYS)

	started := make(chan bool, 1)
	errorChan := make(chan error, 1)
	go func() {
		errorChan <- child.restartAlways(started)
	}()

	for i := 0; i < 20; i++ {
		<-started
		child.Process.Terminate()
	}

	if len(errorChan) != 0 {
		t.Fatal("error chan should be empty")
	}

	child.Terminate()
	t.Log("error chan len: ", len(errorChan))
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
