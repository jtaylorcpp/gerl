package gerl

type ProcessID interface{}

type Pid struct {
	MsgBox chan GerlMsg
}

func NewPid() Pid {
	return Pid{
		MsgBox: make(chan GerlMsg),
	}
}
