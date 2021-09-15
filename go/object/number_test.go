package object

import "testing"

func TestNumber(t *testing.T) {
	tests := []struct {
		number  Number
		inspect string
	}{{
		number:  0,
		inspect: "0",
	}, {
		number:  1,
		inspect: "1",
	}, {
		number:  2,
		inspect: "2",
	}, {
		number:  3,
		inspect: "3",
	}}

	for _, tt := range tests {
		inspect := tt.number.Inspect()
		if inspect != tt.inspect {
			t.Errorf("given %v. want inspect %q. got %q", tt.number, tt.inspect, inspect)
		}
	}
}
