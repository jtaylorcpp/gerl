package process

import (
	"errors"
	"testing"
	"time"

	"gerl/core"
)

var SetValue string

func handler(addr PidAddr, msg core.Message) error {
	switch string(msg.GetRawMsg()) {
	case "term":
		return errors.New("term")
	default:
		SetValue = string(msg.GetRawMsg())
	}

	return nil
}
func TestProcess(t *testing.T) {
	config := &ProcessConfig{
		Scope:   core.LocalScope,
		Handler: handler,
	}
	proc := New(config)

	t.Log("about to start main process")

	procError := make(chan error, 1)
	go func() {
		procError <- proc.Start()
	}()

	for !proc.IsReady() {
		time.Sleep(1 * time.Second)
	}
	t.Log("process started with pid: ", proc.Pid.GetAddr())

	Send(PidAddr(proc.Pid.GetAddr()), "localhost", core.Message{
		RawMsg: []byte("test msg"),
	})

	time.Sleep(50 * time.Millisecond)

	if SetValue != "test msg" {
		t.Fatal("set values didnt match")
	}

	proc.Terminate()
}

func TestProcessError(t *testing.T) {
	config := &ProcessConfig{
		Scope:   core.LocalScope,
		Handler: handler,
	}
	proc := New(config)

	t.Log("about to start main process")

	procError := make(chan error, 1)
	go func() {
		procError <- proc.Start()
	}()

	for !proc.IsReady() {
		time.Sleep(1 * time.Second)
	}
	t.Log("process started with pid: ", proc.Pid.GetAddr())

	Send(PidAddr(proc.Pid.GetAddr()), "localhost", core.Message{
		RawMsg: []byte("test msg"),
	})

	time.Sleep(50 * time.Millisecond)

	if SetValue != "test msg" {
		t.Fatal("set values didnt match")
	}

	Send(PidAddr(proc.Pid.GetAddr()), "localhost", core.Message{
		RawMsg: []byte("term"),
	})

	time.Sleep(50 * time.Millisecond)

	msg := <-procError

	if msg.Error() != "term" {
		t.Fatal("error message propogation failed: ", msg.Error())
	}

	if proc.Pid != nil {
		t.Fatal("process did not close out properly")
	}
}
