package core

import (
	"dabble/eval"
	"dabble/object"
)

func Cons(env *eval.Frame, args ...object.Value) object.Value {
	if err := argsLenError("cons", args, 2); err != nil {
		return err
	}
	car := eval.Eval(env, args[0])
	if car.Type() == object.ERROR {
		return car
	}
	cdr := eval.Eval(env, args[1])
	if cdr.Type() == object.ERROR {
		return cdr
	}
	return object.Cell(car, cdr)
}
