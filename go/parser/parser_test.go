package parser

import (
	"dabble/lexer"
	"dabble/object"
	"strconv"
	"testing"
)

func TestParser(t *testing.T) {
	tests := []struct {
		input   string
		object  object.Value
		wantErr bool
	}{{
		input:  "foo",
		object: object.Symbol("foo"),
	}, {
		input:  "1234",
		object: object.Number(1234),
	}, {
		input:  "()",
		object: object.Null,
	}, {
		input:   "(",
		wantErr: true,
	}, {
		input:   ")",
		wantErr: true,
	}, {
		input:   "())",
		wantErr: true,
	}, {
		input:  "(foo)",
		object: object.Cell(object.Symbol("foo"), nil),
	}, {
		input: "(foo bar)",
		object: object.Cell(object.Symbol("foo"),
			object.Cell(object.Symbol("bar"), nil)),
	}, {
		input: "(foo (bar) baz)",
		object: object.Cell(object.Symbol("foo"),
			object.Cell(object.Cell(object.Symbol("bar"), nil),
				object.Cell(object.Symbol("baz"), nil))),
	}}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			v, err := p.ParseProgram()
			if tt.wantErr {
				if err == nil {
					t.Errorf("wanted error")
				}
				if v != nil {
					t.Errorf("unwanted: %v", v)
				}
			} else {
				if err != nil {
					t.Errorf("unwanted: %v", err)
				}
				if v != tt.object {
					t.Errorf("want %v. got %v", tt.object, v)
				}
			}
		})
	}
}
