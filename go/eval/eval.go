package eval

import (
	"dabble/object"
	"fmt"
)

func Eval(env *object.Binding, value object.Value) object.Value {
	switch value.Type() {
	case object.NUMBER, object.FUNCTION, object.NULL, object.ERROR:
		return value
	case object.SYMBOL:
		r := env.Resolve(value.(object.Symbol))
		return Eval(env, r)
	case object.CELL:
		return call(env, value)
	default:
		return object.Error(fmt.Sprintf("eval: unknown type: %T", value))
	}
}

func call(env *object.Binding, cell object.Value) object.Value {
	first := Eval(env, cell.First())
	if first.Type() == object.ERROR {
		return first
	}
	if first.Type() != object.FUNCTION {
		return object.Error(fmt.Sprintf("calling non-function: %v", first))
	}
	rest := cell.Rest()
	args := []object.Value{}
	for rest.Type() == object.CELL {
		arg := Eval(env, rest.First())
		if arg.Type() == object.ERROR {
			return arg
		}
		args = append(args, arg)
		rest = rest.Rest()
	}

	function := first.(object.Function)
	return function(env, args...)
}
