package core

import (
	"dabble/eval"
	"dabble/object"
)

var Env *eval.Frame

func init() {
	for name, fn := range map[string]eval.Function{
		"atom":    Atom,
		"car":     Car,
		"cdr":     Cdr,
		"cons":    Cons,
		"eq":      Eq,
		"if":      If,
		"label":   Label,
		"lambda":  Lambda,
		"macro":   Macro,
		"quote":   Quote,
		"unquote": Unquote,
		"recur":   Recur,
	} {
		Env = Env.Bind(object.Symbol(name), fn)
	}
}
