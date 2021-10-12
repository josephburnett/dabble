package core

import (
	"dabble/eval"
	"dabble/object"
	"fmt"
)

func Label(env *eval.Frame, args ...object.Value) object.Value {
	if err := argsLenError("label", args, 3); err != nil {
		return err
	}
	symbol := args[0]
	if symbol.Type() != object.SYMBOL {
		return object.Error(fmt.Sprintf("label non-symbol binding: %v", symbol))
	}
	env = env.Bind(symbol.(object.Symbol), args[1])
	return eval.Eval(env, args[2])
}
