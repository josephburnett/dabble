package core

import (
	"dabble/object"
	"testing"
)

func TestQuote(t *testing.T) {

	env := &object.Binding{"quote", object.Function(Quote), nil}

	tests := []coreTest{{
		input: "(quote ())",
		env:   env,
		want:  "()",
	}, {
		input: "(quote (1 2 3 4))",
		env:   env,
		want:  "(1 2 3 4)",
	}}

	testCore(t, tests)
}
