package eval

import (
	"dabble/object"
	"fmt"
)

func resolve(env object.Value, symbol string) object.Value {
	if env == nil || env == object.Null {
		return object.Error(fmt.Sprintf("symbol not bound: %q", symbol))
	}
	binding := env.First()
	if binding.Type() != object.CELL {
		return object.Error(fmt.Sprintf("invalid environment. want cell. got %T (env: %q)", env.First(), env.Inspect()))
	}
	s, ok := binding.First().(object.Symbol)
	if !ok {
		return object.Error(fmt.Sprintf("invalid environment. want symbol binding. got %T (%q)", s, s.Inspect()))
	}
	if string(s) == symbol {
		return binding.Rest()
	}
	return resolve(env.Rest(), symbol)
}
