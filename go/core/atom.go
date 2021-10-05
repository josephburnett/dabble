package core

import (
	"dabble/eval"
	"dabble/object"
)

var _ eval.Function = Atom

func Atom(env *eval.Frame, args ...object.Value) object.Value {
	if err := argsLenError("atom", args, 1); err != nil {
		return err
	}
	value := eval.Eval(env, args[0])
	switch value.Type() {
	case object.ERROR:
		return value
	case object.CELL:
		return object.Nil
	default:
		return object.Cell(value, nil)
	}
}
