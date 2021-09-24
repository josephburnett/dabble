package object

type quoted struct {
	value Value
}

func Quoted(value Value) Value {
	if value == nil {
		value = Nil
	}

	return quoted{value}
}

func (q quoted) First() Value {
	return q.value.First()
}

func (q quoted) Rest() Value {
	return q.value.Rest()
}

func (q quoted) Type() Type {
	return QUOTED
}

func (q quoted) Inspect() string {
	return "'" + q.value.Inspect()
}
