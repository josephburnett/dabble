package core

import (
	"dabble/eval"
	"dabble/object"
	"fmt"
)

var _ eval.Function = Lambda

func Lambda(env *eval.Frame, args ...object.Value) object.Value {
	if err := argsLenError("lambda", args, 2); err != nil {
		return err
	}
	free := []object.Symbol{}
	f := args[0]
	if f.Type() != object.CELL && f.Type() != object.NIL {
		return object.Error(fmt.Sprintf("lambda non-list params: %v", f))
	}
	for f.Type() != object.NIL {
		symbol := f.First()
		if symbol.Type() != object.SYMBOL {
			return object.Error(fmt.Sprintf("lambda non-symbol param: %v", f))
		}
		free = append(free, symbol.(object.Symbol))
		f = f.Rest()
	}
	form := args[1]
	return makeClosure(env, free, form)
}

func makeClosure(env *eval.Frame, free []object.Symbol, form object.Value) eval.Function {
	return func(_ *eval.Frame, args ...object.Value) object.Value {
		if err := argsLenError("lambda args", args, len(free)); err != nil {
			return err
		}
		for i, f := range free {
			value := eval.Eval(env, args[i])
			if value.Type() == object.ERROR {
				return value
			}
			env = env.Bind(f, value)
		}
		return eval.Eval(env, form)
	}
}
