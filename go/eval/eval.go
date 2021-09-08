package eval

import (
	"dabble/object"
	"fmt"
)

func Eval(env, value object.Value) object.Value {
	switch value.Type() {
	case object.NUMBER, object.NULL, object.ERROR:
		return value
	case object.SYMBOL:
		r := resolve(env, value)
		return Eval(env, r)
	case object.CELL:
		first := Eval(env, value.First())
		if first.Type() == object.ERROR {
			return first
		}
		rest := Eval(env, value.Rest())
		if rest.Type() == object.ERROR {
			return rest
		}
		return object.Cell(first, rest)
	default:
		return object.Error(fmt.Sprintf("eval: unknown type: %T", value))
	}
}
