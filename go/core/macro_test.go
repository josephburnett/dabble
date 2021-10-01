package core

import (
	"testing"
)

func TestMacro(t *testing.T) {

	tests := []coreTest{{
		input: "((macro (x) x) 1)",
		want:  "1",
	}, {
		input: "((macro (x y) '(`x `y)) 1 2)",
		want:  "(1 2)",
	}, {
		input: "((macro ((xs)) (cdr '`xs)) 1 2 3)",
		want:  "(2 3)",
	}, {
		input: "((macro (x (xs)) (cons x xs)) 1 2 3)",
		want:  "(1 2 3)",
	}, {
		input: "((macro (x y) '(`y `x)) 1 2)",
		want:  "(2 1)",
	}}

	testCore(t, Env, tests)
}
