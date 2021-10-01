package core

import (
	"dabble/object"
	"testing"
)

func TestCdr(t *testing.T) {

	env := &object.Binding{"cdr", object.Function(Cdr), nil}

	tests := []coreTest{{
		input: "(cdr 2)",
		want:  "1",
	}, {
		input: "(cdr '(1 2 3 4))",
		want:  "(2 3 4)",
	}}

	testCore(t, env, tests)
}
