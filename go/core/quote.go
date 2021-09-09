package core

import "dabble/object"

var _ object.Function = Quote

func Quote(env object.Value, args ...object.Value) object.Value {
	if err := argsLenError("quote", args, 1); err != nil {
		return err
	}
	return args[0]
}
