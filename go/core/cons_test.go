package core

import (
	"dabble/object"
	"testing"
)

func TestCons(t *testing.T) {

	env := &object.Binding{"cons", object.Function(Cons),
		&object.Binding{"quote", object.Function(Quote), nil}}

	tests := []coreTest{{
		input: "(cons 1 (quote (2 3 4)))",
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
		want:  "(1 2)",
	},
		// One of these two ^ v is wrong.
		{
			input: "(cons 1 (quote (2)))",
			env:   env,
			want:  "(1 2)",
		}}

	testCore(t, tests)
}
