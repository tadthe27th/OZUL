package main

import (
	"testing"
)

func TestLexer_Declaration(t *testing.T) {
	source := `Pikachu health is 100`
	lexer := NewLexer(source)
	tokens := lexer.Tokenize()

	expected := []TokenType{PIKACHU, IDENTIFIER, IS, NUMBER, EOF}
	for i, tokType := range expected {
		if tokens[i].Type != tokType {
			t.Errorf("token %d: expected %v, got %v", i, tokType, tokens[i].Type)
		}
	}
}

func TestLexer_Assignment(t *testing.T) {
	source := `health evolves to 150`
	lexer := NewLexer(source)
	tokens := lexer.Tokenize()

	expected := []TokenType{IDENTIFIER, EVOLVES_TO, IDENTIFIER, NUMBER, EOF}
	for i, tokType := range expected {
		if tokens[i].Type != tokType {
			t.Errorf("token %d: expected %v, got %v", i, tokType, tokens[i].Type)
		}
	}
}

func TestLexer_Arithmetic(t *testing.T) {
	source := `Pikachu sum is 10 + 5 * 2`
	lexer := NewLexer(source)
	tokens := lexer.Tokenize()

	expected := []TokenType{PIKACHU, IDENTIFIER, IS, NUMBER, PLUS, NUMBER, MULTIPLY, NUMBER, EOF}
	for i, tokType := range expected {
		if tokens[i].Type != tokType {
			t.Errorf("token %d: expected %v, got %v", i, tokType, tokens[i].Type)
		}
	}
}

func TestLexer_String(t *testing.T) {
	source := `Eevee greeting is "Hello " + "Ash"`
	lexer := NewLexer(source)
	tokens := lexer.Tokenize()

	expected := []TokenType{EEVEE, IDENTIFIER, IS, STRING, PLUS, STRING, EOF}
	for i, tokType := range expected {
		if tokens[i].Type != tokType {
			t.Errorf("token %d: expected %v, got %v", i, tokType, tokens[i].Type)
		}
	}
}

func TestLexer_IO(t *testing.T) {
	source := `catch userInput from trainer
release userInput`
	lexer := NewLexer(source)
	tokens := lexer.Tokenize()

	expected := []TokenType{CATCH, IDENTIFIER, FROM, TRAINER, NEWLINE, RELEASE, IDENTIFIER, EOF}
	for i, tokType := range expected {
		if tokens[i].Type != tokType {
			t.Errorf("token %d: expected %v, got %v", i, tokType, tokens[i].Type)
		}
	}
}
