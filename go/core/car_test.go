package core

import (
	"dabble/object"
	"testing"
)

func TestCar(t *testing.T) {

	env := &object.Binding{"car", object.Function(Car),
		&object.Binding{"quote", object.Function(Quote), nil}}

	tests := []coreTest{{
		input:   "(car)",
		env:     env,
		wantErr: true,
	}, {
		input: "(car 1)",
		env:   env,
		want:  "1",
	}, {
		input: "(car (quote (1 2 3 4)))",
		env:   env,
		want:  "1",
	}, {
		input: "(car ())",
		env:   env,
		want:  "()",
	}}

	testCore(t, tests)
}
