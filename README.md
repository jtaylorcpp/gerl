GERL 
========

Pronounced like gurl ... as in "hey gurl, hey"...

[![CircleCI](https://circleci.com/gh/jtaylorcpp/gerl.svg?style=svg)](https://circleci.com/gh/jtaylorcpp/gerl)

Gerl is an attempt to build out the remarkable parts of the Erlang/OTP for Go
while keeping in spirit with the language.

The vision is to provide a way to build, schedule, and manage both locally and globally
avaialble processes and their ability to communicate.

Gerl provides functionality for:

  - process id\'s

  - generic server (gen_servers)

  - message passing between pids

  - clustered registrar service 

This is mainly done by using:

  - channels

  - go routines

  - grpc

## Why

Golang is a ton of fun. Its reads and writes easily with a good balance of performance and has a robust community around it. Golang is a daily driver.

However, Golang does have a list of issues or features missing that often makes people look at and learn new languages. 

On such journey took me down the path of learning Erlang. Erlangs legedary status of reliable deployments and everest like learning curve made a stark contrast to Golang. Through the hours and days of unraveling the Erlang mystery; there were a few features that, once you got the hang of, made you wonder why they didnt exist elsewhere.

One such set of features is the event driven and message passing process. Message passing is baked into the essence of Erlang as is the spirit of functional programming. Yet the combination of the two into the Erlang gen_server is a piece of magic. Quickly define the handlers for certain types of messages and spin up a hyper-lightweight process. Networking...provided. Event driven messaging...use the handlers. State management...included in the server init. 

MAGIC.

## All the Things

### Process ID (Pid)

Processes IDs (pid) is the main abstraction for communicating
with a running process. The pid contains channels for bidirectional communication 
from the running process and, under the hood, handles the GRPC implemtation. 

All messages sent to a pid are blocking with repsect to the GRPC server implementaiton.
  *Casts* result in an immediately returned empty message and  *Calls* result in the
Pid waiting for the process to return a full message.

### Generic Server (genserver)

Generic servers, genservers, are a concept directly pulled from Erlang/OTP. A genserver
is a process that has a pre-specified set of functionality; mainly a *call* and *cast*.

*call* is a bidirectional action in which a client sends a message to and expects
a result back from a genserver. 

*cast* is a unidirectional action in which the client sends a message to a genserver
and does not expect a message to be returned.

The GenServer constructor accepts handlers for both *Call* and *Cast*.

The client func's for *call* and *cast* only need the Pid address and the message to be 
 passed; no messy business building a GRPC client.

### Process (proc)

Processes are another concept borrowed from Erlang/OTP. In this case, process has a 
handler which is started as a go-routine and the pid, as with a genserver, is the
main way to communicate with the running process.

All messages to processes are implemented like a *cast* in that they nover expect a message directly back.


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
