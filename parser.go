package main

import (
	"errors"
	"fmt"
)

type Token struct {
	tokenType string
	value     string
	line      int
	column    int
}

type Structogram struct {
	name         string
	instructions []string
}

type Parser struct {
	tokenIndex int
	tokens     []Token
}

func (p *Parser) isEof() bool {
	return p.tokenIndex >= len(p.tokens)-1
}

func (p *Parser) next() Token {
	return p.tokens[p.tokenIndex]
}

func (p *Parser) readNext() Token {
	t := p.next()
	p.tokenIndex++
	return t
}

func newTokenValueError(expected string, actual Token) error {
	return errors.New(
		fmt.Sprintf(
			"%d:%d, expected '%s', but got '%s'",
			actual.line,
			actual.column,
			expected,
			actual.tokenType,
		),
	)
}

func parseTokens(tokens []Token) (Structogram, error) {
	p := Parser{
		tokenIndex: 0,
		tokens:     tokens,
	}
	var parsed Structogram
	var err error
	if p.next().tokenType != "name" {
		return parsed, newTokenValueError("name", p.next())
	}
	for !p.isEof() {
		switch p.next().tokenType {
		case "name":
			_ = p.readNext()
			if p.next().tokenType != "openParentheses" {
				return parsed, newTokenValueError("openParentheses", p.next())
			}
			tok := p.readNext()
			if p.next().tokenType != "string" {
				return parsed, newTokenValueError("string", p.next())
			}
			tok = p.readNext()
			parsed.name = tok.value
			if p.next().tokenType != "closeParentheses" {
				return parsed, newTokenValueError("closeParentheses", p.next())
			}
			_ = p.readNext()
		case "instruction":
			_ = p.readNext()
			if p.next().tokenType != "openParentheses" {
				return parsed, newTokenValueError("openParentheses", p.next())
			}
			_ = p.readNext()
			if p.next().tokenType != "string" {
				return parsed, newTokenValueError("string", p.next())
			}
			parsed.instructions = append(parsed.instructions, p.readNext().value)
			if p.next().tokenType != "closeParentheses" {
				return parsed, newTokenValueError("closeParentheses", p.next())
			}
			_ = p.readNext()
		}
	}
	return parsed, err
}
