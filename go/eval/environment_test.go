package eval

import (
	"dabble/object"
	"testing"
)

func TestResolve(t *testing.T) {
	tests := []struct {
		env     object.Value
		symbol  object.Value
		want    string
		wantErr bool
	}{{
		env:     object.Null,
		symbol:  object.Symbol("foo"),
		wantErr: true,
	}, {
		env:     object.Cell(object.Null, object.Null),
		symbol:  object.Symbol("foo"),
		wantErr: true,
	}, {
		env: object.Cell(
			object.Cell(object.Null, object.Null),
			object.Null),
		symbol:  object.Symbol("foo"),
		wantErr: true,
	}, {
		env: object.Cell(
			object.Cell(object.Symbol("foo"), object.Number(1)),
			object.Null),
		symbol: object.Symbol("foo"),
		want:   "1",
	}, {
		env: object.Cell(
			object.Cell(object.Symbol("bar"), object.Number(2)),
			object.Cell(
				object.Cell(object.Symbol("foo"), object.Number(1)),
				object.Null)),
		symbol: object.Symbol("foo"),
		want:   "1",
	}, {
		env: object.Cell(
			object.Cell(object.Symbol("foo"), object.Number(1)),
			object.Cell(
				object.Cell(object.Symbol("foo"), object.Number(2)),
				object.Null)),
		symbol: object.Symbol("foo"),
		want:   "1",
	}}

	for _, tt := range tests {
		got := resolve(tt.env, tt.symbol)
		if tt.wantErr {
			if _, ok := got.(object.Error); !ok {
				t.Errorf("wanted error. got %T (%q)", got, got.Inspect())
			}
		} else {
			if got.Inspect() != tt.want {
				t.Errorf("wanted %q. got %q", tt.want, got.Inspect())
			}
		}
	}
}
