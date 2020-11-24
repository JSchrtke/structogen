package main

type Tokenizer struct{}

func (t *Tokenizer) makeTokens(s string) []Token {
	var tokens []Token
	var runes []rune
	lineNumber := 1
	for i, r := range s {
		runes = append(runes, r)
		switch string(runes) {
		case "name":
			nameToken := Token{
				tokenType: "name",
				value:     "name",
				line:      lineNumber,
				column:    i,
			}
			tokens = append(tokens, nameToken)
		}
	}
	return tokens
}
