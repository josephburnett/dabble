package core

import (
	"dabble/object"
	"testing"
)

func TestMacro(t *testing.T) {

	env := &object.Binding{"macro", object.Function(Macro),
		&object.Binding{"car", object.Function(Car),
			&object.Binding{"cdr", object.Function(Cdr), nil}}}

	tests := []coreTest{{
		input: "((macro (x) `x) 1)",
		env:   env,
		want:  "1",
	}, {
		input: "((macro (x y) '(`x `y)) 1 2)",
		env:   env,
		want:  "(1 2)",
	}, {
		input: "((macro ((xs)) (cdr '`xs)) 1 2 3)",
		env:   env,
		want:  "(2 3)",
	}}

	testCore(t, tests)
}
