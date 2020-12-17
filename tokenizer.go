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
			closeParenthesesToken := Token{
				tokenType: "closeParentheses",
				value:     ")",
				line:      t.currentLineNumber,
				column:    t.currentColumnNumber,
			}
			t.currentColumnNumber++
			tokens = append(tokens, closeParenthesesToken)
			runes = nil
		case `"`, "'":
			quot := string(runes)
			str := ""
			for !t.isEof() {
				if string(t.next()) != quot {
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
		case " ", "\t", "\n":
			for !t.isEof() && t.isNextWhitespace() {
				runes = append(runes, t.readNext())
			}
			whitespace := Token{
				tokenType: "whitespace",
				value:     string(runes),
				line:      t.currentLineNumber,
				column:    t.currentColumnNumber,
			}
			for _, v := range runes {
				if string(v) == "\n" {
					t.currentLineNumber++
					t.currentColumnNumber = 1
				} else {
					t.currentColumnNumber++
				}
			}
			tokens = append(tokens, whitespace)
			runes = nil
		}
	}
	if len(runes) != 0 {
		v := string(runes)
		invalid := Token{
			tokenType: "invalid",
			value:     v,
			line:      t.currentLineNumber,
			column:    t.currentColumnNumber,
		}
		tokens = append(tokens, invalid)
		t.currentColumnNumber += len(v)
	}
	tokens = append(tokens, Token{
		tokenType: "EOF",
		value:     "EOF",
		line:      t.currentLineNumber,
		column:    t.currentColumnNumber,
	})
	return tokens
}
