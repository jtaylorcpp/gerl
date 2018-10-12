package core

import (
	"errors"
	"fmt"
	"log"
	"net"

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
	Addr    string
	Inbox   chan GerlMsg
	Outbox  chan GerlMsg
	Errors  chan error
	Server  *grpc.Server
	Running bool
}

// GRPC function
func (p *Pid) Call(ctx context.Context, in *GerlMsg) (*GerlMsg, error) {
	p.Inbox <- *in
	returnMsg := <-p.Outbox
	return &returnMsg, nil
}

// GRPC function
func (p *Pid) Cast(ctx context.Context, in *GerlMsg) (*Empty, error) {
	p.Inbox <- *in
	return &Empty{}, nil
}

// Generates new Pid to use by process in Gerl
func NewPid(address, port string) *Pid {
	// error chan to elevate to process using pid
	Errors := make(chan error, 10)

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
	npid := &Pid{
		Addr:    lis.Addr().(*net.TCPAddr).String(),
		Inbox:   make(chan GerlMsg, 8),
		Outbox:  make(chan GerlMsg, 8),
		Errors:  Errors,
		Server:  grpcServer,
		Running: false,
	}

	// register pid and grpc server
	RegisterGerlMessagerServer(grpcServer, npid)
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			npid.Errors <- err
		}
	}()

	npid.Running = true
	return npid
}

func (p Pid) GetAddr() string {
	return p.Addr
}

func (p *Pid) Terminate() {
	log.Printf("Pid <%v> terminating\n", p)
	log.Println("closing grpc server")
	p.Server.Stop()
	log.Println("closing channels")
	close(p.Inbox)
	p.Errors <- errors.New("pid terminated")
	close(p.Errors)
	p.Running = false
	log.Printf("pid<%v> terminated\n", p)
}

func newClient(pidAddress string) (*grpc.ClientConn, GerlMessagerClient) {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(pidAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("could not connect to server: ", err)
	}

	client := NewGerlMessagerClient(conn)

	return conn, client

}

func PidCall(toaddr string, fromaddr string, msg Message) Message {
	conn, client := newClient(toaddr)
	defer conn.Close()
	gerlMsg := &GerlMsg{
		Type:        GerlMsg_CALL,
		Processaddr: fromaddr,
		Msg:         &msg,
	}
	returnGerlMsg, err := client.Call(context.Background(), gerlMsg)
	if err != nil {
		log.Printf("error<%v> calling pid<%v> with msg<%v>\n", err, toaddr, msg)
	}
	return *returnGerlMsg.GetMsg()
}

func PidCast(toaddr string, fromaddr string, msg Message) {
	conn, client := newClient(toaddr)
	defer conn.Close()
	gerlMsg := &GerlMsg{
		Type:        GerlMsg_CAST,
		Processaddr: fromaddr,
		Msg:         &msg,
	}
	_, err := client.Cast(context.Background(), gerlMsg)
	if err != nil {
		log.Printf("error<%v> cast pid<%v> with msg<%v>\n", err, toaddr, msg)
	}
}
