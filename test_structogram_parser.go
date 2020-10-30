package main

import "testing"

func TestCanCreateParser(t *testing.T) {
	parser, err := createParser()
	if err != nil {
		t.Errorf("Did not expect any errors")
	}
	_ = parser
}
