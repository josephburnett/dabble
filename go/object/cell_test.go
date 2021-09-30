package object

import (
	"strconv"
	"testing"
)

func TestCell(t *testing.T) {
	tests := []struct {
		cell   Value
		first  string
		rest   string
		string string
	}{{
		cell:   Cell(Nil, Nil),
		first:  "()",
		rest:   "()",
		string: "(())",
	}, {
		cell:   Cell(Symbol("a"), Symbol("b")),
		first:  "a",
		rest:   "b",
		string: "(a b)",
	}, {
		cell:   Cell(Number(1), Number(2)),
		first:  "1",
		rest:   "2",
		string: "(1 2)",
	}, {
		cell: Cell(
			Cell(Nil, Number(1)),
			Cell(Symbol("a"), Nil)),
		first:  "(() 1)",
		rest:   "(a)",
		string: "((() 1) a)",
	}, {
		cell: Cell(Number(1),
			Cell(Number(2),
				Cell(Number(3), Nil))),
		first:  "1",
		rest:   "(2 3)",
		string: "(1 2 3)",
	}}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			first := tt.cell.First().String()
			if first != tt.first {
				t.Errorf("given %v. want first %q. got %q", tt.cell, tt.first, first)
			}
			rest := tt.cell.Rest().String()
			if rest != tt.rest {
				t.Errorf("given %v. want rest %q. got %q", tt.cell, tt.rest, rest)
			}
			got := tt.cell.String()
			if got != tt.string {
				t.Errorf("given %v. want string %q. got %q", tt.cell, tt.string, got)
			}
		})
	}
}
