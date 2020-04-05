package supervisor

import "gerl/core"

type Supervisable interface {
	// Start is ablocking call that returns an error once the
	// Supervisable process stops running; this can be from internal
	// error or by running Terminate
	Start() error
	// IsReady returns whether or not a process is currently
	// ready to accept traffic
	// true - is ready
	// false - is not yet ready
	IsReady() bool
	// GetPid returns the current pid of the supervisable process
	GetPid() *core.Pid
	// Terminate cleanly closes out the child process
	// If a proccess exits by error this should be signalled by the
	//    Start func returning an error
	Terminate()
}
