package parser

// Based on Monkey parser.go.

import (
	"dabble/lexer"
	"dabble/object"
	"dabble/token"
	"fmt"
	"strconv"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	p.nextToken()
	return p
}

func (p *Parser) ParseProgram() object.Value {
	return p.parseValue()
}

func (p *Parser) parseValue() object.Value {
	switch p.curToken.Type {
	case token.SYMBOL:
		return object.Symbol(p.curToken.Literal)
	case token.NUMBER:
		i, err := strconv.ParseUint(p.curToken.Literal, 10, 64)
		if err != nil {
			p.error(err.Error())
			return object.Null
		}
		return object.Number(i)
	case token.EOF:
		p.error("end of file")
		return object.Null
	case token.ILLEGAL:
		p.error("illegal: %v", p.curToken.Literal)
		return object.Null
	case token.LPAREN:
		p.nextToken()
		first := p.parseValue()
		p.nextToken()
		return object.Cell(
			first,
			p.parseValue())
	case token.RPAREN:
		return object.Null
	default:
		p.error("unknown token type: %v", p.curToken.Type)
		return object.Null
	}
}

func (p *Parser) nextToken() {
	p.curToken = p.l.NextToken()
}

func (p *Parser) error(s string, args ...interface{}) {
	p.errors = append(p.errors, fmt.Sprintf(s, args...))
}
