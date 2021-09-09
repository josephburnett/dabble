package core

import (
	"dabble/eval"
	"dabble/object"
)

var _ object.Function = Atom

func Atom(env object.Environment, args ...object.Value) object.Value {
	if err := argsLenError("atom", args, 1); err != nil {
		return err
	}
	value := eval.Eval(env, args[0])
	switch value.Type() {
	case object.ERROR:
		return value
	case object.CELL:
		return object.Null
	default:
		return object.Cell(value, nil)
	}
}
