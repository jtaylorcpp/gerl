package process

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	"gerl/core"
)

func init() {
	log.SetReportCaller(true)
}

type PidAddr = string
type FromAddr = string
type Inbox chan core.GerlMsg

// Handler started as a go-routine to process incoming messages
type ProcHandler func(PidAddr, core.Message) error

type ProcessConfig struct {
	Scope   core.Scope
	Port    string
	Address string
	Handler ProcHandler
}

// Process is the struct designed to be a single threaded loop
// to process core.GerlMsg
type Process struct {
	Config *ProcessConfig
	// Pid use to communicate with the Process
	Pid   *core.Pid
	Scope core.Scope
	// Handler is the Inbox handler that is started as a go-routine
	Handler ProcHandler
}

// Builds a new Process that has yet to be started
func New(config *ProcessConfig) *Process {
	return &Process{
		Config:  config,
		Handler: config.Handler,
	}
}

// Starts a process which is blocking until an error is reported.
// The main thread processes all incoming Ccore.GerlMsg, and errors from both the
// Process and Pid.
func (p *Process) Start() error {
	pid, err := core.NewPid(p.Config.Address, p.Config.Port, p.Config.Scope)
	if err != nil {
		return err
	}
	p.Pid = pid

	log.Println("Process available at addr: ", p.Pid.GetAddr())

	log.Printf("process with pid<%v> entering main loop\n", p.Pid.GetAddr())

	var returnErr error = nil
	for {
		select {
		case msg, ok := <-p.Pid.Inbox:
			if !ok {
				log.Println("inbox has been closed but process not terminated")
				returnErr = errors.New("process inbox has been closed")
				goto terminate
			} else {
				switch msg.Type {
				case core.GerlMsg_PROC:
					err := p.Handler(msg.Fromaddr, *msg.Msg)
					if err != nil {
						log.Errorf("process handler returned error: %s\n", err.Error())
						returnErr = err
						goto terminate
					}
				case core.GerlMsg_TERM:
					log.Println("process recieved terminate signal")
					goto terminate
				default:
					log.Printf("unknown message: %#v\n", msg)
				}
			}
		}
	}
terminate:
	log.Println("process is terminating")
	p.Pid.Terminate()
	p.Pid = nil
	log.Println("process has terminated")
	return returnErr
}

func (p *Process) GetPid() *core.Pid {
	return p.Pid
}

func (p *Process) IsReady() bool {
	if p.Pid == nil {
		return false
	}

	return core.PidHealthCheck(p.GetPid().GetAddr())
}

// Terminates all of the Process side channels. Terminates the Pid and clears
// all resulting errors.
func (p *Process) Terminate() {
	log.Printf("process with pid<%v> terminating\n", p.Pid.GetAddr())
	core.PidTerminate(p.GetPid().GetAddr(), p.GetPid().GetAddr(), core.Message{})
	for p.Pid != nil {
		time.Sleep(10 * time.Nanosecond)
	}
}

// Send sends an arbitrary core.Message to a Process at PidAddr
func Send(to PidAddr, from FromAddr, msg core.Message) {
	core.PidSendProc(to, from, msg)
}
