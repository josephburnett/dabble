package object

import "fmt"

type Number uint64

func (n Number) First() Value {
	return n & 1
}

func (n Number) Rest() Value {
	return n >> 1
}

func (n Number) Inspect() string {
	return fmt.Sprintf("%v", n)
}
