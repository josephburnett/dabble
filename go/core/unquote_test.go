package core

import (
	"testing"
)

func TestUnquote(t *testing.T) {

	tests := []coreTest{{
		input: "(unquote ())",
		want:  "`()",
	}, {
		input: "(unquote (1 2 3 4))",
		want:  "`(1 2 3 4)",
	}}

	testCore(t, Env, tests)
}
