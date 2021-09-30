package object

type Error string

func (e Error) First() Value {
	return e
}

func (e Error) Rest() Value {
	return e
}

func (e Error) Type() Type {
	return ERROR
}

func (e Error) String() string {
	return string("<error: " + e + ">")
}
