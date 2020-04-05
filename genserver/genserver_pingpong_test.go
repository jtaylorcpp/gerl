package genserver

import (
	"testing"
	"time"

	log "github.com/sirupsen/logrus"

	"gerl/core"
)

var PongAddr PidAddr

type Ball struct {
	NextAddr string
	Action   string
}

type PingPongState struct {
	Actions []string
}

func defaultCast(_ core.Pid, _ struct{}, _ FromAddr, s struct{}) struct{} {
	return s
}

func pingCall(pid core.Pid, ball Ball, fromaddr FromAddr, s PingPongState) (Ball, PingPongState) {
	log.Printf("ping description: %#v\n", ball)
	switch ball.Action {
	case "serve":
		//run ping
		s.Actions = append(s.Actions, "server")
		log.Println("ping sending to pong")
		nextBallAction := Ball{
			NextAddr: pid.GetAddr(),
			Action:   "ping",
		}
		log.Printf("fromaddr<%v> pongaddr<%v> msg<%v>\n", fromaddr, ball.NextAddr, nextBallAction)
		pong, err := Call(ball.NextAddr, pid.GetAddr(), nextBallAction)
		if err != nil {
			log.Println("error from call client: ", err.Error())
			return Ball{}, s
		}
		log.Println("ping got msg: ", pong)
		s.Actions = append(s.Actions, "recieved pong")
		return pong.(Ball), s
	default:
		log.Println("ping unknown message: ", ball)
		return Ball{}, s
	}
}

func pongCall(_ core.Pid, ball Ball, fromaddr FromAddr, s PingPongState) (Ball, PingPongState) {
	log.Println("pong description: ", ball.Action)
	switch ball.Action {
	case "ping":
		log.Println("pong got ping")
		nextBallAction := Ball{
			Action:   "pong",
			NextAddr: "",
		}

		s.Actions = append(s.Actions, "recieved ping")
		return nextBallAction, s
	default:
		log.Println("pong unknown message: ", ball)
		return Ball{}, s
	}
}

func TestGenServers(t *testing.T) {
	config1 := &GenServerConfig{
		StartState:  PingPongState{},
		Scope:       core.LocalScope,
		CallHandler: pingCall,
		CastHandler: defaultCast,
	}
	gs1, err := NewGenServer(config1)
	if err != nil {
		t.Fatal(err.Error())
	}

	config2 := &GenServerConfig{
		StartState:  PingPongState{},
		Scope:       core.LocalScope,
		CallHandler: pongCall,
		CastHandler: defaultCast,
	}
	gs2, err := NewGenServer(config2)
	if err != nil {
		t.Fatal(err.Error())
	}

	go func() {
		log.Errorln("error exiting ping server: ", gs1.Start())
	}()

	go func() {
		log.Errorln("error exiting ping server: ", gs2.Start())
	}()

	for !gs1.IsReady() {
		time.Sleep(10 * time.Millisecond)
	}

	for !gs2.IsReady() {
		time.Sleep(10 * time.Millisecond)
	}

	log.Println("ping server addr: ", gs1.pid.GetAddr())
	log.Println("pong server addr: ", gs2.pid.GetAddr())

	log.Println("test pong routine")
	pong1 := Ball{
		Action: "ping",
	}

	returnPong1, err := Call(gs2.pid.GetAddr(), "localhost", pong1)
	if err != nil {
		t.Fatal(err.Error())
	}

	if returnPong1.(Ball).Action != "pong" {
		t.Fatal("didnt get pong back when ping was sent")
	}

	serve1 := Ball{
		Action:   "serve",
		NextAddr: gs2.pid.GetAddr(),
	}

	returnPong2, err := Call(gs1.pid.GetAddr(), "localhost", serve1)
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Logf("%#v\n", returnPong2)

	if returnPong2.(Ball).Action != "pong" {
		t.Fatal("should have gotten a pong back")
	}

	gs1.Terminate()
	gs2.Terminate()
}
