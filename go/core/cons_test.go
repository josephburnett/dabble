package core

import (
	"dabble/eval"
	"testing"
)

func TestCons(t *testing.T) {

	env := eval.NilFrame.Bind("cons", eval.Function(Cons))

	tests := []coreTest{{
		input: "(cons 1 '(2 3 4))",
		want:  "(1 2 3 4)",
	}, {
		input:   "(cons)",
		wantErr: true,
	}, {
		input:   "(cons 1)",
		wantErr: true,
	}, {
		input: "(cons 1 ())",
		want:  "(1)",
	}, {
		input: "(cons 1 2)",
		want:  "(1 2)", // Should be "(1 . 2)".
	}, {
		input: "(cons 1 '(2))",
		want:  "(1 2)",
	}}

	testCore(t, env, tests)
}
