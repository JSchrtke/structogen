package main

type Tokenizer struct {
	input               []rune
	runeIndex           int
	nextRuneIdx         int
	currentLineNumber   int
	currentColumnNumber int
	tokens              []Token
	runes               []rune
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

func (t *Tokenizer) emitToken(tokenType string) {
	v := string(t.runes)
	if tokenType == "EOF" {
		v = "EOF"
	}
	tok := Token{
		tokenType: tokenType,
		value:     v,
		line:      t.currentLineNumber,
		column:    t.currentColumnNumber,
	}
	t.tokens = append(t.tokens, tok)
	t.currentColumnNumber += len(v)
	t.runes = nil
}

func makeTokens(s string) []Token {
	t := Tokenizer{
		runeIndex:           0,
		nextRuneIdx:         0,
		currentLineNumber:   1,
		currentColumnNumber: 1,
	}
	t.nextRuneIdx = 0
	t.input = []rune(s)
	for !t.isEof() {
		r := t.readNext()
		t.runes = append(t.runes, r)
		switch string(t.runes) {
		case "name":
			t.emitToken("name")
		case "(":
			t.emitToken("openParentheses")
		case ")":
			t.emitToken("closeParentheses")
		case "if":
			t.emitToken("if")
		case "{":
			t.emitToken("openBrace")
		case "}":
			t.emitToken("closeBrace")
		case "instruction":
			t.emitToken("instruction")
		case "else":
			t.emitToken("else")
		case "call":
			t.emitToken("call")
		case "while":
			t.emitToken("while")
		case "dowhile":
			t.emitToken("dowhile")
		case "switch":
			t.emitToken("switch")
		case "default":
			t.emitToken("default")
		case "case":
			t.emitToken("case")
		case `"`, "'":
			quot := string(t.runes)
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
			t.tokens = append(t.tokens, stringToken)
			// While we don't want the quotation marks in the value of the
			// string, we do have to make sure the column number is still
			// correct.
			t.currentColumnNumber += len(str) + 2
			t.runes = nil
		case " ", "\t", "\n":
			for !t.isEof() && t.isNextWhitespace() {
				t.runes = append(t.runes, t.readNext())
			}
			whitespace := Token{
				tokenType: "whitespace",
				value:     string(t.runes),
				line:      t.currentLineNumber,
				column:    t.currentColumnNumber,
			}
			for _, v := range t.runes {
				if string(v) == "\n" {
					t.currentLineNumber++
					t.currentColumnNumber = 1
				} else {
					t.currentColumnNumber++
				}
			}
			t.tokens = append(t.tokens, whitespace)
			t.runes = nil
		}
	}
	if len(t.runes) != 0 {
		t.emitToken("invalid")
	}
	t.emitToken("EOF")
	return t.tokens
}
