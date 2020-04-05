package supervisor

import (
	"gerl/core"
	"gerl/genserver"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

type TestStratMessage struct {
	Val int
}

type TestStratState struct {
	Iter int
}

type GS struct {
	gs *genserver.GenServer
}

func NewGS() (*GS, error) {
	config := &genserver.GenServerConfig{
		StartState:  TestStratState{Iter: 0},
		CallHandler: CallTestStrat,
		CastHandler: CastTestStrat,
		Scope:       core.LocalScope,
	}
	gs, err := genserver.NewGenServer(config)
	if err != nil {
		return nil, err
	}

	return &GS{gs: gs}, nil
}

func (g *GS) Start() error {
	return g.gs.Start()
}

func (gs *GS) GetPid() *core.Pid {
	return gs.gs.GetPid()
}

func (gs *GS) IsReady() bool {
	return gs.gs.IsReady()
}

func (g *GS) Terminate() {
	g.gs.Terminate()
}

func CallTestStrat(_ core.Pid, msg TestStratMessage, _ genserver.FromAddr, s TestStratState) (TestStratMessage, TestStratState) {
	log.Println("call test func called")
	val := s.Iter
	s.Iter++
	return TestStratMessage{
		Val: val,
	}, s
}

func CastTestStrat(_ core.Pid, msg TestStratMessage, _ genserver.FromAddr, s TestStratState) TestStratState {
	log.Println("cast test func called")
	return s
}

func TestGS(t *testing.T) {
	gs, err := NewGS()
	if err != nil {
		t.Fatal(err.Error())
	}
	childConfig := &ChildConfig{
		Name:            "test",
		Process:         gs,
		ProcessStrategy: RESTART_ALWAYS,
	}
	child, err := NewChild(childConfig)
	if err != nil {
		t.Fatal(err.Error())
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- child.Start()
	}()

	for !child.IsReady() {
		time.Sleep(100 * time.Millisecond)
	}

	t.Log("child started")

	child.Terminate()

	t.Log(<-errChan)
}

func TestSupervisorStartStopEmpty(t *testing.T) {
	supervisorConfig := &SupervisorConfig{
		Children:         []*Child{},
		ChildrenStrategy: ONE_FOR_ONE,
	}
	sup, _ := NewSupervisor(supervisorConfig)

	errorChan := make(chan error, 1)
	go func() {
		errorChan <- sup.Start()
	}()

	log.Println("waiting for supervisor to start")
	for !sup.IsReady() {
		t.Log("waiting for supervisor to be ready")
		time.Sleep(100 * time.Microsecond)
	}

	log.Println("terminating supervisor")
	sup.Terminate()
}

func TestSupervisorStartStopOneForOne(t *testing.T) {
	gs1, err := NewGS()
	if err != nil {
		t.Fatal(err.Error())
	}
	childConfig := &ChildConfig{
		Name:            "test",
		Process:         gs1,
		ProcessStrategy: RESTART_ALWAYS,
	}
	child, err := NewChild(childConfig)
	if err != nil {
		t.Fatal(err.Error())
	}
	gs2, err := NewGS()
	if err != nil {
		t.Fatal(err.Error())
	}
	child2Config := &ChildConfig{
		Name:            "test2",
		Process:         gs2,
		ProcessStrategy: RESTART_ALWAYS,
	}
	child2, err := NewChild(child2Config)
	if err != nil {
		t.Fatal(err.Error())
	}
	supervisorConfig := &SupervisorConfig{
		Children:         []*Child{child, child2},
		ChildrenStrategy: ONE_FOR_ONE,
	}
	sup, err := NewSupervisor(supervisorConfig)
	if err != nil {
		t.Fatal(err.Error())
	}
	errorChan := make(chan error, 1)
	go func() {
		errorChan <- sup.Start()
	}()

	log.Println("waiting for supervisor to start")
	for !sup.IsReady() {
		t.Log("waiting for supervisor to be ready")
		time.Sleep(100 * time.Microsecond)
	}

	// increment child 1
	testMsg, err := genserver.Call(child.config.Process.GetPid().GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if testMsg.(TestStratMessage).Val != 0 {
		t.Log("chil should have iteration of 0")
	}

	testMsg, err = genserver.Call(child.config.Process.GetPid().GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if testMsg.(TestStratMessage).Val != 1 {
		t.Log("chil should have iteration of 1")
	}

	// increment child 2
	testMsg, err = genserver.Call(child2.config.Process.GetPid().GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if testMsg.(TestStratMessage).Val != 0 {
		t.Log("child should have iteration of 0")
	}

	testMsg, err = genserver.Call(child2.config.Process.GetPid().GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if testMsg.(TestStratMessage).Val != 1 {
		t.Log("chil should have iteration of 1")
	}

	// terminate 1, supervisor should reset
	child.Terminate()
	for i := 0; i < 10; i++ {
		log.Println("testing wait")
		time.Sleep(100 * time.Millisecond)
	}

	for !sup.IsReady() {
		log.Println("waiting for supervisor to be ready")
		time.Sleep(500 * time.Microsecond)
	}

	log.Println("resuming testing")
	// check child 1
	testMsg, err = genserver.Call(child.config.Process.GetPid().GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if testMsg.(TestStratMessage).Val != 0 {
		t.Log("chil should have iteration of 0")
	}

	// check child 2
	testMsg, err = genserver.Call(child2.config.Process.GetPid().GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if testMsg.(TestStratMessage).Val != 2 {
		t.Log("chil should have iteration of 2")
	}

	log.Println("terminating supervisor")
	sup.Terminate()
}

func TestSupervisorStartStopOneForAll(t *testing.T) {
	gs1, err := NewGS()
	if err != nil {
		t.Fatal(err.Error())
	}
	childConfig := &ChildConfig{
		Name:            "test",
		Process:         gs1,
		ProcessStrategy: RESTART_ALWAYS,
	}
	child, err := NewChild(childConfig)
	if err != nil {
		t.Fatal(err.Error())
	}
	gs2, err := NewGS()
	if err != nil {
		t.Fatal(err.Error())
	}
	child2Config := &ChildConfig{
		Name:            "test2",
		Process:         gs2,
		ProcessStrategy: RESTART_ALWAYS,
	}
	child2, err := NewChild(child2Config)
	if err != nil {
		t.Fatal(err.Error())
	}
	supervisorConfig := &SupervisorConfig{
		Children:         []*Child{child, child2},
		ChildrenStrategy: ONE_FOR_ALL,
	}
	sup, err := NewSupervisor(supervisorConfig)
	if err != nil {
		t.Fatal(err.Error())
	}
	errorChan := make(chan error, 1)
	go func() {
		errorChan <- sup.Start()
	}()

	log.Println("waiting for supervisor to start")
	for !sup.IsReady() {
		t.Log("waiting for supervisor to be ready")
		time.Sleep(100 * time.Microsecond)
	}

	// increment both child, child2
	testMsg, err := genserver.Call(child.config.Process.GetPid().GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if testMsg.(TestStratMessage).Val != 0 {
		t.Log("child should have iteration of 0")
	}

	testMsg, err = genserver.Call(child2.config.Process.GetPid().GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if testMsg.(TestStratMessage).Val != 0 {
		t.Log("child2 should have iteration of 0")
	}

	child.Terminate()

	for !sup.IsReady() {
		log.Println("waiting for supervisor to be ready")
		time.Sleep(100 * time.Millisecond)
	}

	// wait
	time.Sleep(200 * time.Millisecond)

	//retest child, child2
	testMsg, err = genserver.Call(child.config.Process.GetPid().GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if testMsg.(TestStratMessage).Val != 0 {
		t.Log("child should have iteration of 0")
	}

	testMsg, err = genserver.Call(child2.config.Process.GetPid().GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if testMsg.(TestStratMessage).Val != 0 {
		t.Log("child2 should have iteration of 0")
	}

	// test terminating child2

	child2.Terminate()

	for !sup.IsReady() {
		log.Println("waiting for supervisor to be ready")
		time.Sleep(100 * time.Millisecond)
	}

	// wait
	time.Sleep(200 * time.Millisecond)

	//retest child, child2
	testMsg, err = genserver.Call(child.config.Process.GetPid().GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if testMsg.(TestStratMessage).Val != 0 {
		t.Fatal("child should have iteration of 0")
	}

	testMsg, err = genserver.Call(child2.config.Process.GetPid().GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if testMsg.(TestStratMessage).Val != 0 {
		t.Fatal("child2 should have iteration of 0")
	}

	sup.Terminate()
}

func TestSupervisorStartStopRestForOne(t *testing.T) {
	gs1, err := NewGS()
	if err != nil {
		t.Fatal(err.Error())
	}
	childConfig := &ChildConfig{
		Name:            "test",
		Process:         gs1,
		ProcessStrategy: RESTART_ALWAYS,
	}
	child, err := NewChild(childConfig)
	if err != nil {
		t.Fatal(err.Error())
	}
	gs2, err := NewGS()
	if err != nil {
		t.Fatal(err.Error())
	}
	child2Config := &ChildConfig{
		Name:            "test2",
		Process:         gs2,
		ProcessStrategy: RESTART_ALWAYS,
	}
	child2, err := NewChild(child2Config)
	if err != nil {
		t.Fatal(err.Error())
	}
	supervisorConfig := &SupervisorConfig{
		Children:         []*Child{child, child2},
		ChildrenStrategy: REST_FOR_ONE,
	}
	sup, err := NewSupervisor(supervisorConfig)
	if err != nil {
		t.Fatal(err.Error())
	}
	errorChan := make(chan error, 1)
	go func() {
		errorChan <- sup.Start()
	}()

	log.Println("waiting for supervisor to start")
	for !sup.IsReady() {
		t.Log("waiting for supervisor to be ready")
		time.Sleep(100 * time.Microsecond)
	}

	// increment both child, child2
	testMsg, err := genserver.Call(child.config.Process.GetPid().GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if testMsg.(TestStratMessage).Val != 0 {
		t.Fatal("child should have iteration of 0")
	}

	testMsg, err = genserver.Call(child2.config.Process.GetPid().GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if testMsg.(TestStratMessage).Val != 0 {
		t.Fatal("child2 should have iteration of 0")
	}

	// terminate child 1
	child.Terminate()

	for !sup.IsReady() {
		log.Println("waiting for supervisor to be ready")
		time.Sleep(100 * time.Millisecond)
	}

	// wait
	time.Sleep(200 * time.Millisecond)

	// terminating child 1 would reset both children
	testMsg, err = genserver.Call(child.config.Process.GetPid().GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if testMsg.(TestStratMessage).Val != 0 {
		t.Fatal("child should have iteration of 0")
	}

	testMsg, err = genserver.Call(child2.config.Process.GetPid().GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if testMsg.(TestStratMessage).Val != 0 {
		t.Fatal("child2 should have iteration of 0")
	}

	// terminate child2
	child2.Terminate()

	for !sup.IsReady() {
		log.Println("waiting for supervisor to be ready")
		time.Sleep(100 * time.Millisecond)
	}

	// wait
	time.Sleep(200 * time.Millisecond)

	// terminating child 2 should only rest child 2
	testMsg, err = genserver.Call(child.config.Process.GetPid().GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if testMsg.(TestStratMessage).Val != 1 {
		t.Fatal("child should have iteration of 0")
	}

	testMsg, err = genserver.Call(child2.config.Process.GetPid().GetAddr(), "localhost", TestStratMessage{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if testMsg.(TestStratMessage).Val != 0 {
		t.Fatal("child2 should have iteration of 0")
	}

	sup.Terminate()
}
