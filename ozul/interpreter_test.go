package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// Helper to capture output from a function
func captureOutput(f func()) string {
	origStdout := os.Stdout
	origStderr := os.Stderr

	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()

	os.Stdout = wOut
	os.Stderr = wErr

	// Run the function
	f()

	wOut.Close()
	wErr.Close()

	var bufOut, bufErr bytes.Buffer
	io.Copy(&bufOut, rOut)
	io.Copy(&bufErr, rErr)

	os.Stdout = origStdout
	os.Stderr = origStderr

	return bufOut.String() + bufErr.String()
}

func runInterpreterWithOutput(program *Program, input string) (string, int) {
	origStdin := os.Stdin
	inR, inW, _ := os.Pipe()
	inW.WriteString(input)
	inW.Close()
	os.Stdin = inR

	var output string
	output = captureOutput(func() {
		it := NewInterpreter()
		it.Run(program)
	})

	os.Stdin = origStdin
	return output, 0
}

func TestInterpreter_DeclarationAndRelease(t *testing.T) {
	source := `Pikachu x is 5
release x`
	lexer := NewLexer(source)
	tokens := lexer.Tokenize()
	parser := NewParser(tokens)
	program := parser.Parse()

	output, _ := runInterpreterWithOutput(program, "")
	if !strings.Contains(output, "5") {
		t.Errorf("Expected output to contain '5', got: %q", output)
	}
}

func TestInterpreter_Arithmetic(t *testing.T) {
	source := `Pikachu sum is 10 + 5 * 2
release sum`
	lexer := NewLexer(source)
	tokens := lexer.Tokenize()
	parser := NewParser(tokens)
	program := parser.Parse()

	output, _ := runInterpreterWithOutput(program, "")
	if !strings.Contains(output, "20") {
		t.Errorf("Expected output to contain '20', got: %q", output)
	}
}

func TestInterpreter_StringConcat(t *testing.T) {
	source := `Eevee greeting is "Hello " + "Ash"
release greeting`
	lexer := NewLexer(source)
	tokens := lexer.Tokenize()
	parser := NewParser(tokens)
	program := parser.Parse()

	output, _ := runInterpreterWithOutput(program, "")
	if !strings.Contains(output, "Hello Ash") {
		t.Errorf("Expected output to contain 'Hello Ash', got: %q", output)
	}
}

func TestInterpreter_CatchInput(t *testing.T) {
	source := `catch age from trainer
release age`
	lexer := NewLexer(source)
	tokens := lexer.Tokenize()
	parser := NewParser(tokens)
	program := parser.Parse()

	output, _ := runInterpreterWithOutput(program, "42\n")
	if !strings.Contains(output, "42") {
		t.Errorf("Expected output to contain '42', got: %q", output)
	}
}

func TestInterpreter_UndefinedVariableError(t *testing.T) {
	source := `release missingno`
	lexer := NewLexer(source)
	tokens := lexer.Tokenize()
	parser := NewParser(tokens)
	program := parser.Parse()

	var output string
	var code int
	func() {
		defer func() {
			if r := recover(); r != nil {
				output = output + r.(string)
				code = 1
			}
		}()
		output, _ = runInterpreterWithOutput(program, "")
	}()
	if code == 0 || !strings.Contains(output, "Undefined variable") {
		t.Errorf("Expected undefined variable error, got: %q", output)
	}
}

func TestInterpreter_DivisionByZeroError(t *testing.T) {
	source := `Pikachu x is 5 / 0
release x`
	lexer := NewLexer(source)
	tokens := lexer.Tokenize()
	parser := NewParser(tokens)
	program := parser.Parse()

	var output string
	var code int
	func() {
		defer func() {
			if r := recover(); r != nil {
				output = output + r.(string)
				code = 1
			}
		}()
		output, _ = runInterpreterWithOutput(program, "")
	}()
	if code == 0 || !strings.Contains(output, "Division by zero") {
		t.Errorf("Expected division by zero error, got: %q", output)
	}
}
