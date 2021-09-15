package object

import (
	"fmt"
)

type Error string

func (e Error) isValue() {}

func (e Error) Inspect() string {
	return fmt.Sprintf("<ERROR: %v>",string(e))
}
