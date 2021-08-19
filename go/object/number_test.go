package object

import "testing"

func TestNumber(t *testing.T) {
	tests := []struct {
		number  Number
		first   Number
		rest    Number
		inspect string
	}{{
		number:  0,
		first:   0,
		rest:    0,
		inspect: "0",
	}, {
		number:  1,
		first:   1,
		rest:    0,
		inspect: "1",
	}, {
		number:  2,
		first:   0,
		rest:    1,
		inspect: "2",
	}, {
		number:  3,
		first:   1,
		rest:    1,
		inspect: "3",
	}}

	for _, tt := range tests {
		first := tt.number.First().(Number)
		if first != tt.first {
			t.Errorf("given %v. want first %v. got %v", tt.number, tt.first, first)
		}
		rest := tt.number.Rest().(Number)
		if rest != tt.rest {
			t.Errorf("given %v. want rest %v. got %v", tt.number, tt.rest, rest)
		}
		inspect := tt.number.Inspect()
		if inspect != tt.inspect {
			t.Errorf("given %v. want inspect %q. got %q", tt.number, tt.inspect, inspect)
		}
	}
}
