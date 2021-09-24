package object

type unquoted struct {
	value Value
}

func Unquoted(value Value) Value {
	if value == nil {
		value = Nil
	}
	return unquoted{value}
}

func (u unquoted) First() Value {
	return u.value.First()
}

func (u unquoted) Rest() Value {
	return u.value.Rest()
}

func (u unquoted) Type() Type {
	return UNQUOTED
}

func (u unquoted) Inspect() string {
	return "`" + u.value.Inspect()
}
