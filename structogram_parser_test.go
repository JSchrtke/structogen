package main

import "testing"

func checkOk(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Did not expect any errors")
	}
}

func TestCanCreateParser(t *testing.T) {
	parser, err := createParser()
	checkOk(t, err)
	_ = parser
}

func TestCanCallParseStructogram(t *testing.T) {
	parser, err := createParser()
	checkOk(t, err)

	diagram, err := parser.parseStructogram("")
	checkOk(t, err)
	_ = diagram
}
