package main

import (
	"errors"
	"fmt"
	"strings"
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

// TODO find a better name for this
type ParsedObject struct {
	name string
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

func parseTokens(tokens []Token) (ParsedObject, error) {
	p := Parser{
		tokenIndex: 0,
		tokens:     tokens,
	}
	var parsed ParsedObject
	var err error
	if p.next().tokenType != "name" {
		return parsed, errors.New(
			fmt.Sprintf(
				"%d:%d, structogram has to start with a name",
				p.next().line,
				p.next().column,
			),
		)
	}
	for !p.isEof() {
		switch p.next().tokenType {
		case "name":
			_ = p.readNext()
			if p.next().tokenType != "openParentheses" {
				return parsed, newTokenValueError("openParentheses", p.next())
			}
			tok := p.readNext()
			if p.next().tokenType == "name" {
				return parsed, errors.New(
					fmt.Sprintf(
						"%d:%d, names can not be nested",
						p.next().line,
						p.next().column,
					),
				)
			}
			if p.next().tokenType != "string" {
				return parsed, errors.New(
					fmt.Sprintf(
						"%d:%d, missing name",
						tok.line,
						tok.column,
					),
				)
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
		}
	}
	return parsed, err
}

func parseToken(s, tokenStart, tokenEnd string) (content, remaining string) {
	tokenStartIndex := strings.Index(s, tokenStart)
	tokenEndIndex := strings.Index(s[tokenStartIndex:], tokenEnd) + tokenStartIndex
	content = s[tokenStartIndex+len(tokenStart) : tokenEndIndex]
	remaining = s[tokenEndIndex:]
	return
}

func parseStructogram(structogram string) (*Structogram, error) {
	if len(structogram) == 0 {
		return nil, errors.New("Parsing error, structogram string is empty!")
	}

	nameToken := "name("
	if strings.Index(structogram, nameToken) != 0 {
		return nil, errors.New("Structogram must have a name!")
	}

	parsed := Structogram{}

	var remaining string
	parsed.name, remaining = parseToken(structogram, nameToken, ")")

	if len(parsed.name) == 0 {
		return nil, errors.New("Structograms can not have empty names!")
	}
	if strings.Contains(parsed.name, nameToken) {
		return nil, errors.New("Structogram names can not be nested!")
	}

	instructionToken := "instruction("
	var instruction string
	for strings.Contains(remaining, instructionToken) {
		instruction, remaining = parseToken(remaining, instructionToken, ")")
		if len(instruction) == 0 {
			return nil, errors.New("Instructions can not be empty!")
		}
		if strings.Contains(instruction, instructionToken) {
			return nil, errors.New("Instructions can not be nested!")
		}
		parsed.instructions = append(
			parsed.instructions,
			instruction,
		)
	}

	return &parsed, nil
}
