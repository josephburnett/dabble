package object

import "fmt"

type Cell [2]Value

func (c Cell) First() Value {
	return c[0]
}

func (c Cell) Rest() Value {
	return c[1]
}

func (c Cell) Inspect() string {
	first, rest := c[0], c[1]
	if first == nil {
		first = Null
	}
	if rest == nil {
		rest = Null
	}
	return fmt.Sprintf("(%v %v)", first.Inspect(), rest.Inspect())
}
