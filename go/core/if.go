package core

import (
	"dabble/eval"
	"dabble/object"
)

var _ eval.Function = If

func If(env *eval.Frame, args ...object.Value) object.Value {
	if err := argsLenError("if", args, 3); err != nil {
		return err
	}

	return args[1]
}
