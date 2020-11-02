package main

import (
	"errors"
	"strings"
)

type Structogram struct {
	name         string
	instructions []string
}

func parseToken(s, tokenStart, tokenEnd string) string {
	tokenStartIndex := strings.Index(s, tokenStart)
	tokenEndIndex := strings.Index(s[tokenStartIndex:], tokenEnd) + tokenStartIndex
	return s[tokenStartIndex+len(tokenStart) : tokenEndIndex]
}

func parseStructogram(structogram string) (*Structogram, error) {
	if len(structogram) == 0 {
		return nil, errors.New("Parsing error, structogram string is empty!")
	}

	nameToken := "name("
	if !strings.Contains(structogram, nameToken) {
		return nil, errors.New("Structogram must have a name!")
	}

	parsed := Structogram{}

	parsed.name = parseToken(structogram, nameToken, ")")

	if len(parsed.name) == 0 {
		return nil, errors.New("Structograms can not have empty names!")
	}

	if strings.Contains(parsed.name, nameToken) {
		return nil, errors.New("Structogram names can not be nested!")
	}

	instructionToken := "instruction("
	if strings.Contains(structogram, instructionToken) {
		instruction := parseToken(structogram, instructionToken, ")")
		if len(instruction) == 0 {
			return nil, errors.New("Instructions can not be empty!")
		}
		parsed.instructions = append(
			parsed.instructions,
			instruction,
		)
	}

	return &parsed, nil
}
