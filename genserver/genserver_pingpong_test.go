package genserver

import (
	"log"
	"testing"
	//"time"

	"gerl/core"
)

var PongAddr PidAddr

type Ball struct {
	NextAddr string
	Action string
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
			Action: "ping",
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
			Action: "pong",
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
	gs1, err := NewGenServer(PingPongState{}, core.LocalScope, pingCall, defaultCast)
	if err != nil {
		t.Fatal(err.Error())
	}
	gs2, err := NewGenServer(PingPongState{}, core.GlobalScope, pongCall, defaultCast)
	if err != nil {
		t.Fatal(err.Error())
	}
	gs1Started := make(chan bool, 1)
	gs2Started := make(chan bool, 1)
	go func() {
		t.Log(gs1.Start(gs1Started))
	}()

	go func() {
		t.Log(gs2.Start(gs2Started))
	}()

	<- gs1Started
	<- gs2Started

	log.Println("ping server addr: ", gs1.Pid.GetAddr())
	log.Println("pong server addr: ", gs2.Pid.GetAddr())
	

	log.Println("test pong routine")
	pong1 := Ball{
		Action: "ping",
	}

	returnPong1, err := Call( gs2.Pid.GetAddr(), "localhost", pong1)
	if err != nil {
		t.Fatal(err.Error())
	}

	if returnPong1.(Ball).Action != "pong" {
		t.Fatal("didnt get pong back when ping was sent")
	}

	serve1 := Ball{
		Action: "serve",
		NextAddr: gs2.Pid.GetAddr(),
	}

	returnPong2, err := Call(gs1.Pid.GetAddr(), "localhost", serve1)
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Logf("%#v\n", returnPong2)

	if returnPong2.(Ball).Action != "pong" {
		t.Fatal("should have gotten a pong back")
	}
	/*rmsg1 := Call(gs2.Pid.GetAddr(), "localhost", Ball {
		Action: "ping",
		Description: "ping",
	})

	t.Log("pong test: ", rmsg1)

	if rmsg1.GetDescription() != "ping pong" {
		t.Fatal("pong test failed")
	}

	rmsg2 := Call(PidAddr(gs1.Pid.GetAddr()), FromAddr("localhost"), core.Message{
		Type:        core.Message_SIMPLE,
		Description: "serve",
	})

	t.Log("ping test: ", rmsg2)

	if rmsg2.GetDescription() != "ping pong" {
		t.Fatal("ping serve test failed")
	}*/

	gs1.Terminate()
	gs2.Terminate()
}

