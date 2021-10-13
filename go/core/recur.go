package core

import (
	"dabble/eval"
	"dabble/object"
)

func Recur(env *eval.Frame, args ...object.Value) object.Value {
	// Is this necessary?
	// for i, arg := range args {
	// 	value := eval.Eval(env, arg)
	// 	if value.Type() == object.ERROR {
	// 		return value
	// 	}
	// 	args[i] = value
	// }
	lastCaller := env.LastCaller()
	eval.T("recurring on %v", lastCaller)
	return lastCaller.Fn(env, args...)
}
