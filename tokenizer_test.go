package main

import (
	"fmt"
	"testing"
)

func TestTokenizingEmptyStringDoesNothing(t *testing.T) {
	var tokenizer Tokenizer
	tokens := tokenizer.makeTokens("")
	if len(tokens) != 0 {
		t.Errorf(fmt.Sprintf("Expected no tokens, but got %d", len(tokens)))
	}
}

func TestCanTokenizeName(t *testing.T) {
	var tokenizer Tokenizer
	tokens := tokenizer.makeTokens("name")
	if len(tokens) != 1 {
		t.Errorf(fmt.Sprintf("Expected 1 token, but got %d", len(tokens)))
	}
	if tokens[0].tokenType != "name" {
		t.Errorf(fmt.Sprintf(
			"Expected token of type name, but got %s", tokens[0].tokenType,
		))
	}
}
