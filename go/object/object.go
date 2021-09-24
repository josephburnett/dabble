package object

type Type string

const (
	SYMBOL Type = "SYMBOL"
	NUMBER      = "NUMBER"
	CELL        = "CELL"
	NIL         = "NIL"

	QUOTED   = "QUOTED"
	UNQUOTED = "UNQUOTED"

	FUNCTION = "FUNCTION"
	CLOSURE  = "CLOSURE"
	ERROR    = "ERROR"
)

type Value interface {
	First() Value
	Rest() Value
	Type() Type
	Inspect() string
}
