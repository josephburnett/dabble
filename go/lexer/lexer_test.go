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
