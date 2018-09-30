package core

var Pid ProcessID

// ProcessID is the piece that connects other pieces of software to a given
// GenericServer.
type ProcessID interface {
	// Allocates new Pid
	NewPid(ProcessBufferSize) ProcessID
	// Sends a gerlMsg to a given ProcessID/GenericServer
	SendToProcess() chan GerlMsg
	SendToTransport() chan GerlMsg
	// Gets a GerlMsg and whether the ProcessID is still open for a
	// GenericServer to consume.
	ProcessReceive() chan GerlMsg
	TransportReceive() chan GerlMsg
	// Address of ProcessID
	GetAddr() ProcessAddr
	// Stops GenericServer from writing messages to ProcessID
	ClosedByGenServer()
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
