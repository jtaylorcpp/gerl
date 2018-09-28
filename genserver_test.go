package gerl

import (
	"testing"
)

func TestGenServer(t *testing.T) {
	t.Log("Starting GenServer test...")

	gs := &GenServer{
		CustomCall: testCustomCall,
		CustomCast: testCustomCast,
		BufferSize: 2,
	}

	gs.Init("test state")

	gsPid := gs.Start()

	t.Log("pid made: ", gsPid)

	t.Log("test call")
	gs.Pid.SendMsg(GerlMsg{0x0, ProcessAddr([]byte("testServer")), "test"})
	t.Log("test cast")
	gs.Pid.SendMsg(GerlMsg{0x1, ProcessAddr([]byte("testServer")), "test"})
	t.Log("test terminate")
	gs.Terminate()

}

//type GenServerCustomCall func(GenericServerMessage, ProcessAddr, GenericServerState) (GenericServerMessage, GenericServerState)
func testCustomCall(gsm GenericServerMessage, pa ProcessAddr, gss GenericServerState) (GenericServerMessage, GenericServerState) {
	return gsm, gss
}

func testCustomCast(gsm GenericServerMessage, pa ProcessAddr, gss GenericServerState) GenericServerState {
	return gss
}
