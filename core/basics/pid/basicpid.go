package basics

import (
	"log"

	"github.com/jtaylorcpp/gerl/core"
)

func init() {
	core.Pid = Pid{}
	log.Println("Setting Pid to type BasicPid")
}

// Pid is an implementation of the ProcessID interface
type Pid struct {
	// Unique address of Pid
	Addr core.ProcessAddr
	// Channel to send GerlMsgs to GenServer
	MsgInbox chan core.GerlMsg
	// Channel to recieve GerlMsgs from GenServer
	MsgOutbox chan core.GerlMsg
	// Channel to recieve TermSig from the GenServer
	TermSig chan bool
	// Keeps track of current GenServer status
	IsRunning bool
}

// Sends GerlMsg to the Pid output channel
/*func (p Pid) SendMsg(msg core.GerlMsg) {
	p.MsgOutbox <- msg
}*/
func (p Pid) SendToProcess() chan core.GerlMsg {
	return p.MsgInbox
}

func (p Pid) SendToTransport() chan core.GerlMsg {
	return p.MsgOutbox
}

// Gets both GerlMsg and Closed/Open bool from input channel
/*func (p Pid) GetMsg() (core.GerlMsg, bool) {
	msg, ok := <-p.MsgInbox
	return msg, ok
}*/

func (p Pid) ProcessReceive() chan core.GerlMsg {
	return p.MsgInbox
}

func (p Pid) TransportReceive() chan core.GerlMsg {
	return p.MsgOutbox
}

// Gets the address of the Pid
func (p Pid) GetAddr() core.ProcessAddr {
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
func (Pid) NewPid(s core.ProcessBufferSize) core.ProcessID {
	return Pid{
		MsgInbox:  make(chan core.GerlMsg, s),
		MsgOutbox: make(chan core.GerlMsg, s),
		TermSig:   make(chan bool, 1),
		IsRunning: true,
	}
}
