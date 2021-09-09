package core

import "dabble/object"

var _ object.Function = Car

func Car(env object.Value, args ...object.Value) object.Value {
	if err := argsLenError("car", args, 1); err != nil {
		return err
	}
	return args[0].First()
}
