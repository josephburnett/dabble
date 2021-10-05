package core

import (
	"dabble/eval"
	"testing"
)

func TestQuote(t *testing.T) {

	env := (*eval.Frame)(nil).Bind("quote", eval.Function(Quote))

	tests := []coreTest{{
		input: "(quote ())",
		want:  "'()",
	}, {
		input: "(quote (1 2 3 4))",
		want:  "'(1 2 3 4)",
	}}

	testCore(t, env, tests)
}
