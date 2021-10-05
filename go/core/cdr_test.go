package core

import (
	"dabble/eval"
	"testing"
)

func TestCdr(t *testing.T) {

	env := (*eval.Frame)(nil).Bind("cdr", eval.Function(Cdr))

	tests := []coreTest{{
		input: "(cdr 2)",
		want:  "1",
	}, {
		input: "(cdr '(1 2 3 4))",
		want:  "(2 3 4)",
	}}

	testCore(t, env, tests)
}
