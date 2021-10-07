package core

import (
	"dabble/eval"
	"dabble/object"
)

var _ eval.Function = Recur

func Recur(env *eval.Frame, args ...object.Value) object.Value {
	return env.LastCaller()(env, args...)
}
