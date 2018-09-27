package gerl

import (
	"testing"
)

func TestGenServer(t *testing.T) {
	t.Log("Starting GenServer test...")

	gs := &GenServer{}
	gsPid := gs.Init("test state")

	t.Log("pid made: ", gsPid)

	gs.Pid.MsgBox <- GerlMsg{0x0, "test"}
	gsPid.MsgBox <- GerlMsg{0x1, "test"}

	gs.Terminate()

}
