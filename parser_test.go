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

func TestInstructionsCanNotBeEmpty(t *testing.T) {
	structogram, err := parseStructogram(
		"name(test structogram)\ninstruction()",
	)
	_ = structogram
	checkErrorMsg(t, err, "Instructions can not be empty!")
}

func TestInstuctionsCanNotBeNested(t *testing.T) {
	structogram, err := parseStructogram(
		"name(test structogram)\ninstruction(instruction())",
	)
	_ = structogram
	checkErrorMsg(t, err, "Instructions can not be nested!")
}

func TestStructogramCanHaveInstructions(t *testing.T) {
	structogram, err := parseStructogram(
		"name(test structogram)\ninstruction(do a thing)",
	)
	checkOk(t, err)
	if structogram.instructions[0] != "do a thing" {
		t.Errorf("Instruction 0 is wrong, expected: %s, but was: %s",
			"do a thing", structogram.instructions[0],
		)
	}
}

func TestStructogramsCanHaveMultipleInstructions(t *testing.T) {
	structogram, err := parseStructogram(
		"name(test structogram)\ninstruction(do a thing)\ninstruction(do another thing)",
	)
	checkOk(t, err)
	if structogram.instructions[0] != "do a thing" {
		t.Errorf("Instruction 0 is wrong, expected: %s, but was: %s",
			"do a thing", structogram.instructions[0],
		)
	} else if structogram.instructions[1] != "do another thing" {
		t.Errorf("Instruction 1 is wrong, expected: %s, but was: %s",
			"do another thing", structogram.instructions[1],
		)
	}
}
