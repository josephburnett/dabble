package object

type null struct{}

var Null null = struct{}{}

func (n null) First() Value {
	return Null
}

func (n null) Rest() Value {
	return Null
}

func (n null) Inspect() string {
	return "()"
}
