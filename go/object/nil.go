package object

type n struct{}

var Nil n = struct{}{}

func (n n) First() Value {
	return Nil
}

func (n n) Rest() Value {
	return Nil
}

func (n n) Type() Type {
	return NIL
}

func (n n) String() string {
	return "()"
}
