package object

import (
	"fmt"
)

type Binding struct {
	Symbol Symbol
	Value  Value
}

type Environment []Binding

func (e Environment) Resolve(env []Binding, symbol Symbol) Value {
	for _, binding := range e {
		if binding.Symbol == symbol {
			return binding.Value
		}
	}
	return Error(fmt.Sprintf("symbol not bound: %q", symbol))
}
