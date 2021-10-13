package core

import (
	"dabble/eval"
	"dabble/object"
	"fmt"
)

func Apply(env *eval.Frame, args ...object.Value) object.Value {
	if err := argsLenError("apply", args, 2); err != nil {
		return err
	}
	function := eval.Eval(env, args[0])
	if function.Type() == object.ERROR {
		return function
	}
	if function.Type() != object.FUNCTION {
		return object.Error(fmt.Sprintf("apply non function %v", function))
	}
	argsList := eval.Eval(env, args[1])
	if argsList.Type() == object.ERROR {
		return argsList
	}
	if argsList.Type() != object.CELL && argsList.Type() != object.NIL {
		return object.Error(fmt.Sprintf("apply to non list %v", argsList))
	}
	flattenedArgs := []object.Value{}
	for argsList.Type() != object.NIL {
		flattenedArgs = append(flattenedArgs, argsList.First())
		argsList = argsList.Rest()
	}
	return function.(*eval.Function).Fn(env, flattenedArgs...)
}
