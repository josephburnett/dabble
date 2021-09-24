package object

import (
	"strconv"
	"testing"
)

func TestQuoted(t *testing.T) {
	tests := []struct {
		quoted  Value
		first   string
		rest    string
		inspect string
	}{{
		quoted:  Quoted(Number(1)),
		first:   "1",
		rest:    "0",
		inspect: "'1",
	}, {
		quoted:  Quoted(Symbol("abc")),
		first:   "a",
		rest:    "bc",
		inspect: "'abc",
	}, {
		quoted: Quoted(Cell(Number(1),
			Cell(Number(2),
				Cell(Number(3), Nil)))),
		first:   "1",
		rest:    "(2 3)",
		inspect: "'(1 2 3)",
	}}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			first := tt.quoted.First().Inspect()
			if first != tt.first {
				t.Errorf("given %v. want first %q. got %q", tt.quoted, tt.first, first)
			}
			rest := tt.quoted.Rest().Inspect()
			if rest != tt.rest {
				t.Errorf("given %v. want rest %q. got %q", tt.quoted, tt.rest, rest)
			}
			inspect := tt.quoted.Inspect()
			if inspect != tt.inspect {
				t.Errorf("given %v. want inspect %q. got %q", tt.quoted, tt.inspect, inspect)
			}
		})
	}
}
