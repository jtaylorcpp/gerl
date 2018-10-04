package genserver

import (
	"reflect"
	"testing"
	"time"

	. "github.com/jtaylorcpp/gerl/core"
	basicpid "github.com/jtaylorcpp/gerl/core/includes/pid"
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

	gsClient := GenServerClient{
		CallHandler: testCustomCallClient,
		CastHandler: testCustomeCastClient,
	}

	t.Log("client made: ", gsClient)

	t.Log("test call")
	//gs.Pid.SendToInbox(GerlMsg{0x0, ProcessAddr([]byte("testServer")), "test"})
	msg1 := GerlMsg{0x0, ProcessAddr([]byte("testServer")), "test"}
	returnMsg1 := gsClient.Call(gsPid, msg1)

	if !reflect.DeepEqual(msg1, returnMsg1) {
		t.Errorf("msg<%v> not same as msg<%v> recieved\n", msg1, returnMsg1)
	}

	t.Log("test cast")
	//gs.Pid.SendToInbox(GerlMsg{0x1, ProcessAddr([]byte("testServer")), "test"})
	msg2 := GerlMsg{0x1, ProcessAddr([]byte("testServer")), "test"}
	gsClient.Cast(gs.Pid, msg2)

	// terminate stops processes before msgs can be processed
	time.Sleep(5 * time.Second)

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

func testCustomCallClient(pid ProcessID, msg GerlMsg) GerlMsg {
	nPid := pid.(basicpid.Pid)
	nPid.MsgChan <- msg
	returnMsg := <-nPid.MsgChan
	return returnMsg
}

func testCustomeCastClient(pid ProcessID, msg GerlMsg) {
	nPid := pid.(basicpid.Pid)
	nPid.MsgChan <- msg
}
