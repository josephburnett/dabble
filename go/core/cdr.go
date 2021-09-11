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
	if value.Type() == object.ERROR {
		return value
	}
	return value.Rest()
}
