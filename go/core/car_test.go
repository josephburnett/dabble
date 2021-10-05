package core

import (
	"dabble/eval"
	"testing"
)

func TestCar(t *testing.T) {

	env := (*eval.Frame)(nil).Bind("car", eval.Function(Car))

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

	testCore(t, env, tests)
}
