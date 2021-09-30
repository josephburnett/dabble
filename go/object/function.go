package object

type Function func(env *Binding, args ...Value) Value

func (f Function) First() Value {
	return Nil
}

func (f Function) Rest() Value {
	return Nil
}

func (f Function) Type() Type {
	return FUNCTION
}

func (f Function) String() string {
	return "<function>"
}
