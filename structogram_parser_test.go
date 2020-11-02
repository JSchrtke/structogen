package main

import "testing"

func checkOk(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Did not expect any errors")
	}
}

func checkErrorMsg(t *testing.T, err error, expectedMsg string) {
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
	structogram, err := parseStructogram("has no name token")
	_ = structogram
	checkErrorMsg(t, err, "Structogram must have a name!")
}

func TestEmptyStructogramNameCausesError(t *testing.T) {
	structogram, err := parseStructogram("name()")
	_ = structogram
	checkErrorMsg(t, err, "Structograms can not have empty names!")
}

func TestStructogramsHaveNames(t *testing.T) {
	expectedName := "test name"
	structogram, err := parseStructogram("name(" + expectedName + ")")
	checkOk(t, err)

	if structogram.name != expectedName {
		t.Errorf(
			"Diagram has wrong name, expected: %s, but was: %s",
			expectedName, structogram.name,
		)
	}
}

func TestNamesCanNotBeNested(t *testing.T) {
	structogram, err := parseStructogram("name(name())")
	_ = structogram
	checkErrorMsg(t, err, "Structogram names can not be nested!")
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
