package registrar

import (
	"testing"

	"github.com/jtaylorcpp/gerl/core"
)

func TestRegistrar(t *testing.T) {
	reg := New("blank state", core.GlobalScope)

	t.Log("new registrar: ", reg)
}
