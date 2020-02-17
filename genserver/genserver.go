package genserver

import (
	"errors"
	"log"
	"encoding/json"
	"reflect"
	"fmt"

	"gerl/core"
)

type State = interface{}
type Message = []byte
type PidAddr string
type FromAddr = string

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
type GenServerCallHandler func(core.Pid, Message, FromAddr, State) (Message, State)
func newCustomCallHandler(handlerFunc interface{}) (GenServerCallHandler, error) {
	if handlerFunc == nil {
		return nil, errors.New("call handler is nil")
	}

	handler := reflect.ValueOf(handlerFunc)

	// validate func signature
	handlerType := reflect.TypeOf(handlerFunc)

	if handlerType.Kind() != reflect.Func {
		return nil, errors.New("call handler is not a func")
	}

	if handlerType.NumIn() != 4 {
		return nil, errors.New("call handler does not take 4 arguments")
	}

	

	pidType := reflect.TypeOf(&core.Pid{}).Elem()
	arg0Type := handlerType.In(0)
	if pidType != arg0Type {
		return nil, errors.New(fmt.Sprintf("call handler param 0 (type: %s) should implement type %s", arg0Type, pidType))
	}

	if handlerType.In(1).Kind() != reflect.Struct {
		return nil, errors.New(fmt.Sprintf("call handler param 1 (type: %s) should implement type struct", handlerType.In(1).Kind()))
	}

	stringType:= reflect.TypeOf((*string)(nil)).Elem()
	
	if handlerType.In(2) != stringType {
		return nil, errors.New(fmt.Sprintf("call handler param 2 (type: %s) should implement type FromAddr(string)", handlerType.In(2).Kind()))
	}

	if handlerType.In(3).Kind() != reflect.Struct {
		return nil, errors.New(fmt.Sprintf("call handler param 3 (type: %s) should implement type struct", handlerType.In(3).Kind()))
	}

	
	if handlerType.NumOut() != 2 {
		return nil, errors.New("call handler does not output 2 parameters")
	}

	if handlerType.Out(0).Kind() != reflect.Struct {
		return nil, errors.New(fmt.Sprintf("call handler output param 0 (type: %s) should implement type struct", handlerType.Out(0).Kind()))
	}
	
	stateType := reflect.TypeOf(handlerType.In(3)).Elem()
	out1Type := reflect.TypeOf(handlerType.Out(1)).Elem()
	if stateType != out1Type {
		return nil, errors.New(fmt.Sprintf("call handler does not use same type for input (type: %s) and output (type: %s) states",stateType, out1Type))
	}

	returnFunc := GenServerCallHandler(func(pid core.Pid, msg Message, faddr FromAddr, s State)(Message, State){
		// args from original func

		var args []reflect.Value

		// arg 0 is always Pid
		args = append(args, reflect.ValueOf(pid))

		// arg 1 is defined message
		msgStruct := reflect.New(handlerType.In(1))
		if err := json.Unmarshal(msg, msgStruct.Interface()); err != nil {
			log.Fatal("unable to unmarshal message into struct: ", handlerType.In(1).Elem())
		}

		args = append(args, msgStruct.Elem())

		// arg 2 is the FromAddr
		args = append(args, reflect.ValueOf(faddr))

		// arg3 is the state struct
		args = append(args, reflect.ValueOf(s))

		response := handler.Call(args)

		// response needs to be converted to msg, state
		returnMsg, err := json.Marshal(response[0].Interface())
		if err != nil {
			log.Println(err.Error())
			log.Fatal("unable to marshal returned message")
		}

		
		

		return returnMsg, response[1].Interface()
	})

	return returnFunc, nil
}


// GenServerCustomCasr acts like HTTP middleware and is wrapped inside the
// GenericServer.CastHandler of the GenServer.
type GenServerCastHandler func(core.Pid, Message, FromAddr, State) State
func newCustomCastHandler(handler interface{}) GenServerCastHandler {
	return nil
} 

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
// call must be a func with a signature of: func(core.Pid, StructIn, FromAddr, StructState) (StructOut, StructState)
// cast must be a func with a signature of: func(core.Pid, StructIn, FromAddr, StructState) StructState
func NewGenServer(state State, scope core.Scope, call , cast interface{}) *GenServer {
	log.Println("Initializing GenServer with state: ", state)
	log.Println("Initializing Call Handler")
	callHandler, err := newCustomCallHandler(call)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Initializing Cast Hander")
	castHandler := newCustomCastHandler(cast)
	return &GenServer{
		Pid:        &core.Pid{},
		State:      state,
		Scope:      scope,
		CustomCall: callHandler,
		CustomCast: castHandler,
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
				var rawMsg []byte
				msgStruct, loopState := gs.CallHandler(*msg.GetMsg(), FromAddr(msg.GetFromaddr()), loopState)
				rawMsg, err := json.Marshal(&msgStruct)
				if err != nil {
					log.Println("error marshalling message to pass: ", err.Error())
					rawMsg = []byte{}
				}
				gs.Pid.Outbox <- core.GerlMsg{
					Type:     core.GerlMsg_CALL,
					Fromaddr: gs.Pid.GetAddr(),
					Msg:      &core.Message{
						RawMsg: rawMsg,
					},
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
func (gs *GenServer) CallHandler(msg core.Message, fa FromAddr, s State) (Message, State) {
	log.Printf("GenServer with pid<%v> calling CustomCaller\n", gs.Pid)
	newMsg, newState := gs.CustomCall(*gs.Pid, msg.GetRawMsg(), fa, s)
	log.Printf("GenServer with pid<%v> has new state<%v>\n", gs.Pid, newState)
	log.Printf("GenServer with pid<%v> call returning msg<%v>\n", gs.Pid, newMsg)
	return newMsg, newState
}

// CastHandler from GenericServer and passes through all variable to
// the GenericServerCustomCast
func (gs *GenServer) CastHandler(msg core.Message, fa FromAddr, s State) State {
	log.Printf("GenServer with pid<%v> calling CustomCaster\n", gs.Pid)
	newState := gs.CustomCast(*gs.Pid, msg.GetRawMsg(), fa, s)
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
