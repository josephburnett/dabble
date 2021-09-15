package object

type Symbol string

func (s Symbol) isValue() {}

func (s Symbol) Inspect() string {
	if s == "" {
		return Nil.Inspect()
	}
	return string(s)
}
