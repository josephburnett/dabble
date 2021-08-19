package object

import "testing"

func TestSymbol(t *testing.T) {
	tests := []struct {
		symbol  Symbol
		first   string
		rest    string
		inspect string
	}{{
		symbol:  "",
		first:   "()",
		rest:    "()",
		inspect: "()",
	}, {
		symbol:  "a",
		first:   "a",
		rest:    "()",
		inspect: "a",
	}, {
		symbol:  "ab",
		first:   "a",
		rest:    "b",
		inspect: "ab",
	}, {
		symbol:  "abc",
		first:   "a",
		rest:    "bc",
		inspect: "abc",
	}}

	for _, tt := range tests {
		first := tt.symbol.First().Inspect()
		if first != tt.first {
			t.Errorf("given %q. want first %q. got %q", tt.symbol, tt.first, first)
		}
		rest := tt.symbol.Rest().Inspect()
		if rest != tt.rest {
			t.Errorf("given %q. want rest %q. got %q", tt.symbol, tt.rest, rest)
		}
		inspect := tt.symbol.Inspect()
		if inspect != tt.inspect {
			t.Errorf("given %q. want inspect %q. got %q", tt.symbol, tt.inspect, inspect)
		}
	}
}
