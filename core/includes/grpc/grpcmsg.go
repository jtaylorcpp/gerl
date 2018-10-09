package grpc

import (
	"fmt"
	//"log"

	"github.com/jtaylorcpp/gerl/core"
)

func init() {
	core.GerlMsg = GRPCGerlMsg{}
}

type GRPCGerlMsg struct {
	Type     core.GerlMsgType
	FromAddr core.ProcessAddr
	Msg      interface{}
}

func (msg GRPCGerlMsg) GetType() core.GerlMsgType {
	return msg.Type
}

func (msg GRPCGerlMsg) GetFromAddr() core.ProcessAddr {
	return msg.FromAddr
}

func (msg GRPCGerlMsg) GetMsg() interface{} {
	return msg.Msg
}

func (GRPCGerlMsg) New(t core.GerlMsgType, pa core.ProcessAddr, msg interface{}) core.GerlPassableMessage {
	return GRPCGerlMsg{
		Type:     t,
		FromAddr: pa,
		Msg:      msg,
	}
}

func ToGRPC(msg core.GerlPassableMessage) GerlMsg {
	var t GerlMsg_Type
	switch msg.GetType() {
	case core.Call:
		t = GerlMsg_CALL
	case core.Cast:
		t = GerlMsg_CAST
	default:
		t = GerlMsg_CALL
	}

	nmsg := &Message{}

	switch msg.GetMsg().(type) {
	case string:
		nmsg.Description = msg.GetMsg().(string)
		nmsg.Type = Message_SIMPLE
	default:
		nmsg.Description = fmt.Sprintf("%v", msg.GetMsg())
		nmsg.Type = Message_SIMPLE
	}

	newMsg := GerlMsg{
		Type:        t,
		Processaddr: string(msg.GetFromAddr()),
		Msg:         nmsg,
	}

	return newMsg

}

func ToGerl(msg GerlMsg) core.GerlPassableMessage {
	var t core.GerlMsgType
	switch msg.GetType() {
	case GerlMsg_CALL:
		t = core.Call
	case GerlMsg_CAST:
		t = core.Cast
	default:
		t = core.Call
	}

	newMsg := core.GerlMsg.New(t, core.ProcessAddr(msg.GetProcessaddr()), msg.GetMsg().GetDescription())

	return newMsg
}
