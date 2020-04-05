package supervisor

import (
	"gerl/core"
	"gerl/genserver"

	log "github.com/sirupsen/logrus"
)

type TestMessage struct {
	Body string
}

type TestState struct {
	Some string
}

func CallTest(_ core.Pid, msg TestMessage, _ genserver.FromAddr, s TestState) (TestMessage, TestState) {
	log.Println("call test func called")
	return msg, s
}

func CastTest(_ core.Pid, msg TestMessage, _ genserver.FromAddr, s TestState) TestState {
	log.Println("cast test func called")
	return s
}
