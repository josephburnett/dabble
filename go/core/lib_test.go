package core

import (
	"dabble/eval"
	"dabble/lexer"
	"dabble/object"
	"dabble/parser"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLib(t *testing.T) {
	err := filepath.Walk("../../tst", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".lisp") {
			return nil
		}
		symbol := info.Name()[:len(info.Name())-len(".lisp")]
		t.Run(symbol, func(t *testing.T) {
			if err != nil {
				t.Errorf(err.Error())
				return
			}
			bytes, err := ioutil.ReadFile(path)
			if err != nil {
				t.Errorf(err.Error())
				return
			}
			l := lexer.New(string(bytes))
			p := parser.New(l)
			program, err := p.ParseProgram()
			if err != nil {
				t.Errorf(err.Error())
				return
			}
			eval.BeginTrace()
			value := eval.Eval(Env, program)
			trace := eval.EndTrace()
			if value.Type() == object.ERROR {
				t.Errorf("%v", value)
				t.Errorf("TRACE:\n%v", trace)
				return
			}
		})
		return nil
	})
	if err != nil {
		t.Errorf(err.Error())
	}
}
