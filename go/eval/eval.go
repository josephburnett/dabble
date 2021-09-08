package eval

import (
	"dabble/object"
	"fmt"
)

func Eval(env, value object.Value) object.Value {
	if value == nil || value == object.Null {
		return object.Null
	}
	switch value.Type() {
	case object.NUMBER:
		return value
	case object.SYMBOL:
		r := resolve(env, value)
		return Eval(env, r)
	case object.CELL:
		first := Eval(env, value.First())
		if _, ok := first.(object.Error); ok {
			return first
		}
		rest := Eval(env, value.Rest())
		if _, ok := rest.(object.Error); ok {
			return rest
		}
		return object.Cell(first, rest)
	default:
		return object.Error(fmt.Sprintf("unknown type: %T", value))
	}
}
