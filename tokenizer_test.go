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
			"Expected token of type: %s, but got type: %s",
			typeString,
			token.tokenType,
		))
	}
}

func checkTokenValue(t *testing.T, token Token, value string) {
	t.Helper()
	if token.value != value {
		t.Errorf(fmt.Sprintf(
			"Expected token with value: %s, but got value: %s",
			value,
			token.value,
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

func TestCanTokenizeString(t *testing.T) {
	var tokenizer Tokenizer
	tokens := tokenizer.makeTokens(`"a test string"`)
	checkTokenCount(t, tokens, 1)
	checkTokenType(t, tokens[0], "string")
	checkTokenValue(t, tokens[0], "a test string")
}

func TestCanTokenizeInstruction(t *testing.T) {
	var tokenizer Tokenizer
	tokens := tokenizer.makeTokens("instruction")
	checkTokenCount(t, tokens, 1)
	checkTokenType(t, tokens[0], "instruction")
}

func TestCanTokenizeSpace(t *testing.T) {
	var tokenizer Tokenizer
	tokens := tokenizer.makeTokens(" ")
	checkTokenCount(t, tokens, 1)
	checkTokenType(t, tokens[0], "whitespace")
	checkTokenValue(t, tokens[0], " ")
}

func TestCanTokenizeMultipleSpaces(t *testing.T) {
	var tokenizer Tokenizer
	tokens := tokenizer.makeTokens("  ")
	checkTokenCount(t, tokens, 1)
	checkTokenType(t, tokens[0], "whitespace")
	checkTokenValue(t, tokens[0], "  ")
}

func TestCanTokenizeMultipleTokens(t *testing.T) {
	var tokenizer Tokenizer
	tokens := tokenizer.makeTokens(`name("a name") instruction("do this")`)
	checkTokenCount(t, tokens, 9)
	checkTokenType(t, tokens[0], "name")
	checkTokenValue(t, tokens[0], "name")
	checkTokenType(t, tokens[1], "openParentheses")
	checkTokenValue(t, tokens[1], "(")
	checkTokenType(t, tokens[2], "string")
	checkTokenValue(t, tokens[2], "a name")
	checkTokenType(t, tokens[3], "closeParentheses")
	checkTokenValue(t, tokens[3], ")")
	checkTokenType(t, tokens[4], "whitespace")
	checkTokenValue(t, tokens[4], " ")
	checkTokenType(t, tokens[5], "instruction")
	checkTokenValue(t, tokens[5], "instruction")
	checkTokenType(t, tokens[6], "openParentheses")
	checkTokenValue(t, tokens[6], "(")
	checkTokenType(t, tokens[7], "string")
	checkTokenValue(t, tokens[7], "do this")
	checkTokenType(t, tokens[8], "closeParentheses")
	checkTokenValue(t, tokens[8], ")")
}
