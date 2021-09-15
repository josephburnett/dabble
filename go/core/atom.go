package core

import (
	"dabble/eval"
	"dabble/object"
)

var _ object.Function = Atom

func Atom(env *object.Binding, args ...object.Value) object.Value {
	if err := argsLenError("atom", args, 1); err != nil {
		return err
	}
	value := eval.Eval(env, args[0])
	switch value.(type) {
	case object.Error:
		return value
	case object.Cell:
		return object.Nil
	default:
		return object.Symbol("t")
	}
}
