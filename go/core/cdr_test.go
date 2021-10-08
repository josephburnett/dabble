package core

import (
	"testing"
)

func TestCdr(t *testing.T) {

	tests := []coreTest{{
		input: "(cdr 2)",
		want:  "1",
	}, {
		input: "(cdr '(1 2 3 4))",
		want:  "(2 3 4)",
	}}

	testCore(t, Env, tests)
}
