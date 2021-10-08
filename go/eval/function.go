package eval

import (
	"dabble/object"
	"fmt"
)

type Function struct {
	Name string
	Fn   func(env *Frame, args ...object.Value) object.Value
}

func (f *Function) First() object.Value {
	return object.Nil
}

func (f *Function) Rest() object.Value {
	return object.Nil
}

func (f *Function) Type() object.Type {
	return object.FUNCTION
}

func (f *Function) String() string {
	return fmt.Sprintf("<function %q>", f.Name)
}
