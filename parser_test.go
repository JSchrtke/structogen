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

func checkNodeCount(t *testing.T, n []Node, count int) {
	t.Helper()
	if len(n) != count {
		t.Errorf("Wrong node count, expected %d, but got %d", count, len(n))
	}
}

func checkNode(t *testing.T, n Node, nodeType string, value string) {
	if n.nodeType != nodeType {
		t.Errorf(
			"Wrong node type, expected %s, but got %s", nodeType, n.nodeType,
		)
	}
	if n.value != value {
		t.Errorf("Wrong node value, expected %s, but got %s", value, n.value)
	}
}

func TestEmptyStructogramNameCausesError(t *testing.T) {
	tokens := makeTokens("name()")
	_, err := parseStructogram(tokens)
	checkErrorMsg(t, err, "1:6, expected 'string', but got 'closeParentheses'")
}

func TestStructogramsHaveNames(t *testing.T) {
	tokens := makeTokens(`name("test name")`)
	structogram, err := parseStructogram(tokens)
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
	_, err := parseStructogram(tokens)
	checkErrorMsg(t, err, "1:6, expected 'string', but got 'name'")
}

func TestNameHasToBeFirstToken(t *testing.T) {
	tokens := makeTokens(`instruction("something")name("a name")`)
	_, err := parseStructogram(tokens)
	checkErrorMsg(t, err, "1:1, expected 'name', but got 'instruction'")
}

func TestNameValueHasToBeEnclosedByParentheses(t *testing.T) {

	tokens := makeTokens(`name"a name"`)
	_, err := parseStructogram(tokens)
	checkErrorMsg(t, err, "1:5, expected 'openParentheses', but got 'string'")

	tokens = makeTokens(`name("a"(`)
	_, err = parseStructogram(tokens)
	checkErrorMsg(
		t, err, "1:9, expected 'closeParentheses', but got 'openParentheses'",
	)
}

func TestInstructionValueHasToBeEnclosedByParentheses(t *testing.T) {
	tokens := makeTokens(`name("some name")instruction"something")`)
	_, err := parseStructogram(tokens)
	checkErrorMsg(t, err, "1:29, expected 'openParentheses', but got 'string'")

	tokens = makeTokens(`name("a")instruction("b"(`)
	_, err = parseStructogram(tokens)
	checkErrorMsg(
		t, err, "1:25, expected 'closeParentheses', but got 'openParentheses'",
	)
}

func TestInstructionsCanNotBeEmpty(t *testing.T) {
	tokens := makeTokens(`name("test structogram")instruction()`)
	_, err := parseStructogram(tokens)
	checkErrorMsg(t, err, "1:37, expected 'string', but got 'closeParentheses'")
}

func TestInstuctionsCanNotBeNested(t *testing.T) {
	tokens := makeTokens(`name("a")instruction(instruction())`)
	_, err := parseStructogram(tokens)
	checkErrorMsg(t, err, "1:22, expected 'string', but got 'instruction'")
}

func TestStructogramCanHaveInstructions(t *testing.T) {
	tokens := makeTokens(`name("a")instruction("something")`)
	structogram, err := parseStructogram(tokens)
	checkOk(t, err)
	checkNodeCount(t, structogram.nodes, 1)
	checkNode(t, structogram.nodes[0], "instruction", "something")
}

func TestStructogramsCanHaveMultipleInstructions(t *testing.T) {
	tokens := makeTokens(`name("a")instruction("b")instruction("c")`)
	structogram, err := parseStructogram(tokens)
	checkOk(t, err)
	checkNodeCount(t, structogram.nodes, 2)
	checkNode(t, structogram.nodes[0], "instruction", "b")
	checkNode(t, structogram.nodes[1], "instruction", "c")
}

func TestParserCanHandleInvalidTokens(t *testing.T) {
	tokens := makeTokens(`name("a")asd`)
	_, err := parseStructogram(tokens)
	checkErrorMsg(t, err, "1:10, expected 'keyword', but got 'invalid'")
}

func TestParserIgnoresWhitespaceTokens(t *testing.T) {
	tokens := makeTokens(`name("a")` + "\n " + `instruction("b")`)
	_, err := parseStructogram(tokens)
	checkOk(t, err)
}

func TestIfTokenValuesAreEnclosedByParentheses(t *testing.T) {
	tokens := makeTokens(`name("a")if"b")`)
	_, err := parseStructogram(tokens)
	checkErrorMsg(t, err, "1:12, expected 'openParentheses', but got 'string'")

	tokens = makeTokens(`name("a")if("b"`)
	_, err = parseStructogram(tokens)
	checkErrorMsg(t, err, "1:16, expected 'closeParentheses', but got 'EOF'")
}

func TestIfTokenValueCanNotBeEmpty(t *testing.T) {
	tokens := makeTokens(`name("a")if()`)
	_, err := parseStructogram(tokens)
	checkErrorMsg(t, err, "1:13, expected 'string', but got 'closeParentheses'")
}

func TestIfTokenHasToHaveBody(t *testing.T) {
	tokens := makeTokens(`name("a")if("b")`)
	_, err := parseStructogram(tokens)
	checkErrorMsg(t, err, "1:17, expected 'openBrace', but got 'EOF'")

	// The only valid tokens inside of an if-body are keywords or whitespace.
	// Whitespace should get entirely ignored, and anything that is not a
	// keyword, so either a string or EOF should cause an error.
	// The only exception are openParentheses, which are legal if they
	// are preceeded by a keyword
	tokens = makeTokens(`name("a")if("b"){`)
	_, err = parseStructogram(tokens)
	checkErrorMsg(t, err, "1:18, expected 'keyword', but got 'EOF'")

	tokens = makeTokens(`name("a")if("b"){"c"`)
	_, err = parseStructogram(tokens)
	checkErrorMsg(t, err, "1:18, expected 'keyword', but got 'string'")

	tokens = makeTokens(`name("a")if("b"){name}`)
	_, err = parseStructogram(tokens)
	checkErrorMsg(t, err, "1:18, expected 'keyword', but got 'name'")

	tokens = makeTokens(`name("a") if("b") {instruction("c")`)
	_, err = parseStructogram(tokens)
	checkErrorMsg(t, err, "1:36, expected 'closeBrace', but got 'EOF'")
}

func TestIfTokenCanHaveWhitespaceBetweenConditionAndBody(t *testing.T) {
	tokens := makeTokens(`name("a")if("b")` + "\n " + `{instruction("c")}`)
	_, err := parseStructogram(tokens)
	checkOk(t, err)
}

func TestInstructionTokenInsideIfBodyBehavesTheSameAsOutside(t *testing.T) {
	tokens := makeTokens(`name("a") if("b") {instruction}`)
	_, err := parseStructogram(tokens)
	checkErrorMsg(
		t, err, "1:31, expected 'openParentheses', but got 'closeBrace'",
	)

	tokens = makeTokens(`name("a") if("b") {instruction(}`)
	_, err = parseStructogram(tokens)
	checkErrorMsg(t, err, "1:32, expected 'string', but got 'closeBrace'")

	tokens = makeTokens(`name("a") if("b") {instruction("c"}`)
	_, err = parseStructogram(tokens)
	checkErrorMsg(
		t, err, "1:35, expected 'closeParentheses', but got 'closeBrace'",
	)
}

func TestCanParseMultipleInstructionsInsideIfBody(t *testing.T) {
	tokens := makeTokens(
		`name("a") if("b") {instruction("c") instruction("d")}`,
	)
	structogram, err := parseStructogram(tokens)
	checkOk(t, err)

	checkNodeCount(t, structogram.nodes, 1)
	ifNode := structogram.nodes[0]
	checkNode(t, ifNode, "if", "b")

	ifBody := ifNode.nodes
	checkNodeCount(t, ifBody, 2)
	instructionNode := ifBody[0]
	checkNode(t, instructionNode, "instruction", "c")

	instructionNode = ifBody[1]
	checkNode(t, instructionNode, "instruction", "d")
}

func TestCanParseNestedIfs(t *testing.T) {
	tokens := makeTokens(`name("a") if("b") {if("c"){instruction("d")}}`)
	structogram, err := parseStructogram(tokens)
	checkOk(t, err)

	checkNodeCount(t, structogram.nodes, 1)
	ifNode := structogram.nodes[0]
	checkNode(t, ifNode, "if", "b")

	ifBody := ifNode.nodes
	checkNodeCount(t, ifBody, 1)

	nestedIf := ifBody[0]
	checkNode(t, nestedIf, "if", "c")

	nestedIfBody := nestedIf.nodes
	checkNodeCount(t, nestedIfBody, 1)
	checkNode(t, nestedIfBody[0], "instruction", "d")
}

func TestElseWithoutIfCausesError(t *testing.T) {
	tokens := makeTokens(`name("a") else {instruction("b")}`)
	_, err := parseStructogram(tokens)
	checkErrorMsg(t, err, "1:11, expected 'statement', but got 'else'")
}

func TestCanParseElse(t *testing.T) {
	tokens := makeTokens(
		`name("a") if("b") {instruction("c")} else {instruction("d")}`,
	)
	structogram, err := parseStructogram(tokens)
	checkOk(t, err)

	checkNodeCount(t, structogram.nodes, 2)

	elseNode := structogram.nodes[1]
	checkNode(t, elseNode, "else", "")
	checkNodeCount(t, elseNode.nodes, 1)

	elseBody := elseNode.nodes
	checkNodeCount(t, elseBody, 1)
	checkNode(t, elseBody[0], "instruction", "d")
}

func TestCanParseCall(t *testing.T) {
	tokens := makeTokens(`name("a") call`)
	_, err := parseStructogram(tokens)
	checkErrorMsg(t, err, "1:15, expected 'openParentheses', but got 'EOF'")

	tokens = makeTokens(`name("a") call(`)
	_, err = parseStructogram(tokens)
	checkErrorMsg(t, err, "1:16, expected 'string', but got 'EOF'")

	tokens = makeTokens(`name("a") call("b"`)
	_, err = parseStructogram(tokens)
	checkErrorMsg(t, err, "1:19, expected 'closeParentheses', but got 'EOF'")

	tokens = makeTokens(`name("a") call("b")`)
	structogram, err := parseStructogram(tokens)
	checkOk(t, err)

	checkNodeCount(t, structogram.nodes, 1)
	checkNode(t, structogram.nodes[0], "call", "b")
}

func TestCanParseCallInsideIfBody(t *testing.T) {
	tokens := makeTokens(`name("a") if("b") {call("c")}`)
	structogram, err := parseStructogram(tokens)
	checkOk(t, err)

	checkNodeCount(t, structogram.nodes, 1)
	ifNode := structogram.nodes[0]
	checkNode(t, ifNode, "if", "b")

	checkNodeCount(t, ifNode.nodes, 1)
	checkNode(t, ifNode.nodes[0], "call", "c")
}

func TestWhileHasToHaveCondition(t *testing.T) {
	tokens := makeTokens(`name("a") while`)
	_, err := parseStructogram(tokens)
	checkErrorMsg(t, err, "1:16, expected 'openParentheses', but got 'EOF'")

	tokens = makeTokens(`name("a") while(`)
	_, err = parseStructogram(tokens)
	checkErrorMsg(t, err, "1:17, expected 'string', but got 'EOF'")

	tokens = makeTokens(`name("a") while("a"`)
	_, err = parseStructogram(tokens)
	checkErrorMsg(t, err, "1:20, expected 'closeParentheses', but got 'EOF'")
}
