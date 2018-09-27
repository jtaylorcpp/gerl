package gerl

type GSType byte

const (
	Call GSType = 0x0
	Cast GSType = 0x1
)

type GerlMsg struct {
	Type    GSType
	Payload interface{}
}
