package supervisor

import (
	"gerl/core"
	"gerl/genserver"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

type TestStratMessage struct {
	Val int
}

type TestStratState struct {
	Iter int
}

var iter int = 0

type GS struct {
	gs *genserver.GenServer
}

func (g *GS) Start(started chan<- bool) error {
	iter = iter + 1
	var err error
	g.gs, err = genserver.NewGenServer(TestStratState{Iter: iter}, core.LocalScope, CallTestStrat, CastTestStrat)
	if err != nil {
		panic(err)
	}
	return g.gs.Start(started)
}

func (gs *GS) GetPid() *core.Pid {
	return gs.gs.GetPid()
}

func (g *GS) Terminate() {
	g.gs.Terminate()
}

func CallTestStrat(_ core.Pid, msg TestStratMessage, _ genserver.FromAddr, s TestStratState) (TestStratMessage, TestStratState) {
	log.Println("call test func called")
	return TestStratMessage{
		Val: s.Iter,
	}, s
}

func CastTestStrat(_ core.Pid, msg TestStratMessage, _ genserver.FromAddr, s TestStratState) TestStratState {
	log.Println("cast test func called")
	return s
}

func TestSupervisorOneForOne(t *testing.T) {
	iter = 0
	gs := &GS{}
	child := NewChild("test", gs, RESTART_ALWAYS)
	child2 := NewChild("test2", &GS{}, RESTART_ALWAYS)

	sup := NewSupervisor([]*Child{child, child2}, ONE_FOR_ONE)

	started := make(chan bool, 1)
	errorChan := make(chan error, 1)

	go func() {
		t.Log("starting supervisor")
		errorChan <- sup.Start(started)
	}()
	<-started
	t.Log("supervisor started")

	returnMsg, err := genserver.Call(gs.gs.Pid.GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal(err.Error())
	}

	currentState := returnMsg.(TestStratMessage).Val
	t.Logf("current process state is: %v\n", currentState)

	returnMsg, err = genserver.Call(gs.gs.Pid.GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if returnMsg.(TestStratMessage).Val != currentState {
		t.Fatal("the state should not have been incremented")
	}

	child.Process.Terminate()

	for child.Process.GetPid() == nil {
		t.Log("waiting for child process to be restarted")
		time.Sleep(1 * time.Second)
	}

	for core.PidHealthCheck(child.Process.GetPid().GetAddr()) == false {
		t.Log("waiting for child process to be healthy")
		time.Sleep(1 * time.Second)
	}

	restartedReturnMsg, err := genserver.Call(child.Process.GetPid().GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if restartedReturnMsg.(TestStratMessage).Val != currentState+2 {
		t.Fatalf("the state should have incremented: %v\n", restartedReturnMsg.(TestStratMessage).Val)
	}

	child2ReturnMsg, err := genserver.Call(child2.Process.GetPid().GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal("failed with error: ", err.Error())
	}

	if child2ReturnMsg.(TestStratMessage).Val != currentState+1 {
		t.Fatal("second child should not have been restarted")
	}

	sup.Terminate()
}

func TestSupervisorOneForAll(t *testing.T) {
	iter = 0
	child := NewChild("test", &GS{}, RESTART_ALWAYS)
	child2 := NewChild("test2", &GS{}, RESTART_ALWAYS)

	sup := NewSupervisor([]*Child{child, child2}, ONE_FOR_ONE)

	started := make(chan bool, 1)
	errorChan := make(chan error, 1)

	go func() {
		t.Log("starting supervisor")
		errorChan <- sup.Start(started)
	}()
	<-started
	t.Log("supervisor started")

	returnMsg, err := genserver.Call(child.Process.GetPid().GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if returnMsg.(TestStratMessage).Val != 1 {
		t.Log("child should be iter number 1")
	}

	returnMsg, err = genserver.Call(child2.Process.GetPid().GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if returnMsg.(TestStratMessage).Val != 2 {
		t.Log("child2 should be iter number 2")
	}

	// since child/child2 will always have states {1,2} * iteration, terminating either should follow the pattern
	child.Terminate()

	for child.Process.GetPid() == nil {
		t.Log("waiting for child to be restarted")
		time.Sleep(1 * time.Second)
	}

	for core.PidHealthCheck(child.Process.GetPid().GetAddr()) == false {
		t.Log("waiting for child process to be healthy")
		time.Sleep(1 * time.Second)
	}

	for child2.Process.GetPid() == nil {
		t.Log("waiting for child to be restarted")
		time.Sleep(1 * time.Second)
	}

	for core.PidHealthCheck(child2.Process.GetPid().GetAddr()) == false {
		t.Log("waiting for child process to be healthy")
		time.Sleep(1 * time.Second)
	}

	returnMsg, err = genserver.Call(child.Process.GetPid().GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if returnMsg.(TestStratMessage).Val != 3 {
		t.Log("child should be iter number 1")
	}

	returnMsg, err = genserver.Call(child2.Process.GetPid().GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if returnMsg.(TestStratMessage).Val != 4 {
		t.Log("child2 should be iter number 2")
	}
}
