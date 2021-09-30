package object

import "testing"

func TestSymbol(t *testing.T) {
	tests := []struct {
		symbol Symbol
		first  string
		rest   string
		string string
	}{{
		symbol: "",
		first:  "()",
		rest:   "()",
		string: "()",
	}, {
		symbol: "a",
		first:  "a",
		rest:   "()",
		string: "a",
	}, {
		symbol: "ab",
		first:  "a",
		rest:   "b",
		string: "ab",
	}, {
		symbol: "abc",
		first:  "a",
		rest:   "bc",
		string: "abc",
	}}

	for _, tt := range tests {
		first := tt.symbol.First().String()
		if first != tt.first {
			t.Errorf("given %q. want first %q. got %q", tt.symbol, tt.first, first)
		}
		rest := tt.symbol.Rest().String()
		if rest != tt.rest {
			t.Errorf("given %q. want rest %q. got %q", tt.symbol, tt.rest, rest)
		}
		got := tt.symbol.String()
		if got != tt.string {
			t.Errorf("given %q. want string %q. got %q", tt.symbol, tt.string, got)
		}
	}
}
