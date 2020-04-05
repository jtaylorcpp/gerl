package supervisor

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetReportCaller(true)
}

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

type SupervisorConfig struct {
	Children         []*Child
	ChildrenStrategy ChildrenStrategy
}

type Supervisor struct {
	config     *SupervisorConfig
	terminated bool
}

func NewSupervisor(config *SupervisorConfig) (*Supervisor, error) {
	s := &Supervisor{
		config: config,
	}

	return s, nil
}

func (s *Supervisor) Start() error {
	s.terminated = false
	switch s.config.ChildrenStrategy {
	case ONE_FOR_ONE:
		return s.oneForOne()
	case ONE_FOR_ALL:
		return s.oneForAll()
	case REST_FOR_ONE:
		return s.restForOne()
	default:
		return errors.New("unknown strtegy provided to supervisor")
	}
	return nil
}

func (s *Supervisor) IsReady() bool {
	for _, child := range s.config.Children {
		log.Printf("checking if child %s is ready\n", child.config.Name)
		if !child.IsReady() {
			log.Printf("child %s is not ready\n", child.config.Name)
			return false
		}
		log.Printf("child %s is ready\n", child.config.Name)
	}

	return true
}

func startChildProcess(c *Child, errorChan chan childError) {
	errorChan <- childError{
		err:       c.Start(),
		childName: c.config.Name,
	}
}

type childError struct {
	err       error
	childName string
}

func (s *Supervisor) oneForOne() error {
	// kick off all of the children
	errorChan := make(chan childError, len(s.config.Children))
	for _, child := range s.config.Children {
		go startChildProcess(child, errorChan)
		for !child.IsReady() {
			log.Printf("waiting for child %s to be ready\n", child.config.Name)
			time.Sleep(10 * time.Nanosecond)
		}
	}

	for !s.terminated {
		select {
		case childErr := <-errorChan:
			log.Printf("child %s has temrinated with error: %#v\n", childErr.childName, childErr.err)
			for _, child := range s.config.Children {
				if child.config.Name == childErr.childName {
					log.Println("restarting child")
					child.Terminate()
					go startChildProcess(child, errorChan)
					for !child.IsReady() {
						log.Printf("waiting for child %s to be ready\n", child.config.Name)
						time.Sleep(100 * time.Nanosecond)
					}
					log.Printf("child %s has been restarted\n", child.config.Name)
				}
			}
		}
	}

	return nil
}

func (s *Supervisor) oneForAll() error {
	errorChan := make(chan childError, 1)
	defer close(errorChan)
	for _, child := range s.config.Children {
		go startChildProcess(child, errorChan)
		for !child.IsReady() {
			log.Printf("waiting for child %s to be ready\n", child.config.Name)
			time.Sleep(10 * time.Nanosecond)
		}
	}

	for !s.terminated {
		select {
		case childErr := <-errorChan:
			log.Printf("child %s has temrinated with error: %#v\n", childErr.childName, childErr.err)
			for _, child := range s.config.Children {
				if child.config.Name == childErr.childName {
					// this child should already have been terminated
					child.Terminate() // run anyways
					go startChildProcess(child, errorChan)
					for !child.IsReady() {
						log.Printf("waiting for child %s to be ready\n", child.config.Name)
						time.Sleep(10 * time.Nanosecond)
					}
				} else {
					child.Terminate()
					currentErr := <-errorChan
					if currentErr.childName != child.config.Name {
						log.Errorf("supervisor ALL_FOR_ONE tryig to restart child %s but got error from child %s\n", child.config.Name, childErr.childName)
						goto terminate
					}
					go startChildProcess(child, errorChan)
					for !child.IsReady() {
						log.Printf("waiting for child %s to be ready\n", child.config.Name)
						time.Sleep(10 * time.Nanosecond)
					}
				}
			}
		}
	}
terminate:
	for _, child := range s.config.Children {
		if child.IsReady() {
			child.Terminate()
			log.Printf("supervisor ALL_FOR_ONE closing child %s with output error %#v\n", child.config.Name, <-errorChan)
		} else {
			// not already running
			child.Terminate()
			log.Printf("supervisor ALL_FOR_ONE child %s already terminated\n", child.config.Name)
		}
	}
	return nil
}

func (s *Supervisor) restForOne() error {
	errorChan := make(chan childError, 1)
	defer close(errorChan)
	for _, child := range s.config.Children {
		go startChildProcess(child, errorChan)
		for !child.IsReady() {
			log.Printf("waiting for child %s to be ready\n", child.config.Name)
			time.Sleep(10 * time.Nanosecond)
		}
	}

	for !s.terminated {
		select {
		case childErr := <-errorChan:
			log.Printf("child %s has temrinated with error: %#v\n", childErr.childName, childErr.err)
			terminateChild := false
			for _, child := range s.config.Children {
				if child.config.Name == childErr.childName {
					// this child should already have been terminated
					child.Terminate() // run anyways
					go startChildProcess(child, errorChan)
					for !child.IsReady() {
						log.Printf("waiting for child %s to be ready\n", child.config.Name)
						time.Sleep(10 * time.Nanosecond)
					}

					terminateChild = true
				}
				// all children after the terminated child
				if terminateChild {
					child.Terminate()
					currentErr := <-errorChan
					if currentErr.childName != child.config.Name {
						log.Errorf("supervisor REST_FOR_ONE trying to restart child %s but got error from child %s\n", child.config.Name, childErr.childName)
						goto terminate
					}
					go startChildProcess(child, errorChan)
					for !child.IsReady() {
						log.Printf("waiting for child %s to be ready\n", child.config.Name)
						time.Sleep(10 * time.Nanosecond)
					}
				}
			}
		}
	}
terminate:
	for _, child := range s.config.Children {
		if child.IsReady() {
			child.Terminate()
			log.Printf("supervisor REST_FOR_ONE closing child %s with output error %#v\n", child.config.Name, <-errorChan)
		} else {
			// not already running
			child.Terminate()
			log.Printf("supervisor REST_FOR_ONE child %s already terminated\n", child.config.Name)
		}
	}
	return nil
}

func (s *Supervisor) Terminate() {
	s.terminated = true

	for _, child := range s.config.Children {
		child.Terminate()
		for child.IsReady() {
			log.Printf("waiting for child %s to close out\n", child.config.Name)
			time.Sleep(10 * time.Nanosecond)
		}
	}
}
