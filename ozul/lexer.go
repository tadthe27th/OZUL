package main

import (
	"unicode"
)

type Lexer struct {
	source string
	pos    int
	line   int
	column int
	ch     rune
}

func NewLexer(input string) *Lexer {
	l := &Lexer{source: input, line: 1, column: 0}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.pos >= len(l.source) {
		l.ch = 0
	} else {
		l.ch = rune(l.source[l.pos])
	}
	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
	l.pos++
}

func (l *Lexer) Tokenize() []Token {
	tokens := []Token{}
	for {
		tok := l.nextToken()
		tokens = append(tokens, tok)
		if tok.Type == EOF {
			break
		}
	}
	return tokens
}

func (l *Lexer) nextToken() Token {
	l.skipWhitespace()

	tok := Token{Line: l.line, Column: l.column}

	switch {
	case l.ch == 0:
		tok.Type = EOF
		tok.Value = ""
	case l.ch == '+':
		tok.Type = PLUS
		tok.Value = "+"
		l.readChar()
	case l.ch == '-':
		tok.Type = MINUS
		tok.Value = "-"
		l.readChar()
	case l.ch == '*':
		tok.Type = MULTIPLY
		tok.Value = "*"
		l.readChar()
	case l.ch == '/':
		tok.Type = DIVIDE
		tok.Value = "/"
		l.readChar()
	case l.ch == '\n':
		tok.Type = NEWLINE
		tok.Value = "\n"
		l.readChar()
	case l.ch == '"':
		tok.Type = STRING
		tok.Value = l.readString()
	case unicode.IsLetter(l.ch):
		startCol := l.column
		ident := l.readIdentifier()
		tok.Value = ident
		tok.Column = startCol
		tok.Type = lookupIdentOrKeyword(ident)
	case unicode.IsDigit(l.ch):
		num, isFloat := l.readNumber()
		if isFloat {
			tok.Type = FLOAT
		} else {
			tok.Type = NUMBER
		}
		tok.Value = num
	default:
		tok.Type = EOF
		l.readChar()
	}
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	start := l.pos - 1
	for unicode.IsLetter(l.ch) || unicode.IsDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.source[start : l.pos-1]
}

func (l *Lexer) readNumber() (string, bool) {
	start := l.pos - 1
	isFloat := false
	for unicode.IsDigit(l.ch) {
		l.readChar()
	}
	if l.ch == '.' {
		isFloat = true
		l.readChar()
		for unicode.IsDigit(l.ch) {
			l.readChar()
		}
	}
	return l.source[start : l.pos-1], isFloat
}

func (l *Lexer) readString() string {
	l.readChar() // skip opening quote
	start := l.pos - 1
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	str := l.source[start : l.pos-1]
	l.readChar() // skip closing quote
	return str
}

func lookupIdentOrKeyword(ident string) TokenType {
	switch ident {
	case "Pikachu":
		return PIKACHU
	case "Psyduck":
		return PSYDUCK
	case "Eevee":
		return EEVEE
	case "is":
		return IS
	case "evolves":
		return EVOLVES_TO
	case "catch":
		return CATCH
	case "release":
		return RELEASE
	case "from":
		return FROM
	case "wildgrass":
		return WILDGRASS
	default:
		return IDENTIFIER
	}
}
