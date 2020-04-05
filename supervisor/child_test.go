package supervisor

import (
	"gerl/core"
	"gerl/genserver"
	"gerl/process"
	"testing"
	"time"
)

func TestChildEmptyNone(t *testing.T) {
	config := &ChildConfig{}

	child, _ := NewChild(config)

	if child.checkIfProcessNil() != true {
		t.Fatal("child has no process and should be nil")
	}
}

func TestChildEmptyGenServer(t *testing.T) {
	config := &ChildConfig{
		Process: &genserver.GenServerV2{},
	}

	child, _ := NewChild(config)

	if child.checkIfProcessNil() == true {
		t.Fatal("child was given empty genserver")
	}

	if child.config.Process.IsReady() == true {
		t.Fatal("child genserver has no pid config")
	}
}

func TestChildEmptyProcess(t *testing.T) {
	config := &ChildConfig{
		Process: &process.Process{},
	}

	child, _ := NewChild(config)

	if child.checkIfProcessNil() == true {
		t.Fatal("child was given empty process")
	}

	if child.config.Process.IsReady() == true {
		t.Fatal("child process has no pid config")
	}
}

func TestChildGenServer(t *testing.T) {
	gsConfig := genserver.GenServerV2Config{
		StartState:  TestState{},
		CallHandler: CallTest,
		CastHandler: CastTest,
		Scope:       core.LocalScope,
	}

	gs, _ := genserver.NewGenServerV2(&gsConfig)

	childConfig := &ChildConfig{
		Name:            "test gs",
		Process:         gs,
		ProcessStrategy: RESTART_NEVER,
	}

	child, err := NewChild(childConfig)

	if err != nil {
		t.Fatal("error should not have been thrown since GS is a pointer")
	}

	if child.IsReady() == true {
		t.Fatal("child has not been started")
	}

	errorChan := make(chan error, 1)
	go func() {
		errorChan <- child.Start()
	}()

	for !child.IsReady() {
		t.Log("waiting for child/genserver to start")
		time.Sleep(1 * time.Second)
	}

	child.config.Process.Terminate()

	exitError := <-errorChan
	if exitError != nil {
		t.Fatal("genserver was closed out properly")
	}

	child.Terminate()

	go func() {
		errorChan <- child.Start()
	}()

	for !child.IsReady() {
		t.Log("waiting for child/genserver to start")
		time.Sleep(1 * time.Second)
	}

	child.Terminate()

	exitError = <-errorChan

	if exitError != nil {
		t.Fatal("genserver was closed out properly")
	}
}

func TestChildRestartNever(t *testing.T) {
	gsConfig := genserver.GenServerV2Config{
		StartState:  TestState{},
		CallHandler: CallTest,
		CastHandler: CastTest,
		Scope:       core.LocalScope,
	}

	gs, _ := genserver.NewGenServerV2(&gsConfig)

	childConfig := &ChildConfig{
		Name:            "test gs",
		Process:         gs,
		ProcessStrategy: RESTART_NEVER,
	}

	child, _ := NewChild(childConfig)

	errorChan := make(chan error, 2)
	go func() {
		errorChan <- child.Start()
	}()

	for !child.IsReady() {
		t.Log("waiting for child/genserver to start")
		time.Sleep(50 * time.Microsecond)
	}

	child.config.Process.GetPid().Terminate()

	for child.IsReady() {
		t.Log("waiting for process to exit")
		time.Sleep(50 * time.Microsecond)
	}

	if len(errorChan) != 1 {
		t.Log("error chan len: ", len(errorChan))
		t.Fatal("child should have already exited")
	}

	child.config.Process.GetPid().Terminate()

	for child.IsReady() {
		t.Log("waiting for process to exit")
		time.Sleep(50 * time.Microsecond)
	}

	if len(errorChan) != 1 {
		t.Fatal("child should have already exited")
	}

	t.Log(<-errorChan)
}

func TestChildRestartOnce(t *testing.T) {
	gsConfig := genserver.GenServerV2Config{
		StartState:  TestState{},
		CallHandler: CallTest,
		CastHandler: CastTest,
		Scope:       core.LocalScope,
	}

	gs, _ := genserver.NewGenServerV2(&gsConfig)

	childConfig := &ChildConfig{
		Name:            "test gs",
		Process:         gs,
		ProcessStrategy: RESTART_ONCE,
	}

	child, _ := NewChild(childConfig)

	errorChan := make(chan error, 2)
	go func() {
		errorChan <- child.Start()
	}()

	for !child.IsReady() {
		t.Log("waiting for child/genserver to start")
		time.Sleep(50 * time.Microsecond)
	}

	child.config.Process.GetPid().Terminate()

	for !child.IsReady() {
		t.Log("waiting for process to restart")
		time.Sleep(50 * time.Microsecond)
	}

	if len(errorChan) != 0 {
		t.Fatal("child should not have exited")
	}

	child.config.Process.GetPid().Terminate()

	for child.IsReady() {
		t.Log("waiting for process to exit")
		time.Sleep(50 * time.Microsecond)
	}

	// there is a race condition in this check on whether or not
	// the message propogates fast enough
	time.Sleep(1 * time.Second)
	if len(errorChan) != 1 {
		t.Log("len of chan: ", len(errorChan))
		t.Fatal("child should have exited")
	}

	child.config.Process.GetPid().Terminate()

	for child.IsReady() {
		t.Log("waiting for process to exit")
		time.Sleep(100 * time.Microsecond)
	}

	// weird case where the IsReady shows its closed but
	// we check for error before message is propogated
	time.Sleep(100 * time.Microsecond)

	if len(errorChan) != 1 {
		t.Log("len of chan: ", len(errorChan))
		t.Fatal("child should have exited")
	}

	t.Log(<-errorChan)
}

func TestChildRestartAlways(t *testing.T) {
	gsConfig := genserver.GenServerV2Config{
		StartState:  TestState{},
		CallHandler: CallTest,
		CastHandler: CastTest,
		Scope:       core.LocalScope,
	}

	gs, _ := genserver.NewGenServerV2(&gsConfig)

	childConfig := &ChildConfig{
		Name:            "test gs",
		Process:         gs,
		ProcessStrategy: RESTART_ALWAYS,
	}

	child, _ := NewChild(childConfig)

	errorChan := make(chan error, 2)
	go func() {
		errorChan <- child.Start()
	}()

	for !child.IsReady() {
		t.Log("waiting for child/genserver to start")
		time.Sleep(50 * time.Microsecond)
	}

	for i := 0; i < 10; i++ {
		child.config.Process.GetPid().Terminate()

		for !child.IsReady() {
			t.Log("waiting for process to restart")
			time.Sleep(50 * time.Microsecond)
		}

		if len(errorChan) != 0 {
			t.Fatal("child should not have exited")
		}
	}

	child.Terminate()

	t.Log(<-errorChan)
}
