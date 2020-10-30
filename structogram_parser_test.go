package main

import "testing"

func checkOk(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Did not expect any errors")
	}
}

func TestParsingEmptyStringCausesError(t *testing.T) {
	parser, err := createParser()
	checkOk(t, err)

	diagram, err := parser.parseStructogram("")
	_ = diagram
	if err == nil {
		t.Errorf("Expected error but was nil")
	}
	expectedMsg := "Parsing error, structogram string is empty!"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error with msg %s, but got %s", expectedMsg, err.Error())
	}
}
