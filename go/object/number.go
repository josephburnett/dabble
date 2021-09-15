package object

import "fmt"

type Number int64

func (n Number) isValue() {}

func (n Number) Inspect() string {
	return fmt.Sprintf("%v", n)
}
