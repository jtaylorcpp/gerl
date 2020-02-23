package process

import (
	"errors"
	"log"

	"gerl/core"
)

type PidAddr = string
type FromAddr = string
type Inbox chan core.GerlMsg

// Handler started as a go-routine to process incoming messages
type ProcHandler func(PidAddr, Inbox) error

// Process is the struct designed to be a single threaded loop
// to process core.GerlMsg
type Process struct {
	// Pid use to communicate with the Process
	Pid   *core.Pid
	Scope core.Scope
	// Handler is the Inbox handler that is started as a go-routine
	Handler ProcHandler
	// Error channel that forces a termination when an error is sent
	Errors chan error
	// Terminate channel that forces the main loop to terminate with termination error
	Terminated chan bool
}

// Builds a new Process that has yet to be started
func New(scope core.Scope, handler ProcHandler) *Process {
	return &Process{
		Pid:        &core.Pid{},
		Scope:      scope,
		Handler:    handler,
		Errors:     make(chan error, 2),
		Terminated: make(chan bool, 1),
	}
}

// Starts a process which is blocking until an error is reported.
// The main thread processes all incoming Ccore.GerlMsg, and errors from both the
// Process and Pid.
func (p *Process) Start(started chan<- bool) error {
	var err error
	p.Pid, err = core.NewPid("", "", p.Scope)
	if err != nil {
		started <- false
		return err
	}

	log.Println("Process available at addr: ", p.Pid.GetAddr())

	in := make(Inbox, 1)

	go func() {
		p.Errors <- p.Handler(PidAddr(p.Pid.GetAddr()), in)
	}()

	log.Printf("process with pid<%v> entering main loop\n", p.Pid.GetAddr())

	started <- true
	for {
		select {
		case err := <-p.Pid.Errors:
			log.Println("process pid error: ", err)
			return errors.New("pid error, close process")
		case err := <-p.Errors:
			log.Println("process error, close process")
			return err
		case <-p.Terminated:
			log.Println("process terminated")
			return nil
		case msg, ok := <-p.Pid.Inbox:
			log.Println("process message from inbox")
			if !ok {
				return errors.New("process inbox closed")
			}
			switch msg.GetType() {
			case core.GerlMsg_PROC:
				log.Println("process recieved proc message")
				in <- msg
			default:
				log.Println("process recieved unknown type")
			}
		}
	}

	return nil
}

// Terminates all of the Process side channels. Terminates the Pid and clears
// all resulting errors.
func (p *Process) Terminate() {
	log.Printf("process with pid<%v> terminating\n", p.Pid.GetAddr())
	p.Terminated <- true
	close(p.Terminated)
	p.Pid.Terminate()
	for {
		err, ok := <-p.Pid.Errors
		if !ok {
			break
		}
		log.Println("process clearing pid errors: ", err)
	}
	log.Printf("process with pid<%v> temrinated\n", p.Pid.GetAddr())
}

// Send sends an arbitrary core.Message to a Process at PidAddr
func Send(to PidAddr, from FromAddr, msg core.Message) {
	core.PidSendProc(to, from, msg)
}
