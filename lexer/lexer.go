package lexer

import (
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

func (lexer *Lexer) peekChar() byte {
	if lexer.readPosition >= len(lexer.input) {
		return 0
	} else {
		return lexer.input[lexer.readPosition]
	}
}

// NextToken detects and returns new token
func (lexer *Lexer) NextToken() token.Token {
	var tok token.Token
	lexer.skipWhitespace()
	switch lexer.chr {
	case '=':
		if lexer.peekChar() == '=' {
			ch := lexer.chr
			lexer.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(lexer.chr)}
		} else {
			tok = newToken(token.ASSIGN, lexer.chr)
		}
	case '!':

		if lexer.peekChar() == '=' {
			ch := lexer.chr
			lexer.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(lexer.chr)}
		} else {
			tok = newToken(token.BANG, lexer.chr)
		}
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
	case '-':
		tok = newToken(token.MINUS, lexer.chr)

	case '/':
		tok = newToken(token.SLASH, lexer.chr)
	case '*':
		tok = newToken(token.ASTERISK, lexer.chr)
	case '<':
		tok = newToken(token.LT, lexer.chr)
	case '>':
		tok = newToken(token.GT, lexer.chr)
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
	for lexer.chr == ' ' || lexer.chr == '\t' || lexer.chr == '\r' || lexer.chr == '\n' {
		lexer.readChar()
	}
}

func (lexer *Lexer) readNumber() string {
	position := lexer.position
	for isDigit(lexer.chr) {
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
