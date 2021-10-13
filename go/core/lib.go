package core

import (
	"dabble/eval"
	"dabble/lexer"
	"dabble/object"
	"dabble/parser"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var Env *eval.Frame

func init() {

	Env = Env.Bind("t", object.Symbol("t"))

	var builtins *eval.Frame
	for name, fn := range map[string]func(*eval.Frame, ...object.Value) object.Value{
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
		"error":   Error,
		"apply":   Apply,
	} {
		function := &eval.Function{
			Name: name,
			Fn:   fn,
		}
		builtins = builtins.Bind(object.Symbol(name), function)
	}
	Env = Env.BindAll(builtins)

	// TODO: pack core and library code and read from in-binary.
	var lib *eval.Frame
	err := filepath.Walk("../../src", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".lisp") {
			return nil
		}
		symbol := info.Name()[:len(info.Name())-len(".lisp")]
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		l := lexer.New(string(bytes))
		p := parser.New(l)
		program, err := p.ParseProgram()
		if err != nil {
			return err
		}
		value := eval.Eval(Env, program)
		if value.Type() == object.ERROR {
			panic(fmt.Sprintf("%q %v", info.Name(), string(value.(object.Error))))
		}
		lib = lib.Bind(object.Symbol(symbol), value)
		return nil
	})
	if err != nil {
		panic(err)
	}
	Env = Env.BindAll(lib)
}
