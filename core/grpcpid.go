package core

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// global var which allows pids to know what ip address to bind to
const (
	localIP string = "127.0.0.1"
)

var GlobalIP string

var MessageTimeout time.Duration
var HealthTimeout time.Duration

var DefaultLogger *logrus.Entry = nil

type Scope = byte

const (
	LocalScope  Scope = 0x01
	GlobalScope Scope = 0x02
)

// initializes the pid environment
//  - IPAddress to bind all pids to
func init() {
	MessageTimeout = 500 * time.Millisecond
	HealthTimeout = 50 * time.Millisecond

	if DefaultLogger == nil {
		logrus.SetReportCaller(true)
		DefaultLogger = logrus.WithField("component", "core/pid")
	} else {
		DefaultLogger = DefaultLogger.WithField("component", "core/pid")
	}
}

func getPublicIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, iface := range ifaces {
		DefaultLogger.Println("checking iface: ", iface)

		switch iface.Name {
		case "docker0", "lo":
			continue
		default:
			addrs, err := iface.Addrs()
			if err != nil {
				panic(err)
			}

			for _, addr := range addrs {
				DefaultLogger.Println("checking addr: ", addr)
				switch t := addr.(type) {
				case *net.IPNet:
					DefaultLogger.Println("ipnet: ", t)
					if t.IP.To4() == nil {
						continue
					}
					return t.IP.To4().String()
				case *net.IPAddr:
					DefaultLogger.Println("ipaddr: ", t)
					if t.IP.To4() == nil {
						continue
					}
					return t.IP.To4().String()
				}
			}
		}
	}

	return ""
}

type PidError struct {
	Reasone string
	Error   error
}

type PidTerm struct {
	Reason string
	Error  error
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
	Outbox          chan GerlMsg
	serverErrorChan chan error
	Scope           Scope
	Logger          *logrus.Entry
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
// If address is empty an address based on scope is assigned
// If port is left empty a random one will be used
func NewPid(address, port string, scope Scope) (*Pid, error) {
	// get default addresses to use
	var ipaddress string
	if address != "" {
		ipaddress = address
	} else {
		switch scope {
		case LocalScope:
			ipaddress = localIP
		case GlobalScope:
			if GlobalIP == "" {
				ipaddress = getPublicIP()
			} else {
				ipaddress = GlobalIP
			}
		}
	}

	// generate tcp listener
	// no ip address will use 0.0.0.0
	// no port number(string) will result in one being assigned
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", ipaddress, port))
	if err != nil {
		DefaultLogger.Println("error genreating listener: ", err.Error())
		return nil, err
	}

	if lis == nil {
		DefaultLogger.Println("generated listener is nil")
		return nil, errors.New("pid listener is nil")
	}

	if lis.Addr() == nil {
		DefaultLogger.Println("generated listener addr is nil")
		return nil, errors.New("pid listener addr is nil")
	}

	// new grpc server constructor
	grpcServer := grpc.NewServer()
	// create pid to return
	pidLogger := DefaultLogger.WithField("listenerAddr", lis.Addr().String)
	npid := &Pid{
		Listener:        lis,
		Addr:            lis.Addr().String(),
		Inbox:           make(chan GerlMsg, 1),
		Outbox:          make(chan GerlMsg, 1),
		Server:          grpcServer,
		serverErrorChan: make(chan error, 1),
		Scope:           scope,
		Logger:          pidLogger,
	}

	// register pid and grpc server
	RegisterGerlMessagerServer(grpcServer, npid)
	reflection.Register(grpcServer)

	go func(npid *Pid) {
		npid.serverErrorChan <- npid.Server.Serve(lis)
		close(npid.serverErrorChan)
	}(npid)

	return npid, nil
}

// Getter for address in format ip:port
func (p Pid) GetAddr() string {
	return p.Addr
}

// Terminate terminates the Pid and closes all of the Pid side components
func (p *Pid) Terminate() error {
	if p == nil {
		p.Logger.Println("pid has already been terminated")
		return nil
	}
	p.Logger.Printf("Pid <%v> terminating\n", p)
	p.Logger.Println("closing inbox")
	close(p.Inbox)
	close(p.Outbox)
	p.Logger.Println("closing grpc listner")
	p.Server.GracefulStop()
	p.Logger.Println("closing listener")
	p.Listener.Close()
	p.Logger.Println("closing grpc server")
	p.Server.Stop()
	p.Logger.Println("closing channels")
	// blocking since the listener close out may generate an error
	serverError := <-p.serverErrorChan
	if serverError != nil {
		p.Logger.Errorln("error recieved when closing out grpc server: ", serverError.Error())
	}
	p.Logger.Printf("pid<%v> terminated\n", p)
	return nil
}

// Creates GRPC client with only an address string
func newClient(pidAddress string) (*grpc.ClientConn, GerlMessagerClient) {
	var conn *grpc.ClientConn

	// gets connection to remote GRPC server
	conn, err := grpc.Dial(pidAddress, grpc.WithInsecure())
	if err != nil {
		DefaultLogger.Fatalln("could not connect to server: ", err)
	}

	client := NewGerlMessagerClient(conn)

	// return connection and messager client
	return conn, client

}

// Sends a Call message to the Pid from a Message struct
// Constructs both the client and GerlMessage needed
func PidCall(toaddr string, fromaddr string, msg Message) Message {
	conn, client := newClient(toaddr)
	defer conn.Close()
	gerlMsg := &GerlMsg{
		Type:     GerlMsg_CALL,
		Fromaddr: fromaddr,
		Msg:      &msg,
	}
	deadline := time.Now().Add(MessageTimeout)
	ctx, _ := context.WithDeadline(context.Background(), deadline)
	returnGerlMsg, err := client.Call(ctx, gerlMsg)
	if err != nil {
		DefaultLogger.Printf("error<%v> calling pid<%v> with msg<%v>\n", err, toaddr, msg)
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
	deadline := time.Now().Add(MessageTimeout)
	ctx, _ := context.WithDeadline(context.Background(), deadline)
	_, err := client.Cast(ctx, gerlMsg)
	if err != nil {
		DefaultLogger.Printf("error<%v> cast pid<%v> with msg<%v>\n", err, toaddr, msg)
	}
}

func PidTerminate(toaddr string, fromaddr string, msg Message) error {
	conn, client := newClient(toaddr)
	defer conn.Close()
	gerlMsg := &GerlMsg{
		Type:     GerlMsg_TERM,
		Fromaddr: fromaddr,
		Msg:      &msg,
	}
	deadline := time.Now().Add(MessageTimeout)
	ctx, _ := context.WithDeadline(context.Background(), deadline)
	_, err := client.Cast(ctx, gerlMsg)
	if err != nil {
		DefaultLogger.Printf("error<%v> cast pid<%v> with msg<%v>\n", err, toaddr, msg)
		return err
	}

	return nil
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
	deadline := time.Now().Add(MessageTimeout)
	ctx, _ := context.WithDeadline(context.Background(), deadline)
	_, err := client.Cast(ctx, gerlMsg)
	if err != nil {
		DefaultLogger.Printf("error<%v> cast pid<%v> with msg<%v>\n", err, toaddr, msg)
	}
}

func PidHealthCheck(toaddr string) bool {
	conn, client := newClient(toaddr)
	defer conn.Close()
	deadline := time.Now().Add(HealthTimeout)
	ctx, _ := context.WithDeadline(context.Background(), deadline)
	health, err := client.RUOK(ctx, &Empty{})
	if err != nil {
		DefaultLogger.Printf("error<%v> getting pid<%v> health\n", err, toaddr)
		return false
	}

	if health.GetStatus() == Health_ALIVE {
		return true
	}

	return false
}
