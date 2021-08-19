package eval

import (
	"dabble/object"
	"testing"
)

func TestEval(t *testing.T) {
	tests := []struct {
		value   object.Value
		env     object.Value
		want    string
		wantErr bool
	}{{
		value: object.Number(1),
		env:   object.Cell{},
		want:  "1",
	}, {
		value:   object.Symbol("a"),
		env:     object.Cell{},
		wantErr: true,
	}, {
		value: object.Cell{},
		env:   object.Cell{},
		want:  "(() ())",
	}, {
		value: object.Cell{
			object.Cell{object.Number(1), object.Number(2)},
			object.Null},
		env:  object.Cell{},
		want: "((1 2) ())",
	}, {
		value: object.Cell{object.Symbol("foo"), object.Symbol("bar")},
		env: object.Cell{
			object.Cell{object.Symbol("foo"), object.Number(1)},
			object.Cell{
				object.Cell{object.Symbol("bar"), object.Number(2)}}},
		want: "(1 2)",
	}, {
		value:   object.Cell{object.Symbol("foo")},
		env:     object.Cell{},
		wantErr: true,
	}, {
		value: object.Null,
		env:   object.Cell{},
		want:  "()",
	}}

	for _, tt := range tests {
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
	}
}
