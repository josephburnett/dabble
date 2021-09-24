package core

import "dabble/object"

var _ object.Function = Unquote

func Unquote(env *object.Binding, args ...object.Value) object.Value {
	if err := argsLenError("unquote", args, 1); err != nil {
		return err
	}
	return object.Unquoted(args[0])
}
