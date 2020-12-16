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

func checkTokenLineNumber(t *testing.T, token Token, lineNumber int) {
	t.Helper()
	if token.line != lineNumber {
		t.Errorf(fmt.Sprintf(
			"Expected token with line number: %d, but got line number: %d",
			lineNumber,
			token.line,
		))
	}
}

func checkTokenColumnNumber(t *testing.T, token Token, columnNumber int) {
	t.Helper()
	if token.column != columnNumber {
		t.Errorf(fmt.Sprintf(
			"Expected token with column number: %d, but got column number: %d",
			columnNumber,
			token.column,
		))
	}
}

func checkToken(
	t *testing.T, token Token, tokenType string,
	tokenValue string, lineNumber int, columnNumber int,
) {
	t.Helper()
	checkTokenType(t, token, tokenType)
	checkTokenValue(t, token, tokenValue)
	checkTokenLineNumber(t, token, lineNumber)
	checkTokenColumnNumber(t, token, columnNumber)
}

func TestTokenizingEmptyStringDoesNothing(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens("")
	checkTokenCount(t, tokens, 0)
}

func TestCanTokenizeName(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens("name")
	checkTokenCount(t, tokens, 1)
	checkToken(t, tokens[0], "name", "name", 1, 1)
}

func TestCanTokenizeOpenParentheses(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens("(")
	checkTokenCount(t, tokens, 1)
	checkToken(t, tokens[0], "openParentheses", "(", 1, 1)
}

func TestCanTokenizeCloseParentheses(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens(")")
	checkTokenCount(t, tokens, 1)
	checkToken(t, tokens[0], "closeParentheses", ")", 1, 1)
}

func TestCanTokenizeString(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens(`"a test string"`)
	checkTokenCount(t, tokens, 1)
	checkToken(t, tokens[0], "string", "a test string", 1, 1)
}

func TestCanTokenizeInstruction(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens("instruction")
	checkTokenCount(t, tokens, 1)
	checkToken(t, tokens[0], "instruction", "instruction", 1, 1)
}

func TestCanTokenizeSpace(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens(" ")
	checkTokenCount(t, tokens, 1)
	checkToken(t, tokens[0], "whitespace", " ", 1, 1)
}

func TestCanTokenizeMultipleSpaces(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens("  ")
	checkTokenCount(t, tokens, 1)
	checkToken(t, tokens[0], "whitespace", "  ", 1, 1)
}

func TestCanTokenizeTabs(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens("\t")
	checkTokenCount(t, tokens, 1)
	checkToken(t, tokens[0], "whitespace", "\t", 1, 1)
}

func TestCanTokenizeMultipleTabs(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens("\t\t")
	checkTokenCount(t, tokens, 1)
	checkToken(t, tokens[0], "whitespace", "\t\t", 1, 1)
}

func TestCanTokenizeNewlines(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens("\n")
	checkTokenCount(t, tokens, 1)
	checkToken(t, tokens[0], "whitespace", "\n", 1, 1)
}

func TestCanTokenizeMultipleNewlines(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens("\n\n")
	checkTokenCount(t, tokens, 1)
	checkToken(t, tokens[0], "whitespace", "\n\n", 1, 1)
}

func TestTokenizingNewlineAdvancesLineNumber(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens("name\ninstruction\n\ninstruction")
	checkTokenCount(t, tokens, 5)
	checkToken(t, tokens[0], "name", "name", 1, 1)
	checkToken(t, tokens[1], "whitespace", "\n", 1, 5)
	checkToken(t, tokens[2], "instruction", "instruction", 2, 1)
	checkToken(t, tokens[3], "whitespace", "\n\n", 2, 12)
	checkToken(t, tokens[4], "instruction", "instruction", 4, 1)
}

func TestDifferentWhitespacesAreOneToken(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens("\t \nname")
	checkTokenCount(t, tokens, 2)
	checkToken(t, tokens[0], "whitespace", "\t \n", 1, 1)
	checkToken(t, tokens[1], "name", "name", 2, 1)
}

func TestCanTokenizeMultipleTokens(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens(`name("a name") instruction("do this")`)
	checkTokenCount(t, tokens, 9)

	checkToken(t, tokens[0], "name", "name", 1, 1)
	checkToken(t, tokens[1], "openParentheses", "(", 1, 5)
	checkToken(t, tokens[2], "string", "a name", 1, 6)
	checkToken(t, tokens[3], "closeParentheses", ")", 1, 14)
	checkToken(t, tokens[4], "whitespace", " ", 1, 15)
	checkToken(t, tokens[5], "instruction", "instruction", 1, 16)
	checkToken(t, tokens[6], "openParentheses", "(", 1, 27)
	checkToken(t, tokens[7], "string", "do this", 1, 28)
	checkToken(t, tokens[8], "closeParentheses", ")", 1, 37)
}

func TestTokenizerCanHandleInvalidStrings(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens("some invalid string")
	checkTokenCount(t, tokens, 1)
	checkToken(t, tokens[0], "invalid", "some invalid string", 1, 1)

	tokenizer = makeTokenizer()
	tokens = tokenizer.makeTokens(`name("some name")invalid`)
	checkTokenCount(t, tokens, 5)
	checkToken(t, tokens[0], "name", "name", 1, 1)
	checkToken(t, tokens[1], "openParentheses", "(", 1, 5)
	checkToken(t, tokens[2], "string", "some name", 1, 6)
	checkToken(t, tokens[3], "closeParentheses", ")", 1, 17)
	checkToken(t, tokens[4], "invalid", "invalid", 1, 18)
}
