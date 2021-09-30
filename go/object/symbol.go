package object

type Symbol string

func (s Symbol) First() Value {
	if s == "" {
		return Nil
	}
	return s[0:1]
}

func (s Symbol) Rest() Value {
	if len(s) < 2 {
		return Nil
	}
	return s[1:]
}

func (s Symbol) Type() Type {
	return SYMBOL
}

func (s Symbol) String() string {
	if s == "" {
		return Nil.String()
	}
	return string(s)
}
