package main

import (
	"fmt"
	"testing"
)

func checkTokenCount(t *testing.T, tokens []Token, count int) {
	t.Helper()
	if len(tokens) != count {
		t.Errorf(
			fmt.Sprintf("Expected %d tokens, but got %d", count, len(tokens)),
		)
	}
}

func checkTokenType(t *testing.T, token Token, typeString string) {
	t.Helper()
	if token.tokenType != typeString {
		t.Errorf(fmt.Sprintf(
			"Expected token of type %s, but got %s",
			typeString,
			token.tokenType,
		))
	}
}

func TestTokenizingEmptyStringDoesNothing(t *testing.T) {
	var tokenizer Tokenizer
	tokens := tokenizer.makeTokens("")
	checkTokenCount(t, tokens, 0)
}

func TestCanTokenizeName(t *testing.T) {
	var tokenizer Tokenizer
	tokens := tokenizer.makeTokens("name")
	checkTokenCount(t, tokens, 1)
	checkTokenType(t, tokens[0], "name")
}

func TestCanTokenizeOpenParentheses(t *testing.T) {
	var tokenizer Tokenizer
	tokens := tokenizer.makeTokens("(")
	checkTokenCount(t, tokens, 1)
	checkTokenType(t, tokens[0], "openParentheses")
}

func TestCanTokenizeCloseParentheses(t *testing.T) {
	var tokenizer Tokenizer
	tokens := tokenizer.makeTokens(")")
	checkTokenCount(t, tokens, 1)
	checkTokenType(t, tokens[0], "closeParentheses")
}

func TestCanTokenizeMultipleTokens(t *testing.T) {
	var tokenizer Tokenizer
	tokens := tokenizer.makeTokens("name()")
	checkTokenCount(t, tokens, 3)
	checkTokenType(t, tokens[0], "name")
	checkTokenType(t, tokens[1], "openParentheses")
	checkTokenType(t, tokens[2], "closeParentheses")
}

func TestCanTokenizeString(t *testing.T) {
	var tokenizer Tokenizer
	tokens := tokenizer.makeTokens(`"a test string"`)
	checkTokenCount(t, tokens, 1)
	checkTokenType(t, tokens[0], "string")
	if tokens[0].value != "a test string" {
		t.Errorf(fmt.Sprintf("Expected token with value 'a test string', but got %s", tokens[0].value))
	}
}

func TestCanTokenizeInstruction(t *testing.T) {
	var tokenizer Tokenizer
	tokens := tokenizer.makeTokens("instruction")
	checkTokenCount(t, tokens, 1)
	checkTokenType(t, tokens[0], "instruction")
}
