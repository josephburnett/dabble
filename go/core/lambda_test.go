package core

import (
	"dabble/eval"
	"dabble/object"
	"fmt"
	"testing"
)

func TestLambda(t *testing.T) {

	adder := func(env *eval.Frame, args ...object.Value) object.Value {
		if err := argsLenError("adder", args, 2); err != nil {
			return err
		}
		first := eval.Eval(env, args[0])
		if first.Type() != object.NUMBER {
			return object.Error(fmt.Sprintf("not a number: %v", first))
		}
		second := eval.Eval(env, args[1])
		if second.Type() != object.NUMBER {
			return object.Error(fmt.Sprintf("not a number: %v", second))
		}
		return first.(object.Number) + second.(object.Number)
	}

	env := eval.NilFrame.Bind("lambda", eval.Function(Lambda)).
		Bind("+", eval.Function(adder))

	tests := []coreTest{{
		input: "((lambda () 1))",
		want:  "1",
	}, {
		input: "((lambda (a) a) 1)",
		want:  "1",
	}, {
		input: "((lambda (a b) (+ a b)) 1 2)",
		want:  "3",
	}, {
		input: "((lambda (a) (+ 4 a)) 1)",
		want:  "5",
	}, {
		input:   "((lambda () 1) 2)",
		wantErr: true,
	}, {
		input:   "((lambda (a) a))",
		wantErr: true,
	}, {
		input:   "((lambda (a) a) 1 2 )",
		wantErr: true,
	}}

	testCore(t, env, tests)
}
