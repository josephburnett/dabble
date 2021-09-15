package core

import (
	"dabble/eval"
	"dabble/object"
)

var _ object.Function = Cons

func Cons(env *object.Binding, args ...object.Value) object.Value {
	if err := argsLenError("cons", args, 2); err != nil {
		return err
	}
	car := eval.Eval(env, args[0])
	if _, ok := car.(object.Error); ok {
		return car
	}
	cdr := eval.Eval(env, args[1])
	if _, ok := cdr.(object.Error); ok {
		return cdr
	}
	return object.Cell{car, cdr}
}
