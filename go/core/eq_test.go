package core

import (
	"dabble/object"
	"testing"
)

func TestEq(t *testing.T) {

	env := &object.Binding{"eq", object.Function(Eq), nil}

	tests := []coreTest{{
		input: "(eq 1 1)",
		want:  "t",
	}, {
		input: "(eq 1 2)",
		want:  "()",
	}, {
		input: "(eq () ())",
		want:  "t",
	}, {
		input: "(eq '(1 2) '(1 2))",
		want:  "t",
	}, {
		input: "(eq 'abc 'abc)",
		want:  "t",
	}, {
		input: "(eq 'abc 'cba)",
		want:  "()",
	}}

	testCore(t, env, tests)
}
