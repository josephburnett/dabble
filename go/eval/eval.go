package eval

import (
	"dabble/object"
	"fmt"
)

func Eval(env *object.Binding, value object.Value) object.Value {
	return eval(env, false, nil, value)
}

func eval(env *object.Binding, quoted bool, trace *Trace, value object.Value) object.Value {
	switch value.Type() {
	case object.NUMBER, object.FUNCTION, object.NIL, object.ERROR:
		trace.T("self evaluation of %v", value)
		return value
	case object.SYMBOL:
		if quoted {
			trace.T("quoted symbol %v", value)
			return value
		} else {
			r := env.Resolve(value.(object.Symbol))
			trace.T("resolved symbol %v to %v", value, r)
			return eval(env, quoted, trace, r)
		}
	case object.CELL:
		if quoted {
			first := eval(env, quoted, trace, value.First())
			if first.Type() == object.ERROR {
				return first
			}
			rest := eval(env, quoted, trace, value.Rest())
			if rest.Type() == object.ERROR {
				return rest
			}
			return object.Cell(first, rest)
		} else {
			return call(env, quoted, trace, value)
		}
	case object.QUOTED:
		return eval(env, true, trace, value.First())
	case object.UNQUOTED:
		return eval(env, false, trace, value.First())
	default:
		return object.Error(fmt.Sprintf("eval: unknown type: %T", value))
	}
}

func call(env *object.Binding, quoted bool, trace *Trace, cell object.Value) object.Value {
	first := eval(env, quoted, trace, cell.First())
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
