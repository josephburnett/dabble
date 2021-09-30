package object

import "fmt"

type Number int64

func (n Number) First() Value {
	return n & 1
}

func (n Number) Rest() Value {
	return n >> 1
}

func (n Number) Type() Type {
	return NUMBER
}

func (n Number) String() string {
	return fmt.Sprintf("%v", int64(n))
}
