package core

import "dabble/object"

func Car(args ...object.Value) object.Value {
	if err := argsLenError("car", args, 1); err != nil {
		return err
	}
	return args[0].First()
}
