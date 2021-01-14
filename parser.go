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

func (p *Parser) parseParentheses() (string, error) {
	if p.next().tokenType != "openParentheses" {
		return "", newTokenValueError("openParentheses", p.next())
	}
	p.readNext()

	if p.next().tokenType != "string" {
		return "", newTokenValueError("string", p.next())
	}
	content := p.readNext().value

	if p.next().tokenType != "closeParentheses" {
		return "", newTokenValueError("closeParentheses", p.next())
	}
	p.readNext()
	return content, nil
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

			n.value, err = p.parseParentheses()
			if err != nil {
				return nodes, err
			}

			nodes = append(nodes, n)
		case "if":
			ifNode, err := p.parseConditional()
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
			whileNode, err := p.parseConditional()
			if err != nil {
				return nodes, err
			}
			nodes = append(nodes, whileNode)
		}
	}
	if p.next().tokenType != delimiter {
		return nodes, newTokenValueError(delimiter, p.next())
	}
	p.readNext()

	return nodes, err
}

func (p *Parser) parseConditional() (Node, error) {
	var node Node

	node.nodeType = p.readNext().value

	v, err := p.parseParentheses()
	node.value = v
	if err != nil {
		return node, err
	}

	if p.next().tokenType == "whitespace" {
		p.readNext()
	}

	if p.next().tokenType != "openBrace" {
		return node, newTokenValueError("openBrace", p.next())
	}
	p.readNext()

	if p.next().tokenType == "whitespace" {
		p.readNext()
	}

	if !isKeyword(p.next().tokenType) {
		return node, newTokenValueError("keyword", p.next())
	}
	body, err := p.parseTokensUntil("closeBrace")
	node.nodes = body

	return node, err
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
	parsed.name, err = p.parseParentheses()
	if err != nil {
		return parsed, err
	}

	nodes, err := p.parseTokensUntil("EOF")
	parsed.nodes = nodes
	return parsed, err
}
