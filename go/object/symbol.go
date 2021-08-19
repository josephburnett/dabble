package object

type Symbol string

func (s Symbol) First() Value {
	if s == "" {
		return Null
	}
	return s[0:1]
}

func (s Symbol) Rest() Value {
	if len(s) < 2 {
		return Null
	}
	return s[1:]
}

func (s Symbol) Inspect() string {
	if s == "" {
		return Null.Inspect()
	}
	return string(s)
}
