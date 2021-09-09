package core

import (
	"dabble/object"
	"fmt"
)

func argsLenError(name string, args []object.Value, want int) object.Value {
	if len(args) != want {
		return object.Error(fmt.Sprintf("%v wants 1 arg. got %v", name, len(args)))
	}
	return nil
}
