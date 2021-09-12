package object

import (
	"strconv"
	"testing"
)

func TestCell(t *testing.T) {
	tests := []struct {
		cell    Value
		first   string
		rest    string
		inspect string
	}{{
		cell:    Cell(Nil, Nil),
		first:   "()",
		rest:    "()",
		inspect: "(())",
	}, {
		cell:    Cell(Symbol("a"), Symbol("b")),
		first:   "a",
		rest:    "b",
		inspect: "(a b)",
	}, {
		cell:    Cell(Number(1), Number(2)),
		first:   "1",
		rest:    "2",
		inspect: "(1 2)",
	}, {
		cell: Cell(
			Cell(Nil, Number(1)),
			Cell(Symbol("a"), Nil)),
		first:   "(() 1)",
		rest:    "(a)",
		inspect: "((() 1) a)",
	}, {
		cell: Cell(Number(1),
			Cell(Number(2),
				Cell(Number(3), Nil))),
		first:   "1",
		rest:    "(2 3)",
		inspect: "(1 2 3)",
	}}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			first := tt.cell.First().Inspect()
			if first != tt.first {
				t.Errorf("given %v. want first %q. got %q", tt.cell, tt.first, first)
			}
			rest := tt.cell.Rest().Inspect()
			if rest != tt.rest {
				t.Errorf("given %v. want rest %q. got %q", tt.cell, tt.rest, rest)
			}
			inspect := tt.cell.Inspect()
			if inspect != tt.inspect {
				t.Errorf("given %v. want inspect %q. got %q", tt.cell, tt.inspect, inspect)
			}
		})
	}
}
