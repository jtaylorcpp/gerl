package registrar

import (
	"fmt"
	"log"
	"time"

	"github.com/jtaylorcpp/gerl/core"
	"github.com/jtaylorcpp/gerl/genserver"
)

var REFRESH_TIMER time.Duration

type State genserver.State
type CallHander genserver.GenServerCallHandler
type CastHandler genserver.GenServerCastHandler

func init() {
	REFRESH_TIMER = 10 * time.Second
}

func registrarCallHander(pid core.Pid, in core.Message, from genserver.FromAddr, state genserver.State) (core.Message, genserver.State) {
	log.Printf("Registrar call handler msg: <%v>\n", in)
	switch in.GetType() {
	case core.Message_REGISTER:
		switch in.GetSubtype() {
		case core.Message_SET, core.Message_PUT:
			log.Println("Registar handling REGISTER_[SET|PUT]")
			log.Println("Registrar handling values: ", in.GetValues())
			if len(in.GetValues())%3 != 0 {
				log.Println("register message has values that are not same length as records: ", in.Values)
			} else if len(in.GetValues()) == 0 {
				log.Println("no records to register in len 0 values")
			} else {
				log.Println("registering records: ", in.GetValues())
				reg := state.(register)
				for idx := 0; idx < len(in.GetValues()); idx += 3 {
					reg.addRecords(NewRecord(in.GetValues()[idx], in.GetValues()[idx+1], core.Scope(in.GetValues()[idx+2])))
				}
				log.Println("register: ", reg)
				return core.Message{
					Type:        core.Message_SIMPLE,
					Description: "register",
					Values:      []string{fmt.Sprintf("%d", len(in.GetValues())/3)},
				}, reg
			}

		case core.Message_GET:
			// handle register of new genserver
			log.Println("Registrar handling SIMPLE message")
			log.Println("Registrar handling SIMPLE get for name: ", in.GetValues())
			reg := state.(register)
			records := reg.getRecords(in.GetValues()[0])
			values := make([]string, 0)
			for _, record := range records {
				values = append(values, record.Name, record.Address, string(record.Scope))
			}
			return core.Message{
				Type:        core.Message_SIMPLE,
				Description: "get",
				Values:      values,
			}, reg

		default:
			log.Println("unknonw register message: ", in)

		}

	case core.Message_SYNC:
		log.Println("Registrar SYNC recieved: ", in)
		switch in.GetSubtype() {
		case core.Message_JOIN:
			log.Println("Registrar joining: ", in)
			reg := state.(register)
			if _, ok := reg.registrarmap[in.GetDescription()]; ok {
				log.Println("Registrar already registered node: ", in.GetDescription())
			} else {
				reg.registrarmap[in.GetDescription()] = true
			}
			/*
				TODO:
				Add other synched nodes to message
			*/
			return core.Message{
				Type:        core.Message_SYNC,
				Subtype:     core.Message_JOIN,
				Description: "registered",
			}, reg
		default:
			log.Println("Registrar unhandled SYNC: ", in)
		}

	case core.Message_SIMPLE:
		// handle simple messages
		log.Println("Registrar SIMPLE recieved: ", in)
	default:
		log.Println("unknonw message handled: ", in)

	}

	return core.Message{}, state
}

func registrarCastHander(pid core.Pid, in core.Message, from genserver.FromAddr, state genserver.State) genserver.State {
	log.Println("registrar handling cast message: ", in)
	reg := state.(register)
	switch in.GetType() {
	case core.Message_SYNC:
		log.Println("recieved sync message: ", in)
		switch in.GetSubtype() {
		case core.Message_REFRESH:

			log.Println("handling refresh with state: ", reg)
			for node, _ := range reg.registrarmap {
				if core.PidHealthCheck(node) {
					log.Println("node is still available: ", node)
				} else {
					log.Println("node is no longer available: ", node)
					delete(reg.registrarmap, node)
				}
			}

			for group, nodes := range reg.recordmap {
				log.Println("checking availability for pid group: ", group)
				for addr, _ := range nodes {
					log.Printf("checking availability of pid %v in group %v\n", addr, group)
					if core.PidHealthCheck(addr) {
						log.Printf("pid %v in group %v still available\n", addr, group)
					} else {
						log.Printf("pid %v in group %v no longer available\n", addr, group)
						delete(reg.recordmap[group], addr)
					}
				}
			}

			timer := time.NewTimer(REFRESH_TIMER)
			go func() {
				<-timer.C
				log.Println("registrar sending keep alives")
				genserver.Cast(genserver.PidAddr(pid.GetAddr()),
					genserver.PidAddr(pid.GetAddr()),
					core.Message{
						Type:    core.Message_SYNC,
						Subtype: core.Message_REFRESH,
					},
				)
			}()
		}

	default:
		log.Println("unknown call msg: ", in)

	}

	return reg
}

func NewRegistrar(scope core.Scope) *genserver.GenServer {
	gensvr := genserver.NewGenServer(newRegister(), scope, registrarCallHander, registrarCastHander)

	genserver.Cast(genserver.PidAddr(gensvr.Pid.GetAddr()),
		genserver.PidAddr(gensvr.Pid.GetAddr()),
		core.Message{
			Type:    core.Message_SYNC,
			Subtype: core.Message_REFRESH,
		},
	)

	return gensvr
}

func JoinRegistrar(from genserver.FromAddr, to genserver.PidAddr) bool {
	msg := genserver.Call(to, from, core.Message{
		Type:        core.Message_SYNC,
		Subtype:     core.Message_JOIN,
		Description: string(from),
	})

	if (msg.GetType() == core.Message_SYNC) &&
		(msg.GetSubtype() == core.Message_JOIN) &&
		(msg.GetDescription() == "registered") {
		log.Printf("registrar <%v> registered with registrar <%v>\n", from, to)

		return true
	}

	return false
}
