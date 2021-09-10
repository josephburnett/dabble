package eval

import (
	"dabble/lexer"
	"dabble/object"
	"dabble/parser"
	"strconv"
	"testing"
)

func TestEval(t *testing.T) {

	passFunction := object.Function(func(env *object.Binding, args ...object.Value) object.Value {
		return object.Symbol("pass")
	})

	tests := []struct {
		input   string
		env     *object.Binding
		want    string
		wantErr bool
	}{{
		input: "1",
		env:   nil,
		want:  "1",
	}, {
		input:   "a",
		env:     nil,
		wantErr: true,
	}, {
		input:   "(foo)",
		env:     nil,
		wantErr: true,
	}, {
		input: "()",
		env:   nil,
		want:  "()",
	}, {
		input: "(bar)",
		env:   &object.Binding{"bar", passFunction, nil},
		want:  "pass",
	}}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			value, err := p.ParseProgram()
			if err != nil {
				t.Fatalf(err.Error())
			}
			got := Eval(tt.env, value)
			if tt.wantErr {
				if _, ok := got.(object.Error); !ok {
					t.Errorf("given value %v env %v. want err. got %v", value, tt.env, got)
				}
			} else {
				if got.Inspect() != tt.want {
					t.Errorf("given value %v env %v. want %v. got %v", value, tt.env, tt.want, got.Inspect())
				}
			}
		})
	}
}
