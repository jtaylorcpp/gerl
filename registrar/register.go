package registrar

import (
	"log"
	"time"

	"github.com/jtaylorcpp/gerl/core"
)

type register struct {
	recordmap    map[string]map[string]Record
	registrarmap map[string]bool
	refresh      time.Timer
}

func newRegister() register {
	reg := register{
		recordmap:    make(map[string]map[string]Record),
		registrarmap: make(map[string]bool),
	}

	return reg
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
