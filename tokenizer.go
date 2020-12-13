package main

type Tokenizer struct {
	input               []rune
	runeIndex           int
	nextRuneIdx         int
	currentLineNumber   int
	currentColumnNumber int
}

func makeTokenizer() Tokenizer {
	return Tokenizer{
		runeIndex:           0,
		nextRuneIdx:         0,
		currentLineNumber:   1,
		currentColumnNumber: 1,
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
		column:    t.currentColumnNumber,
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
				column:    t.currentColumnNumber,
			}
			tokens = append(tokens, nameToken)
			t.currentColumnNumber += len("name")
			runes = nil
		case "(":
			openParethesesToken := Token{
				tokenType: "openParentheses",
				value:     "(",
				line:      t.currentLineNumber,
				column:    t.currentColumnNumber,
			}
			t.currentColumnNumber++
			tokens = append(tokens, openParethesesToken)
			runes = nil
		case ")":
			openParethesesToken := Token{
				tokenType: "closeParentheses",
				value:     ")",
				line:      t.currentLineNumber,
				column:    t.currentColumnNumber,
			}
			t.currentColumnNumber++
			tokens = append(tokens, openParethesesToken)
			runes = nil
		case `"`:
			str := ""
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
				column:    t.currentColumnNumber,
			}
			tokens = append(tokens, stringToken)
			// While we don't want the quotation marks in the value of the
			// string, we do have to make sure the column number is still
			// correct.
			t.currentColumnNumber += len(str) + 2
			runes = nil
		case "instruction":
			instructionToken := Token{
				tokenType: "instruction",
				value:     "instruction",
				line:      t.currentLineNumber,
				column:    t.currentColumnNumber,
			}
			tokens = append(tokens, instructionToken)
			t.currentColumnNumber += len("instruction")
			runes = nil
		case " ":
			whitespaceToken := t.makeWhitespaceToken(" ")
			tokens = append(tokens, whitespaceToken)
			t.currentColumnNumber++
			runes = nil
		case "\t":
			whitespaceToken := t.makeWhitespaceToken("\t")
			tokens = append(tokens, whitespaceToken)
			t.currentColumnNumber++
			runes = nil
		case "\n":
			tokenValue := "\n"
			for !t.isEof() && t.isNextWhitespace() {
				tokenValue += string(t.readNext())
			}
			whitespaceToken := Token{
				tokenType: "whitespace",
				value:     tokenValue,
				line:      t.currentLineNumber,
				column:    t.currentColumnNumber,
			}
			tokens = append(tokens, whitespaceToken)
			t.currentLineNumber += len(tokenValue)
			t.currentColumnNumber = 1
			runes = nil
		}
	}
	return tokens
}
