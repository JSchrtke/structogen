package main

type parser struct{}

func createParser() (*parser, error) {
	p := parser{}
	return &p, nil
}
