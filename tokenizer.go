package main

type Tokenizer struct {
	input             []rune
	runeIndex         int
	nextRuneIdx       int
	currentLineNumber int
}

func makeTokenizer() Tokenizer {
	return Tokenizer{
		runeIndex:         0,
		nextRuneIdx:       0,
		currentLineNumber: 1,
	}
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

func (t *Tokenizer) isNextWhitespace() bool {
	return string(t.next()) == " " ||
		string(t.next()) == "\t" ||
		string(t.next()) == "\n"
}

func (t *Tokenizer) makeWhitespaceToken(tokenValue string) Token {
	for !t.isEof() && t.isNextWhitespace() {
		tokenValue += string(t.readNext())
	}
	return Token{
		tokenType: "whitespace",
		value:     tokenValue,
		line:      t.currentLineNumber,
		column:    t.runeIndex,
	}
}

func (t *Tokenizer) makeTokens(s string) []Token {
	var tokens []Token
	var runes []rune
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
				line:      t.currentLineNumber,
				column:    t.runeIndex,
			}
			tokens = append(tokens, nameToken)
			runes = nil
		case "(":
			openParethesesToken := Token{
				tokenType: "openParentheses",
				value:     "(",
				line:      t.currentLineNumber,
				column:    t.runeIndex,
			}
			tokens = append(tokens, openParethesesToken)
			runes = nil
		case ")":
			openParethesesToken := Token{
				tokenType: "closeParentheses",
				value:     ")",
				line:      t.currentLineNumber,
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
					// We don't want the quotation marks in the string, so when
					// one is found, read and discard it.
					t.readNext()
					break
				}
			}
			stringToken := Token{
				tokenType: "string",
				value:     str,
				line:      t.currentLineNumber,
				column:    startIdx,
			}
			tokens = append(tokens, stringToken)
			runes = nil
		case "instruction":
			instructionToken := Token{
				tokenType: "instruction",
				value:     "instruction",
				line:      t.currentLineNumber,
				column:    t.runeIndex,
			}
			tokens = append(tokens, instructionToken)
			runes = nil
		case " ":
			whitespaceToken := t.makeWhitespaceToken(" ")
			tokens = append(tokens, whitespaceToken)
			runes = nil
		case "\t":
			whitespaceToken := t.makeWhitespaceToken("\t")
			tokens = append(tokens, whitespaceToken)
			runes = nil
		case "\n":
			whitespaceToken := t.makeWhitespaceToken("\n")
			tokens = append(tokens, whitespaceToken)
			runes = nil
			t.currentLineNumber++
		}
	}
	return tokens
}
