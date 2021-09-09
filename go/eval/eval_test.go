package eval

import (
	"dabble/object"
	"strconv"
	"testing"
)

func TestEval(t *testing.T) {
	tests := []struct {
		value   object.Value
		env     *object.Binding
		want    string
		wantErr bool
	}{{
		value: object.Number(1),
		env:   nil,
		want:  "1",
	}, {
		value:   object.Symbol("a"),
		env:     nil,
		wantErr: true,
	}, {
		value:   object.Cell(object.Symbol("foo"), nil),
		env:     nil,
		wantErr: true,
	}, {
		value: object.Null,
		env:   nil,
		want:  "()",
	}}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := Eval(tt.env, tt.value)
			if tt.wantErr {
				if _, ok := got.(object.Error); !ok {
					t.Errorf("given value %v env %v. want err. got %v", tt.value, tt.env, got)
				}
			} else {
				if got.Inspect() != tt.want {
					t.Errorf("given value %v env %v. want %v. got %v", tt.value, tt.env, tt.want, got.Inspect())
				}
			}
		})
	}
}
