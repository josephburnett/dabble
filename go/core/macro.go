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
	var rest bool
	for f.Type() != object.NIL {
		symbol := f.First()
		if symbol.Type() == object.CELL && symbol.First().Type() == object.SYMBOL && symbol.Rest().Type() == object.NIL {
			rest = true
			symbol = symbol.First()
		}
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
	return makeMacro(env, free, rest, form)
}

func makeMacro(macroEnv *object.Binding, free []object.Symbol, haveRest bool, form object.Value) object.Function {
	return func(env *object.Binding, args ...object.Value) object.Value {
		if len(args) < len(free) {
			return object.Error("not enough arguments to macro")
		}
		if !haveRest && len(args) != len(free) {
			return object.Error("wrong number of arguments to macro")
		}
		var i int
		for i = 0; i < len(free)-1; i++ {
			macroEnv = &object.Binding{
				Symbol: free[i],
				Value:  args[i],
				Next:   macroEnv,
			}

		}
		var rest object.Value
		if haveRest {
			rest = object.Nil
			for j := i; j < len(args); j++ {
				rest = object.Cell(args[j], rest)
			}
			macroEnv = &object.Binding{
				Symbol: free[i],
				Value:  rest,
				Next:   macroEnv,
			}
		} else {
			macroEnv = &object.Binding{
				Symbol: free[i],
				Value:  args[i],
				Next:   macroEnv,
			}
		}
		expandedForm := eval.Eval(macroEnv, form)
		eval.T("expanded macro form: %v", expandedForm)
		if expandedForm.Type() == object.ERROR {
			return expandedForm
		}
		return eval.Eval(env, expandedForm)
	}
}
