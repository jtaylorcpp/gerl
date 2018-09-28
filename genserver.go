package gerl

import (
	"log"
)

// GenericServer is an implementation of the Erlang OTP gen_server.
// It is intended to be a "single threaded" go routine in which all
// communication with the server happens through the ProcessID (pid).
// Messages are passed through the pid to the GenericServer and are
// processed sequentially by use of the CallHandler and CastHandler.
type GenericServer interface {
	// Builds the GenericServer with intial state
	Init(GenericServerState)
	// Starts the GenericServer and returns the ProcessID associated with it
	Start() ProcessID
	// Processes a synchronous message passed to the GenericServer
	CallHandler(GenericServerMessage, ProcessAddr, GenericServerState) GenericServerMessage
	// Processes an asynchronous message passed to the GenericServer
	CastHandler(GenericServerMessage, ProcessAddr, GenericServerState)
	// Terminate closes the ProcessID and clears out the GenericServer
	Terminate()
}

// GenericServerMessage is an emtpy interface which allows arbitrary objects to
// be passed through the pid to the GenericServer.
type GenericServerMessage interface{}

// GenericServerState is an empty interface which allows the GenericServer to
// manage arbitrary objects inside the GenericServer state
type GenericServerState interface{}

// GenServer is an implementation of the GenericServer.
// It serves both as a reference implemntation and
// easy way to build and use the GenericServer pattern.
type GenServer struct {
	// Pid (type) is a struct defined in pid.go which
	// fullfills the ProcessID interface and allows outside
	// code to interact with the GenServer.
	Pid ProcessID
	// State uses the empty interface GenericServerState to handle arbitrary
	// state information.
	State GenericServerState
	// CustomCall is a func ran inside of this implementations CallHandler.
	// This allows a user defined call routine to be ran within the
	// GenericServer interface.
	CustomCall GenServerCustomCall
	// CustomCast is a func ran inside of this implementations CastHandler.
	// This allows a user defined cast routine to be ran within the
	// GenericServer interface.
	CustomCast GenServerCustomCast
	// BufferSize (uint64) sets the initial Pid GerlMessage buffer size.
	BufferSize ProcessBufferSize
}

// GenServerCustomCall acts like HTTP middleware and is wrapped inside the
// GenericServer.CallHandler of the GenServer.
type GenServerCustomCall func(GenericServerMessage, ProcessAddr, GenericServerState) (GenericServerMessage, GenericServerState)

// GenServerCustomCasr acts like HTTP middleware and is wrapped inside the
// GenericServer.CastHandler of the GenServer.
type GenServerCustomCast func(GenericServerMessage, ProcessAddr, GenericServerState) GenericServerState

// Initializes the GenServer with the intial state
func (gs *GenServer) Init(state GenericServerState) {
	log.Println("Initializing GenServer with state: ", state)
	gs.State = state
}

// Starts the GenServer main loop in which messages are read from the
// "inbox" and then passed to either the CallHandler or the CastHandler.
// This main loop is ran in a go-routine in which when the "inbox" is closed
// the loop sends a termination signal and closes.
// All messages output by the CallHandler are sent to the "outbox" and processed
// by the pid. This "outbox" is also closed when the GenServer main loop is broken.
func (gs *GenServer) Start() Pid {
	// generate a new pid
	pid := NewPid(gs.BufferSize)
	log.Println("GenServer available at pid: ", pid)
	// assign pid to the genserver
	gs.Pid = pid

	// main loop
	go func() {
		log.Printf("GenServer with pid<%v> started\n", gs.Pid)
		for {
			// reads GerlMsg from the inbox
			nextMessage, open := gs.Pid.GetMsg()
			if !open {
				break
			}
			log.Printf("GenServer with pid<%v> got type<%v>\n",
				gs.Pid, nextMessage.Type)
			switch nextMessage.Type {
			case Call:
				log.Printf("GenServer with pid<%v> got call with payload<%v>\n",
					gs.Pid, nextMessage.Msg)
				msg := gs.CallHandler(nextMessage.Msg, nextMessage.FromAddr, gs.State)
				// builds new GerlMsg to send to the outbox
				outMsg := GerlMsg{
					Type:     nextMessage.Type,
					FromAddr: gs.Pid.GetAddr(),
					Msg:      msg,
				}
				log.Printf("GenServer with pid<%v> replied with msg<%v>\n",
					gs.Pid, outMsg)
				gs.Pid.SendMsg(outMsg)
			case Cast:
				log.Printf("pid <%v> got cast with payload<%v>\n",
					gs.Pid, nextMessage.Msg)
				gs.CastHandler(nextMessage.Msg, nextMessage.FromAddr, gs.State)
			}
		}
		log.Printf("GenServer with pid<%v> message box closed\n", gs.Pid)
		gs.Pid.ClosedByGenServer()
	}()

	pid.IsRunning = true
	return pid
}

// CallHandler from GenericServer and passes through all variables to
// the GenServerCustomCall.
func (gs *GenServer) CallHandler(gsm GenericServerMessage, pa ProcessAddr, gss GenericServerState) GenericServerMessage {
	log.Printf("GenServer with pid<%v> calling CustomCaller\n", gs.Pid)
	newMsg, newState := gs.CustomCall(gsm, pa, gss)
	log.Printf("GenServer with pid<%v> has new state<%v>\n", gs.Pid, newState)
	gs.State = newState
	log.Printf("GenServer with pid<%v> call returning msg<%v>\n", gs.Pid, newMsg)
	return newMsg
}

// CastHandler from GenericServer and passes through all variable to
// the GenericServerCustomCast
func (gs *GenServer) CastHandler(gsm GenericServerMessage, pa ProcessAddr, gss GenericServerState) {
	log.Printf("GenServer with pid<%v> calling CustomCaster\n", gs.Pid)
	newState := gs.CustomCast(gsm, pa, gss)
	log.Printf("GenServer with pid<%v> has new state<%v>\n", gs.Pid, newState)
	gs.State = newState
}

// Terminate Calls the Pid.Terminate function to close out both the
// Pid and GenServer
func (gs *GenServer) Terminate() {
	log.Printf("Genserver with pid<%v> terminating\n", gs.Pid)
	gs.Pid.Terminate()
	log.Printf("Genserver with pid<%v> terminated\n", gs.Pid)
}

type GenericServerClient interface {
	Init()
	Call()
	Cast()
	Terminate()
}
