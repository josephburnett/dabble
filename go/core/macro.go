package core

import (
	"dabble/eval"
	"dabble/object"
	"fmt"
)

var _ object.Function = Macro

func Macro(env *object.Binding, args ...object.Value) object.Value {
	if err := argsLenError("macro", args, 2); err != nil {
		return err
	}
	free := []object.Symbol{}
	f := args[0]
	if f.Type() != object.CELL && f.Type() != object.NIL {
		return object.Error(fmt.Sprintf("macro non-list params: %v", f))
	}
	for f.Type() != object.NIL {
		symbol := f.First()
		if symbol.Type() != object.SYMBOL {
			return object.Error(fmt.Sprintf("macro non-symbol param: %v", f))
		}
		free = append(free, symbol.(object.Symbol))
		f = f.Rest()
	}
	form := args[1]
	if len(free) == 0 {
		return object.Error("macro requires at least one free variable")
	}
	return makeMacro(env, free, form)
}

func makeMacro(macroEnv *object.Binding, free []object.Symbol, form object.Value) object.Function {
	return func(env *object.Binding, args ...object.Value) object.Value {
		if len(args) < len(free) {
			return object.Error("no enough arguments to macro")
		}
		var i int
		for i = 0; i < len(free)-1; i++ {
			env = &object.Binding{
				Symbol: free[i],
				Value:  args[i],
				Next:   env,
			}
		}
		var rest object.Value = object.Nil
		for j := i; j < len(args); j++ {
			rest = object.Cell(args[j], rest)
		}
		env = &object.Binding{
			Symbol: free[i],
			Value:  rest,
			Next:   env,
		}
		expandedForm := eval.Eval(macroEnv, form)
		if expandedForm.Type() == object.ERROR {
			return expandedForm
		}
		return eval.Eval(env, expandedForm)
	}
}
