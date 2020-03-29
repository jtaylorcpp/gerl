package supervisor

import (
	"errors"
	"gerl/core"

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
	// returns the current pid of the supervisable process
	GetPid() *core.Pid
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
	// restart the one that fails
	ONE_FOR_ONE ChildrenStrategy = iota
	// restart all
	ONE_FOR_ALL
	// restart all after sequentially
	REST_FOR_ONE
)

type Supervisor struct {
	Children     []*Child
	Strategy     ChildrenStrategy
	ErrorIn      chan ChildError
	TerminateIn  chan bool
	TerminateOut chan bool
}

func NewSupervisor(children []*Child, childStrategy ChildrenStrategy) *Supervisor {
	s := &Supervisor{}
	s.Children = children
	s.Strategy = childStrategy
	s.TerminateIn = make(chan bool, 1)
	s.TerminateOut = make(chan bool, 1)
	s.ErrorIn = make(chan ChildError, len(children))
	return s
}

func (s *Supervisor) Start(started chan<- bool) error {
	switch s.Strategy {
	case ONE_FOR_ONE:
		s.oneForOne(started)
	case ONE_FOR_ALL:
		s.oneForAll(started)
	case REST_FOR_ONE:
		s.restForOne(started)
	}
	return nil
}

func startChildProcess(c *Child, s *Supervisor, start chan<- bool) {
	s.ErrorIn <- ChildError{
		Name:  c.Name,
		Error: c.Start(start),
	}
}

func (s *Supervisor) oneForOne(started chan<- bool) error {
	// kick off all of the children
	for _, child := range s.Children {
		childStarted := make(chan bool, 1)
		childRestarted := make(chan bool, 1)
		log.Printf("starting child %s\n", child.Name)
		go startChildProcess(child, s, childRestarted)
		go child.RestartHandler(childRestarted, childStarted, child)
		<-childStarted

		log.Printf("child %s started\n", child.Name)
	}
	started <- true

	// wait for child errors or term sigs
	log.Println("supervisor entering main loop")
	for {
		select {
		case childErr := <-s.ErrorIn:
			log.Printf("child %s terminated with error %#v\n", childErr.Name, childErr.Error)
			var childToRestart *Child
			for _, child := range s.Children {
				if child.Name == childErr.Name {
					childToRestart = child
					break
				}
			}
			log.Printf("restarting child %s\n", childToRestart.Name)
			childStarted := make(chan bool, 1)
			childRestarted := make(chan bool, 1)
			go startChildProcess(childToRestart, s, childRestarted)
			go childToRestart.RestartHandler(childRestarted, childStarted, childToRestart)
			<-childStarted
			log.Printf("child %s has been restarted\n", childToRestart.Name)

		case <-s.TerminateIn:
			log.Println("supervisor has been manually terminated")
			for _, child := range s.Children {
				child.Terminate()
			}
			if len(s.ErrorIn) != len(s.Children) {
				s.TerminateOut <- true
				close(s.TerminateOut)
				return errors.New("not all child processes have terminated")
			}

			close(s.ErrorIn)
			for cErr := range s.ErrorIn {
				if cErr.Error != nil {
					log.Printf("error from child process %s: %s\n", cErr.Name, cErr.Error.Error())
					s.TerminateOut <- true
					close(s.TerminateOut)
					return errors.New("child process did not exit cleanly")
				}
			}

			s.TerminateOut <- true
			close(s.TerminateOut)
			return nil
		}
	}
}

func (s *Supervisor) oneForAll(started chan<- bool) error {
	// kick off all of the children
	for _, child := range s.Children {
		childStarted := make(chan bool, 1)
		childRestarted := make(chan bool, 1)
		log.Printf("starting child %s\n", child.Name)
		go startChildProcess(child, s, childRestarted)
		go child.RestartHandler(childRestarted, childStarted, child)
		<-childStarted
		log.Printf("child %s started\n", child.Name)
	}
	started <- true

	// wait for child errors or term sigs
	log.Println("supervisor entering main loop")
	for {
		select {
		case childErr := <-s.ErrorIn:
			log.Printf("child %s terminated with error %#v\n", childErr.Name, childErr.Error)
			for _, child := range s.Children {
				log.Printf("restarting child %s\n", child.Name)
				log.Printf("terminating child %s\n", child.Name)

				if child.Name == childErr.Name {
					log.Printf("child %s has already been terminated\n", child.Name)
				} else if child.Process.GetPid() != nil {
					child.Terminate()
				}

				log.Printf("starting child %s\n", child.Name)
				childStarted := make(chan bool, 1)
				childRestarted := make(chan bool, 1)
				go startChildProcess(child, s, childRestarted)
				go child.RestartHandler(childRestarted, childStarted, child)
				<-childStarted
				log.Printf("child %s has been restarted\n", child.Name)
			}

			log.Println("all children processes have been restarted")

		case <-s.TerminateIn:
			log.Println("supervisor has been manually terminated")
			for _, child := range s.Children {
				child.Terminate()
			}
			if len(s.ErrorIn) != len(s.Children) {
				s.TerminateOut <- true
				close(s.TerminateOut)
				return errors.New("not all child processes have terminated")
			}

			close(s.ErrorIn)
			for cErr := range s.ErrorIn {
				if cErr.Error != nil {
					log.Printf("error from child process %s: %s\n", cErr.Name, cErr.Error.Error())
					s.TerminateOut <- true
					close(s.TerminateOut)
					return errors.New("child process did not exit cleanly")
				}
			}

			s.TerminateOut <- true
			close(s.TerminateOut)
			break
		}
	}
}

func (s *Supervisor) restForOne(started chan<- bool) error {
	// kick off all of the children
	for _, child := range s.Children {
		childStarted := make(chan bool, 1)
		childRestarted := make(chan bool, 1)
		log.Printf("starting child %s\n", child.Name)
		go startChildProcess(child, s, childRestarted)
		go child.RestartHandler(childRestarted, childStarted, child)
		<-childStarted
		log.Printf("child %s started\n", child.Name)
	}
	started <- true

	// wait for child errors or term sigs
	log.Println("supervisor entering main loop")
	select {
	case childErr := <-s.ErrorIn:
		log.Printf("child %s terminated with error %#v\n", childErr.Name, childErr.Error)
		foundErroredChild := false
		for _, child := range s.Children {
			if child.Name == childErr.Name {
				foundErroredChild = true
			}
			switch foundErroredChild {
			case true:
				log.Printf("restarting child %s\n", child.Name)
				if child.Process.GetPid() != nil {
					child.Terminate()
				}
				childStarted := make(chan bool, 1)
				childRestarted := make(chan bool, 1)
				go startChildProcess(child, s, childRestarted)
				go child.RestartHandler(childRestarted, childStarted, child)
				<-childStarted
				log.Printf("child %s has been restarted\n", child.Name)
			case false:
				// do nothing as we have not either reached the child that needs
				// to be restarted or passed it
			}
		}
		log.Println("all children processes have been restarted")

	case <-s.TerminateIn:
		log.Println("supervisor has been manually terminated")
		for _, child := range s.Children {
			child.Terminate()
		}
		if len(s.ErrorIn) != len(s.Children) {
			return errors.New("not all child processes have terminated")
		}

		close(s.ErrorIn)
		for cErr := range s.ErrorIn {
			if cErr.Error != nil {
				log.Printf("error from child process %s: %s\n", cErr.Name, cErr.Error.Error())
				return errors.New("child process did not exit cleanly")
			}
		}

		s.TerminateOut <- true
		close(s.TerminateOut)
		break
	}

	return nil
}

func (s *Supervisor) Terminate() {
	log.Println("terminating supervisor")
	s.TerminateIn <- true
	log.Println("terminating sigterm in chan")
	close(s.TerminateIn)
	<-s.TerminateOut
	log.Println("supervisor has terminated")
}

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
	RestartHandler  ChildRestartHandler
	TerminateIn     chan bool
	TerminateOut    chan bool
}

type ChildError struct {
	Name  string
	Error error
}

// ChildRestartHandler will be called as a go routine and will process
// the start events as the process restarts following its strategy
type ChildRestartHandler func(restart chan bool, firstStart chan bool, child *Child)

func defaultRestartHandler(restart chan bool, firstStart chan bool, c *Child) {
	var started bool = false
	for {
		event, ok := <-restart
		if !ok {
			log.Warningf("restarts chan for child %s has been closed", c.Name)
		}
		if !started {
			log.Infof("child %s start event recieved", c.Name)
			firstStart <- event
			started = true
		} else {
			log.Infof("child %s restart event recieved", c.Name)
		}
	}
}

func NewChild(name string, proc Supervisable, strat ProcessStrategy) *Child {
	c := &Child{
		Name:            name,
		Process:         proc,
		RestartStrategy: strat,
		TerminateIn:     make(chan bool, 1),
		TerminateOut:    make(chan bool, 1),
		RestartHandler:  defaultRestartHandler,
	}

	return c
}

func (c *Child) SetResethandler(handler ChildRestartHandler) {
	c.RestartHandler = handler
}

func (c *Child) Start(started chan<- bool) error {
	c.TerminateIn = make(chan bool, 1)
	c.TerminateOut = make(chan bool, 1)
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
		close(c.TerminateIn)
		c.Process.Terminate()
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

	if hasStarted := <-childStarted; !hasStarted {
		if restarted {
			started <- false
			return errors.New("child has restarted too many times")
		} else {
			restarted = true
			started <- false
			goto restart
		}
	}

	// child process has started
	started <- true

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
		close(c.TerminateIn)
		c.Process.Terminate()
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
			//started <- false
			started <- false
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
