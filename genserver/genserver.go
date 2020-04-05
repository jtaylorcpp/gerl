package genserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"gerl/core"
)

func init() {
	log.SetReportCaller(true)
}

type State = interface{}
type Message = []byte
type PidAddr = string
type FromAddr = string

// GenericServer is an implementation of the Erlang OTP gen_server.
// It is intended to be a "single threaded" go routine in which all
// communication with the server happens through the ProcessID (pid).
// Messages are passed through the pid to the GenericServer and are
// processed sequentially by use of the CallHandler and CastHandler.
type GenericServer interface {
	// Starts the GenericServer and returns the ProcessID associated with it
	Start() error
	GetPid() *core.Pid
	IsReady() bool
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

	stringType := reflect.TypeOf((*string)(nil)).Elem()

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

	stateType := handlerType.In(3)
	out1Type := handlerType.Out(1)
	if stateType != out1Type {
		return nil, errors.New(fmt.Sprintf("call handler does not use same type for input (type: %s) and output (type: %s) states", stateType, out1Type))
	}

	returnFunc := GenServerCallHandler(func(pid core.Pid, msg Message, faddr FromAddr, s State) (Message, State) {
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

func newCustomCastHandler(handlerFunc interface{}) (GenServerCastHandler, error) {
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

	stringType := reflect.TypeOf((*string)(nil)).Elem()

	if handlerType.In(2) != stringType {
		return nil, errors.New(fmt.Sprintf("call handler param 2 (type: %s) should implement type FromAddr(string)", handlerType.In(2).Kind()))
	}

	if handlerType.In(3).Kind() != reflect.Struct {
		return nil, errors.New(fmt.Sprintf("call handler param 3 (type: %s) should implement type struct", handlerType.In(3).Kind()))
	}

	if handlerType.NumOut() != 1 {
		return nil, errors.New("call handler does not output 1 parameter")
	}

	stateType := handlerType.In(3)
	out0Type := handlerType.Out(0)
	if stateType != out0Type {
		return nil, errors.New(fmt.Sprintf("call handler does not use same type for input (type: %s) and output (type: %s) states", stateType, out0Type))
	}

	returnFunc := GenServerCastHandler(func(pid core.Pid, msg Message, faddr FromAddr, s State) State {
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

		// response needs to be converted to state

		return response[0].Interface()
	})

	return returnFunc, nil
}

type GenServerConfig struct {
	StartState  State
	CallHandler interface{}
	CastHandler interface{}
	Scope       core.Scope
	Port        string
	Address     string
}

type GenServer struct {
	config *GenServerConfig
	// Pid (type) is a struct defined in pid.go which
	// fullfills the ProcessID interface and allows outside
	// code to interact with the GenServer.
	pid *core.Pid
	// State uses the empty interface GenericServerState to handle arbitrary
	// state information.
	state State
	scope core.Scope
	// CustomCall is a func ran inside of this implementations CallHandler.
	// This allows a user defined call routine to be ran within the
	// GenericServer interface.
	customCall GenServerCallHandler
	// CustomCast is a func ran inside of this implementations CastHandler.
	// This allows a user defined cast routine to be ran within the
	// GenericServer interface.
	customCast GenServerCastHandler
	mutex      *sync.Mutex
}

func NewGenServer(config *GenServerConfig) (*GenServer, error) {
	log.Println("Initializing Call Hander")
	callHandler, err := newCustomCallHandler(config.CallHandler)
	if err != nil {
		return nil, err
	}
	log.Println("Initializing Cast Hander")
	castHandler, err := newCustomCastHandler(config.CastHandler)
	if err != nil {
		return nil, err
	}
	return &GenServer{
		config:     config,
		customCall: callHandler,
		customCast: castHandler,
		mutex:      &sync.Mutex{},
	}, nil
}

func (gs *GenServer) Start() error {
	gs.mutex.Lock()
	pid, err := core.NewPid(gs.config.Address, gs.config.Port, gs.config.Scope)
	if err != nil {
		return err
	}
	gs.pid = pid
	gs.state = gs.config.StartState

	log.Println("GenServer available at pid: ", gs.pid.GetAddr())
	log.Printf("genserver: %#v\n", gs)

	var exitError error = nil
	gs.mutex.Unlock()
	for {
		select {
		case msg, ok := <-gs.pid.Inbox:
			log.Printf("recieved message <%#v> from inbox<status:%#v>\n", msg, ok)
			if !ok {
				log.Println("pid inbox has been closed without genserver termination")
				exitError = errors.New("genserver inbox has been closed")
				goto pid_terminated
			} else {
				switch msg.Type {
				case core.GerlMsg_CALL:
					log.Println("genserver recieved call")
					var returnMsg core.Message
					returnMsg, gs.state = gs.CallHandler(*msg.GetMsg(), msg.GetFromaddr(), gs.state)
					gs.pid.Outbox <- core.GerlMsg{
						Type:     core.GerlMsg_CALL,
						Fromaddr: gs.pid.GetAddr(),
						Msg:      &returnMsg,
					}
					log.Println("state after call: ", gs.state)
				case core.GerlMsg_CAST:
					log.Println("genserver recieved cast")
					gs.state = gs.CastHandler(*msg.GetMsg(), msg.GetFromaddr(), gs.state)
					log.Println("state after cast: ", gs.state)
				case core.GerlMsg_ERR:
					log.Errorf("genserver recieved error message: <%#v>\n", msg)
				case core.GerlMsg_TERM:
					log.Printf("genserver recieved terminate message: <%#v>\n", msg)
					goto terminate
				}
			}
		}
	}
terminate:
	log.Printf("genserver terminating")
	gs.pid.Terminate()
pid_terminated:
	log.Printf("Genserver with pid<%#v> terminated\n", gs.pid)
	gs.pid = nil
	log.Println("genserver end state: ", gs.state)
	return exitError
}

func (gs *GenServer) GetPid() *core.Pid {
	return gs.pid
}

func (gs *GenServer) IsReady() bool {
	if gs.pid == nil {
		return false
	}

	return core.PidHealthCheck(gs.GetPid().GetAddr())
}

// CallHandler from GenericServer and passes through all variables to
// the GenServerCustomCall.
func (gs *GenServer) CallHandler(msg core.Message, fa FromAddr, s State) (core.Message, State) {
	log.Printf("GenServer with pid<%#v> calling CustomCaller\n", gs.pid)
	newMsg, newState := gs.customCall(*gs.pid, msg.GetRawMsg(), fa, s)
	log.Printf("GenServer with pid<%#v> has new state<%v>\n", gs.pid, newState)
	log.Printf("GenServer with pid<%#v> call returning msg<%v>\n", gs.pid, newMsg)
	return core.Message{RawMsg: newMsg}, newState
}

// CastHandler from GenericServer and passes through all variable to
// the GenericServerCustomCast
func (gs *GenServer) CastHandler(msg core.Message, fa FromAddr, s State) State {
	log.Printf("GenServer with pid<%v> calling CustomCaster\n", gs.pid)
	newState := gs.customCast(*gs.pid, msg.GetRawMsg(), fa, s)
	log.Printf("GenServer with pid<%v> has new state<%v>\n", gs.pid, newState)
	return newState
}

func (gs *GenServer) Terminate() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	if gs.pid == nil {
		return
	}

	core.PidTerminate(gs.pid.GetAddr(), gs.pid.GetAddr(), core.Message{})
	for gs.pid != nil {
		time.Sleep(10 * time.Nanosecond)
	}
}

// Call sends an arbitrary core.Message to the GenServer at address PidAddr
// and includes the FromAddr
// This is desigend to send Call messages specifically to GenServers
func Call(to PidAddr, from FromAddr, msg interface{}) (interface{}, error) {
	rawMsg, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	returnMsg := core.PidCall(string(to), string(from), core.Message{RawMsg: rawMsg})

	returnStruct := reflect.New(reflect.TypeOf(msg))
	log.Printf("returned message: %#v\n", returnMsg.GetRawMsg())
	err = json.Unmarshal(returnMsg.GetRawMsg(), returnStruct.Interface())
	if err != nil {
		return nil, err
	}

	return returnStruct.Elem().Interface(), nil
}

// Cast sends an arbitrary core.Message to the GenServer at address PidAddr
// and includes the FromAddr
// This is desigend to send Cast messages specifically to GenServers
func Cast(to PidAddr, from PidAddr, msg core.Message) {
	core.PidCast(string(to), string(from), msg)
}
