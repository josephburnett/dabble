package core

import (
	"testing"
)

func TestQuote(t *testing.T) {

	tests := []coreTest{{
		input: "(quote ())",
		want:  "'()",
	}, {
		input: "(quote (1 2 3 4))",
		want:  "'(1 2 3 4)",
	}}

	testCore(t, Env, tests)
}
