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
	return fmt.Sprintf("(%v %v)", c[0].Inspect(), c[1].Inspect())
}
