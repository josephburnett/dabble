package core

import (
	"dabble/object"
	"testing"
)

func TestCdr(t *testing.T) {

	env := &object.Binding{"cdr", object.Function(Cdr),
		&object.Binding{"quote", object.Function(Quote), nil}}

	tests := []coreTest{{
		input: "(cdr (quote (1 2 3 4)))",
		env:   env,
		want:  "(2 3 4)",
	}}

	testCore(t, tests)
}
