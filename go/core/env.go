package core

import "dabble/object"

var Env *object.Binding

func init() {
	for name, fn := range map[string]object.Function{
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
	} {
		Env = &object.Binding{
			Symbol: object.Symbol(name),
			Value:  fn,
			Next:   Env,
		}
	}
}
