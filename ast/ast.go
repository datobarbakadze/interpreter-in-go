package ast

import "interpreter/token"

type Node interface {
	TokenLiteral() string
}

// Statement represents statements in AST
// for e.g let x = 5, 'let x' is a statement
type Statement interface {
	Node
	statementNode()
}

// Expression represents expressions in AST
// for e.g let x = 5, '5' is an expression
type Expression interface {
	Node
	expressionNode()
}

// Identifier stores value and the name of the identifiers
// like variables.
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// LetStatement stores and expresses `let` statement
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// Program is a root node in AST
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}
