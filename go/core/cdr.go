package core

import "dabble/object"

var _ object.Function = Cdr

func Cdr(env object.Value, args ...object.Value) object.Value {
	if err := argsLenError("cdr", args, 1); err != nil {
		return err
	}
	return args[0].Rest()
}
