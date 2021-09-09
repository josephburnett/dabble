package core

import "dabble/object"

func Atom(args ...object.Value) object.Value {
	if err := argsLenError("atom", args, 1); err != nil {
		return err
	}
	if args[0].Type() == object.CELL {
		return object.Null
	}
	return object.Symbol("t")
}
