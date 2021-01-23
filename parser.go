package main

import (
	"errors"
	"fmt"
)

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

type Token struct {
	tokenType string
	value     string
	line      int
	column    int
}

func parseStructogram(tokens []Token) (Structogram, error) {
	p := Parser{
		tokenIndex: 0,
		tokens:     tokens,
	}
	var parsed Structogram
	var err error
	if p.next().tokenType != "name" {
		return parsed, newTokenTypeError("name", p.next())
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

func (p *Parser) next() Token {
	return p.tokens[p.tokenIndex]
}

func (p *Parser) readNext() Token {
	t := p.next()
	p.tokenIndex++
	return t
}

func (p *Parser) parseParentheses() (string, error) {
	if p.next().tokenType != "openParentheses" {
		return "", newTokenTypeError("openParentheses", p.next())
	}
	p.readNext()

	if p.next().tokenType != "string" {
		return "", newTokenTypeError("string", p.next())
	}
	content := p.readNext().value

	if p.next().tokenType != "closeParentheses" {
		return "", newTokenTypeError("closeParentheses", p.next())
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
			return nodes, newTokenTypeError(delimiter, p.next())
		case "whitespace":
			p.readNext()
		case "invalid":
			return nodes, newTokenTypeError("keyword", p.next())
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
			return nodes, newTokenTypeError("statement", p.next())
		case "while":
			whileNode, err := p.parseConditional()
			if err != nil {
				return nodes, err
			}
			nodes = append(nodes, whileNode)
		}
	}
	if p.next().tokenType != delimiter {
		return nodes, newTokenTypeError(delimiter, p.next())
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
		return node, newTokenTypeError("openBrace", p.next())
	}
	p.readNext()

	if p.next().tokenType == "whitespace" {
		p.readNext()
	}

	if !isKeyword(p.next().tokenType) {
		return node, newTokenTypeError("keyword", p.next())
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
		return elseNode, newTokenTypeError("openBrace", p.next())
	}
	p.readNext()
	if !isKeyword(p.next().tokenType) {
		return elseNode, newTokenTypeError("keyword", p.next())
	}
	elseBody, err := p.parseTokensUntil("closeBrace")
	elseNode.nodes = elseBody

	return elseNode, err
}

func newTokenTypeError(expected string, actual Token) error {
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
