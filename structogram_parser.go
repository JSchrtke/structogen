package main

import (
	"errors"
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

func parseTokens(tokens []Token) (ParsedObject, error) {
	return ParsedObject{}, errors.New("1:5, missing name")
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
