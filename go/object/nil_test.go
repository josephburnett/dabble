package object

import "testing"

func TestNil(t *testing.T) {
	want := "()"
	first := Nil.First().String()
	if first != want {
		t.Errorf("want %q. got %q", want, first)
	}
	rest := Nil.Rest().String()
	if rest != want {
		t.Errorf("want %q. got %q", want, rest)
	}
	got := Nil.String()
	if got != want {
		t.Errorf("want %q. got %q", want, got)
	}
}
