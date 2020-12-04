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
	checkTokenType(t, tokens[0], "name")
}

func TestCanTokenizeOpenParentheses(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens("(")
	checkTokenCount(t, tokens, 1)
	checkTokenType(t, tokens[0], "openParentheses")
}

func TestCanTokenizeCloseParentheses(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens(")")
	checkTokenCount(t, tokens, 1)
	checkTokenType(t, tokens[0], "closeParentheses")
}

func TestCanTokenizeString(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens(`"a test string"`)
	checkTokenCount(t, tokens, 1)
	checkTokenType(t, tokens[0], "string")
	checkTokenValue(t, tokens[0], "a test string")
}

func TestCanTokenizeInstruction(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens("instruction")
	checkTokenCount(t, tokens, 1)
	checkTokenType(t, tokens[0], "instruction")
}

func TestCanTokenizeSpace(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens(" ")
	checkTokenCount(t, tokens, 1)
	checkTokenType(t, tokens[0], "whitespace")
	checkTokenValue(t, tokens[0], " ")
}

func TestCanTokenizeMultipleSpaces(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens("  ")
	checkTokenCount(t, tokens, 1)
	checkTokenType(t, tokens[0], "whitespace")
	checkTokenValue(t, tokens[0], "  ")
}

func TestCanTokenizeTabs(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens("\t")
	checkTokenCount(t, tokens, 1)
	checkTokenType(t, tokens[0], "whitespace")
	checkTokenValue(t, tokens[0], "\t")
}

func TestCanTokenizeMultipleTabs(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens("\t\t")
	checkTokenCount(t, tokens, 1)
	checkTokenType(t, tokens[0], "whitespace")
	checkTokenValue(t, tokens[0], "\t\t")
}

func TestCanTokenizeNewlines(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens("\n")
	checkTokenCount(t, tokens, 1)
	checkTokenType(t, tokens[0], "whitespace")
	checkTokenValue(t, tokens[0], "\n")
}

func TestCanTokenizeMultipleNewlines(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens("\n\n")
	checkTokenCount(t, tokens, 1)
	checkTokenType(t, tokens[0], "whitespace")
	checkTokenValue(t, tokens[0], "\n\n")
}

func TestTokenizingNewlineAdvancesLineNumber(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens("name\ninstruction")
	checkTokenCount(t, tokens, 3)
	checkTokenLineNumber(t, tokens[0], 1)
	checkTokenLineNumber(t, tokens[1], 1)
	checkTokenLineNumber(t, tokens[2], 2)
}

func TestColumnNumberResetsAfterEncounteringNewline(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens("name\ninstruction")
	checkTokenCount(t, tokens, 3)
	checkToken(t, tokens[0], "name", "name", 1, 1)
	checkToken(t, tokens[1], "whitespace", "\n", 1, 5)
	checkToken(t, tokens[2], "instruction", "instruction", 2, 1)
}

func TestCanTokenizeMultipleTokens(t *testing.T) {
	tokenizer := makeTokenizer()
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
