package core

import "dabble/object"

var _ object.Function = If

func If(env object.Value, args ...object.Value) object.Value {
	if err := argsLenError("if", args, 3); err != nil {
		return err
	}

	return args[1]
}
