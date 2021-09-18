package object

type Closure func(args ...Value) Value

func (c Closure) Type() Type {
	return CLOSURE
}

func (c Closure) Inspect() string {
	return "<closure>"
}

func (c Closure) First() Value {
	return Nil
}

func (c Closure) Rest() Value {
	return Nil
}
