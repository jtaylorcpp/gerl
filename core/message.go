package core

// GerlMsgType is a Enum used to designate types of messages
// as they are passed between GenericServers
type GerlMsgType byte

const (
	Call GerlMsgType = 0x0
	Cast GerlMsgType = 0x1
)

// GerlMsg provides the structure of messages to be passed between
// GenericServers
type GerlMsg struct {
	// Enum type expressed as Byte
	Type GerlMsgType
	// Address of ProcesID being sent from
	FromAddr ProcessAddr
	// Msg is a blank interface that is the actual data sent between processes
	Msg interface{}
}
