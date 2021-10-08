package core

import (
	"dabble/eval"
	"dabble/object"
	"testing"
)

func TestRecur(t *testing.T) {

	env := Env.Bind("-", &eval.Function{Fn: func(env *eval.Frame, args ...object.Value) object.Value {
		a := eval.Eval(env, args[0])
		b := eval.Eval(env, args[1])
		return object.Number(a.(object.Number) - b.(object.Number))
	}})

	tests := []coreTest{{
		input: `
((lambda (repeat n)
  (if (eq 0 n)
    ()
    (cons repeat (recur repeat (- n 1)))))
 3 9)`,
		want: "(3 3 3 3 3 3 3 3 3)",
	}}

	testCore(t, env, tests)
}
