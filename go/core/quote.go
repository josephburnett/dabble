package core

import (
	"dabble/eval"
	"dabble/object"
)

func Quote(env *eval.Frame, args ...object.Value) object.Value {
	if err := argsLenError("quote", args, 1); err != nil {
		return err
	}
	return object.Quoted(args[0])
}
