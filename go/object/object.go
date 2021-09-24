package object

type Type string

const (
	SYMBOL   Type = "SYMBOL"
	NUMBER        = "NUMBER"
	CELL          = "CELL"
	QUOTED        = "QUOTED"
	UNQUOTED      = "UNQUOTED"
	FUNCTION      = "FUNCTION"
	CLOSURE       = "CLOSURE"
	NIL           = "NIL"
	ERROR         = "ERROR"
)

type Value interface {
	First() Value
	Rest() Value
	Type() Type
	Inspect() string
}
