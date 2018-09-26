package gerl

import (
	"log"
)

type GeneralServer interface {
	Init(GerlState) ProcessID
	Start() ProcessID
	CallHandler()
	CastHandler()
	Terminate()
}

type GeneralServerClient interface {
	Init()
	Call()
	Cast()
	Terminate()
}

type GenServer struct {
	Pid   ProcessID
	State interface{}
}

func (gs *GenServer) Init(state GerlState) Pid {
	log.Println("Initializing GenServer with state: ", state)
	gs.State = state
	pid := gs.Start()
	log.Println("GenServer available at pid: ", pid)
	return pid
}

func (gs *GenServer) Start() Pid {
	//do things to init pid
	pid := Pid{}
	return pid
}
