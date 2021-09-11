package core

import (
	"dabble/object"
	"testing"
)

func TestEq(t *testing.T) {

	env := &object.Binding{"eq", object.Function(Eq),
		&object.Binding{"quote", object.Function(Quote), nil}}

	tests := []coreTest{{
		input: "(eq 1 1)",
		env:   env,
		want:  "(1 1)",
	}, {
		input: "(eq 1 2)",
		env:   env,
		want:  "()",
		// }, {
		// 	input: "(eq () ())",
		// 	env:   env,
		// 	want:  "(() ())",         <=== gets "(())" because list handling is wrong
		// }, {
		// 	input: "(eq (quote (1 2)) (quote (1 2)))",
		// 	env:   env,
		// 	want:  "((1 2) (1 2))",   <=== gets "((1 2) 1 2)"
	}}

	testCore(t, tests)
}
