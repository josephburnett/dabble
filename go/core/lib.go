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
)

var Env *eval.Frame

func init() {

	var builtins *eval.Frame
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
		builtins = builtins.Bind(object.Symbol(name), fn)
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
		value := eval.Eval(builtins, program)
		if value.Type() == object.ERROR {
			panic(fmt.Sprintf("%q %v", info.Name(), string(value.(object.Error))))
		}
		lib = lib.Bind(object.Symbol(info.Name()), value)
		return nil
	})
	if err != nil {
		panic(err)
	}
	Env = Env.BindAll(lib)

	// lib_test.go will walk dabble/tst/*
	// and eval each .lisp file, asserting t.
}
