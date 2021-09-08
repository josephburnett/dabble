package eval

import (
	"dabble/object"
	"fmt"
)

func resolve(env object.Value, symbol object.Value) object.Value {
	if env == object.Null {
		return object.Error(fmt.Sprintf("symbol not bound: %q", symbol))
	}
	if env.Type() != object.CELL {
		return object.Error(fmt.Sprintf("invalid environment: %v", env))
	}
	if symbol.Type() != object.SYMBOL {
		return object.Error(fmt.Sprintf("invalid symbol: %v", symbol))
	}
	binding := env.First()
	if binding.Type() != object.CELL {
		return object.Error(fmt.Sprintf("invalid environment binding: %v", binding))
	}
	s := binding.First()
	if s.Type() != object.SYMBOL {
		return object.Error(fmt.Sprintf("invalid environment binding symbol: %v", s))
	}
	if s == symbol {
		return binding.Rest()
	}
	return resolve(env.Rest(), symbol)
}
