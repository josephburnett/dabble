package token

// Based on Monkey token.go.

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	SYMBOL = "SYMBOL"
	NUMBER = "NUMBER"

	LPAREN  = "("
	RPAREN  = ")"
	DOT     = "."
	QUOTE   = "'"
	UNQUOTE = "`"
)

type Token struct {
	Type    TokenType
	Literal string
}
