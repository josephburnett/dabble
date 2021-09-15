package object

import "testing"

func TestSymbol(t *testing.T) {
	tests := []struct {
		symbol  Symbol
		inspect string
	}{{
		symbol:  "",
		inspect: "()",
	}, {
		symbol:  "a",
		inspect: "a",
	}, {
		symbol:  "ab",
		inspect: "ab",
	}, {
		symbol:  "abc",
		inspect: "abc",
	}}

	for _, tt := range tests {
		inspect := tt.symbol.Inspect()
		if inspect != tt.inspect {
			t.Errorf("given %q. want inspect %q. got %q", tt.symbol, tt.inspect, inspect)
		}
	}
}
