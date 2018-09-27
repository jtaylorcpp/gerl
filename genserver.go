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
			nextMessage, open := <-gs.Pid.MsgBox
			if !open {
				break
			}
			log.Printf("GenServer with pid<%v> got type<%v>\n",
				gs.Pid, nextMessage.Type)
			switch nextMessage.Type {
			case Call:
				log.Printf("pid <%v> got call with payload<%v>\n",
					gs.Pid, nextMessage.Payload)
			case Cast:
				log.Printf("pid <%v> got cast with payload<%v>\n",
					gs.Pid, nextMessage.Payload)
			}
		}
		log.Printf("GenServer with pid<%v> message box closed\n", gs.Pid)
		gs.Pid.TermSig <- true
	}()

	return pid
}

func (gs *GenServer) Terminate() {
	gs.IsRunning = false
	close(gs.Pid.MsgBox)
	log.Printf("Genserver with pid<%v> terminating\n", gs.Pid)
	<-gs.Pid.TermSig
	log.Printf("GenServer with pid<%v> termianted\n", gs.Pid)
}

type GeneralServerClient interface {
	Init()
	Call()
	Cast()
	Terminate()
}
