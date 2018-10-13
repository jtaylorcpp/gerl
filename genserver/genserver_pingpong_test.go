package genserver

import (
	"log"
	"testing"
	"time"

	"github.com/jtaylorcpp/gerl/core"
)

var PongAddr PidAddr

func defaultCast(_ PidAddr, _ core.Message, s State) State {
	return s
}

func pingCall(myAddr PidAddr, msg core.Message, s State) (core.Message, State) {
	log.Println("ping description: ", msg.GetDescription())
	switch msg.GetDescription() {
	case "serve":
		//run ping
		log.Println("ping sending to pong")
		log.Printf("myaddr<%v> pongaddr<%v> msg<%v>\n", myAddr, PongAddr, core.Message{Type: core.Message_SIMPLE, Fromaddr: string(myAddr), Description: "ping"})
		pong := Call(PongAddr, myAddr, core.Message{
			Type:        core.Message_SIMPLE,
			Fromaddr:    string(myAddr),
			Description: "ping",
		})
		log.Println("ping go msg: ", pong)
		return pong, s
	default:
		log.Println("ping unknown message: ", msg)
		return core.Message{}, s
	}
}

func pongCall(myAddr PidAddr, msg core.Message, s State) (core.Message, State) {
	log.Println("pong description: ", msg.GetDescription())
	switch msg.GetDescription() {
	case "ping":
		log.Println("pong got ping")
		desc := msg.GetDescription() + " pong"
		return core.Message{Type: core.Message_SIMPLE, Fromaddr: string(myAddr), Description: desc}, s
	default:
		log.Println("pong unknown message: ", msg)
		return core.Message{}, s
	}
}

func TestGenServers(t *testing.T) {
	gs1 := NewGenServer("genserver 1", pingCall, defaultCast)
	gs2 := NewGenServer("genserver 2", pongCall, defaultCast)

	go func() {
		t.Log(gs1.Start())
	}()

	go func() {
		t.Log(gs2.Start())
	}()

	time.Sleep(50 * time.Millisecond)

	PongAddr = PidAddr(gs2.Pid.GetAddr())

	log.Println("ping server addr: ", gs1.Pid.GetAddr())
	log.Println("pong server addr: ", gs2.Pid.GetAddr())
	log.Println("var pong server addr: ", PongAddr)

	rmsg1 := Call(PidAddr(gs2.Pid.GetAddr()), PidAddr("localhost"), core.Message{
		Type:        core.Message_SIMPLE,
		Fromaddr:    "localhost",
		Description: "ping",
	})

	t.Log("pong test: ", rmsg1)

	if rmsg1.GetDescription() != "ping pong" {
		t.Fatal("pong test failed")
	}

	rmsg2 := Call(PidAddr(gs1.Pid.GetAddr()), PidAddr("localhost"), core.Message{
		Type:        core.Message_SIMPLE,
		Fromaddr:    "localhost",
		Description: "serve",
	})

	t.Log("ping test: ", rmsg2)

	if rmsg2.GetDescription() != "ping pong" {
		t.Fatal("ping serve test failed")
	}

	gs1.Terminate()
	gs2.Terminate()
}
