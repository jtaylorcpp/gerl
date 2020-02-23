package supervisor

type Supervisable interface{
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

const (
	// restart values assigned to a child
	RESTART_ALWAYS uint8 = 0
	RESTART_ONCE uint8 = 1
	RESTART_NEVER uint8 = 2
)

type Supervisor struct {
	Children []Child
}

func (s *Supervisor) New(children []Child) {
	s.Children = children
}

func (s *Supervisor) Start(chan<- bool) error {
	return nil
}

func (s *Supervisor) Terminate() {}

type Child struct {
	Name string
	Process Supervisable
	RestartStrategy uint8
}