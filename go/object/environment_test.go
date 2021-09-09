package object

import (
	"testing"
)

func TestResolve(t *testing.T) {
	tests := []struct {
		env     Environment
		symbol  Symbol
		want    string
		wantErr bool
	}{{
		env:     nil,
		symbol:  Symbol("foo"),
		wantErr: true,
	}, {
		env:     []Binding{},
		symbol:  Symbol("foo"),
		wantErr: true,
	}, {
		env: []Binding{
			{"foo", Number(1)},
		},
		symbol: Symbol("foo"),
		want:   "1",
	}, {
		env: []Binding{
			{"bar", Number(2)},
			{"foo", Number(1)},
		},
		symbol: Symbol("foo"),
		want:   "1",
	}, {
		env: []Binding{
			{"foo", Number(1)},
			{"foo", Number(2)},
		},
		symbol: Symbol("foo"),
		want:   "1",
	}}

	for _, tt := range tests {
		got := tt.env.Resolve(tt.env, tt.symbol)
		if tt.wantErr {
			if _, ok := got.(Error); !ok {
				t.Errorf("wanted error. got %T (%q)", got, got.Inspect())
			}
		} else {
			if got.Inspect() != tt.want {
				t.Errorf("wanted %q. got %q", tt.want, got.Inspect())
			}
		}
	}
}
