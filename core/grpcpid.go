package core

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// global var which allows pids to know what ip address to bind to
var IPAddress string

// initializes the pid environment
//  - IPAddress to bind all pids to
func init() {
	if IPAddress == "" {
		IPAddress = "127.0.0.1"
	}
}

// ProcessID (Pid) is the struct used to keep track of the main
//  communication method to a running process
type Pid struct {
	// GRPC server
	Server   *grpc.Server
	Listener net.Listener
	// Address of the currently running Pid
	Addr string
	// Inbox for messages passed to a process
	Inbox chan GerlMsg
	// Outbox for messages passed from a process
	Outbox chan GerlMsg
	// Error chan to be monitored by the process using the Pid
	Errors chan error
	// Listener termination
	LisTerm chan bool
	// Running check
	Running bool
}

// GRPC function for interface GerlMessager
func (p *Pid) Call(ctx context.Context, in *GerlMsg) (*GerlMsg, error) {
	p.Inbox <- *in
	outMsg := <-p.Outbox
	returnMsg := &outMsg
	returnMsg.Fromaddr = p.GetAddr()
	return returnMsg, nil
}

// GRPC function for interface GerlMessager
func (p *Pid) Cast(ctx context.Context, in *GerlMsg) (*Empty, error) {
	p.Inbox <- *in
	return &Empty{}, nil
}

//GRPC function for interface GerlMessager
func (p *Pid) RUOK(ctx context.Context, _ *Empty) (*Health, error) {
	return &Health{Status: Health_ALIVE}, nil
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
	// if port is unset use 0
	ipport := "0"
	if port != "" {
		ipport = port
	}

	// generate tcp listener
	// no ip address will use 0.0.0.0
	// no port number(string) will result in one being assigned
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", ipaddress, ipport))
	if err != nil {
		Errors <- err
	}

	// new grpc server constructor
	grpcServer := grpc.NewServer()

	// create pid to return
	npid := &Pid{
		Listener: lis,
		Addr:     lis.Addr().(*net.TCPAddr).String(),
		Inbox:    make(chan GerlMsg, 8),
		Outbox:   make(chan GerlMsg, 8),
		Errors:   Errors,
		Server:   grpcServer,
		Running:  false,
		LisTerm:  make(chan bool, 1),
	}

	// register pid and grpc server
	RegisterGerlMessagerServer(grpcServer, npid)
	reflection.Register(grpcServer)

	// go routine to run grpc server in the background
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			npid.Errors <- err
			npid.LisTerm <- true
		} else {
			npid.LisTerm <- true
		}
	}()

	npid.Running = true
	return npid
}

// Getter for address in format ip:port
func (p Pid) GetAddr() string {
	return p.Addr
}

// Terminates the Pid and closes all of the Pid side components
func (p *Pid) Terminate() {
	log.Printf("Pid <%v> terminating\n", p)
	log.Println("closing listener")
	p.Listener.Close()
	log.Println("closing grpc server")
	p.Server.Stop()
	log.Println("closing channels")
	close(p.Inbox)
	p.Errors <- errors.New("pid terminated")
	// blocking since the listener close out may generate an error
	<-p.LisTerm
	close(p.Errors)
	p.Running = false
	log.Printf("pid<%v> terminated\n", p)
}

// Creates GRPC client with only an address string
func newClient(pidAddress string) (*grpc.ClientConn, GerlMessagerClient) {
	var conn *grpc.ClientConn

	// gets connection to remote GRPC server
	conn, err := grpc.Dial(pidAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("could not connect to server: ", err)
	}

	client := NewGerlMessagerClient(conn)

	// return connection and messager client
	return conn, client

}

// Sends a Call message to the Pid from a Message struct
// Constructs both the client and GerlMessagee needed
func PidCall(toaddr string, fromaddr string, msg Message) Message {
	conn, client := newClient(toaddr)
	defer conn.Close()
	gerlMsg := &GerlMsg{
		Type:     GerlMsg_CALL,
		Fromaddr: fromaddr,
		Msg:      &msg,
	}
	returnGerlMsg, err := client.Call(context.Background(), gerlMsg)
	if err != nil {
		log.Printf("error<%v> calling pid<%v> with msg<%v>\n", err, toaddr, msg)
	}
	return *returnGerlMsg.GetMsg()
}

// Sends a Cast message to the Pid from a Message struct
// Constructs both the client and GerlMessagee needed
func PidCast(toaddr string, fromaddr string, msg Message) {
	conn, client := newClient(toaddr)
	defer conn.Close()
	gerlMsg := &GerlMsg{
		Type:     GerlMsg_CAST,
		Fromaddr: fromaddr,
		Msg:      &msg,
	}
	_, err := client.Cast(context.Background(), gerlMsg)
	if err != nil {
		log.Printf("error<%v> cast pid<%v> with msg<%v>\n", err, toaddr, msg)
	}
}

// Sends a Cast message with type PROC to the Pid from a Message struct
// Constructs both the client and GerlMessagee needed
func PidSendProc(toaddr string, fromaddr string, msg Message) {
	conn, client := newClient(toaddr)
	defer conn.Close()
	gerlMsg := &GerlMsg{
		Type:     GerlMsg_PROC,
		Fromaddr: fromaddr,
		Msg:      &msg,
	}
	_, err := client.Cast(context.Background(), gerlMsg)
	if err != nil {
		log.Printf("error<%v> cast pid<%v> with msg<%v>\n", err, toaddr, msg)
	}
}

func PidHealthCheck(toaddr string) bool {
	conn, client := newClient(toaddr)
	defer conn.Close()
	deadline := time.Now().Add(10 * time.Millisecond)
	ctx, _ := context.WithDeadline(context.Background(), deadline)
	health, err := client.RUOK(ctx, &Empty{})
	if err != nil {
		log.Printf("error<%v> getting pid<%v> health\n", err, toaddr)
		return false
	}

	if health.GetStatus() == Health_ALIVE {
		return true
	}

	return false
}
