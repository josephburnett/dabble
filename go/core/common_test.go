package core

import (
	"dabble/eval"
	"dabble/lexer"
	"dabble/object"
	"dabble/parser"
	"strconv"
	"testing"
)

type coreTest struct {
	input   string
	want    string
	wantErr bool
}

func testCore(t *testing.T, env *eval.Frame, tests []coreTest) {
	t.Helper()
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			value, err := p.ParseProgram()
			if err != nil {
				t.Fatalf(err.Error())
			}
			eval.BeginTrace()
			got := eval.Eval(env, value)
			trace := eval.EndTrace()
			var printTrace bool
			if tt.wantErr {
				if _, ok := got.(object.Error); !ok {
					t.Errorf("given value %v env %+v. want err. got %v", value.String(), env, got.String())
					printTrace = true
				}
			} else {
				if got.String() != tt.want {
					t.Errorf("given value %v env %+v. want %v. got %v", value.String(), env, tt.want, got.String())
					printTrace = true
				}
			}
			if printTrace {
				t.Log("TRACE:\n" + trace)
			}
		})
	}
}
