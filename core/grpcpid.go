package core

import (
	"net"

	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// ip address var to set for pid to use
var IPAddress string

func init() {
	if IPAddress == "" {
		IPAddress = "127.0.0.1"
	}
}

type Pid struct {
	Addr   string
	Inbox  chan Message
	Outbox chan Message
	Errors chan error
	Server *grpc.Server
}

// GRPC function
func (p *Pid) Call(ctx context.Context, in *GerlMsg) (*GerlMsg, error) {
	p.Inbox <- *in.GetMsg()
	returnMsg := <-p.Outbox
	returnGerl := &GerlMsg{
		Type:        GerlMsg_CALL,
		Processaddr: returnMsg.GetFromaddr(),
		Msg:         &returnMsg,
	}
	return returnGerl, nil
}

// GRPC function
func (p *Pid) Cast(ctx context.Context, in *GerlMsg) (*Empty, error) {
	p.Inbox <- *in.GetMsg()
	return &Empty{}, nil
}

// Generates new Pid to use by process in Gerl
func NewPid(address, port string) Pid {
	// error chan to elevate to process using pid
	Errors := make(chan error)

	// get default addresses to use
	ipaddress := IPAddress
	if address != "" {
		ipaddress = address
	}
	ipport := "0"
	if port != "" {
		ipport = port
	}

	// generate tcp listener
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", ipaddress, ipport))
	if err != nil {
		Errors <- err
	}

	// new grpc server constructor
	grpcServer := grpc.NewServer()

	// create pid to return
	npid := Pid{
		Addr:   lis.Addr().(*net.TCPAddr).String(),
		Inbox:  make(chan Message, 8),
		Outbox: make(chan Message, 8),
		Errors: Errors,
		Server: grpcServer,
	}

	// register pid and grpc server
	RegisterGerlMessagerServer(grpcServer, &npid)
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			npid.Errors <- err
		}
	}()

	return npid
}

func (p Pid) Read() (Message, bool) {
	msg, open := <-p.Inbox
	return msg, open
}

func (p Pid) Write(msg Message) {
	p.Outbox <- msg
}

func (p Pid) GetAddr() string {
	return p.Addr
}

func (p Pid) Terminate() {
	p.Server.Stop()
	close(p.Inbox)
}
