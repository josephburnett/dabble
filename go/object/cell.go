package object

import (
	"fmt"
	"strings"
)

type cell [2]Value

func Cell(v1, v2 Value) Value {
	if v1 == nil {
		v1 = Nil
	}
	if v2 == nil {
		v2 = Nil
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

func (c cell) String() string {
	first, rest := c[0], c[1]
	if first == nil {
		first = Nil
	}
	if rest == nil {
		rest = Nil
	}
	var b strings.Builder
	fmt.Fprintf(&b, "(%v", first.String())
	for rest.Type() == CELL {
		fmt.Fprintf(&b, " %v", rest.First().String())
		rest = rest.Rest()
	}
	if rest.Type() != NIL {
		fmt.Fprintf(&b, " %v", rest.String())
	}
	fmt.Fprintf(&b, ")")
	return b.String()
}
