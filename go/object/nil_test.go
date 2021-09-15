package object

import "testing"

func TestNil(t *testing.T) {
	want := "()"
	inspect := Nil.Inspect()
	if inspect != want {
		t.Errorf("want %q. got %q", want, inspect)
	}
}
