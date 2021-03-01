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

	nodes, err := p.parseUntil("EOF")
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

func (p *Parser) parseUntil(delimiter string) ([]Node, error) {
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
		case "while", "dowhile":
			conditionalNode, err := p.parseConditional()
			if err != nil {
				return nodes, err
			}
			nodes = append(nodes, conditionalNode)
		case "switch":
			// discard the switch token
			p.readNext()
			// next token should be open parentheses
			if p.next().tokenType != "openParentheses" {
				return nodes, newTokenTypeError("openParentheses", p.next())
			}
			// discard the open parentheses
			p.readNext()
			// next token should be a string
			if p.next().tokenType != "string" {
				return nodes, newTokenTypeError("string", p.next())
			}
			// discard the string for now
			p.readNext()
			// next token should be closeParentheses
			if p.next().tokenType != "closeParentheses" {
				return nodes, newTokenTypeError("closeParentheses", p.next())
			}
			// discard the closeParentheses
			p.readNext()
			// the next token should be openBrace
			if p.next().tokenType != "openBrace" {
				return nodes, newTokenTypeError("openBrace", p.next())
			}
			// discard the openBrace
			p.readNext()
			// the next token should be closeBrace
			if p.next().tokenType != "closeBrace" {
				return nodes, newTokenTypeError("closeBrace", p.next())
			}
		case "default":
			var defaultNode Node
			defaultNode.nodeType = p.readNext().tokenType
			defaultNode.value = ""

			if p.next().tokenType == "whitespace" {
				p.readNext()
			}
			if p.next().tokenType != "openBrace" {
				return nodes, newTokenTypeError("openBrace", p.next())
			}
			p.readNext()
			if p.next().tokenType == "whitespace" {
				p.readNext()
			}

			defaultNode.nodes, err = p.parseUntil("closeBrace")
			if err != nil {
				return nodes, err
			}

			nodes = append(nodes, defaultNode)
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
	body, err := p.parseUntil("closeBrace")
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
	elseBody, err := p.parseUntil("closeBrace")
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
