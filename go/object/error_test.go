package object

import "testing"

func TestError(t *testing.T) {
	err := Error("something went wrong")
	first := err.First()
	if first != err {
		t.Errorf("given %q. want first %q. got %q", err, err, first)
	}
	rest := err.Rest()
	if rest != err {
		t.Errorf("given %q. want rest %q. got %q", err, err, rest)
	}
	inspect := err.Rest()
	if inspect != err {
		t.Errorf("given %q. want inspect %q. got %q", err, err, inspect)
	}
}
