package registrar

import (
	"fmt"
	"gerl/core"
	"gerl/genserver"
	"time"

	log "github.com/sirupsen/logrus"
)

type RegistarAction = string

const (
	ADD    RegistarAction = "ADD"
	REMOVE RegistarAction = "REMOVE"
	QUERY  RegistarAction = "QUERY"
)

type RegistrarMessage struct {
	Action RegistarAction
	Record RegistrarRecord
}

type RegistrarReply struct {
	Error   error
	Record  RegistrarRecord
	Records []ProcessAddress
}

type RegistrarState struct {
	Registrar     *Registrar
	RecordTimeOut time.Duration
}

type RegistrarGenServer struct {
	*genserver.GenServer
}

func NewRegistrarGenServer(state *RegistrarState) (*genserver.GenServer, error) {
	config := &genserver.GenServerConfig{
		StartState:  state,
		CallHandler: RegistrarCallHandler,
		CastHandler: RegistrarCastHandler,
		Scope:       core.GlobalScope,
	}

	return genserver.NewGenServer(config)
}

func RegistrarCallHandler(_ core.Pid, msg RegistrarMessage, _ genserver.FromAddr, s RegistrarState) (RegistrarReply, RegistrarState) {
	log.Printf("registrar call handler with message: %#v", msg)
	returnMsg := &RegistrarReply{}
	switch msg.Action {
	case ADD:
		returnMsg.Error = s.Registrar.AddServiceRecord(msg.Record)
	case REMOVE:
		returnMsg.Error = s.Registrar.RemoveServiceRecord(msg.Record)
	case QUERY:
		returnMsg.Records, returnMsg.Error = s.Registrar.GetServiceAddresses(msg.Record)
	default:
		log.Printf("registrar does not recognize action: %#v", msg)
		returnMsg.Error = fmt.Errorf("registrar does not recognize the action %s", msg.Action)
	}
	return *returnMsg, s
}

func RegistrarCastHandler(_ core.Pid, msg RegistrarMessage, _ genserver.FromAddr, s RegistrarState) RegistrarState {
	log.Printf("registrar cast handler with message: %#v", msg)

	return s
}
