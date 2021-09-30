package core

import (
	"dabble/object"
	"testing"
)

func TestMacro(t *testing.T) {

	env := &object.Binding{"macro", object.Function(Macro),
		&object.Binding{"car", object.Function(Car), nil}}

	tests := []coreTest{{
		input: "((macro (x xs) '(car `x)) 1 2)",
		env:   env,
		want:  "1",
	}}

	testCore(t, tests)
}
