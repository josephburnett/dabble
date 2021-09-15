package object

import "testing"

func TestError(t *testing.T) {
	err := Error("something went wrong")
	inspect := err.Inspect()
	if inspect != "something went wrong" {
		t.Errorf("given %q. want inspect %q. got %q", err, err, inspect)
	}
}
