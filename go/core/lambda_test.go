package core

import (
	"dabble/eval"
	"dabble/object"
	"fmt"
	"testing"
)

func TestLambda(t *testing.T) {

	adder := func(env *object.Binding, args ...object.Value) object.Value {
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

	env := &object.Binding{"lambda", object.Function(Lambda),
		&object.Binding{"quote", object.Function(Quote),
			&object.Binding{"+", object.Function(adder), nil}}}

	tests := []coreTest{{
		input: "((lambda () 1))",
		env:   env,
		want:  "1",
	}, {
		input: "((lambda (a) a) 1)",
		env:   env,
		want:  "1",
	}, {
		input: "((lambda (a b) (+ a b)) 1 2)",
		env:   env,
		want:  "3",
	}}

	testCore(t, tests)
}
