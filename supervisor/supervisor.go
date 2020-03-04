package supervisor

import (
	"errors"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetReportCaller(true)
}

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

func (c *Child) Start(started chan<- bool) error {
	switch c.RestartStrategy {
	case RESTART_NEVER:
		return c.restartNever(started)
	case RESTART_ONCE:
		return c.restartOnce(started)
	case RESTART_ALWAYS:
		return c.restartAlways(started)
	default:
		started <- false
		return errors.New("unkown strategy provided to child")
	}
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

func (c *Child) restartAlways(started chan<- bool) error {
	for {
		childStarted := make(chan bool, 1)
		defer close(childStarted)
		childError := make(chan error, 1)
		defer close(childError)

		go func() {
			log.Println("starting child process")
			childError <- c.Process.Start(childStarted)
			log.Println("exiting child process")
		}()

		if hasStarted := <-childStarted; !hasStarted {
			continue
		}

		started <- true

		select {
		case err := <-childError:
			if err != nil {
				log.Printf("child process has errored out: %s\n", err.Error())
			}
			log.Println("restarting")
		case <-c.TerminateIn:
			log.Println("child directly terminated")
			c.Process.Terminate()
			log.Println("child process terminated")
			close(c.TerminateIn)
			c.TerminateOut <- true
			return nil
		}
	}
}

func (c *Child) Terminate() {
	log.Println("terminating child process")
	log.Println("sending termination signal")
	c.TerminateIn <- true
	log.Println("waiting for termination confirmation")
	<-c.TerminateOut
	log.Println("termination confirmed")
	close(c.TerminateOut)
}
