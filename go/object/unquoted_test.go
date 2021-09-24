package object

import (
	"strconv"
	"testing"
)

func TestUnquoted(t *testing.T) {
	tests := []struct {
		unquoted Value
		first    string
		rest     string
		inspect  string
	}{{
		unquoted: Unquoted(Number(1)),
		first:    "1",
		rest:     "0",
		inspect:  "`1",
	}, {
		unquoted: Unquoted(Symbol("abc")),
		first:    "a",
		rest:     "bc",
		inspect:  "`abc",
	}, {
		unquoted: Unquoted(Cell(Number(1),
			Cell(Number(2),
				Cell(Number(3), Nil)))),
		first:   "1",
		rest:    "(2 3)",
		inspect: "`(1 2 3)",
	}}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			first := tt.unquoted.First().Inspect()
			if first != tt.first {
				t.Errorf("given %v. want first %q. got %q", tt.unquoted, tt.first, first)
			}
			rest := tt.unquoted.Rest().Inspect()
			if rest != tt.rest {
				t.Errorf("given %v. want rest %q. got %q", tt.unquoted, tt.rest, rest)
			}
			inspect := tt.unquoted.Inspect()
			if inspect != tt.inspect {
				t.Errorf("given %v. want inspect %q. got %q", tt.unquoted, tt.inspect, inspect)
			}
		})
	}
}
