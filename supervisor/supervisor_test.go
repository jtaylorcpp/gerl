package supervisor

import (
	"gerl/core"
	"gerl/genserver"
	"log"
	"testing"
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
	child := Child{
		Name:            "test",
		Process:         gserver,
		RestartStrategy: RESTART_ONCE,
	}

	started := make(chan bool, 1)
	errorChan := make(chan error, 1)
	go func() {
		errorChan <- child.restartOnce(started)
	}()
	<-started

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

/*func TestChildTerminate(t *testing.T) {
	gserver, err := genserver.NewGenServer(TestState{"test state"}, core.LocalScope, CallTest, CastTest)
	if err != nil {
		t.Fatal(err.Error())
	}

	child := Child{
		Name:            "test",
		Process:         gserver,
		RestartStrategy: RESTART_NEVER,
	}

	started := make(chan bool, 1)
	stopped := make(chan bool, 1)
	log.Println("starting child")
	go func() {
		childFail := child.Start(started)
		if childFail.Error != nil {
			t.Fatal(childFail.Error.Error())
		}
		log.Println("child closed out in test")
		stopped <- true
	}()
	<-started
	log.Println("terminating child")
	child.Terminate()
	<-stopped
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
	childStart := make(chan bool, 1)
	go func() {
		childErr <- child.Start(childStart)
	}()
	<-childStart

	if len(childErr) != 0 {
		t.Fatal("child exited prematurely")
	}

	child.Process.Terminate()

	// give time to propogate terminate to child
	<-child.TerminateOut

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
	childStarted := make(chan bool, 1)
	child := NewChild("test2", gserver, RESTART_ONCE)
	go func() {
		childErr <- child.Start(childStarted)
	}()
	<-childStarted

	if len(childErr) != 0 {
		t.Fatal("child closed out prematurely")
	}

	child.Process.Terminate()
	time.Sleep(100 * time.Millisecond)
	t.Log(len(childStarted))
	/*
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
}*/

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
