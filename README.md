GERL 
========

Pronounced like gurl ... as in "hey gurl, hey"...


Gerl is an attempt to build out the remarkable parts of the Erlang/OTP for Go
while keeping in spirit with the language.

The vision is to provide a way to build, schedule, and manage both locally and globally
avaialble processes and their ability to communicate.

Gerl provides functionality for:

  - process id\'s

  - generic server (gen_servers)

  - message passing between pids

This is mainly done by using:

  - channels

  - go routines

  - grpc

## Basic Concepts

### Process ID (Pid)

Processes IDs (pid) is the main abstraction for communicating
with a running process. The pid contains channels for bidirectional communication 
from the running process and, under the hood, handles the GRPC implemtation. 

The pid has both an *inbox* and *outbox* which are channels used to pass messages
into a processes handler/main loop. This allows for the go-routine running the process
to get to the message once it is done handling other messages.

All messages sent to a pid are blocking with repsect to the GRPC server implementaiton.
 For *casts*, which to a process appear to be non-blocking, have an empty message returned
at the GRPC layer which forces both pid to be able to confirm a message was passed.
*Calls* return a new message at the GRPC layer and will wait until a process gets to and 
processes the message.

### Generic Server (genserver)

Generic servers, genservers, are a concept directly pulled from Erlang/OTP. A genserver
is a process that has a pre-specified set of functionality; mainly a *call* and *cast*.


*call* is a bidirectional action in which a client sends a message to and expects
a result back from a genserver. The genserver has a specific function dedicated
to handling *call* actions.

*cast* is a unidirectional action in which the client sends a message to a genserver
and moves on.

The genserver client builds the GRCP client necessary to make the calls and needs the 
address of the pid of the genserver to send messages back and forth.

### Process (proc)

Processes are another concept borrowed from Erlang/OTP. In this case, process has a 
handler which is started as a go-routine and the pid, as with a genserver, is the
main way to communicate with the running process.

All messages to processes are intentionally unidirection and processes must be designed
to allow for bi-directional communicate. Although less featureful, the genserver is the
child of the process in which the genserver implements opinionated and strict constraits on
the process idea.


## Getting started

### GenServer Ping Pong

```go
import (
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

``` 