package main

import "testing"

func checkOk(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("Did not expect any errors")
	}
}

func checkErrorMsg(t *testing.T, err error, expectedMsg string) {
	t.Helper()
	if err == nil {
		t.Errorf("Expected error but was nil")
	}
	if err.Error() != expectedMsg {
		t.Errorf(
			"Expected error with msg %s, but got %s",
			expectedMsg,
			err.Error(),
		)
	}
}

func TestStructogramHasToHaveAName(t *testing.T) {
	// TODO Is this test still needed in this form if we are parsing tokens?
	structogram, err := parseStructogram("has no name token")
	_ = structogram
	checkErrorMsg(t, err, "Structogram must have a name!")
}

func TestEmptyStructogramNameCausesError(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens("name()")
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:5, missing name")
}

func TestStructogramsHaveNames(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens(`name("test name")`)
	structogram, err := parseTokens(tokens)
	checkOk(t, err)

	if structogram.name != "test name" {
		t.Errorf(
			"Diagram has wrong name, expected: %s, but was: %s",
			"test name", structogram.name,
		)
	}
}

func TestNamesCanNotBeNested(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens("name(name())")
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:6, names can not be nested")
}

func TestNameHasToBeFirstToken(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens(`instruction("something")name("a name")`)
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:1, structogram has to start with a name")
}

func TestNameValueHasToBeEnclosedByParentheses(t *testing.T) {
	tokenizer := makeTokenizer()

	tokens := tokenizer.makeTokens(`name"a name"`)
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:5, expected 'openParentheses', but got 'string'")

	tokenizer = makeTokenizer()
	tokens = tokenizer.makeTokens(`name("a"(`)
	_, err = parseTokens(tokens)
	checkErrorMsg(t, err, "1:9, expected 'closeParentheses', but got 'openParentheses'")
}

func TestInstructionValueHasToBeEnclosedByParentheses(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens(`name("some name")instruction"something")`)
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:29, expected 'openParentheses', but got 'string'")

	tokenizer = makeTokenizer()
	tokens = tokenizer.makeTokens(`name("a")instruction("b"(`)
	_, err = parseTokens(tokens)
	checkErrorMsg(t, err, "1:25, expected 'closeParentheses', but got 'openParentheses'")
}

func TestInstructionsCanNotBeEmpty(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens(`name("test structogram")instruction()`)
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:37, expected 'string', but got 'closeParentheses'")
}

func TestInstuctionsCanNotBeNested(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens(`name("a")instruction(instruction())`)
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:22, expected 'string', but got 'instruction'")
}

func TestStructogramCanHaveInstructions(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens(`name("a")instruction("something")`)
	structogram, err := parseTokens(tokens)
	checkOk(t, err)
	if structogram.instructions[0] != "something" {
		t.Errorf("Instruction 0 is wrong, expected: '%s', but was: '%s'",
			"something", structogram.instructions[0],
		)
	}
}

func TestStructogramsCanHaveMultipleInstructions(t *testing.T) {
	tokenizer := makeTokenizer()
	tokens := tokenizer.makeTokens(`name("a")instruction("b")instruction("c")`)
	structogram, err := parseTokens(tokens)
	checkOk(t, err)
	if len(structogram.instructions) != 2 {
		t.Errorf("Wrong instruction count, expected %d, but was %d",
			2, len(structogram.instructions),
		)
	}
	if structogram.instructions[0] != "b" {
		t.Errorf("Instruction 0 is wrong, expected: %s, but was: %s",
			"b", structogram.instructions[0],
		)
	} else if structogram.instructions[1] != "c" {
		t.Errorf("Instruction 1 is wrong, expected: %s, but was: %s",
			"c", structogram.instructions[1],
		)
	}
}
