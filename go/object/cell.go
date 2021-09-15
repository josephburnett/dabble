package object

import (
	"fmt"
	"strings"
)

type Cell [2]Value

func (c Cell) isValue() {}

func (c Cell) Inspect() string {
	first, rest := c[0], c[1]
	if first == nil {
		first = Nil
	}
	if rest == nil {
		rest = Nil
	}
	var b strings.Builder
	fmt.Fprintf(&b, "(%v", first.Inspect())
	for _, ok := rest.(Cell); ok; _, ok = rest.(Cell) {
		fmt.Fprintf(&b, " %v", rest.(Cell).Car().Inspect())
		rest = rest.(Cell).Cdr()
	}
	if rest != Nil {
		fmt.Fprintf(&b, " %v", rest.Inspect())
	}
	fmt.Fprintf(&b, ")")
	return b.String()
}

func (c Cell) Car() Value {
	return c[0]
}

func (c Cell) Cdr() Value {
	return c[1]
}
