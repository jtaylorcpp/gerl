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
	// Channel GerlMsgs to/from GenServer
	MsgChan chan core.GerlMsg
	// Keeps track of current GenServer status
	IsRunning bool
}

func (p Pid) Read() (core.GerlMsg, bool) {
	msg, ok := <-p.MsgChan
	return msg, ok
}

func (p Pid) Write(msg core.GerlMsg) {
	p.MsgChan <- msg
}

// Gets the address of the Pid
func (p Pid) GetAddr() core.ProcessAddr {
	return p.Addr
}

// Terminates both the Pid and sends the close signal to let the GenServer close.
func (p Pid) Terminate() {
	log.Printf("pid<%v> closing inbox\n", p)
	// closing MsgInbox forces the genserever main go-routine
	// to exit
	close(p.MsgChan)
	p.IsRunning = false
}

// Builds a new pid of type Pid
func (Pid) NewPid(s core.ProcessBufferSize) core.ProcessID {
	return Pid{
		MsgChan:   make(chan core.GerlMsg, s),
		IsRunning: true,
	}
}
