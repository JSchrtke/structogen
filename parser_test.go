package main

import "testing"

func checkOk(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("Did not expect any errors, but got %s", err.Error())
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

func TestEmptyStructogramNameCausesError(t *testing.T) {
	tokens := makeTokens("name()")
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:6, expected 'string', but got 'closeParentheses'")
}

func TestStructogramsHaveNames(t *testing.T) {
	tokens := makeTokens(`name("test name")`)
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
	tokens := makeTokens("name(name())")
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:6, expected 'string', but got 'name'")
}

func TestNameHasToBeFirstToken(t *testing.T) {
	tokens := makeTokens(`instruction("something")name("a name")`)
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:1, expected 'name', but got 'instruction'")
}

func TestNameValueHasToBeEnclosedByParentheses(t *testing.T) {

	tokens := makeTokens(`name"a name"`)
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:5, expected 'openParentheses', but got 'string'")

	tokens = makeTokens(`name("a"(`)
	_, err = parseTokens(tokens)
	checkErrorMsg(t, err, "1:9, expected 'closeParentheses', but got 'openParentheses'")
}

func TestInstructionValueHasToBeEnclosedByParentheses(t *testing.T) {
	tokens := makeTokens(`name("some name")instruction"something")`)
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:29, expected 'openParentheses', but got 'string'")

	tokens = makeTokens(`name("a")instruction("b"(`)
	_, err = parseTokens(tokens)
	checkErrorMsg(t, err, "1:25, expected 'closeParentheses', but got 'openParentheses'")
}

func TestInstructionsCanNotBeEmpty(t *testing.T) {
	tokens := makeTokens(`name("test structogram")instruction()`)
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:37, expected 'string', but got 'closeParentheses'")
}

func TestInstuctionsCanNotBeNested(t *testing.T) {
	tokens := makeTokens(`name("a")instruction(instruction())`)
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:22, expected 'string', but got 'instruction'")
}

func TestStructogramCanHaveInstructions(t *testing.T) {
	tokens := makeTokens(`name("a")instruction("something")`)
	structogram, err := parseTokens(tokens)
	checkOk(t, err)
	if structogram.instructions[0] != "something" {
		t.Errorf("Instruction 0 is wrong, expected: '%s', but was: '%s'",
			"something", structogram.instructions[0],
		)
	}
}

func TestStructogramsCanHaveMultipleInstructions(t *testing.T) {
	tokens := makeTokens(`name("a")instruction("b")instruction("c")`)
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

func TestParserCanHandleInvalidTokens(t *testing.T) {
	tokens := makeTokens(`name("a")asd`)
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:10, expected 'identifier', but got 'invalid'")
}

func TestParserIgnoresWhitespaceTokens(t *testing.T) {
	tokens := makeTokens(`name("a")` + "\n " + `instruction("b")`)
	_, err := parseTokens(tokens)
	checkOk(t, err)
}

func TestIfTokenValuesAreEnclosedByParentheses(t *testing.T) {
	tokens := makeTokens(`name("a")if"b")`)
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:12, expected 'openParentheses', but got 'string'")

	tokens = makeTokens(`name("a")if("b"`)
	_, err = parseTokens(tokens)
	checkErrorMsg(t, err, "1:16, expected 'closeParentheses', but got 'EOF'")
}

func TestIfTokenValueCanNotBeEmpty(t *testing.T) {
	tokens := makeTokens(`name("a")if()`)
	_, err := parseTokens(tokens)
	checkErrorMsg(t, err, "1:13, expected 'string', but got 'closeParentheses'")
}
