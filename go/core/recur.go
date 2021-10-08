package core

import (
	"dabble/eval"
	"dabble/object"
)

func Recur(env *eval.Frame, args ...object.Value) object.Value {
	return env.LastCaller().Fn(env, args...)
}
