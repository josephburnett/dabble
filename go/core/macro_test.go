package core

import (
	"dabble/object"
	"testing"
)

func XTestMacro(t *testing.T) {

	env := &object.Binding{"macro", object.Function(Macro),
		&object.Binding{"car", object.Function(Car), nil}}

	tests := []coreTest{{
		input: "((macro (x) '(car `x)) 1)",
		env:   env,
		want:  "1",
	}}

	testCore(t, tests)
}
