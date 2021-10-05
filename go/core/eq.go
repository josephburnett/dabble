package core

import (
	"dabble/eval"
	"dabble/object"
)

var _ eval.Function = Eq

func Eq(env *eval.Frame, args ...object.Value) object.Value {
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
		return object.Nil
	}
	if a.Type() != object.CELL {
		if a == b {
			return object.Symbol("t")
		} else {
			return object.Nil
		}
	}
	if Eq(env, a.First(), b.First()).Type() == object.NIL {
		return object.Nil
	}
	if Eq(env, a.Rest(), b.Rest()).Type() == object.NIL {
		return object.Nil
	}
	return object.Symbol("t")
}
