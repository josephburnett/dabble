package parser

// Based on Monkey parser.go.

import (
	"dabble/lexer"
	"dabble/object"
	"dabble/token"
	"fmt"
	"strconv"
	"strings"
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

func (p *Parser) ParseProgram() (object.Value, error) {
	v := p.parseValue()
	p.nextToken()
	if p.curToken.Type != token.EOF {
		p.error("unexpected: %v", p.curToken.Literal)
	}
	if len(p.errors) != 0 {
		errString := strings.Join(p.errors, "\n")
		return nil, fmt.Errorf("error parsing:\n%v", errString)
	}
	return v, nil
}

func (p *Parser) parseValue() object.Value {
	switch p.curToken.Type {
	case token.RPAREN:
		p.error("unexpected: %v", p.curToken.Literal)
		return object.Nil
	case token.SYMBOL:
		return object.Symbol(p.curToken.Literal)
	case token.NUMBER:
		i, err := strconv.ParseUint(p.curToken.Literal, 10, 64)
		if err != nil {
			p.error("invalid number: %v", err.Error())
			return object.Nil
		}
		return object.Number(i)
	case token.EOF:
		p.error("end of file")
		return object.Nil
	case token.ILLEGAL:
		p.error("illegal: %v", p.curToken.Literal)
		return object.Nil
	case token.LPAREN:
		p.nextToken()
		return p.parseCell()
	default:
		p.error("unknown token type: %v", p.curToken.Type)
		return object.Nil
	}
}

func (p *Parser) parseCell() object.Value {
	switch p.curToken.Type {
	case token.RPAREN:
		return object.Nil
	case token.EOF:
		p.error("end of file")
		return object.Nil
	case token.ILLEGAL:
		p.error("illegal: %v", p.curToken.Literal)
		return object.Nil
	case token.DOT:
		p.error("expected value before dot")
		return object.Nil
	default:
		first := p.parseValue()
		if _, ok := first.(object.Error); ok {
			return first
		}
		p.nextToken()
		if p.curToken.Type == token.DOT {
			p.nextToken()
			return p.parseDottedList(first)
		}
		rest := p.parseList()
		if _, ok := rest.(object.Error); ok {
			return rest
		}
		return object.Cell{first, rest}
	}
}

func (p *Parser) parseDottedList(first object.Value) object.Value {
	rest := p.parseValue()
	if _, ok := rest.(object.Error); ok {
		return rest
	}
	p.nextToken()
	if p.curToken.Type != token.RPAREN {
		p.error("expecting ) after dot construction")
		return object.Nil
	}
	return object.Cell{first, rest}
}

func (p *Parser) parseList() object.Value {
	switch p.curToken.Type {
	case token.RPAREN:
		return object.Nil
	case token.EOF:
		p.error("end of file")
		return object.Nil
	case token.ILLEGAL:
		p.error("illegal: %v", p.curToken.Literal)
		return object.Nil
	default:
		first := p.parseValue()
		p.nextToken()
		rest := p.parseList()
		return object.Cell{first, rest}
	}
}

func (p *Parser) nextToken() {
	p.curToken = p.l.NextToken()
}

func (p *Parser) error(s string, args ...interface{}) {
	p.errors = append(p.errors, fmt.Sprintf(s, args...))
}
