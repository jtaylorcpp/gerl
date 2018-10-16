package process

import (
	"errors"
	"testing"
	"time"

	"github.com/jtaylorcpp/gerl/core"
)

var SetValue string

func TestProcess(t *testing.T) {
	proc := New(func(pa PidAddr, in Inbox) error {
		for {
			msg, ok := <-in
			if !ok {
				t.Log("inbox closed")
			}
			t.Log("recieved desc: ", msg.GetMsg().GetDescription())
			if msg.GetMsg().GetDescription() != SetValue {
				t.Fatalf("msg value<%v> not equal to <%v>\n", msg.GetMsg().GetDescription(), SetValue)
				return errors.New("not the set value")
			}
		}

		return nil
	})

	t.Log("about to start main process")

	go func() {
		t.Log(proc.Start())
	}()

	t.Log("process started with pid: ", proc.Pid.GetAddr())

	time.Sleep(25 * time.Millisecond)

	SetValue = "test1"

	Send(PidAddr(proc.Pid.GetAddr()), "localhost", core.Message{
		Type:        core.Message_SIMPLE,
		Description: SetValue,
	})

	time.Sleep(50 * time.Millisecond)

	proc.Terminate()

}
