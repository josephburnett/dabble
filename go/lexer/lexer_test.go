package lexer

import (
	"dabble/token"
	"testing"
)

// Based on Monkey lexer_test.go

func TestNextToken(t *testing.T) {
	input := `
()
foo
123
(foo)
(foo bar)
"foo"
"foo bar"
("foo bar" "baz bam")
("foo""bar")
((foo))
(())
"+"
"("
" "
(let (a 1) (+ 1 a))
""
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.SYMBOL, "foo"},
		{token.NUMBER, "123"},
		{token.LPAREN, "("},
		{token.SYMBOL, "foo"},
		{token.RPAREN, ")"},
		{token.LPAREN, "("},
		{token.SYMBOL, "foo"},
		{token.SYMBOL, "bar"},
		{token.RPAREN, ")"},
		{token.SYMBOL, "foo"},
		{token.SYMBOL, "foo bar"},
		{token.LPAREN, "("},
		{token.SYMBOL, "foo bar"},
		{token.SYMBOL, "baz bam"},
		{token.RPAREN, ")"},
		{token.LPAREN, "("},
		{token.SYMBOL, "foo"},
		{token.SYMBOL, "bar"},
		{token.RPAREN, ")"},
		{token.LPAREN, "("},
		{token.LPAREN, "("},
		{token.SYMBOL, "foo"},
		{token.RPAREN, ")"},
		{token.RPAREN, ")"},
		{token.LPAREN, "("},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.RPAREN, ")"},
		{token.SYMBOL, "+"},
		{token.SYMBOL, "("},
		{token.SYMBOL, " "},
		{token.LPAREN, "("},
		{token.SYMBOL, "let"},
		{token.LPAREN, "("},
		{token.SYMBOL, "a"},
		{token.NUMBER, "1"},
		{token.RPAREN, ")"},
		{token.LPAREN, "("},
		{token.SYMBOL, "+"},
		{token.NUMBER, "1"},
		{token.SYMBOL, "a"},
		{token.RPAREN, ")"},
		{token.RPAREN, ")"},
		{token.SYMBOL, ""},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
