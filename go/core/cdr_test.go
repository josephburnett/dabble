package core

import (
	"dabble/object"
	"testing"
)

func TestCdr(t *testing.T) {

	env := &object.Binding{"cdr", object.Function(Cdr), nil}

	tests := []coreTest{{
		input: "(cdr 2)",
		env:   env,
		want:  "1",
	}, {
		input: "(cdr '(1 2 3 4))",
		env:   env,
		want:  "(2 3 4)",
	}}

	testCore(t, tests)
}
