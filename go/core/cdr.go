package core

import "dabble/object"

func Cdr(args ...object.Value) object.Value {
	if err := argsLenError("cdr", args, 1); err != nil {
		return err
	}
	return args[0].Rest()
}
