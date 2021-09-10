package core

import (
	"dabble/object"
	"testing"
)

func TestCar(t *testing.T) {

	env := &object.Binding{"car", object.Function(Car), nil}

	tests := []coreTest{{
		input:   "(car)",
		env:     env,
		wantErr: true,
	}, {
		input: "(car 1)",
		env:   env,
		want:  "1",
	}}

	testCore(t, tests)
}
