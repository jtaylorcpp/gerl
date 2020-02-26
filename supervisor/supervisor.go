package supervisor

import (
	"errors"
	"log"
)

type Supervisable interface {
	// Start is a blocking call that puts  a true/false on the channel
	// to allow the supervisor to know when the process is able to recieve requests
	// true - supervised process is ready to recieve/send messages
	// false - there was an error on startup
	// In the case of a false by the supervised process the Start func should end with an error returned
	Start(chan<- bool) error
	// Terminate cleanly closes out the child process
	// If a proccess exits by error this should be signalled by the
	//    Start func returning an error
	Terminate()
}

type ProcessStrategy uint8

const (
	// restart values assigned to a child
	RESTART_ALWAYS ProcessStrategy = iota
	RESTART_ONCE
	RESTART_NEVER
)

type ChildrenStrategy uint8

const (
	// startegy use for restarting children as they terminate/fail
	ONE_FOR_ONE ChildrenStrategy = iota
	ONE_FOR_ALL
	REST_FOR_ONE
)

type Supervisor struct {
	Children []Child
	Strategy ChildrenStrategy
}

func NewSupervisor(children []Child, childStrategy ChildrenStrategy) *Supervisor {
	s := &Supervisor{}
	s.Children = children
	s.Strategy = childStrategy

	return s
}

func (s *Supervisor) Start(chan<- bool) error {
	return nil
}

func (s *Supervisor) Terminate() {}

// Child struct describes a supervised process for a supervisor
//   Name is a human readbale label to apply to a supervised instance of a process
//   Process is the supervised process (genserver.Start, process.Start)
//   Restart strategy is one of the ito values prefixed by RESTART
//     RESTART_ONCE   - only restart a given process once (if exited by error)
//     RESTART_NEVER  - never restart a given process (if exited by error)
//     RESTART_ALWAYS - always restart a given process (if exited by error)
type Child struct {
	Name            string
	Process         Supervisable
	RestartStrategy ProcessStrategy
	Terminate       chan bool
	Termianted      chan bool
}

type ChildFailure struct {
	Name  string
	Error error
}

func NewChild(name string, proc Supervisable, strat ProcessStrategy) *Child {
	c := &Child{
		Name:            name,
		Process:         proc,
		RestartStrategy: strat,
		Terminate:       make(chan bool, 1),
		Termianted:      make(chan bool, 1),
	}

	return c
}

func (c *Child) Start() ChildFailure {
	restarted := false

restartProcess:
	log.Println("starting child process")
	childStarted := make(chan bool, 1)
	childError := make(chan error, 1)
	go func() {
		childError <- c.Process.Start(childStarted)
	}()

	for {
		select {
		case started := <-childStarted:
			if !started {
				log.Println("child process failed ot start")
				switch c.RestartStrategy {
				case RESTART_ALWAYS:
					goto restartProcess
				case RESTART_ONCE:
					if !restarted {
						restarted = true
						goto restartProcess
					} else {
						return ChildFailure{
							Name:  c.Name,
							Error: errors.New("Process has failed and already been restarted"),
						}
					}
				case RESTART_NEVER:
					return ChildFailure{
						Name:  c.Name,
						Error: errors.New("Process has failed"),
					}
				}
			}
		case err := <-childError:
			log.Println("child process has terminated")
			switch c.RestartStrategy {
			case RESTART_ALWAYS:
				log.Println("restarting child process - always")
				goto restartProcess
			case RESTART_ONCE:
				log.Println("restarting child process - once")
				if !restarted {
					restarted = true
					log.Println("first restart of child process")
					goto restartProcess
				} else {
					return ChildFailure{
						Name:  c.Name,
						Error: err,
					}
				}
			case RESTART_NEVER:
				log.Println("restarting child process - never")
				return ChildFailure{
					Name:  c.Name,
					Error: err,
				}
			}
		case _ = <-c.Terminate:
			log.Println("child has been terminated")
			switch c.RestartStrategy {
			case RESTART_ALWAYS:
				// terminate the process but forcefully close child since it
				//  will try to restart the process
				c.Process.Terminate()
				close(c.Terminate)
				c.Termianted <- true
				close(c.Termianted)
				return ChildFailure{
					Name:  c.Name,
					Error: nil,
				}
			case RESTART_ONCE:
				// set the restarted flag to allow it to naturally close out once
				//  the process is termianted
				restarted = true
				c.Process.Terminate()
				close(c.Terminate)
				c.Termianted <- true
				close(c.Termianted)
			case RESTART_NEVER:
				// since the process will never be restarted; terminate it
				c.Process.Terminate()
				close(c.Terminate)
				c.Termianted <- true
				close(c.Termianted)
			}
			c.Process.Terminate()
			close(c.Terminate)
			c.Termianted <- true
			close(c.Termianted)
		}
	}

	return ChildFailure{
		Name:  c.Name,
		Error: nil,
	}
}
