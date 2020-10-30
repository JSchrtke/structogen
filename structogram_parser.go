package main

type parser struct{}

func createParser() (*parser, error) {
	p := parser{}
	return &p, nil
}

type parsedObject struct{}

func (p *parser) parseStructogram(s string) (*parsedObject, error) {
	return &parsedObject{}, nil
}
