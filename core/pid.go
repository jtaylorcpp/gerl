package core

var Pid ProcessID

// ProcessID is the piece that connects other pieces of software to a given
// GenericServer.
type ProcessID interface {
	// Allocates new Pid
	NewPid(ProcessBufferSize) ProcessID
	Read() (GerlPassableMessage, bool)
	Write(GerlPassableMessage)
	// Address of ProcessID
	GetAddr() ProcessAddr
	// Closes out ProcessID and stops messages going to the GenericServer
	Terminate()
}

// ProcessAddr used to uniquely identify a GenericServer and where it exists
type ProcessAddr []byte

// ProcessBufferSize used to set size of GerlMsg buffin in ProcessID
type ProcessBufferSize uint64

func NewPid(s ProcessBufferSize) ProcessID {
	return Pid.NewPid(s)
}
