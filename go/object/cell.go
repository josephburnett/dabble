package object

import (
	"fmt"
	"strings"
)

type cell [2]Value

func Cell(v1, v2 Value) Value {
	if v1 == nil {
		v1 = Null
	}
	if v2 == nil {
		v2 = Null
	}
	return cell{v1, v2}
}

func (c cell) First() Value {
	return c[0]
}

func (c cell) Rest() Value {
	return c[1]
}

func (c cell) Type() Type {
	return CELL
}

func (c cell) Inspect() string {
	first, rest := c[0], c[1]
	if first == nil {
		first = Null
	}
	if rest == nil {
		rest = Null
	}
	var b strings.Builder
	fmt.Fprintf(&b, "(%v", first.Inspect())
	for rest.Type() == CELL {
		fmt.Fprintf(&b, " %v", rest.First().Inspect())
		rest = rest.Rest()
	}
	if rest.Type() != NULL {
		fmt.Fprintf(&b, " %v", rest.Inspect())
	}
	fmt.Fprintf(&b, ")")
	return b.String()
}
