package object

type Value interface {
	First() Value
	Rest() Value
	Inspect() string
}
