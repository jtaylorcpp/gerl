package supervisor

import (
	"errors"
	"gerl/core"
	"gerl/genserver"
	"gerl/process"
	"sync"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetReportCaller(true)
}

type ProcessStrategy uint8

const (
	// restart values assigned to a child
	RESTART_ALWAYS ProcessStrategy = iota
	RESTART_ONCE
	RESTART_NEVER
)

// Child struct describes a supervised process for a supervisor
//   Name is a human readbale label to apply to a supervised instance of a process
//   Process is the supervised process (genserver.Start, process.Start)
//   Restart strategy is one of the ito values prefixed by RESTART
//     RESTART_ONCE   - only restart a given process once (if exited by error)
//     RESTART_NEVER  - never restart a given process (if exited by error)
//     RESTART_ALWAYS - always restart a given process (if exited by error)
type ChildConfig struct {
	Name            string
	Process         Supervisable
	ProcessStrategy ProcessStrategy
}

type Child struct {
	config     *ChildConfig
	terminated bool
	mutex      *sync.Mutex
}

/*
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
*/

func NewChild(config *ChildConfig) (*Child, error) {
	c := &Child{
		config: config,
		mutex:  &sync.Mutex{},
	}

	if c.checkIfProcessNil() {
		return c, errors.New("Process(type: Supervisable) should be of type: *gerl/genserver.GenServerV2, *gerl/process.Process")
	}

	return c, nil
}

func (c *Child) Start() error {
	c.mutex.Lock()
	c.terminated = false

	switch c.config.ProcessStrategy {
	case RESTART_NEVER:
		c.mutex.Unlock()
		return c.restartNever()
	case RESTART_ONCE:
		c.mutex.Unlock()
		return c.restartOnce()
	case RESTART_ALWAYS:
		c.mutex.Unlock()
		return c.restartAlways()
	default:
		c.mutex.Unlock()
		return errors.New("unkown strategy provided to child")
	}
}

func (c *Child) GetPid() *core.Pid {
	return c.config.Process.GetPid()
}

func (c *Child) IsReady() bool {
	if c.checkIfProcessNil() {
		return false
	}

	if c.config.Process.GetPid() == nil {
		return false
	}

	return core.PidHealthCheck(c.config.Process.GetPid().GetAddr())
}

func (c *Child) checkIfProcessNil() bool {
	if c.config.Process == nil {
		return true
	}

	switch c.config.Process.(type) {
	case *genserver.GenServerV2:
		if gs, ok := c.config.Process.(*genserver.GenServerV2); ok {
			if gs == nil {
				return true
			}
		} else {
			return true
		}
	case *process.Process:
		if p, ok := c.config.Process.(*process.Process); ok {
			if p == nil {
				return true
			}
		} else {
			return true
		}
	}

	return false
}

func (c *Child) restartNever() error {
	return c.config.Process.Start()
}

func (c *Child) restartOnce() error {
	restarts := 0
	var err error = nil
	for restarts < 2 && !c.terminated {
		err = c.config.Process.Start()
		if err != nil {
			log.Errorf("error recived from child process RESTART_ONCE: %s\n", err.Error())
		}

		restarts++
	}

	log.Println("Child process with RESTART_ONCE has restarted too many times")
	return err
}

func (c *Child) restartAlways() error {
	var err error = nil
	for !c.terminated {
		err = c.config.Process.Start()
		if err != nil {
			log.Printf("error recieved from child process RESTART_ALWAYS: %s\n", err.Error())
		}
	}

	return err
}

func (c *Child) Terminate() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// terminated is check in all startegy runs except for restart never
	c.terminated = true
	c.config.Process.Terminate()
}
