package eval

import (
	"dabble/object"
	"fmt"
)

func Eval(env *object.Binding, value object.Value) object.Value {
	if value == object.Nil {
		return value
	}
	switch value.(type) {
	case object.Number, object.Function, object.Error:
		return value
	case object.Symbol:
		r := env.Resolve(value.(object.Symbol))
		return Eval(env, r)
	case object.Cell:
		return call(env, value)
	default:
		return object.Error(fmt.Sprintf("eval: unknown type: %T", value))
	}
}

func call(env *object.Binding, cell object.Value) object.Value {
	first := Eval(env, cell.(object.Cell).Car())
	if _, ok := first.(object.Error); ok {
		return first
	}
	if _, ok := first.(object.Function); !ok {
		return object.Error(fmt.Sprintf("calling non-function: %v", first))
	}
	rest := cell.(object.Cell).Cdr()
	args := []object.Value{}
	if _, ok := rest.(object.Cell); ok {
		args = append(args, rest.(object.Cell).Car())
		rest = rest.(object.Cell).Cdr()
	}

	function := first.(object.Function)
	return function(env, args...)
}
