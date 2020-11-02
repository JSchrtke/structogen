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
		t.Errorf("Expected error with msg %s, but got %s", expectedMsg, err.Error())
	}
}

func TestParsingEmptyStringCausesError(t *testing.T) {
	structogram, err := parseStructogram("")
	_ = structogram
	checkErrorMsg(t, err, "Parsing error, structogram string is empty!")
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
