package registrar

import (
	"errors"
	"gerl/core"
	"time"

	log "github.com/sirupsen/logrus"
)

type ProcessName = string
type ProcessAddress = string
type ProcessRecord struct {
	address     ProcessAddress
	lastUpdated time.Time
}

type Registrar struct {
	registry map[ProcessName]map[core.Scope][]ProcessRecord
}

func NewRegistrar() *Registrar {
	return &Registrar{
		registry: map[ProcessName]map[core.Scope][]ProcessRecord{},
	}
}

type RegistrarRecord struct {
	Name    ProcessName
	Scope   core.Scope
	Address ProcessAddress
}

func (r *Registrar) AddServiceRecord(addRecord RegistrarRecord) error {
	name, scope, addr := addRecord.Name, addRecord.Scope, addRecord.Address
	timeAdded := time.Now()
	if _, nameok := r.registry[name]; nameok {
		if knownAddrs, scopeok := r.registry[name][scope]; scopeok {
			addrIsKnown := false
			for _, knownProcess := range knownAddrs {
				if knownProcess.address == addr {
					addrIsKnown = true
					break
				}
			}
			if !addrIsKnown {
				r.registry[name][scope] = append(r.registry[name][scope], ProcessRecord{addr, timeAdded})
			}
		} else {
			if scope == core.LocalScope || scope == core.GlobalScope {
				r.registry[name][scope] = []ProcessRecord{{addr, timeAdded}}
			} else {
				return errors.New("scope provided for registrar is not core.LocalScope or core.GlobalScope")
			}
		}
	} else {
		if scope == core.LocalScope || scope == core.GlobalScope {
			if addr != "" {
				r.registry[name] = map[core.Scope][]ProcessRecord{
					scope: []ProcessRecord{{addr, timeAdded}},
				}
			} else {
				return errors.New("no address prvided to registrar")
			}
		} else {
			return errors.New("scope provided for registrar is not core.LocalScope or core.GlobalScope")
		}
	}

	return nil
}

// RemoveServiceRecord will remove records from the registrar using the following logic:
// If no name is included the method will return an error
// if a name but no scope is included all entries will be removed for that name
// if a name and scope are included all entries for that scope will be removed
// if a name, scope, and address are included only that address will be removed
func (r *Registrar) RemoveServiceRecord(removeRecord RegistrarRecord) error {
	log.Printf("removing records from registrar: %#v\n", removeRecord)
	name, scope, addr := removeRecord.Name, removeRecord.Scope, removeRecord.Address

	switch name {
	case "":
		return errors.New("name needed to remove service records")
	default:
		switch scope {
		case core.LocalScope, core.GlobalScope:
			if addr != "" {
				newAddrs := []ProcessRecord{}
				for _, currentProcRecord := range r.registry[name][scope] {
					if currentProcRecord.address != addr {
						newAddrs = append(newAddrs, currentProcRecord)
					}
				}
				r.registry[name][scope] = newAddrs
			} else {
				delete(r.registry[name], scope)
			}
		default:
			// delete all records for service
			delete(r.registry, name)
		}
	}

	return nil
}

func (r *Registrar) GetServiceAddresses(record RegistrarRecord) ([]ProcessAddress, error) {
	name, scope, addr := record.Name, record.Scope, record.Address
	switch name {
	case "":
		return []ProcessAddress{}, errors.New("a service name must be included when getting process addresses")
	default:
		switch scope {
		case core.LocalScope, core.GlobalScope:
			switch addr {
			case "":
				returnAddrs := []ProcessAddress{}
				for _, procRecord := range r.registry[name][scope] {
					returnAddrs = append(returnAddrs, procRecord.address)
				}
				return returnAddrs, nil
			default:
				returnAddresses := []ProcessAddress{}
				for _, knownProc := range r.registry[name][scope] {
					if knownProc.address == addr {
						returnAddresses = append(returnAddresses, addr)
					}
				}

				return returnAddresses, nil
			}
		default:
			returnAddresses := []ProcessAddress{}
			for _, procRecords := range r.registry[name] {
				for _, procRecord := range procRecords {
					returnAddresses = append(returnAddresses, procRecord.address)
				}
			}

			return returnAddresses, nil
		}
	}
}
