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
