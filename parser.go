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
	tokenIndex     int
	tokens         []Token
	isInSwitchBody bool
	isInCaseBody   bool
}

type Token struct {
	tokenType string
	value     string
	line      int
	column    int
}

func parseStructogram(tokens []Token) (Structogram, error) {
	// we do not need whitespace for anything, so they just get discarded
	var cleanTokens []Token
	for _, t := range tokens {
		if t.tokenType != "whitespace" {
			cleanTokens = append(cleanTokens, t)
		}
	}

	p := Parser{
		tokenIndex: 0,
		tokens:     cleanTokens,
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
			switchNode := Node{}
			switchNode.nodeType = p.readNext().tokenType
			switchNode.value, err = p.parseParentheses()
			if err != nil {
				return nodes, err
			}
			switchNode.nodes, err = p.parseSwitchBody()
			if err != nil {
				return nodes, err
			}
			nodes = append(nodes, switchNode)
		case "default":
			if p.isInSwitchBody {
				return nodes, newTokenTypeError("keyword", p.next())
			}
			var defaultNode Node
			defaultNode.nodeType = p.readNext().tokenType
			defaultNode.value = ""

			defaultNode.nodes, err = p.parseBraces()
			if err != nil {
				return nodes, err
			}

			nodes = append(nodes, defaultNode)
		case "case":
			if p.isInCaseBody {
				return nodes, newTokenTypeError("keyword", p.next())
			}
			caseNode, err := p.parseConditional()
			if err != nil {
				return nodes, err
			}
			nodes = append(nodes, caseNode)
		}
	}
	if p.next().tokenType != delimiter {
		return nodes, newTokenTypeError(delimiter, p.next())
	}
	p.readNext()

	return nodes, err
}

func (p *Parser) parseBraces() ([]Node, error) {
	var body []Node
	if p.next().tokenType != "openBrace" {
		return body, newTokenTypeError("openBrace", p.next())
	}
	p.readNext()
	if !isKeyword(p.next().tokenType) {
		return body, newTokenTypeError("keyword", p.next())
	}
	body, err := p.parseUntil("closeBrace")
	if err != nil {
		return body, err
	}
	return body, nil
}

func (p *Parser) parseSwitchBody() ([]Node, error) {
	p.isInSwitchBody = true
	var switchBody []Node
	if p.next().tokenType != "openBrace" {
		return switchBody, newTokenTypeError("openBrace", p.next())
	}
	p.readNext()
	for p.next().tokenType == "case" {
		p.isInCaseBody = true
		caseNode, err := p.parseConditional()
		if err != nil {
			return switchBody, err
		}
		switchBody = append(switchBody, caseNode)
	}
	p.isInCaseBody = false
	if p.next().tokenType != "default" {
		return switchBody, newTokenTypeError("default", p.next())
	}
	var defaultNode Node
	defaultNode.nodeType = p.readNext().tokenType
	defaultNode.value = ""
	defaultBody, err := p.parseBraces()
	if err != nil {
		return switchBody, err
	}
	if p.next().tokenType != "closeBrace" {
		return switchBody, newTokenTypeError("closeBrace", p.next())
	}
	p.readNext()
	defaultNode.nodes = defaultBody
	switchBody = append(switchBody, defaultNode)
	p.isInSwitchBody = false
	return switchBody, nil
}

func (p *Parser) parseConditional() (Node, error) {
	var node Node

	node.nodeType = p.readNext().value

	v, err := p.parseParentheses()
	node.value = v
	if err != nil {
		return node, err
	}

	body, err := p.parseBraces()
	node.nodes = body
	return node, err
}

func (p *Parser) parseElse() (Node, error) {
	var elseNode Node

	elseNode.nodeType = p.readNext().tokenType

	elseBody, err := p.parseBraces()
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
	return s == "instruction" || s == "if" || s == "call" || s == "default" || s == "switch"
}
