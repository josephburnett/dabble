package eval

import "dabble/object"

type Function func(env *Frame, args ...object.Value) object.Value

func (f Function) First() object.Value {
	return object.Nil
}

func (f Function) Rest() object.Value {
	return object.Nil
}

func (f Function) Type() object.Type {
	return object.FUNCTION
}

func (f Function) String() string {
	return "<function>"
}
