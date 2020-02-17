package genserver

/*import (
	"log"
	"testing"
	"time"

	"github.com/jtaylorcpp/gerl/core"
)

var PongAddr PidAddr

func defaultCast(_ core.Pid, _ core.Message, _ FromAddr, s State) State {
	return s
}

func pingCall(pid core.Pid, msg core.Message, fromaddr FromAddr, s State) (core.Message, State) {
	log.Println("ping description: ", msg.GetDescription())
	switch msg.GetDescription() {
	case "serve":
		//run ping
		log.Println("ping sending to pong")
		log.Printf("fromaddr<%v> pongaddr<%v> msg<%v>\n", fromaddr, PongAddr, core.Message{Type: core.Message_SIMPLE, Description: "ping"})
		pong := Call(PongAddr, FromAddr(pid.GetAddr()), core.Message{
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

func pongCall(_ core.Pid, msg core.Message, fromaddr FromAddr, s State) (core.Message, State) {
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

func TestGenServers(t *testing.T) {
	gs1 := NewGenServer("genserver 1", core.LocalScope, pingCall, defaultCast)
	gs2 := NewGenServer("genserver 2", core.GlobalScope, pongCall, defaultCast)

	go func() {
		t.Log(gs1.Start())
	}()

	go func() {
		t.Log(gs2.Start())
	}()

	for !core.PidHealthCheck(gs1.Pid.GetAddr()) || !core.PidHealthCheck(gs2.Pid.GetAddr()) {
		time.Sleep(25 * time.Microsecond)
		t.Log("waiting for genserver to start")
	}

	PongAddr = PidAddr(gs2.Pid.GetAddr())

	log.Println("ping server addr: ", gs1.Pid.GetAddr())
	log.Println("pong server addr: ", gs2.Pid.GetAddr())
	log.Println("var pong server addr: ", PongAddr)

	rmsg1 := Call(PidAddr(gs2.Pid.GetAddr()), FromAddr("localhost"), core.Message{
		Type:        core.Message_SIMPLE,
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
	}

	gs1.Terminate()
	gs2.Terminate()
}
*/