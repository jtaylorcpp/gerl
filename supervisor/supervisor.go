package supervisor

import (
	"errors"
	log "github.com/sirupsen/logrus"
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
	TerminateIn     chan bool
	TerminateOut    chan bool
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
		TerminateIn:     make(chan bool, 1),
		TerminateOut:    make(chan bool, 1),
	}

	return c
}

func (c *Child) restartNever(started chan<- bool) error {
	childStarted := make(chan bool, 1)
	defer close(childStarted)
	ChildError := make(chan error, 1)
	defer close(ChildError)
	go func() {
		ChildError <- c.Process.Start(childStarted)
	}()

	started <- <-childStarted

	select {
	case err := <-ChildError:
		log.Println("child process terminated")
		return err
	case <-c.TerminateIn:
		log.Println("terminating child process")
		c.Process.Terminate()
		close(c.TerminateIn)
		c.TerminateOut <- true
		return nil
	}
}

func (c *Child) restartOnce(started chan<- bool) error {
	restarted := false
restart:
	childStarted := make(chan bool, 1)
	defer close(childStarted)
	childError := make(chan error, 1)
	defer close(childError)

	go func() {
		childError <- c.Process.Start(childStarted)
	}()

	started <- true

	if hasStarted := <-childStarted; !hasStarted {
		if restarted {
			return errors.New("child has restarted too many times")
		} else {
			restarted = true
			goto restart
		}
	}

	select {
	case err := <-childError:
		if err != nil {
			log.Printf("child process returned error: %s\n", err.Error())
		}

		if restarted {
			log.Println("child has restarted too many times")
			return err
		} else {
			restarted = true
			goto restart
		}
	case <-c.TerminateIn:
		log.Println("child directly terminated")
		c.Process.Terminate()
		close(c.TerminateIn)
		c.TerminateOut <- true
		return nil
	}
}

func (c *Child) restartAlways(started chan<- bool) {}

/*func (c *Child) Start(started chan<- bool) ChildFailure {
	restarted := false

restartProcess:
	c.TerminateIn = make(chan bool, 1)
	c.TerminateOut = make(chan bool, 1)
	childStarted := make(chan bool, 1)
	childError := make(chan error, 1)
	log.Println("starting child process")

	go func() {
		childError <- c.Process.Start(childStarted)
	}()

	if !<-childStarted {
		log.Println("child process failed to start")
		started <- false
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
	started <- true

	log.Println("child started")
	for {
		select {
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
				close(c.TerminateIn)
				c.TerminateOut <- true
				return ChildFailure{
					Name:  c.Name,
					Error: err,
				}
			}
		case <-c.TerminateIn:
			log.Println("child has been terminated")
			goto closeout
		}
	}
closeout:
	log.Println("child closing out")
	c.Process.Terminate()
	close(c.TerminateIn)
	c.TerminateOut <- true
	return ChildFailure{
		Name:  c.Name,
		Error: nil,
	}
}

func (c *Child) Terminate() {
	c.TerminateIn <- true
	<-c.TerminateOut
	close(c.TerminateOut)
}
*/
