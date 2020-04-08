package registrar

import (
	"errors"
	"gerl/core"

	log "github.com/sirupsen/logrus"
)

type ProcessName = string
type ProcessAddress = string

type Registrar struct {
	registry map[ProcessName]map[core.Scope][]ProcessAddress
}

func NewRegistrar() *Registrar {
	return &Registrar{
		registry: map[ProcessName]map[core.Scope][]ProcessAddress{},
	}
}

type RegistrarRecord struct {
	Name    string
	Scope   core.Scope
	Address ProcessAddress
}

func (r *Registrar) AddServiceRecord(addRecord RegistrarRecord) error {
	name, scope, addr := addRecord.Name, addRecord.Scope, addRecord.Address
	if _, nameok := r.registry[name]; nameok {
		if knownAddrs, scopeok := r.registry[name][scope]; scopeok {
			addrIsKnown := false
			for _, knownAddr := range knownAddrs {
				if knownAddr == addr {
					addrIsKnown = true
					break
				}
			}
			if !addrIsKnown {
				r.registry[name][scope] = append(r.registry[name][scope], addr)
			}
		} else {
			if scope == core.LocalScope || scope == core.GlobalScope {
				r.registry[name][scope] = []ProcessAddress{addr}
			} else {
				return errors.New("scope provided for registrar is not core.LocalScope or core.GlobalScope")
			}
		}
	} else {
		if scope == core.LocalScope || scope == core.GlobalScope {
			if addr != "" {
				r.registry[name] = map[core.Scope][]ProcessAddress{
					scope: []ProcessAddress{addr},
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
				newAddrs := []ProcessAddress{}
				for _, currentAddr := range r.registry[name][scope] {
					if currentAddr != addr {
						newAddrs = append(newAddrs, currentAddr)
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
				return r.registry[name][scope], nil
			default:
				returnAddresses := []ProcessAddress{}
				for _, knownAddr := range r.registry[name][scope] {
					if knownAddr == addr {
						returnAddresses = append(returnAddresses, addr)
					}
				}

				return returnAddresses, nil
			}
		default:
			returnAddresses := []ProcessAddress{}
			for _, addrs := range r.registry[name] {
				returnAddresses = append(returnAddresses, addrs...)
			}

			return returnAddresses, nil
		}
	}
}
