package core

import (
	"dabble/eval"
	"dabble/object"
	"fmt"
)

func Error(_ *eval.Frame, args ...object.Value) object.Value {
	if err := argsLenError("error", args, 1); err != nil {
		return err
	}
	if args[0].Type() != object.SYMBOL {
		return object.Error(fmt.Sprintf("non-symbol error: %v", args[0]))
	}
	return object.Error(args[0].(object.Symbol))
}
