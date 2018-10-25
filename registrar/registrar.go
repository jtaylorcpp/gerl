package registrar

import (
	"log"

	"github.com/jtaylorcpp/gerl/core"
	"github.com/jtaylorcpp/gerl/genserver"
)

type Registrar genserver.GenServer
type State genserver.State
type CallHander genserver.GenServerCallHandler
type CastHandler genserver.GenServerCastHandler

func registrarCallHander(pid core.Pid, in core.Message, from genserver.FromAddr, state genserver.State) (core.Message, genserver.State) {

	switch in.GetType() {
	case core.Message_SIMPLE:

	case core.Message_SYNC:

	default:

	}

	return core.Message{}, state
}

func registrarCastHander(pid core.Pid, in core.Message, from genserver.FromAddr, state genserver.State) genserver.State {

	return state
}

func New(scope core.Scope) *Registrar {
	gensvr := genserver.NewGenServer(newRegister(), scope, registrarCallHander, registrarCastHander)

	reg := Registrar(*gensvr)

	return &reg
}

type register struct {
	recordmap map[string]map[string]record
}

type record struct {
	name    string
	address string
	scope   core.Scope
}

func newRegister() register {
	return register{recordmap: make(map[string]map[string]record)}
}

func (r register) addRecords(records ...record) register {
	log.Println("adding records to register: ", records)
	for _, rec := range records {
		log.Println("adding record: ", rec)
		if _, svc := r.recordmap[rec.name]; !svc {
			log.Println("adding service: ", rec.name)
			r.recordmap[rec.name] = make(map[string]record)
		}
		r.recordmap[rec.name][rec.address] = rec
		log.Println("record added: ", r)
	}
	log.Println("new register state: ", r)
	return r
}
