package core

import (
	"dabble/eval"
	"dabble/object"
)

var _ object.Function = Car

func Car(env *object.Binding, args ...object.Value) object.Value {
	if err := argsLenError("car", args, 1); err != nil {
		return err
	}
	value := eval.Eval(env, args[0])
	switch v := value.(type) {
	case object.Cell:
		return v.Car()
	case object.Symbol:
		if len(v) == 0 {
			return object.Nil
		}
		return v[0:1]
	default:
		return v
	}
}
