package core

import (
	"dabble/object"
	"testing"
)

func TestUnquote(t *testing.T) {

	env := &object.Binding{"unquote", object.Function(Unquote), nil}

	tests := []coreTest{{
		input: "(unquote ())",
		env:   env,
		want:  "`()",
	}, {
		input: "(unquote (1 2 3 4))",
		env:   env,
		want:  "`(1 2 3 4)",
	}}

	testCore(t, tests)
}
