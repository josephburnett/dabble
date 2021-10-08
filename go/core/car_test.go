package core

import (
	"testing"
)

func TestCar(t *testing.T) {

	tests := []coreTest{{
		input:   "(car)",
		wantErr: true,
	}, {
		input: "(car 1)",
		want:  "1",
	}, {
		input: "(car '(1 2 3 4))",
		want:  "1",
	}, {
		input: "(car ())",
		want:  "()",
	}}

	testCore(t, Env, tests)
}
