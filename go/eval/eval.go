package eval

import (
	"dabble/object"
	"fmt"
)

func Eval(env *object.Binding, value object.Value) object.Value {
	return eval(env, false, value)
}

func eval(env *object.Binding, quoted bool, value object.Value) object.Value {
	switch value.Type() {
	case object.NUMBER, object.FUNCTION, object.NIL, object.ERROR:
		return value
	case object.SYMBOL:
		if quoted {
			return value
		} else {
			r := env.Resolve(value.(object.Symbol))
			return eval(env, quoted, r)
		}
	case object.CELL:
		if quoted {
			first := eval(env, quoted, value.First())
			if first.Type() == object.ERROR {
				return first
			}
			rest := eval(env, quoted, value.Rest())
			if rest.Type() == object.ERROR {
				return rest
			}
			return object.Cell(first, rest)
		} else {
			return call(env, value)
		}
	case object.QUOTED:
		return eval(env, true, value)
	case object.UNQUOTED:
		return eval(env, false, value)
	default:
		return object.Error(fmt.Sprintf("eval: unknown type: %T", value))
	}
}

func call(env *object.Binding, cell object.Value) object.Value {
	first := Eval(env, cell.First())
	if first.Type() == object.ERROR {
		return first
	}
	if first.Type() != object.FUNCTION && first.Type() != object.CLOSURE {
		return object.Error(fmt.Sprintf("calling non-function: %v", first.Inspect()))
	}
	rest := cell.Rest()
	args := []object.Value{}
	for rest.Type() == object.CELL {
		args = append(args, rest.First())
		rest = rest.Rest()
	}

	if first.Type() == object.FUNCTION {
		function := first.(object.Function)
		return function(env, args...)
	} else {
		closure := first.(object.Closure)
		return closure(args...)
	}
}
