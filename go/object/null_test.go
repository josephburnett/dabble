package object

import "testing"

func TestNull(t *testing.T) {
	want := "()"
	first := Null.First().Inspect()
	if first != want {
		t.Errorf("want %q. got %q", want, first)
	}
	rest := Null.Rest().Inspect()
	if rest != want {
		t.Errorf("want %q. got %q", want, rest)
	}
	inspect := Null.Inspect()
	if inspect != want {
		t.Errorf("want %q. got %q", want, inspect)
	}
}
