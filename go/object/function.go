package object

type Function func(env *Binding, args ...Value) Value

func (f Function) isValue() {}

func (f Function) Inspect() string {
	return "<function>"
}
