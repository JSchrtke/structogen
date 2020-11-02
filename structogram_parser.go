package main

import (
	"errors"
	"strings"
)

type Structogram struct {
	name string
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
	nameTokenStart := strings.Index(structogram, nameToken)
	nameTokenEnd := strings.Index(structogram[nameTokenStart:], ")")
	parsed.name = structogram[nameTokenStart+len(nameToken) : nameTokenEnd]

	if len(parsed.name) == 0 {
		return nil, errors.New("Structograms can not have empty names!")
	}

	if strings.Contains(parsed.name, nameToken) {
		return nil, errors.New("Structogram names can not be nested!")
	}

	return &parsed, nil
}
