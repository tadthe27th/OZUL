package main

import (
	"testing"
)

func TestParser_Declaration(t *testing.T) {
	source := `Pikachu health is 100`
	lexer := NewLexer(source)
	tokens := lexer.Tokenize()

	parser := NewParser(tokens)
	program := parser.Parse()

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	decl, ok := program.Statements[0].(*DeclarationStmt)
	if !ok {
		t.Fatalf("Expected DeclarationStmt, got %T", program.Statements[0])
	}

	if decl.PokemonType != "Pikachu" {
		t.Errorf("Expected PokemonType 'Pikachu', got '%s'", decl.PokemonType)
	}
	if decl.Name != "health" {
		t.Errorf("Expected Name 'health', got '%s'", decl.Name)
	}

	num, ok := decl.Value.(*NumberLiteral)
	if !ok {
		t.Fatalf("Expected NumberLiteral, got %T", decl.Value)
	}
	if num.Value != 100 {
		t.Errorf("Expected Value 100, got %d", num.Value)
	}
}

func TestParser_Assignment(t *testing.T) {
	source := `health evolves to 150`
	lexer := NewLexer(source)
	tokens := lexer.Tokenize()

	parser := NewParser(tokens)
	program := parser.Parse()

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	assign, ok := program.Statements[0].(*AssignmentStmt)
	if !ok {
		t.Fatalf("Expected AssignmentStmt, got %T", program.Statements[0])
	}

	if assign.Name != "health" {
		t.Errorf("Expected Name 'health', got '%s'", assign.Name)
	}

	num, ok := assign.Value.(*NumberLiteral)
	if !ok {
		t.Fatalf("Expected NumberLiteral, got %T", assign.Value)
	}
	if num.Value != 150 {
		t.Errorf("Expected Value 150, got %d", num.Value)
	}
}

func TestParser_Arithmetic(t *testing.T) {
	source := `Pikachu sum is 10 + 5 * 2`
	lexer := NewLexer(source)
	tokens := lexer.Tokenize()

	parser := NewParser(tokens)
	program := parser.Parse()

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	decl, ok := program.Statements[0].(*DeclarationStmt)
	if !ok {
		t.Fatalf("Expected DeclarationStmt, got %T", program.Statements[0])
	}

	binary, ok := decl.Value.(*BinaryExpr)
	if !ok {
		t.Fatalf("Expected BinaryExpr, got %T", decl.Value)
	}

	if binary.Operator != "+" {
		t.Errorf("Expected operator '+', got '%s'", binary.Operator)
	}

	// Check left side: 10
	left, ok := binary.Left.(*NumberLiteral)
	if !ok {
		t.Fatalf("Expected NumberLiteral for left, got %T", binary.Left)
	}
	if left.Value != 10 {
		t.Errorf("Expected left value 10, got %d", left.Value)
	}

	// Check right side: 5 * 2
	right, ok := binary.Right.(*BinaryExpr)
	if !ok {
		t.Fatalf("Expected BinaryExpr for right, got %T", binary.Right)
	}
	if right.Operator != "*" {
		t.Errorf("Expected right operator '*', got '%s'", right.Operator)
	}
}

func TestParser_String(t *testing.T) {
	source := `Eevee greeting is "Hello " + "Ash"`
	lexer := NewLexer(source)
	tokens := lexer.Tokenize()

	parser := NewParser(tokens)
	program := parser.Parse()

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	decl, ok := program.Statements[0].(*DeclarationStmt)
	if !ok {
		t.Fatalf("Expected DeclarationStmt, got %T", program.Statements[0])
	}

	if decl.PokemonType != "Eevee" {
		t.Errorf("Expected PokemonType 'Eevee', got '%s'", decl.PokemonType)
	}

	binary, ok := decl.Value.(*BinaryExpr)
	if !ok {
		t.Fatalf("Expected BinaryExpr, got %T", decl.Value)
	}

	if binary.Operator != "+" {
		t.Errorf("Expected operator '+', got '%s'", binary.Operator)
	}

	left, ok := binary.Left.(*StringLiteral)
	if !ok {
		t.Fatalf("Expected StringLiteral for left, got %T", binary.Left)
	}
	if left.Value != "Hello " {
		t.Errorf("Expected left value 'Hello ', got '%s'", left.Value)
	}

	right, ok := binary.Right.(*StringLiteral)
	if !ok {
		t.Fatalf("Expected StringLiteral for right, got %T", binary.Right)
	}
	if right.Value != "Ash" {
		t.Errorf("Expected right value 'Ash', got '%s'", right.Value)
	}
}

func TestParser_IO(t *testing.T) {
	source := `catch userInput from trainer
release userInput`
	lexer := NewLexer(source)
	tokens := lexer.Tokenize()

	parser := NewParser(tokens)
	program := parser.Parse()

	if len(program.Statements) != 2 {
		t.Fatalf("Expected 2 statements, got %d", len(program.Statements))
	}

	// Check catch statement
	catch, ok := program.Statements[0].(*CatchStmt)
	if !ok {
		t.Fatalf("Expected CatchStmt, got %T", program.Statements[0])
	}
	if catch.Variable != "userInput" {
		t.Errorf("Expected variable 'userInput', got '%s'", catch.Variable)
	}

	// Check release statement
	release, ok := program.Statements[1].(*ReleaseStmt)
	if !ok {
		t.Fatalf("Expected ReleaseStmt, got %T", program.Statements[1])
	}

	ident, ok := release.Value.(*Identifier)
	if !ok {
		t.Fatalf("Expected Identifier, got %T", release.Value)
	}
	if ident.Name != "userInput" {
		t.Errorf("Expected identifier 'userInput', got '%s'", ident.Name)
	}
}
