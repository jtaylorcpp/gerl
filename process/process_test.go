package process

import (
	"errors"
	"testing"
	"time"

	"gerl/core"
)

var SetValue string

func TestProcess(t *testing.T) {
	proc := New(core.LocalScope, func(pa PidAddr, in Inbox) error {
		for {
			msg, ok := <-in
			if !ok {
				t.Log("inbox closed")
				break
			}
			t.Log("recieved: ", string(msg.GetMsg().GetRawMsg()))
			if string(msg.GetMsg().GetRawMsg()) != SetValue {
				t.Fatalf("msg value<%v> not equal to <%v>\n", string(msg.GetMsg().GetRawMsg()), SetValue)
				return errors.New("not the set value")
			}
		}

		return nil
	})

	t.Log("about to start main process")

	go func() {
		t.Log(proc.Start())
	}()

	time.Sleep(25 * time.Millisecond)
	t.Log("process started with pid: ", proc.Pid.GetAddr())

	time.Sleep(25 * time.Millisecond)

	SetValue = "test1"

	Send(PidAddr(proc.Pid.GetAddr()), "localhost", core.Message{
		RawMsg: []byte(SetValue),
	})

	time.Sleep(50 * time.Millisecond)

	proc.Terminate()
}
