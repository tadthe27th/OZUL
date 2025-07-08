package main

import "fmt"

// Token types for OZUL
//go:generate stringer -type=TokenType

type TokenType int

const (
	// Data types
	PIKACHU TokenType = iota // int
	PSYDUCK                  // float64
	EEVEE                    // string

	// Keywords
	IS         // =
	HAS_LEVEL // = (reassignment)
	CATCH      // input
	RELEASE    // print
	FROM       // from
	WILDGRASS    // wildgrass

	// Literals
	NUMBER     // 123
	FLOAT      // 3.14
	STRING     // "hello"
	IDENTIFIER // variable names

	// Operators
	PLUS     // +
	MINUS    // -
	MULTIPLY // *
	DIVIDE   // /

	// Delimiters
	NEWLINE
	EOF
)

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

// AST Node interfaces and structs

type ASTNode interface {
	String() string
}

type Program struct {
	Statements []Statement
}

type Statement interface {
	ASTNode
}

type Expression interface {
	ASTNode
}

// Variable declaration: "Pikachu health is 100"
type DeclarationStmt struct {
	PokemonType string     // "Pikachu", "Psyduck", or "Eevee"
	Name        string     // variable name
	Value       Expression // initial value
}

func (d *DeclarationStmt) String() string {
	return fmt.Sprintf("%s %s is %s", d.PokemonType, d.Name, d.Value.String())
}

// Assignment: "health evolves to 150"
type AssignmentStmt struct {
	Name  string
	Value Expression
}

func (a *AssignmentStmt) String() string {
	return fmt.Sprintf("%s has level %s", a.Name, a.Value.String())
}

// Output: "release health"
type ReleaseStmt struct {
	Value Expression
}

func (r *ReleaseStmt) String() string {
	return fmt.Sprintf("release %s", r.Value.String())
}

// Input: "catch userInput from wildgrass"
type CatchStmt struct {
	Variable string
}

func (c *CatchStmt) String() string {
	return fmt.Sprintf("catch %s from wildgrass", c.Variable)
}

// Binary operations: "10 + 5"
type BinaryExpr struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (b *BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Left.String(), b.Operator, b.Right.String())
}

// Literals
type NumberLiteral struct {
	Value int
}

func (n *NumberLiteral) String() string {
	return fmt.Sprintf("%d", n.Value)
}

type FloatLiteral struct {
	Value float64
}

func (f *FloatLiteral) String() string {
	return fmt.Sprintf("%f", f.Value)
}

type StringLiteral struct {
	Value string
}

func (s *StringLiteral) String() string {
	return fmt.Sprintf("\"%s\"", s.Value)
}

type Identifier struct {
	Name string
}

func (i *Identifier) String() string {
	return i.Name
}

// Pokemon-themed error type
type PokemonError struct {
	Message string
	Line    int
	Column  int
}

func (e PokemonError) Error() string {
	return fmt.Sprintf("Professor Oak says: %s (line %d)", e.Message, e.Line)
}
