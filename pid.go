package gerl

type ProcessID interface{}

type Pid struct {
	MsgBox  chan GerlMsg
	TermSig chan bool
}

func NewPid() Pid {
	return Pid{
		MsgBox:  make(chan GerlMsg),
		TermSig: make(chan bool, 1),
	}
}
