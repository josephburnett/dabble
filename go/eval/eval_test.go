package eval

import (
	"dabble/lexer"
	"dabble/object"
	"dabble/parser"
	"fmt"
	"strconv"
	"testing"
)

func TestEval(t *testing.T) {

	passFunction := &Function{Fn: func(env *Frame, args ...object.Value) object.Value {
		return object.Symbol("pass")
	}}

	identityFunction := &Function{Fn: func(env *Frame, args ...object.Value) object.Value {
		if len(args) != 1 {
			return object.Error(fmt.Sprintf("wrong args: %v", args))
		}
		return args[0]
	}}

	addingFunction := &Function{Fn: func(env *Frame, args ...object.Value) object.Value {
		if len(args) != 1 {
			return object.Error(fmt.Sprintf("wrong args: %v", args))
		}
		value := Eval(env, args[0])
		if value.Type() == object.ERROR {
			return value
		}
		if value.Type() != object.NUMBER {
			return object.Error(fmt.Sprintf("wrong type: %v", value))
		}
		return object.Number(value.(object.Number) + 1)
	}}

	tests := []struct {
		input   string
		env     *Frame
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
		env:   NilFrame.Bind("bar", passFunction),
		want:  "pass",
	}, {
		input: "(baz 123)",
		env:   NilFrame.Bind("baz", identityFunction),
		want:  "123",
	}, {
		input: "(+ (+ 1))",
		env:   NilFrame.Bind("+", addingFunction),
		want:  "3",
	}, {
		input: "'a",
		env:   NilFrame.Bind("a", object.Number(1)),
		want:  "a",
	}, {
		input: "'(1 `b 3)",
		env:   NilFrame.Bind("b", object.Number(2)),
		want:  "(1 2 3)",
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
				if got.String() != tt.want {
					t.Errorf("given value %v env %v. want %v. got %v", value, tt.env, tt.want, got.String())
				}
			}
		})
	}
}
