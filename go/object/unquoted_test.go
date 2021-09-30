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
		string   string
	}{{
		unquoted: Unquoted(Number(1)),
		first:    "1",
		rest:     "()",
		string:   "`1",
	}, {
		unquoted: Unquoted(Symbol("abc")),
		first:    "abc",
		rest:     "()",
		string:   "`abc",
	}, {
		unquoted: Unquoted(Cell(Number(1),
			Cell(Number(2),
				Cell(Number(3), Nil)))),
		first:  "(1 2 3)",
		rest:   "()",
		string: "`(1 2 3)",
	}}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			first := tt.unquoted.First().String()
			if first != tt.first {
				t.Errorf("given %v. want first %q. got %q", tt.unquoted, tt.first, first)
			}
			rest := tt.unquoted.Rest().String()
			if rest != tt.rest {
				t.Errorf("given %v. want rest %q. got %q", tt.unquoted, tt.rest, rest)
			}
			got := tt.unquoted.String()
			if got != tt.string {
				t.Errorf("given %v. want string %q. got %q", tt.unquoted, tt.string, got)
			}
		})
	}
}
