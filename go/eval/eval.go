package eval

import (
	"dabble/object"
	"fmt"
)

func Eval(env, value object.Value) object.Value {
	if value == nil || value == object.Null {
		return object.Null
	}
	switch v := value.(type) {
	case object.Number:
		return v
	case object.Symbol:
		r := resolve(env, string(v))
		return Eval(env, r)
	case object.Cell:
		first := Eval(env, v[0])
		if _, ok := first.(object.Error); ok {
			return first
		}
		rest := Eval(env, v[1])
		if _, ok := rest.(object.Error); ok {
			return rest
		}
		return object.Cell{first, rest}
	default:
		return object.Error(fmt.Sprintf("unknown type: %T", value))
	}
}
