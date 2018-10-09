package core

var GerlMsg GerlPassableMessage

type GerlPassableMessage interface {
	GetType() GerlMsgType
	GetFromAddr() ProcessAddr
	GetMsg() interface{}
	New(GerlMsgType, ProcessAddr, interface{}) GerlPassableMessage
}

// GerlMsgType is a Enum used to designate types of messages
// as they are passed between GenericServers
type GerlMsgType byte

const (
	Call GerlMsgType = 0x0
	Cast GerlMsgType = 0x1
)
