package core

import (
	"dabble/object"
	"testing"
)

func TestEq(t *testing.T) {

	env := &object.Binding{"eq", object.Function(Eq), nil}

	tests := []coreTest{{
		input: "(eq 1 1)",
		env:   env,
		want:  "t",
	}, {
		input: "(eq 1 2)",
		env:   env,
		want:  "()",
	}, {
		input: "(eq () ())",
		env:   env,
		want:  "t",
	}, {
		input: "(eq '(1 2) '(1 2))",
		env:   env,
		want:  "t",
	}, {
		input: "(eq 'abc 'abc)",
		env:   env,
		want:  "t",
	}, {
		input: "(eq 'abc 'cba)",
		env:   env,
		want:  "()",
	}}

	testCore(t, tests)
}
