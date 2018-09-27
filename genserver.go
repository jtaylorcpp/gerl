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

type GenServer struct {
	Pid       Pid
	State     GerlState
	IsRunning bool
}

func (gs *GenServer) Init(state GerlState) Pid {
	log.Println("Initializing GenServer with state: ", state)
	gs.State = state
	gs.IsRunning = false
	pid := gs.Start()
	log.Println("GenServer available at pid: ", pid)
	return pid
}

func (gs *GenServer) Start() Pid {
	//do things to init pid
	pid := NewPid()
	gs.Pid = pid

	go func() {
		log.Printf("GenServer with pid<%v> started\n", gs.Pid)
		for {
			nextMessage := <-gs.Pid.MsgBox
			log.Printf("GenServer with pid<%v> got type<%v>\n",
				gs.Pid, nextMessage.Type)
			switch nextMessage.Type {
			case Call:
				log.Printf("pid <%v> got call\n", gs.Pid)
			case Cast:
				log.Printf("pid <%v> got cast\n", gs.Pid)
			}
		}
	}()

	return pid
}

type GeneralServerClient interface {
	Init()
	Call()
	Cast()
	Terminate()
}
