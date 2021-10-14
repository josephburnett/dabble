package core

import (
	"dabble/eval"
	"dabble/object"
)

func Recur(env *eval.Frame, args ...object.Value) object.Value {
	lastCaller := env.LastCaller()
	eval.T("recurring on %v", lastCaller)
	return lastCaller.Fn(env, args...)
}
