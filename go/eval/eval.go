package eval

import (
	"dabble/object"
	"fmt"
)

func Eval(env *object.Binding, value object.Value) object.Value {
	return eval(env, false, value)
}

func eval(env *object.Binding, quoted bool, value object.Value) (ret object.Value) {
	t.In()
	defer t.Out()
	defer func() {
		T("returning %v", ret)
	}()
	switch value.Type() {
	case object.NUMBER, object.FUNCTION, object.NIL, object.ERROR:
		T("self evaluation of %v", value)
		return value
	case object.SYMBOL:
		if quoted {
			T("quoted symbol %v", value)
			return value
		} else {
			r := env.Resolve(value.(object.Symbol))
			if r.Type() == object.ERROR {
				T("error resolving symbol %v in environment %v", value, env)
				return r
			}
			T("resolved symbol %v to %v", value, r)
			return r
		}
	case object.CELL:
		if quoted {
			T("eval first %v", value.First())
			first := eval(env, quoted, value.First())
			if first.Type() == object.ERROR {
				return first
			}
			T("eval rest %v", value.Rest())
			rest := eval(env, quoted, value.Rest())
			if rest.Type() == object.ERROR {
				return rest
			}
			return object.Cell(first, rest)
		} else {
			T("calling cell %v", value)
			return call(env, quoted, value)
		}
	case object.QUOTED:
		if quoted {
			T("looking for unquotes in quoted value")
			q := eval(env, true, value.First())
			if q.Type() == object.ERROR {
				return q
			}
			return object.Quoted(q)
		} else {
			T("unwrapping quoted %v", value)
			return eval(env, true, value.First())
		}
	case object.UNQUOTED:
		T("evaluating within unquoted %v", value)
		return eval(env, false, value.First())
	default:
		return object.Error(fmt.Sprintf("eval: unknown type: %T", value))
	}
}

func call(env *object.Binding, quoted bool, cell object.Value) (ret object.Value) {
	t.In()
	defer t.Out()
	defer func() {
		T("returning %v", ret)
	}()
	T("evaluting %v", cell.First())
	first := eval(env, quoted, cell.First())
	if first.Type() == object.ERROR {
		return first
	}
	if first.Type() != object.FUNCTION && first.Type() != object.CLOSURE {
		return object.Error(fmt.Sprintf("calling non-function: %v", first.String()))
	}
	rest := cell.Rest()
	args := []object.Value{}
	for rest.Type() == object.CELL {
		args = append(args, rest.First())
		rest = rest.Rest()
	}

	T("calling %v with args %v", first, cell.Rest())
	if first.Type() == object.FUNCTION {
		function := first.(object.Function)
		return function(env, args...)
	} else {
		closure := first.(object.Closure)
		return closure(args...)
	}
}
