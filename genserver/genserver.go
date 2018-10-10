package genserver

import (
	"errors"
	"log"

	"github.com/jtaylorcpp/gerl/core"
)

type State interface{}
type PidAddr string

// GenericServer is an implementation of the Erlang OTP gen_server.
// It is intended to be a "single threaded" go routine in which all
// communication with the server happens through the ProcessID (pid).
// Messages are passed through the pid to the GenericServer and are
// processed sequentially by use of the CallHandler and CastHandler.
type GenericServer interface {
	// Starts the GenericServer and returns the ProcessID associated with it
	Start() error
	// Processes a synchronous message passed to the GenericServer
	CallHandler(core.Message, PidAddr, State) (core.Message, State)
	// Processes an asynchronous message passed to the GenericServer
	CastHandler(core.Message, PidAddr, State) State
	// Terminate closes the ProcessID and clears out the GenericServer
	Terminate()
}

// GenServerCustomCall acts like HTTP middleware and is wrapped inside the
// GenericServer.CallHandler of the GenServer.
type GenServerCallHandler func(core.Message, PidAddr, State) (core.Message, State)

// GenServerCustomCasr acts like HTTP middleware and is wrapped inside the
// GenericServer.CastHandler of the GenServer.
type GenServerCastHandler func(core.Message, PidAddr, GenericServerState) State

// GenServer is an implementation of the GenericServer.
// It serves both as a reference implemntation and
// easy way to build and use the GenericServer pattern.
type GenServer struct {
	// Pid (type) is a struct defined in pid.go which
	// fullfills the ProcessID interface and allows outside
	// code to interact with the GenServer.
	Pid core.Pid
	// State uses the empty interface GenericServerState to handle arbitrary
	// state information.
	State State
	// CustomCall is a func ran inside of this implementations CallHandler.
	// This allows a user defined call routine to be ran within the
	// GenericServer interface.
	CustomCall GenServerCustomCall
	// CustomCast is a func ran inside of this implementations CastHandler.
	// This allows a user defined cast routine to be ran within the
	// GenericServer interface.
	CustomCast GenServerCustomCast
	Errors     chan error
	// BufferSize (uint64) sets the initial Pid GerlMessage buffer size.
	Terminated chan bool
}

// Initializes the GenServer with the intial state
func New(state State, call GenServerCustomCall, cast GenServerCustomCast) *GenServer {
	log.Println("Initializing GenServer with state: ", state)
	return &GenServer{
		Pid:        nil,
		State:      state,
		CustomCall: call,
		CustomCast: cast,
		Errors:     make(chan error),
		Terminated: make(chan bool),
	}
}

// Starts the GenServer main loop in which messages are read from the
// "inbox" and then passed to either the CallHandler or the CastHandler.
// This main loop is ran in a go-routine in which when the "inbox" is closed
// the loop sends a termination signal and closes.
// All messages output by the CallHandler are sent to the "outbox" and processed
// by the pid. This "outbox" is also closed when the GenServer main loop is broken.
func (gs *GenServer) Start() error {
	// generate a new pid
	gs.Pid = core.NewPid("", "")
	log.Println("GenServer available at pid: ", gs.Pid.GetAddr())

	for {
		select {
		case err := <-gs.Pid.Errors:
			log.Println("genserver pid error: ", err)
			gs.Errors <- errors.New("pid error, close genserver")
		case term := <-gs.Terminated:
			log.Println("genserver terminated")
			return errors.New("genserver terminated")
		case msg, ok := <-gs.Pid.Inbox:
			if !ok {
				return errors.New("genserver message inbox closed")
			}
			switch msg.GetType() {
			case core.GerlMsg_CALL:
				log.Println("genserver recieved call")
			case core.GerlMsg_CAST:
				log.Println("genserver recieved cast")
			default:
				log.Println("genserver recieved unknown type: ", msg.GetType())
			}
		default:

		}
	}

	/*
		// main loop
		ready := make(chan bool)
		go func() {
			log.Printf("GenServer with pid<%v> started\n", gs.Pid)
			var loopState GenericServerState = gs.State
			ready <- true
			for {
				// reads GerlMsg from the inbox
				nextMessage, ok := gs.Pid.Read()
				if !ok {
					break
				}
				log.Printf("GenServer with pid<%v> got type<%v>\n",
					gs.Pid, nextMessage.GetType())
				switch nextMessage.GetType() {
				case Call:
					log.Printf("GenServer with pid<%v> got call with payload<%v>\n",
						gs.Pid, nextMessage.GetMsg())
					msg, newState := gs.CallHandler(nextMessage.GetMsg(), nextMessage.GetFromAddr(), loopState)
					// builds new GerlMsg to send to the outbox
					outMsg := GerlMsg.New(nextMessage.GetType(), gs.Pid.GetAddr(), msg)
					loopState = newState
					log.Printf("GenServer with pid<%v> replied with msg<%v>\n",
						gs.Pid, outMsg)
					gs.Pid.Write(outMsg)
				case Cast:
					log.Printf("pid <%v> got cast with payload<%v>\n",
						gs.Pid, nextMessage.GetMsg())
					loopState = gs.CastHandler(nextMessage.GetMsg(), nextMessage.GetFromAddr(), loopState)
				}
				log.Printf("GenServer with pid<%v> new state<%v>\n", gs.Pid, loopState)
			}
			log.Printf("GenServer with pid<%v> message box closed\n", gs.Pid)
			gs.terminated <- true
		}()
		<-ready
		return pid
	*/
}

// CallHandler from GenericServer and passes through all variables to
// the GenServerCustomCall.
func (gs *GenServer) CallHandler(gsm GenericServerMessage, pa ProcessAddr, gss GenericServerState) (GenericServerMessage, GenericServerState) {
	log.Printf("GenServer with pid<%v> calling CustomCaller\n", gs.Pid)
	newMsg, newState := gs.CustomCall(gsm, pa, gss)
	log.Printf("GenServer with pid<%v> has new state<%v>\n", gs.Pid, newState)
	log.Printf("GenServer with pid<%v> call returning msg<%v>\n", gs.Pid, newMsg)
	return newMsg, newState
}

// CastHandler from GenericServer and passes through all variable to
// the GenericServerCustomCast
func (gs *GenServer) CastHandler(gsm GenericServerMessage, pa ProcessAddr, gss GenericServerState) GenericServerState {
	log.Printf("GenServer with pid<%v> calling CustomCaster\n", gs.Pid)
	newState := gs.CustomCast(gsm, pa, gss)
	log.Printf("GenServer with pid<%v> has new state<%v>\n", gs.Pid, newState)
	return newState
}

// Terminate Calls the Pid.Terminate function to close out both the
// Pid and GenServer
func (gs *GenServer) Terminate() {
	log.Printf("Genserver with pid<%v> terminating\n", gs.Pid)
	gs.Pid.Terminate()
	<-gs.terminated
	log.Printf("Genserver with pid<%v> terminated\n", gs.Pid)
}

type GenericServerClient interface {
	Call(ProcessID, GerlPassableMessage) GerlPassableMessage
	Cast(ProcessID, GerlPassableMessage)
}

type GenServerClient struct {
	CallHandler GenServerClientCall
	CastHandler GenServerClientCast
}

type GenServerClientCall func(ProcessID, GerlPassableMessage) GerlPassableMessage
type GenServerClientCast func(ProcessID, GerlPassableMessage)

func (gsc GenServerClient) Call(pid ProcessID, msg GerlPassableMessage) GerlPassableMessage {
	log.Printf("client calling pid<%v> with msg<%v>\n", pid, msg)
	return gsc.CallHandler(pid, msg)
}

func (gsc GenServerClient) Cast(pid ProcessID, msg GerlPassableMessage) {
	log.Printf("client casting pid<%v> with msg<%v>\n", pid, msg)
	gsc.CastHandler(pid, msg)
}
