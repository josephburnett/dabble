package core

import (
	"dabble/eval"
	"dabble/object"
)

var _ object.Function = Cdr

func Cdr(env *object.Binding, args ...object.Value) object.Value {
	if err := argsLenError("cdr", args, 1); err != nil {
		return err
	}
	value := eval.Eval(env, args[0])
	switch v := value.(type) {
	case object.Cell:
		return v.Cdr()
	case object.Symbol:
		if len(v) < 2 {
			return object.Nil
		}
		return v[1:]
	default:
		return v
	}
}
