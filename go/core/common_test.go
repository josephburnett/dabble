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
	env     *object.Binding
	want    string
	wantErr bool
}

func testCore(t *testing.T, tests []coreTest) {
	t.Helper()
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			value, err := p.ParseProgram()
			if err != nil {
				t.Fatalf(err.Error())
			}
			got := eval.Eval(tt.env, value)
			if tt.wantErr {
				if _, ok := got.(object.Error); !ok {
					t.Errorf("given value %v env %+v. want err. got %v", value.Inspect(), tt.env, got.Inspect())
				}
			} else {
				if got.Inspect() != tt.want {
					t.Errorf("given value %v env %+v. want %v. got %v", value.Inspect(), tt.env, tt.want, got.Inspect())
				}
			}
		})
	}
}
