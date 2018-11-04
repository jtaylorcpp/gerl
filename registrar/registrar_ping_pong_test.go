package registrar

import (
	"log"
	"testing"

	"github.com/jtaylorcpp/gerl/core"
	gs "github.com/jtaylorcpp/gerl/genserver"
)

var REGADDR string

func defaultCast(_ core.Pid, _ core.Message, _ gs.FromAddr, s gs.State) gs.State {
	return s
}

func pingCall(pid core.Pid, msg core.Message, fromaddr gs.FromAddr, s gs.State) (core.Message, gs.State) {
	log.Println("ping description: ", msg.GetDescription())
	switch msg.GetDescription() {
	case "serve":
		//run ping
		log.Println("ping sending to pong")
		pongRecord := GetRecords(REGADDR, pid.GetAddr(), "pong")[0]
		log.Println("pong record: ", pongRecord)
		log.Printf("fromaddr<%v> pongaddr<%v> msg<%v>\n", fromaddr, pongRecord.Address, core.Message{Type: core.Message_SIMPLE, Description: "ping"})
		pong := gs.Call(gs.PidAddr(pongRecord.Address), gs.FromAddr(pid.GetAddr()), core.Message{
			Type:        core.Message_SIMPLE,
			Description: "ping",
		})
		log.Println("ping go msg: ", pong)
		return pong, s
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

	log.Println("ping server addr: ", gs1.Pid.GetAddr())
	log.Println("pong server addr: ", gs2.Pid.GetAddr())
	log.Println("registrar sever addr: ", reg.Pid.GetAddr())

	REGADDR = string(reg.Pid.GetAddr())

	pingRec := NewRecord("ping", gs1.Pid.GetAddr(), core.GlobalScope)
	t.Log("ping record to add: ", pingRec)
	pongRec := NewRecord("pong", gs2.Pid.GetAddr(), core.GlobalScope)
	t.Log("pong record to add:", pongRec)

	if AddRecords(reg.Pid.GetAddr(), "localhost", pingRec, pongRec) {
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

	if pingMsg.GetDescription() != "ping pong" {
		t.Log("ping pong and register not working")
	}

	gs1.Terminate()
	gs2.Terminate()
	reg.Terminate()
}
