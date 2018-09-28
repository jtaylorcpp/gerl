package gerl

import (
	"log"
)

// ProcessID is the piece that connects other pieces of software to a given
// GenericServer.
type ProcessID interface {
	// Sends a gerlMsg to a given ProcessID/GenericServer
	SendMsg(GerlMsg)
	// Gets a GerlMsg and whether the ProcessID is still open for a
	// GenericServer to consume.
	GetMsg() (GerlMsg, bool)
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

// Pid is an implementation of the ProcessID interface
type Pid struct {
	// Unique address of Pid
	Addr ProcessAddr
	// Channel to send GerlMsgs to GenServer
	MsgInbox chan GerlMsg
	// Channel to recieve GerlMsgs from GenServer
	MsgOutbox chan GerlMsg
	// Channel to recieve TermSig from the GenServer
	TermSig chan bool
	// Keeps track of current GenServer status
	IsRunning bool
}

// Sends GerlMsg to the Pid output channel
func (p Pid) SendMsg(msg GerlMsg) {
	p.MsgOutbox <- msg
}

// Gets both GerlMsg and Closed/Open bool from input channel
func (p Pid) GetMsg() (GerlMsg, bool) {
	msg, ok := <-p.MsgInbox
	return msg, ok
}

// Gets the address of the Pid
func (p Pid) GetAddr() ProcessAddr {
	return p.Addr
}

// Closes the channels MsgOutbox and TermSig both of which the GenServer writes to.
func (p Pid) ClosedByGenServer() {
	p.TermSig <- true
	close(p.MsgOutbox)
	close(p.TermSig)
}

// Terminates both the Pid and sends the close signal to let the GenServer close.
func (p Pid) Terminate() {
	log.Printf("pid<%v> closing inbox\n", p)
	// closing MsgInbox forces the genserever main go-routine
	// to exit
	close(p.MsgInbox)
	for {
		// if there are messages in the outbox drop them on the floor
		outMsg, open := <-p.MsgOutbox
		if !open {
			break
		} else {
			log.Printf("pid<%v> closing: extra msg<%v>\n", p, outMsg)
		}
	}
	// block until TermSig is recieved
	<-p.TermSig
	// pid/genserver is no longer running
	p.IsRunning = false
}

// Builds a new pid of type Pid
func NewPid(s ProcessBufferSize) Pid {
	return Pid{
		MsgInbox:  make(chan GerlMsg, s),
		MsgOutbox: make(chan GerlMsg, s),
		TermSig:   make(chan bool, 1),
		IsRunning: false,
	}
}
