package main

import "testing"

func checkOk(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Did not expect any errors")
	}
}

func TestParsingEmptyStringCausesError(t *testing.T) {
	structogram, err := parseStructogram("")
	_ = structogram
	if err == nil {
		t.Errorf("Expected error but was nil")
	}
	expectedMsg := "Parsing error, structogram string is empty!"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error with msg %s, but got %s", expectedMsg, err.Error())
	}
}

func TestStructogramHasToHaveAName(t *testing.T) {
	structogram, err := parseStructogram("has no name token")
	_ = structogram
	if err == nil {
		t.Errorf("Expected error but was nil")
	}
	expectedMsg := "Structogram must have a name!"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error with msg %s, but got %s", expectedMsg, err.Error())
	}
}

func TestEmptyStructogramNameCausesError(t *testing.T) {
	structogram, err := parseStructogram("name()")
	_ = structogram
	if err == nil {
		t.Errorf("Expected error but was nil")
	}
	expectedMsg := "Structograms can not have empty names!"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error with msg %s, but got %s", expectedMsg, err.Error())
	}
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
