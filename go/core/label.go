package core

import (
	"dabble/eval"
	"dabble/object"
	"fmt"
)

var _ object.Function = Label

func Label(env *object.Binding, args ...object.Value) object.Value {
	if err := argsLenError("car", args, 3); err != nil {
		return err
	}
	symbol := args[0]
	if symbol.Type() != object.SYMBOL {
		return object.Error(fmt.Sprintf("label non-symbol binding: %v", symbol))
	}
	env = &object.Binding{
		Symbol: symbol.(object.Symbol),
		Value:  args[1],
		Next:   env,
	}
	return eval.Eval(env, args[2])
}
