package core

import (
	"dabble/eval"
	"testing"
)

func TestUnquote(t *testing.T) {

	env := (*eval.Frame)(nil).Bind("unquote", eval.Function(Unquote))

	tests := []coreTest{{
		input: "(unquote ())",
		want:  "`()",
	}, {
		input: "(unquote (1 2 3 4))",
		want:  "`(1 2 3 4)",
	}}

	testCore(t, env, tests)
}
