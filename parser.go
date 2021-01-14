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
	return s == "instruction" || s == "if" || s == "call"
}

func (p *Parser) parseTokensUntil(delimiter string) ([]Node, error) {
	var nodes []Node
	var err error

	for p.next().tokenType != delimiter {
		switch p.next().tokenType {
		case "EOF":
			return nodes, newTokenValueError(delimiter, p.next())
		case "whitespace":
			p.readNext()
		case "invalid":
			return nodes, newTokenValueError("keyword", p.next())
		case "instruction", "call":
			var n Node
			n.nodeType = p.readNext().tokenType

			if p.next().tokenType != "openParentheses" {
				return nodes, newTokenValueError("openParentheses", p.next())
			}
			p.readNext()

			if p.next().tokenType != "string" {
				return nodes, newTokenValueError("string", p.next())
			}
			n.value = p.readNext().value

			if p.next().tokenType != "closeParentheses" {
				return nodes, newTokenValueError("closeParentheses", p.next())
			}
			p.readNext()

			nodes = append(nodes, n)
		case "if":
			ifNode, err := p.parseIf()
			if err != nil {
				return nodes, err
			}
			nodes = append(nodes, ifNode)

			if p.next().tokenType == "whitespace" {
				p.readNext()
			}
			if p.next().tokenType == "else" {
				elseNode, err := p.parseElse()
				if err != nil {
					return nodes, err
				}
				nodes = append(nodes, elseNode)
			}
		case "else":
			return nodes, newTokenValueError("statement", p.next())
		case "while":
			err := p.parseWhile()
			if err != nil {
				return nodes, err
			}
		}
	}
	if p.next().tokenType != delimiter {
		return nodes, newTokenValueError(delimiter, p.next())
	}
	p.readNext()

	return nodes, err
}

func (p *Parser) parseIf() (Node, error) {
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

	// Parsing of the if body
	if p.next().tokenType != "openBrace" {
		return ifNode, newTokenValueError("openBrace", p.next())
	}
	p.readNext()

	if !isKeyword(p.next().tokenType) {
		return ifNode, newTokenValueError("keyword", p.next())
	}
	ifBody, err := p.parseTokensUntil("closeBrace")
	ifNode.nodes = ifBody

	return ifNode, err
}

func (p *Parser) parseElse() (Node, error) {
	var elseNode Node

	elseNode.nodeType = p.readNext().tokenType
	if p.next().tokenType == "whitespace" {
		p.readNext()
	}
	if p.next().tokenType != "openBrace" {
		return elseNode, newTokenValueError("openBrace", p.next())
	}
	p.readNext()
	if !isKeyword(p.next().tokenType) {
		return elseNode, newTokenValueError("keyword", p.next())
	}
	elseBody, err := p.parseTokensUntil("closeBrace")
	elseNode.nodes = elseBody

	return elseNode, err
}

func (p *Parser) parseWhile() error {
	// discard the while token
	p.readNext()

	// TODO Maybe refactor this to something like p.checkNextType. To
	// be able to do that, I need a global error in the parser and then
	// probably check if the error is nil each parsing iteration.
	if p.next().tokenType != "openParentheses" {
		return newTokenValueError("openParentheses", p.next())
	}
	// discard the openParentheses
	p.readNext()

	if p.next().tokenType != "string" {
		return newTokenValueError("string", p.next())
	}
	// TODO discard the while condition for now
	p.readNext()

	if p.next().tokenType != "closeParentheses" {
		return newTokenValueError("closeParentheses", p.next())
	}

	return nil
}

func parseStructogram(tokens []Token) (Structogram, error) {
	p := Parser{
		tokenIndex: 0,
		tokens:     tokens,
	}
	var parsed Structogram
	var err error
	if p.next().tokenType != "name" {
		return parsed, newTokenValueError("name", p.next())
	}
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
	p.readNext()

	nodes, err := p.parseTokensUntil("EOF")
	parsed.nodes = nodes
	return parsed, err
}
