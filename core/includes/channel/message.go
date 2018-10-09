package channel

import (
	"log"

	"github.com/jtaylorcpp/gerl/core"
)

/*
type GerlPassableMessage interface {
	GetType() GerlMsgType
	GetFromAddr() ProcessAddr
	GetMsg() interface{}
}
*/

func init() {
	core.GerlMsg = GerlMsg{}
	log.Printf("Setting GerlMsg to StructGerlMsg")
}

// GerlMsg provides the structure of messages to be passed between
// GenericServers
type GerlMsg struct {
	// Enum type expressed as Byte
	Type core.GerlMsgType
	// Address of ProcesID being sent from
	FromAddr core.ProcessAddr
	// Msg is a blank interface that is the actual data sent between processes
	Msg interface{}
}

func (gm GerlMsg) GetType() core.GerlMsgType {
	return gm.Type
}

func (gm GerlMsg) GetFromAddr() core.ProcessAddr {
	return gm.FromAddr
}

func (gm GerlMsg) GetMsg() interface{} {
	return gm.Msg
}

func (GerlMsg) New(t core.GerlMsgType, addr core.ProcessAddr, msg interface{}) core.GerlPassableMessage {
	return GerlMsg{
		Type:     t,
		FromAddr: addr,
		Msg:      msg,
	}
}
