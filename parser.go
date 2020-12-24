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
	name  string
	nodes []Node
}

type Node struct {
	nodeType string
	value    string
	nodes    []Node
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

func (p *Parser) parseIf() (Node, error) {
	var nodes []Node
	var ifNode Node

	ifNode.nodeType = p.readNext().value
	if p.next().tokenType != "openParentheses" {
		return ifNode, newTokenValueError("openParentheses", p.next())
	}
	// Discard the openParentheses
	p.readNext()

	if p.next().tokenType != "string" {
		return ifNode, newTokenValueError("string", p.next())
	}
	ifNode.value = p.readNext().value

	if p.next().tokenType != "closeParentheses" {
		return ifNode, newTokenValueError("closeParentheses", p.next())
	}
	p.readNext()
	if p.next().tokenType == "whitespace" {
		p.readNext()
	}
	if p.next().tokenType != "openBrace" {
		return ifNode, newTokenValueError("openBrace", p.next())
	}
	p.readNext()

	if !isKeyword(p.next().tokenType) {
		return ifNode, newTokenValueError("keyword", p.next())
	}

	// Parsing of the if body
	var err error
	for p.next().tokenType != "closeBrace" {
		switch p.next().tokenType {
		case "whitespace":
			p.readNext()
		case "instruction":
			var instructionNode Node
			instructionNode.nodeType = p.readNext().tokenType

			if p.next().tokenType != "openParentheses" {
				return ifNode, newTokenValueError(
					"openParentheses", p.next(),
				)
			}
			p.readNext()

			if p.next().tokenType != "string" {
				return ifNode, newTokenValueError("string", p.next())
			}
			instructionNode.value = p.readNext().value

			if p.next().tokenType != "closeParentheses" {
				return ifNode, newTokenValueError(
					"closeParentheses", p.next(),
				)
			}
			p.readNext()
			ifNode.nodes = append(ifNode.nodes, instructionNode)
		case "if":
			n, err := p.parseIf()
			if err != nil {
				return n, err
			}
			ifNode.nodes = append(ifNode.nodes, n)
		}
	}

	if p.next().tokenType != "closeBrace" {
		return ifNode, newTokenValueError("closeBrace", p.next())
	}
	p.readNext()
	nodes = append(nodes, ifNode)
	return ifNode, err
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
			var instructionNode Node
			instructionNode.nodeType = p.readNext().tokenType

			if p.next().tokenType != "openParentheses" {
				return parsed, newTokenValueError("openParentheses", p.next())
			}
			p.readNext()

			if p.next().tokenType != "string" {
				return parsed, newTokenValueError("string", p.next())
			}
			instructionNode.value = p.readNext().value

			if p.next().tokenType != "closeParentheses" {
				return parsed, newTokenValueError("closeParentheses", p.next())
			}
			p.readNext()

			parsed.nodes = append(parsed.nodes, instructionNode)
		case "whitespace":
			// TODO Maybe have a function that runs once that strips all the
			// whitespace out of the tokens?
			// Whitespace should be completely ignored
			_ = p.readNext()
		case "if":
			n, err := p.parseIf()
			if err != nil {
				return parsed, err
			}
			parsed.nodes = append(parsed.nodes, n)
		case "invalid":
			return parsed, newTokenValueError("identifier", p.next())
		case "EOF":
			return parsed, err
		}
	}
}
