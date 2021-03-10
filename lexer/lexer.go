package lexer

import (
	"fmt"
	"interpreter/token"
)

// Lexer is a type for storing input and iteration
// char by char
type Lexer struct {
	input        string
	position     int  // current char pos
	readPosition int  // after current char pos
	chr          byte // current char
}

func (lexer *Lexer) readChar() {
	if lexer.readPosition >= len(lexer.input) {
		lexer.chr = 0
	} else {
		lexer.chr = lexer.input[lexer.readPosition]
	}
	lexer.position = lexer.readPosition
	lexer.readPosition++
}

// NextToken detects and returns new token
func (lexer *Lexer) NextToken() token.Token {
	var tok token.Token
	lexer.skipWhitespace()
	switch lexer.chr {
	case '=':
		tok = newToken(token.ASSIGN, lexer.chr)
	case ';':
		tok = newToken(token.SEMICOLON, lexer.chr)
	case '(':
		tok = newToken(token.LPAREN, lexer.chr)
	case ')':
		tok = newToken(token.RPAREN, lexer.chr)
	case ',':
		tok = newToken(token.COMMA, lexer.chr)
	case '+':
		tok = newToken(token.PLUS, lexer.chr)
	case '{':
		tok = newToken(token.LBRACE, lexer.chr)
	case '}':
		tok = newToken(token.RBRACE, lexer.chr)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:

		if isLetter(lexer.chr) {
			tok.Literal = lexer.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(lexer.chr) {
			tok.Type = token.INT
			tok.Literal = lexer.readNumber()
			return tok
		} else {
			fmt.Println("Asfas")
			tok = newToken(token.ILLEGAL, lexer.chr)
		}
	}
	lexer.readChar()
	return tok
}
func (lexer *Lexer) readIdentifier() string {
	position := lexer.position
	for isLetter(lexer.chr) {
		lexer.readChar()
	}
	return lexer.input[position:lexer.position]
}

func (lexer *Lexer) skipWhitespace() {
	for lexer.chr == ' ' || lexer.chr == '\t' || lexer.chr == '\r' {
		lexer.readChar()
	}
}

func (lexer *Lexer) readNumber() string {
	position := lexer.position
	for isDigit(lexer.input[position]) {
		lexer.readChar()
	}
	return lexer.input[position:lexer.position]
}

func isDigit(chr byte) bool {
	return '0' <= chr && chr <= '9'
}

func newToken(tokenType token.TokenType, chr byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(chr)}
}

func isLetter(chr byte) bool {
	return 'a' <= chr && chr <= 'z' || 'A' <= chr && chr <= 'Z' || chr == '_'
}

// New creates new lexer
func New(input string) *Lexer {
	lexer := &Lexer{input: input}
	lexer.readChar()
	return lexer
}
