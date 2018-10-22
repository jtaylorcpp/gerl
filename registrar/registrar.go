package registrar

import (
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

func New(state State, scope core.Scope) *Registrar {
	gensvr := genserver.NewGenServer(state, scope, registrarCallHander, registrarCastHander)

	reg := Registrar(*gensvr)

	return &reg
}
