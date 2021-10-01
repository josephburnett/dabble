package core

import (
	"dabble/object"
	"testing"
)

func TestCar(t *testing.T) {

	env := &object.Binding{"car", object.Function(Car), nil}

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
