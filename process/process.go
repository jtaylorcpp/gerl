package process

import (
	"errors"
	"log"

	"github.com/jtaylorcpp/gerl/core"
)

/*
type Process interface {
	// Starts the GenericServer and returns the ProcessID associated with it
	Start() error
	// Processes a synchronous message passed to the GenericServer
	Handler(core.Message, FromAddr, State) (core.Message, State)
	// Terminate closes the ProcessID and clears out the Process
	Terminate()
}
*/

type PidAddr string
type FromAddr string
type Inbox chan core.GerlMsg

type ProcHandler func(PidAddr, Inbox) error

type Process struct {
	Pid        *core.Pid
	Handler    ProcHandler
	Errors     chan error
	Terminated chan bool
}

func New(handler ProcHandler) *Process {
	return &Process{
		Pid:        &core.Pid{},
		Handler:    handler,
		Errors:     make(chan error, 2),
		Terminated: make(chan bool, 1),
	}
}

func (p *Process) Start() error {

	p.Pid = core.NewPid("", "")

	log.Println("Process available at addr: ", p.Pid.GetAddr())

	in := make(Inbox, 32)

	go func() {
		p.Errors <- p.Handler(PidAddr(p.Pid.GetAddr()), in)
	}()

	log.Printf("process with pid<%v> entering main loop\n", p.Pid.GetAddr())

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
			return errors.New("process terminated")
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
				log.Println("process recieved unknonw type")
			}
		default:
			continue
		}
	}

	return errors.New("process ended in error")
}

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

func Send(to PidAddr, from FromAddr, msg core.Message) {
	core.PidSendProc(string(to), string(from), msg)
}
