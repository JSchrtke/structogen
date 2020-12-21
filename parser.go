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

func isKeyword(s string) bool {
	return s == "instruction" || s == "if"
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
	for {
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
		case "whitespace":
			// TODO Maybe have a function that runs once that strips all the
			// whitespace out of the tokens?
			// Whitespace should be completely ignored
			_ = p.readNext()
		case "if":
			p.readNext()
			if p.next().tokenType != "openParentheses" {
				return parsed, newTokenValueError("openParentheses", p.next())
			}
			// Discard the openParentheses
			p.readNext()
			if p.next().tokenType != "string" {
				return parsed, newTokenValueError("string", p.next())
			}
			// this is the value of the if, it get's discarded for now. Will be
			// implemented properly in a later test
			p.readNext()
			if p.next().tokenType != "closeParentheses" {
				return parsed, newTokenValueError("closeParentheses", p.next())
			}
			p.readNext()
			if p.next().tokenType == "whitespace" {
				p.readNext()
			}
			if p.next().tokenType != "openBrace" {
				return parsed, newTokenValueError("openBrace", p.next())
			}
			p.readNext()

			if !isKeyword(p.next().tokenType) {
				return parsed, newTokenValueError("keyword", p.next())
			}
			// Discard the keyword token for now
			p.readNext()

			if p.next().tokenType != "closeBrace" {
				return parsed, newTokenValueError("closeBrace", p.next())
			}
			p.readNext()
		case "invalid":
			return parsed, newTokenValueError("identifier", p.next())
		case "EOF":
			return parsed, err
		}
	}
}
