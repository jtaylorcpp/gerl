package genserver

import (
	"errors"
	"log"

	"github.com/jtaylorcpp/gerl/core"
)

type State interface{}
type PidAddr string
type FromAddr string

// GenericServer is an implementation of the Erlang OTP gen_server.
// It is intended to be a "single threaded" go routine in which all
// communication with the server happens through the ProcessID (pid).
// Messages are passed through the pid to the GenericServer and are
// processed sequentially by use of the CallHandler and CastHandler.
type GenericServer interface {
	// Starts the GenericServer and returns the ProcessID associated with it
	Start() error
	// Processes a synchronous message passed to the GenericServer
	CallHandler(core.Message, FromAddr, State) (core.Message, State)
	// Processes an asynchronous message passed to the GenericServer
	CastHandler(core.Message, FromAddr, State) State
	// Terminate closes the ProcessID and clears out the GenericServer
	Terminate()
}

// GenServerCustomCall acts like HTTP middleware and is wrapped inside the
// GenericServer.CallHandler of the GenServer.
type GenServerCallHandler func(core.Pid, core.Message, FromAddr, State) (core.Message, State)

// GenServerCustomCasr acts like HTTP middleware and is wrapped inside the
// GenericServer.CastHandler of the GenServer.
type GenServerCastHandler func(core.Pid, core.Message, FromAddr, State) State

// GenServer is an implementation of the GenericServer.
// It serves both as a reference implemntation and
// easy way to build and use the GenericServer pattern.
type GenServer struct {
	// Pid (type) is a struct defined in pid.go which
	// fullfills the ProcessID interface and allows outside
	// code to interact with the GenServer.
	Pid *core.Pid
	// State uses the empty interface GenericServerState to handle arbitrary
	// state information.
	State State
	Scope core.Scope
	// CustomCall is a func ran inside of this implementations CallHandler.
	// This allows a user defined call routine to be ran within the
	// GenericServer interface.
	CustomCall GenServerCallHandler
	// CustomCast is a func ran inside of this implementations CastHandler.
	// This allows a user defined cast routine to be ran within the
	// GenericServer interface.
	CustomCast GenServerCastHandler
	Errors     chan error
	// BufferSize (uint64) sets the initial Pid GerlMessage buffer size.
	Terminated chan bool
}

// Initializes the GenServer with the intial state
// takes in both the Call handler and Cast handler to be used in the main loop
func NewGenServer(state State, scope core.Scope, call GenServerCallHandler, cast GenServerCastHandler) *GenServer {
	log.Println("Initializing GenServer with state: ", state)
	return &GenServer{
		Pid:        &core.Pid{},
		State:      state,
		Scope:      scope,
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
	gs.Pid = core.NewPid("", "", gs.Scope)
	log.Println("GenServer available at pid: ", gs.Pid.GetAddr())
	var loopState State
	loopState = gs.State
	for {
		select {
		case err := <-gs.Pid.Errors:
			log.Println("genserver pid error: ", err)
			return errors.New("pid error, close genserver")
		case <-gs.Terminated:
			log.Println("genserver terminated")
			return errors.New("genserver terminated")
		case msg, ok := <-gs.Pid.Inbox:
			if !ok {
				return errors.New("genserver message inbox closed")
			}
			switch msg.GetType() {
			case core.GerlMsg_CALL:
				log.Println("genserver recieved call")
				var newMsg core.Message
				newMsg, loopState = gs.CallHandler(*msg.GetMsg(), FromAddr(msg.GetFromaddr()), loopState)
				gs.Pid.Outbox <- core.GerlMsg{
					Type:     core.GerlMsg_CALL,
					Fromaddr: gs.Pid.GetAddr(),
					Msg:      &newMsg,
				}
				log.Println("state after call: ", loopState)
			case core.GerlMsg_CAST:
				log.Println("genserver recieved cast")
				loopState = gs.CastHandler(*msg.GetMsg(), FromAddr(msg.GetFromaddr()), loopState)
				log.Println("state after cast: ", loopState)
			default:
				log.Println("genserver recieved unknown type: ", msg.GetType())
			}
		default:
			continue
			//log.Println("genserver no matching cases")
		}
	}

	log.Println("genserver end state: ", loopState)
	return errors.New("genserver EOP")
}

// CallHandler from GenericServer and passes through all variables to
// the GenServerCustomCall.
func (gs *GenServer) CallHandler(msg core.Message, fa FromAddr, s State) (core.Message, State) {
	log.Printf("GenServer with pid<%v> calling CustomCaller\n", gs.Pid)
	newMsg, newState := gs.CustomCall(*gs.Pid, msg, fa, s)
	log.Printf("GenServer with pid<%v> has new state<%v>\n", gs.Pid, newState)
	log.Printf("GenServer with pid<%v> call returning msg<%v>\n", gs.Pid, newMsg)
	return newMsg, newState
}

// CastHandler from GenericServer and passes through all variable to
// the GenericServerCustomCast
func (gs *GenServer) CastHandler(msg core.Message, fa FromAddr, s State) State {
	log.Printf("GenServer with pid<%v> calling CustomCaster\n", gs.Pid)
	newState := gs.CustomCast(*gs.Pid, msg, fa, s)
	log.Printf("GenServer with pid<%v> has new state<%v>\n", gs.Pid, newState)
	return newState
}

// Terminate Calls the Pid.Terminate function to close out both the
// Pid and GenServer
func (gs *GenServer) Terminate() {
	log.Printf("Genserver with pid<%v> terminating\n", gs.Pid)
	gs.Terminated <- true
	close(gs.Terminated)
	gs.Pid.Terminate()
	for {
		err, ok := <-gs.Pid.Errors
		if !ok {
			break
		}
		log.Println("genserver clearing pid errors: ", err)
	}
	log.Printf("Genserver with pid<%v> terminated\n", gs.Pid)
}

// Call sends an arbitrary core.Message to the GenServer at address PidAddr
// and includes the FromAddr
// This is desigend to send Call messages specifically to GenServers
func Call(to PidAddr, from FromAddr, msg core.Message) core.Message {
	return core.PidCall(string(to), string(from), msg)
}

// Cast sends an arbitrary core.Message to the GenServer at address PidAddr
// and includes the FromAddr
// This is desigend to send Cast messages specifically to GenServers
func Cast(to PidAddr, from PidAddr, msg core.Message) {
	core.PidCast(string(to), string(from), msg)
}
