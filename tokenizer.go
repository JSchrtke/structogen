package main

type Tokenizer struct {
	input       []rune
	runeIndex   int
	nextRuneIdx int
}

func (t *Tokenizer) isEof() bool {
	return t.nextRuneIdx > len(t.input)-1
}

func (t *Tokenizer) readNext() rune {
	t.runeIndex = t.nextRuneIdx
	t.nextRuneIdx++
	r := t.input[t.runeIndex]
	return r
}

func (t *Tokenizer) next() rune {
	return t.input[t.nextRuneIdx]
}

func (t *Tokenizer) makeTokens(s string) []Token {
	var tokens []Token
	var runes []rune
	lineNumber := 1
	t.nextRuneIdx = 0
	t.input = []rune(s)
	for !t.isEof() {
		r := t.readNext()
		runes = append(runes, r)
		switch string(runes) {
		case "name":
			nameToken := Token{
				tokenType: "name",
				value:     "name",
				line:      lineNumber,
				column:    t.runeIndex,
			}
			tokens = append(tokens, nameToken)
			runes = nil
		case "(":
			openParethesesToken := Token{
				tokenType: "openParentheses",
				value:     "(",
				line:      lineNumber,
				column:    t.runeIndex,
			}
			tokens = append(tokens, openParethesesToken)
			runes = nil
		case ")":
			openParethesesToken := Token{
				tokenType: "closeParentheses",
				value:     ")",
				line:      lineNumber,
				column:    t.runeIndex,
			}
			tokens = append(tokens, openParethesesToken)
			runes = nil
		case `"`:
			str := ""
			startIdx := t.runeIndex
			for !t.isEof() {
				if string(t.next()) != `"` {
					str += string(t.readNext())
				} else {
					t.readNext()
				}
			}
			stringToken := Token{
				tokenType: "string",
				value:     str,
				line:      lineNumber,
				column:    startIdx,
			}
			tokens = append(tokens, stringToken)
		case "instruction":
			instructionToken := Token{
				tokenType: "instruction",
				value:     "",
				line:      lineNumber,
				column:    t.runeIndex,
			}
			tokens = append(tokens, instructionToken)
		}
	}
	return tokens
}
