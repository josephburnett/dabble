package object

type Value interface {
	isValue()
	Inspect() string
}
