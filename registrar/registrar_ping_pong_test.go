package registrar

import (
	"log"
	"testing"
	"time"

	"github.com/jtaylorcpp/gerl/core"
	gs "github.com/jtaylorcpp/gerl/genserver"
)

func defaultCast(_ core.Pid, _ core.Message, _ gs.FromAddr, s gs.State) gs.State {
	return s
}

func pingCall(pid core.Pid, msg core.Message, fromaddr gs.FromAddr, s gs.State) (core.Message, gs.State) {
	log.Println("ping description: ", msg.GetDescription())
	switch msg.GetDescription() {
	/*case "serve":

	//run ping
	log.Println("ping sending to pong")
	log.Printf("fromaddr<%v> pongaddr<%v> msg<%v>\n", fromaddr, PongAddr, core.Message{Type: core.Message_SIMPLE, Description: "ping"})
	pong := Call(PongAddr, FromAddr(pid.GetAddr()), core.Message{
		Type:        core.Message_SIMPLE,
		Description: "ping",
	})
	log.Println("ping go msg: ", pong)
	return pong, s
	*/
	default:
		log.Println("ping unknown message: ", msg)
		return core.Message{}, s
	}
}

func pongCall(_ core.Pid, msg core.Message, fromaddr gs.FromAddr, s gs.State) (core.Message, gs.State) {
	log.Println("pong description: ", msg.GetDescription())
	switch msg.GetDescription() {
	case "ping":
		log.Println("pong got ping")
		desc := msg.GetDescription() + " pong"
		return core.Message{Type: core.Message_SIMPLE, Description: desc}, s
	default:
		log.Println("pong unknown message: ", msg)
		return core.Message{}, s
	}
}

func TestRegistrarPingPong(t *testing.T) {
	gs1 := gs.NewGenServer("genserver 1", core.LocalScope, pingCall, defaultCast)
	gs2 := gs.NewGenServer("genserver 2", core.GlobalScope, pongCall, defaultCast)
	reg := NewRegistrar(core.LocalScope)

	go func() {
		t.Log(gs1.Start())
	}()

	go func() {
		t.Log(gs2.Start())
	}()

	go func() {
		t.Log(reg.Start())
	}()

	for !core.PidHealthCheck(gs1.Pid.GetAddr()) || !core.PidHealthCheck(gs2.Pid.GetAddr()) || !core.PidHealthCheck(reg.Pid.GetAddr()) {
		time.Sleep(25 * time.Microsecond)
		t.Log("waiting for genserver to start")
	}

	log.Println("ping server addr: ", gs1.Pid.GetAddr())
	log.Println("pong server addr: ", gs2.Pid.GetAddr())
	log.Println("registrar sever addr: ", reg.Pid.GetAddr())

	if AddRecords(reg.Pid.GetAddr(), "localhost",
		NewRecord("ping", gs1.Pid.GetAddr(), core.GlobalScope),
		NewRecord("pong", gs2.Pid.GetAddr(), core.GlobalScope)) {
		t.Log("added ping and pong to registrar")
	} else {
		t.Fatal("ping and poing not added to registrar")
	}

	pingRecord := GetRecords(reg.Pid.GetAddr(), "localhost", "ping")[0]

	t.Log("address for ping from registrar: ", pingRecord)

	pingMsg := gs.Call(gs.PidAddr(pingRecord.Address), "localhost", core.Message{
		Type:        core.Message_SIMPLE,
		Description: "serve",
	})

	t.Log("msg recieved from serve: ", pingMsg)

	/*

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
		}
	*/

	gs1.Terminate()
	gs2.Terminate()
	reg.Terminate()
}
