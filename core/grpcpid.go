package grpc

import (
	"net"

	"github.com/jtaylorcpp/gerl/core"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var GRPCAddr string

func init() {
	GRPCAddr = ":8081"
	core.Pid = Pid{}
}

type Pid struct {
	Addr   string
	Inbox  chan core.GerlPassableMessage
	Outbox chan core.GerlPassableMessage
	Server grpc.Server
}

/*
Call(context.Context, *GerlMsg) (*GerlMsg, error)
Cast(context.Context, *GerlMsg) (*Empty, error)
*/
func (p *Pid) Call(ctx *context.Context, in *GerlMsg) (*GerlMsg, error) {
	p.Inbox <- *in
	returnMsg := <-p.Outbox
	return returnMsg
}

func (p *Pid) Cast(ctx *context.Context, in *GerlMsg) (*Empty, error) {
	p.Inbox <- *in
	return &Empty{}, nil
}

func (Pid) NewPid(pbs core.ProcessBufferSize) Pid {
	grpcServer := grpc.NewServer()

	npid := Pid{
		Addr:   GRPCAddr,
		Inbox:  make(chan core.GerlPassableMessage, pbs),
		Outbox: make(chan core.GerlPassableMessage, pbs),
		Server: *grpcServer,
	}
	lis, err := net.Listen("tcp", npid.Addr)
	if err != nil {
		panic(err)
	}

	RegisterGerlMessagerServer(grpcServer, &npid)
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	return npid
}

func (p Pid) Read() (core.GerlPassableMessage, bool) {
	msg, open := <-p.Inbox
	return msg, open
}

func (p Pid) Write(msg core.GerlPassableMessage) {
	p.Outbox <- msg
}

func (p Pid) GetAddr() core.ProcessAddr {
	return core.ProcessAddr(p.Addr)
}

func (p Pid) Terminate() {
	close(p.Inbox)
}
