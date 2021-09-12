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
		object: object.Nil,
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
	}, {
		input: `("""")`,
		object: object.Cell(object.Symbol(""),
			object.Cell(object.Symbol(""), nil)),
	}, {
		input:  "(1 . 2)",
		object: object.Cell(object.Number(1), object.Number(2)),
	}, {
		input: "((1 . 2) . (3 . 4))",
		object: object.Cell(
			object.Cell(object.Number(1), object.Number(2)),
			object.Cell(object.Number(3), object.Number(4))),
	}, {
		input:   "(1 . 2 . 3)",
		wantErr: true,
	}, {
		input:   "(. 2)",
		wantErr: true,
	}, {
		input:   "(1 .)",
		wantErr: true,
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
