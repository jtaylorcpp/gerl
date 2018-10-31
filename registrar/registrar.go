package registrar

import (
	"fmt"
	"log"
	"time"

	"github.com/jtaylorcpp/gerl/core"
	"github.com/jtaylorcpp/gerl/genserver"
)

type State genserver.State
type CallHander genserver.GenServerCallHandler
type CastHandler genserver.GenServerCastHandler

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
		// handle syncing of other registrars
		log.Println("Registrar SYNC recieved: ", in)
	case core.Message_SIMPLE:
		// handle simple messages
		log.Println("Registrar SIMPLE recieved: ", in)
	default:
		log.Println("unknonw message handled: ", in)

	}

	return core.Message{}, state
}

func registrarCastHander(pid core.Pid, in core.Message, from genserver.FromAddr, state genserver.State) genserver.State {

	return state
}

func NewRegistrar(scope core.Scope) *genserver.GenServer {
	gensvr := genserver.NewGenServer(newRegister(), scope, registrarCallHander, registrarCastHander)

	return gensvr
}

type register struct {
	recordmap map[string]map[string]Record
	ticker    *time.Ticker
}

func newRegister() register {
	return register{recordmap: make(map[string]map[string]Record)}
}

type Record struct {
	Name    string
	Address string
	Scope   core.Scope
}

func NewRecord(name, address string, scope core.Scope) Record {
	return Record{
		Name:    name,
		Address: address,
		Scope:   scope,
	}
}

func (r register) addRecords(records ...Record) register {
	log.Println("adding records to register: ", records)
	for _, rec := range records {
		log.Println("adding record: ", rec)
		if _, svc := r.recordmap[rec.Name]; !svc {
			log.Println("adding service: ", rec.Name)
			r.recordmap[rec.Name] = make(map[string]Record)
		}
		r.recordmap[rec.Name][rec.Address] = rec
		log.Println("record added: ", r)
	}
	log.Println("new register state: ", r)
	return r
}

func AddRecords(regaddr string, fromaddr string, records ...Record) bool {
	log.Println("Registrar call to add records: ", records)
	msg := core.Message{
		Type:        core.Message_REGISTER,
		Subtype:     core.Message_PUT,
		Description: "register",
		Values:      make([]string, 0),
	}

	for _, rec := range records {
		msg.Values = append(msg.Values, rec.Name, rec.Address, string(rec.Scope))
	}

	log.Println("Adding records: ", msg.Values)

	returnMsg := core.PidCall(regaddr, fromaddr, msg)

	log.Println("recieved register message back: ", returnMsg)

	return true
}

func (r register) getRecords(name string) []Record {
	records := make([]Record, 0)
	if namedRecords, ok := r.recordmap[name]; ok {
		for _, record := range namedRecords {
			records = append(records, record)
		}

		return records
	}

	return records
}

func GetRecords(regaddr string, fromaddr string, name string) []Record {
	log.Println("Registrar call to get records with name: ", name)
	msg := core.Message{
		Type:        core.Message_REGISTER,
		Subtype:     core.Message_GET,
		Description: "get",
		Values:      []string{name},
	}

	returnMsg := core.PidCall(regaddr, fromaddr, msg)

	log.Println("recieved register message back: ", returnMsg)

	records := make([]Record, 0)

	values := returnMsg.GetValues()

	for idx := 0; idx < len(values); idx += 3 {
		records = append(records, NewRecord(values[idx], values[idx+1], core.Scope(values[idx+2])))
	}

	return records
}
