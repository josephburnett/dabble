package core

import (
	"dabble/eval"
	"testing"
)

func TestEq(t *testing.T) {

	env := (*eval.Frame)(nil).Bind("eq", eval.Function(Eq))

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
