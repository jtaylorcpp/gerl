package registrar

import (
	"github.com/jtaylorcpp/gerl/core"
	"github.com/jtaylorcpp/gerl/genserver"
)

type Registrar genserver.GenServer

func registrarCallHander(pid core.Pid, in core.Message, from genserver.FromAddr, genserver.State) (core.Message, genserver.State){}

func registrarCallHander(pid core.Pid, in core.Message, from genserver.FromAddr, genserver.State) genserver.State {}