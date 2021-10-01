package core

import (
	"dabble/object"
	"testing"
)

func TestUnquote(t *testing.T) {

	env := &object.Binding{"unquote", object.Function(Unquote), nil}

	tests := []coreTest{{
		input: "(unquote ())",
		want:  "`()",
	}, {
		input: "(unquote (1 2 3 4))",
		want:  "`(1 2 3 4)",
	}}

	testCore(t, env, tests)
}
