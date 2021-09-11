package core

import (
	"dabble/eval"
	"dabble/object"
)

var _ object.Function = Eq

func Eq(env *object.Binding, args ...object.Value) object.Value {
	if err := argsLenError("eq", args, 2); err != nil {
		return err
	}
	a := eval.Eval(env, args[0])
	if a.Type() == object.ERROR {
		return a
	}
	b := eval.Eval(env, args[1])
	if b.Type() == object.ERROR {
		return b
	}
	if a.Type() != b.Type() {
		return object.Null
	}
	if a.Type() != object.CELL {
		if a == b {
			return object.Cell(a, b)
		} else {
			return object.Null
		}
	}
	if Eq(env, a.First(), b.First()).Type() == object.NULL {
		return object.Null
	}
	if Eq(env, a.Rest(), b.Rest()).Type() == object.NULL {
		return object.Null
	}
	return object.Cell(a, b)
}
