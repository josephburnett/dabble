package eval

import (
	"dabble/object"
	"strconv"
	"testing"
)

func TestResolve(t *testing.T) {
	tests := []struct {
		env     *Frame
		symbol  object.Symbol
		want    string
		wantErr bool
	}{{
		env:     nil,
		symbol:  object.Symbol("foo"),
		wantErr: true,
	}, {
		env:    (*Frame)(nil).Bind("foo", object.Number(1)),
		symbol: object.Symbol("foo"),
		want:   "1",
	}, {
		env: (*Frame)(nil).Bind("bar", object.Number(2)).
			Bind("foo", object.Number(1)),
		symbol: object.Symbol("foo"),
		want:   "1",
	}, {
		env: (*Frame)(nil).Bind("foo", object.Number(2)).
			Bind("foo", object.Number(1)),
		symbol: object.Symbol("foo"),
		want:   "1",
	}}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := tt.env.Resolve(tt.symbol)
			if tt.wantErr {
				if _, ok := got.(object.Error); !ok {
					t.Errorf("wanted error. got %T (%q)", got, got.String())
				}
			} else {
				if got.String() != tt.want {
					t.Errorf("wanted %q. got %q", tt.want, got.String())
				}
			}
		})
	}
}
