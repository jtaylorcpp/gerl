package registrar

import (
	"gerl/core"
	"testing"
)

func TestRegistrarAddRecords(t *testing.T) {
	r := NewRegistrar()

	err := r.AddServiceRecord(RegistrarRecord{
		Name: "test",
	})

	if err == nil {
		t.Fatal("should have errored since there is no scope")
	}

	err = r.AddServiceRecord(RegistrarRecord{
		Name:  "test",
		Scope: 0x00,
	})

	if err == nil {
		t.Fatal("provided scope is not usable")
	}

	err = r.AddServiceRecord(RegistrarRecord{
		Name:  "test",
		Scope: core.LocalScope,
	})

	if err == nil {
		t.Fatal("no address provided")
	}

	err = r.AddServiceRecord(RegistrarRecord{
		Name:    "test",
		Scope:   core.LocalScope,
		Address: "1",
	})

	if err != nil {
		t.Fatal("given record should be able to be stored")
	}

	if r.registry["test"][core.LocalScope][0] != "1" {
		t.Fatal("address not properly stored")
	}

	err = r.AddServiceRecord(RegistrarRecord{
		Name:    "test",
		Scope:   0x00,
		Address: "1",
	})

	if err == nil {
		t.Fatal("complete record but usable scope")
	}

	err = r.AddServiceRecord(RegistrarRecord{
		Name:    "test",
		Scope:   core.GlobalScope,
		Address: "1",
	})

	if err != nil {
		t.Fatal("should be able to store global record")
	}

	if len(r.registry["test"][core.LocalScope]) != len(r.registry["test"][core.GlobalScope]) {
		t.Fatal("should have 1 record in local and global")
	}

	err = r.AddServiceRecord(RegistrarRecord{
		Name:    "test",
		Scope:   core.LocalScope,
		Address: "1",
	})

	if err != nil {
		t.Fatal("given record should be able to be stored")
	}

	if len(r.registry["test"][core.LocalScope]) != 1 {
		t.Fatal("should not have duplicated a known record")
	}
}

func TestRegistrarRemoveRecords(t *testing.T) {
	r := NewRegistrar()

	records := []RegistrarRecord{
		RegistrarRecord{
			Name:    "testA",
			Scope:   core.LocalScope,
			Address: "1",
		},
		RegistrarRecord{
			Name:    "testA",
			Scope:   core.LocalScope,
			Address: "2",
		},
		RegistrarRecord{
			Name:    "testB",
			Scope:   core.LocalScope,
			Address: "3",
		},
		RegistrarRecord{
			Name:    "testB",
			Scope:   core.LocalScope,
			Address: "4",
		},
		RegistrarRecord{
			Name:    "testA",
			Scope:   core.GlobalScope,
			Address: "a",
		},
		RegistrarRecord{
			Name:    "testA",
			Scope:   core.GlobalScope,
			Address: "b",
		},
		RegistrarRecord{
			Name:    "testB",
			Scope:   core.GlobalScope,
			Address: "c",
		},
		RegistrarRecord{
			Name:    "testB",
			Scope:   core.GlobalScope,
			Address: "d",
		},
	}

	// load registrar
	for _, record := range records {
		r.AddServiceRecord(record)
	}

	// remove complete service
	err := r.RemoveServiceRecord(RegistrarRecord{
		Name: "testA",
	})

	if err != nil {
		t.Fatal("should be able to remove all records for one service")
	}

	if _, ok := r.registry["testA"]; ok {
		t.Fatal("there should be no records for testA service")
	}

	err = r.RemoveServiceRecord(RegistrarRecord{
		Name:  "testB",
		Scope: core.LocalScope,
	})

	if err != nil {
		t.Fatal("should be able to remove service testB local scope addresses")
	}

	if _, ok := r.registry["testB"][core.LocalScope]; ok {
		t.Fatal("should havea prunes testB scope.LocalScope")
	}

	err = r.RemoveServiceRecord(RegistrarRecord{
		Name:    "testB",
		Scope:   core.GlobalScope,
		Address: "a",
	})

	if err != nil {
		t.Fatal("should be no error for deleting record that doesnt exist")
	}

	if len(r.registry["testB"][core.GlobalScope]) != 2 {
		t.Fatal("issue with saved records for testB core.GlobalScope")
	}

	err = r.RemoveServiceRecord(RegistrarRecord{
		Name:    "testB",
		Scope:   core.GlobalScope,
		Address: "c",
	})

	if err != nil {
		t.Fatal("should be able to remove record that exists")
	}

	if r.registry["testB"][core.GlobalScope][0] != "d" {
		t.Logf("%#v\n", *r)
		t.Fatal("remaining record is not the one expected")
	}

}

func TestRegistrarGetServiceAddresses(t *testing.T) {
	r := NewRegistrar()

	records := []RegistrarRecord{
		RegistrarRecord{
			Name:    "testA",
			Scope:   core.LocalScope,
			Address: "1",
		},
		RegistrarRecord{
			Name:    "testA",
			Scope:   core.LocalScope,
			Address: "2",
		},
		RegistrarRecord{
			Name:    "testA",
			Scope:   core.GlobalScope,
			Address: "a",
		},
		RegistrarRecord{
			Name:    "testA",
			Scope:   core.GlobalScope,
			Address: "b",
		},
	}

	for _, record := range records {
		r.AddServiceRecord(record)
	}

	allRecords, err := r.GetServiceAddresses(RegistrarRecord{
		Name: "testA",
	})

	if err != nil {
		t.Fatal(err.Error())
	}

	if allRecords[0] != "1" &&
		allRecords[1] != "2" &&
		allRecords[2] != "a" &&
		allRecords[3] != "b" {
		t.Logf("%#v\n", allRecords)
		t.Fatal("registrar did not return all records")
	}

	localRecords, err := r.GetServiceAddresses(RegistrarRecord{
		Name:  "testA",
		Scope: core.LocalScope,
	})

	if err != nil {
		t.Fatal(err.Error())
	}

	if localRecords[0] != "1" &&
		localRecords[1] != "2" {
		t.Fatal("Local Records were not returned")
	}

	globalRecord, err := r.GetServiceAddresses(RegistrarRecord{
		Name:    "testA",
		Scope:   core.GlobalScope,
		Address: "a",
	})

	if err != nil {
		t.Fatal(err.Error())
	}

	if globalRecord[0] != "a" {
		t.Fatal("incorrect global record returned")
	}

	globalDoesntExist, err := r.GetServiceAddresses(RegistrarRecord{
		Name:    "testA",
		Scope:   core.GlobalScope,
		Address: "c",
	})

	if err != nil {
		t.Fatal(err.Error())
	}

	if len(globalDoesntExist) != 0 {
		t.Fatal("should have returned an empty record set")
	}

	_, err = r.GetServiceAddresses(RegistrarRecord{})

	if err == nil {
		t.Fatal("should fail with a record with no name")
	}
}
