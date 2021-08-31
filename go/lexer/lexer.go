package lexer

import (
	"dabble/token"
	"unicode"
	"unicode/utf8"
)

// Based on Monkey lexer.go.

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '\'':
		tok.Type = token.SYMBOL
		l.readChar()
		tok.Literal = l.readQuotedSymbol()
	case 0:
		tok.Type = token.EOF
	default:
		if isDigit(l.ch) {
			tok.Type = token.NUMBER
			tok.Literal = l.readNumber()
		} else if isSymbolChar(l.ch) {
			tok.Type = token.SYMBOL
			tok.Literal = l.readSymbol()
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for isSpace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readQuotedSymbol() string {
	position := l.position
	for isQuotedSymbolChar(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readSymbol() string {
	position := l.position
	for isSymbolChar(l.peekChar()) {
		l.readChar()
	}
	return l.input[position:l.readPosition]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.peekChar()) {
		l.readChar()
	}
	return l.input[position:l.readPosition]
}

func isSymbolChar(ch byte) bool {
	if isParenChar(ch) || isSpace(ch) {
		return false
	}
	r, _ := utf8.DecodeRune([]byte{ch})
	return unicode.IsPrint(r)
}

func isParenChar(ch byte) bool {
	return ch == '(' || ch == ')'
}

func isQuotedSymbolChar(ch byte) bool {
	if ch == '\'' {
		return false
	}
	r, _ := utf8.DecodeRune([]byte{ch})
	return unicode.IsPrint(r)
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isSpace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
