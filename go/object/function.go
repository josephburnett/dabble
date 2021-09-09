package object

type Function func(env Environment, args ...Value) Value

func (f Function) First() Value {
	return Null
}

func (f Function) Rest() Value {
	return Null
}

func (f Function) Type() Type {
	return FUNCTION
}

func (f Function) Inspect() string {
	return "<function>"
}
