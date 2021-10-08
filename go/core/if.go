package core

import (
	"dabble/eval"
	"dabble/object"
)

func If(env *eval.Frame, args ...object.Value) object.Value {
	if err := argsLenError("if", args, 3); err != nil {
		return err
	}
	cond := eval.Eval(env, args[0])
	if cond.Type() == object.ERROR {
		return cond
	}
	eval.T("if condition evaluted to %v", cond)
	if cond.Type() == object.NIL {
		return eval.Eval(env, args[2])
	} else {
		return eval.Eval(env, args[1])
	}
}
