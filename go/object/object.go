package object

type Type string

const (
	SYMBOL   Type = "SYMBOL"
	NUMBER        = "NUMBER"
	CELL          = "CELL"
	FUNCTION      = "FUNCTION"
	NIL           = "NIL"
	ERROR         = "ERROR"
)

type Value interface {
	First() Value
	Rest() Value
	Type() Type
	Inspect() string
}
