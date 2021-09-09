package core

import (
	"dabble/eval"
	"dabble/object"
)

var _ object.Function = Car

func Car(env object.Value, args ...object.Value) object.Value {
	if err := argsLenError("car", args, 1); err != nil {
		return err
	}
	value := eval.Eval(env, args[0])
	if value.Type() == object.ERROR {
		return value
	}
	return value.First()
}
