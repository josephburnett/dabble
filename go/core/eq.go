package core

import (
	"reflect"
	"dabble/eval"
	"dabble/object"
)

var _ object.Function = Eq

func Eq(env *object.Binding, args ...object.Value) object.Value {
	if err := argsLenError("eq", args, 2); err != nil {
		return err
	}
	a := eval.Eval(env, args[0])
	if _, ok := a.(object.Error); ok {
		return a
	}
	b := eval.Eval(env, args[1])
	if _, ok := b.(object.Error); ok {
		return b
	}
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return object.Nil
	}
	if _, ok := a.(object.Cell); !ok {
		if a == b {
			return object.Symbol("t")
		} else {
			return object.Nil
		}
	}
	if Eq(env, a.(object.Cell).Car(), b.(object.Cell).Car()) == object.Nil {
		return object.Nil
	}
	if Eq(env, a.(object.Cell).Cdr(), b.(object.Cell).Cdr()) == object.Nil { 
		return object.Nil
	}
	return object.Symbol("t")
}
