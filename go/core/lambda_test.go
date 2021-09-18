package core

import (
	"dabble/object"
	"testing"
)

func TestLambda(t *testing.T) {

	env := &object.Binding{"lambda", object.Function(Lambda),
		&object.Binding{"quote", object.Function(Quote), nil}}

	tests := []coreTest{{
		input: "((lambda () 1))",
		env:   env,
		want:  "1",
	}}

	testCore(t, tests)
}
