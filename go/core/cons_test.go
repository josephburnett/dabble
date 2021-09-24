package core

import (
	"dabble/object"
	"testing"
)

func TestCons(t *testing.T) {

	env := &object.Binding{"cons", object.Function(Cons), nil}

	tests := []coreTest{{
		input: "(cons 1 '(2 3 4))",
		env:   env,
		want:  "(1 2 3 4)",
	}, {
		input:   "(cons)",
		env:     env,
		wantErr: true,
	}, {
		input:   "(cons 1)",
		env:     env,
		wantErr: true,
	}, {
		input: "(cons 1 ())",
		env:   env,
		want:  "(1)",
	}, {
		input: "(cons 1 2)",
		env:   env,
		want:  "(1 2)", // Should be "(1 . 2)".
	}, {
		input: "(cons 1 '(2))",
		env:   env,
		want:  "(1 2)",
	}}

	testCore(t, tests)
}
