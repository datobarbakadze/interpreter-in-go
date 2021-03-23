package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"interpreter/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > OR <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	lexer     *lexer.Lexer // lexer/tokenizer
	curToken  token.Token  // current token
	peekToken token.Token  // token after curToken
	errors    []string     // slice to register errors

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// New creates, new Parser.
// Registers lexer, and errors slice.
// Prceeds tokens with two steps:
// step 1: curToken = nil, peekToken = someToken1
// step 2: curToken = someToken1, peekToken = someToken2
func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer, errors: []string{}}
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)

	return p
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// Loops throug all of the tokens by using nextToken(),
// parses statements on each step via parseStatement(),
// Loops until it comes across end of file token: EOF
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// Proceeds curToken and peekToken with one step.
// by using lexer's NextToken() method
//
// step 1: curToken = nil, peekToken = someToken1
// step 2: curToken = someToken1, peekToken = someToken2
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

// Begin of parsing statements.
// Switch statement checks token types and associates
// relevant parsing function to them in case of matching.
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// Parse let statement
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	for !p.curTokIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// Parsing of return statement
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()
	for !p.curTokIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// Parse the whole statement of expressions
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// Parse basic expression
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()
	return leftExp
}

// Determines if curent token of the Parse(the one on which it
// stands on in the specific situatuation) is the same as the
// one passed to the function curTokIs.
//
// used just for checking
func (p *Parser) curTokIs(token token.TokenType) bool {
	return p.curToken.Type == token
}

// Almost same as an curTokIs() method but instead of checking
// current token it checks if token passed to the function is
// same as token that appears after current token.
func (p *Parser) peekTokenIs(token token.TokenType) bool {
	return p.peekToken.Type == token
}

// Proceed to the next token if the next token quals passed
// one. If not register the error in the slice of `errors`.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// Get errors slice.
func (p *Parser) Errors() []string {
	return p.errors
}

// Register the error.
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected newxt token to be %s. got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
