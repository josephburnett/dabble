package core

import (
	"dabble/eval"
	"dabble/object"
)

var _ eval.Function = Unquote

func Unquote(env *eval.Frame, args ...object.Value) object.Value {
	if err := argsLenError("unquote", args, 1); err != nil {
		return err
	}
	return object.Unquoted(args[0])
}
