package core

import (
	"dabble/object"
	"testing"
)

func TestMacro(t *testing.T) {

	env := &object.Binding{"macro", object.Function(Macro),
		&object.Binding{"car", object.Function(Car),
			&object.Binding{"quote", object.Function(Quote), nil}}}

	tests := []coreTest{{
		input: "((macro (x) `x) 1)",
		env:   env,
		want:  "1",
	}, {
		input: "((macro (x y) ''(`x `y)) 1 2)",
		env:   env,
		want:  "(1 2)",
	}}

	testCore(t, tests)
}
