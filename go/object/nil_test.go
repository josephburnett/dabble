package object

import "testing"

func TestNil(t *testing.T) {
	want := "()"
	first := Nil.First().Inspect()
	if first != want {
		t.Errorf("want %q. got %q", want, first)
	}
	rest := Nil.Rest().Inspect()
	if rest != want {
		t.Errorf("want %q. got %q", want, rest)
	}
	inspect := Nil.Inspect()
	if inspect != want {
		t.Errorf("want %q. got %q", want, inspect)
	}
}
