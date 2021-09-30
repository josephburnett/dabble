package object

import (
	"strconv"
	"testing"
)

func TestQuoted(t *testing.T) {
	tests := []struct {
		quoted Value
		first  string
		rest   string
		string string
	}{{
		quoted: Quoted(Number(1)),
		first:  "1",
		rest:   "()",
		string: "'1",
	}, {
		quoted: Quoted(Symbol("abc")),
		first:  "abc",
		rest:   "()",
		string: "'abc",
	}, {
		quoted: Quoted(Cell(Number(1),
			Cell(Number(2),
				Cell(Number(3), Nil)))),
		first:  "(1 2 3)",
		rest:   "()",
		string: "'(1 2 3)",
	}}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			first := tt.quoted.First().String()
			if first != tt.first {
				t.Errorf("given %v. want first %q. got %q", tt.quoted, tt.first, first)
			}
			rest := tt.quoted.Rest().String()
			if rest != tt.rest {
				t.Errorf("given %v. want rest %q. got %q", tt.quoted, tt.rest, rest)
			}
			got := tt.quoted.String()
			if got != tt.string {
				t.Errorf("given %v. want string %q. got %q", tt.quoted, tt.string, got)
			}
		})
	}
}
