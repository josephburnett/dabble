package parser

import (
	"dabble/lexer"
	"dabble/object"
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	tests := []struct {
		input  string
		object object.Value
	}{{
		input:  "foo",
		object: object.Symbol("foo"),
	}, {
		input:  "1234",
		object: object.Number(1234),
	}, {
		input:  "()",
		object: object.Cell(nil, nil),
	}, {
		input:  "(foo)",
		object: object.Cell(object.Symbol("foo"), nil),
	}}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("[%v]", i), func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			v := p.ParseProgram()
			if v != tt.object {
				t.Errorf("want %v. got %v", tt.object, v)
			}
		})
	}
}
