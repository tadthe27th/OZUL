package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestCodeGen_Declaration(t *testing.T) {
	// Test Pikachu (int) declaration
	program := &Program{
		Statements: []Statement{
			&DeclarationStmt{
				Name:        "x",
				PokemonType: "Pikachu",
				Value:       &NumberLiteral{Value: 42},
			},
		},
	}

	cg := NewCodeGen()
	cg.GenerateProgram(program)
	code := cg.GetCode()

	if !strings.Contains(code, "int x = 42;") {
		t.Errorf("Expected 'int x = 42;' in generated code, got: %s", code)
	}
}

func TestCodeGen_PsyduckDeclaration(t *testing.T) {
	// Test Psyduck (float) declaration
	program := &Program{
		Statements: []Statement{
			&DeclarationStmt{
				Name:        "y",
				PokemonType: "Psyduck",
				Value:       &FloatLiteral{Value: 3.14},
			},
		},
	}

	cg := NewCodeGen()
	cg.GenerateProgram(program)
	code := cg.GetCode()

	if !strings.Contains(code, "double y = 3.140000;") {
		t.Errorf("Expected 'double y = 3.140000;' in generated code, got: %s", code)
	}
}

func TestCodeGen_EeveeDeclaration(t *testing.T) {
	// Test Eevee (string) declaration
	program := &Program{
		Statements: []Statement{
			&DeclarationStmt{
				Name:        "msg",
				PokemonType: "Eevee",
				Value:       &StringLiteral{Value: "Hello, Pokemon!"},
			},
		},
	}

	cg := NewCodeGen()
	cg.GenerateProgram(program)
	code := cg.GetCode()

	if !strings.Contains(code, "char* msg = \"Hello, Pokemon!\";") {
		t.Errorf("Expected 'char* msg = \"Hello, Pokemon!\";' in generated code, got: %s", code)
	}
}

func TestCodeGen_Assignment(t *testing.T) {
	// Test variable assignment
	program := &Program{
		Statements: []Statement{
			&DeclarationStmt{
				Name:        "x",
				PokemonType: "Pikachu",
				Value:       &NumberLiteral{Value: 10},
			},
			&AssignmentStmt{
				Name:  "x",
				Value: &NumberLiteral{Value: 20},
			},
		},
	}

	cg := NewCodeGen()
	cg.GenerateProgram(program)
	code := cg.GetCode()

	if !strings.Contains(code, "x = 20;") {
		t.Errorf("Expected 'x = 20;' in generated code, got: %s", code)
	}
}

func TestCodeGen_ArithmeticOperations(t *testing.T) {
	// Test arithmetic operations
	program := &Program{
		Statements: []Statement{
			&DeclarationStmt{
				Name:        "result",
				PokemonType: "Pikachu",
				Value: &BinaryExpr{
					Left:     &NumberLiteral{Value: 10},
					Operator: "+",
					Right:    &NumberLiteral{Value: 5},
				},
			},
		},
	}

	cg := NewCodeGen()
	cg.GenerateProgram(program)
	code := cg.GetCode()

	if !strings.Contains(code, "int result = (10 + 5);") {
		t.Errorf("Expected 'int result = (10 + 5);' in generated code, got: %s", code)
	}
}

func TestCodeGen_ComplexArithmetic(t *testing.T) {
	// Test complex arithmetic expression
	program := &Program{
		Statements: []Statement{
			&DeclarationStmt{
				Name:        "result",
				PokemonType: "Pikachu",
				Value: &BinaryExpr{
					Left: &BinaryExpr{
						Left:     &NumberLiteral{Value: 10},
						Operator: "*",
						Right:    &NumberLiteral{Value: 2},
					},
					Operator: "+",
					Right: &BinaryExpr{
						Left:     &NumberLiteral{Value: 5},
						Operator: "-",
						Right:    &NumberLiteral{Value: 3},
					},
				},
			},
		},
	}

	cg := NewCodeGen()
	cg.GenerateProgram(program)
	code := cg.GetCode()

	expected := "int result = ((10 * 2) + (5 - 3));"
	if !strings.Contains(code, expected) {
		t.Errorf("Expected '%s' in generated code, got: %s", expected, code)
	}
}

func TestCodeGen_StringConcatenation(t *testing.T) {
	// Test string concatenation
	program := &Program{
		Statements: []Statement{
			&DeclarationStmt{
				Name:        "greeting",
				PokemonType: "Eevee",
				Value: &BinaryExpr{
					Left:     &StringLiteral{Value: "Hello"},
					Operator: "+",
					Right:    &StringLiteral{Value: "World"},
				},
			},
		},
	}

	cg := NewCodeGen()
	cg.GenerateProgram(program)
	code := cg.GetCode()

	if !strings.Contains(code, "char str_buffer_0[256];") {
		t.Errorf("Expected buffer declaration for string concatenation, got: %s", code)
	}
	if !strings.Contains(code, "strcpy(str_buffer_0, \"Hello\");") {
		t.Errorf("Expected strcpy for string concatenation, got: %s", code)
	}
	if !strings.Contains(code, "strcat(str_buffer_0, \"World\");") {
		t.Errorf("Expected strcat for string concatenation, got: %s", code)
	}
	if !strings.Contains(code, "char* greeting = str_buffer_0;") {
		t.Errorf("Expected assignment of buffer to greeting, got: %s", code)
	}
}

func TestCodeGen_ReleaseInteger(t *testing.T) {
	// Test releasing an integer
	program := &Program{
		Statements: []Statement{
			&ReleaseStmt{
				Value: &NumberLiteral{Value: 42},
			},
		},
	}

	cg := NewCodeGen()
	cg.GenerateProgram(program)
	code := cg.GetCode()

	if !strings.Contains(code, "printf(\"%d\\n\", 42);") {
		t.Errorf("Expected 'printf(\"%%d\\n\", 42);' in generated code, got: %s", code)
	}
}

func TestCodeGen_ReleaseFloat(t *testing.T) {
	// Test releasing a float
	program := &Program{
		Statements: []Statement{
			&ReleaseStmt{
				Value: &FloatLiteral{Value: 3.14},
			},
		},
	}

	cg := NewCodeGen()
	cg.GenerateProgram(program)
	code := cg.GetCode()

	if !strings.Contains(code, "printf(\"%f\\n\", 3.140000);") {
		t.Errorf("Expected 'printf(\"%%f\\n\", 3.140000);' in generated code, got: %s", code)
	}
}

func TestCodeGen_ReleaseString(t *testing.T) {
	// Test releasing a string
	program := &Program{
		Statements: []Statement{
			&ReleaseStmt{
				Value: &StringLiteral{Value: "Hello, Pokemon!"},
			},
		},
	}

	cg := NewCodeGen()
	cg.GenerateProgram(program)
	code := cg.GetCode()

	if !strings.Contains(code, "printf(\"%s\\n\", \"Hello, Pokemon!\");") {
		t.Errorf("Expected 'printf(\"%%s\\n\", \"Hello, Pokemon!\");' in generated code, got: %s", code)
	}
}

func TestCodeGen_ReleaseVariable(t *testing.T) {
	// Test releasing a variable
	program := &Program{
		Statements: []Statement{
			&DeclarationStmt{
				Name:        "x",
				PokemonType: "Pikachu",
				Value:       &NumberLiteral{Value: 42},
			},
			&ReleaseStmt{
				Value: &Identifier{Name: "x"},
			},
		},
	}

	cg := NewCodeGen()
	cg.GenerateProgram(program)
	code := cg.GetCode()

	if !strings.Contains(code, "printf(\"%d\\n\", x);") {
		t.Errorf("Expected 'printf(\"%%d\\n\", x);' in generated code, got: %s", code)
	}
}

func TestCodeGen_Catch(t *testing.T) {
	// Test catch statement
	program := &Program{
		Statements: []Statement{
			&CatchStmt{
				Variable: "input",
			},
		},
	}

	cg := NewCodeGen()
	cg.GenerateProgram(program)
	code := cg.GetCode()

	if !strings.Contains(code, "int input;") {
		t.Errorf("Expected 'int input;' in generated code, got: %s", code)
	}
	if !strings.Contains(code, "scanf(\"%d\", &input);") {
		t.Errorf("Expected 'scanf(\"%%d\", &input);' in generated code, got: %s", code)
	}
}

func TestCodeGen_CompleteProgram(t *testing.T) {
	// Test a complete program with multiple statements
	program := &Program{
		Statements: []Statement{
			&DeclarationStmt{
				Name:        "x",
				PokemonType: "Pikachu",
				Value:       &NumberLiteral{Value: 10},
			},
			&DeclarationStmt{
				Name:        "y",
				PokemonType: "Psyduck",
				Value:       &FloatLiteral{Value: 3.14},
			},
			&DeclarationStmt{
				Name:        "msg",
				PokemonType: "Eevee",
				Value:       &StringLiteral{Value: "Hello"},
			},
			&AssignmentStmt{
				Name: "x",
				Value: &BinaryExpr{
					Left:     &Identifier{Name: "x"},
					Operator: "*",
					Right:    &NumberLiteral{Value: 2},
				},
			},
			&ReleaseStmt{
				Value: &Identifier{Name: "x"},
			},
			&ReleaseStmt{
				Value: &Identifier{Name: "y"},
			},
			&ReleaseStmt{
				Value: &Identifier{Name: "msg"},
			},
		},
	}

	cg := NewCodeGen()
	cg.GenerateProgram(program)
	code := cg.GetCode()

	// Check for required elements
	expectedElements := []string{
		"#include <stdio.h>",
		"#include <stdlib.h>",
		"#include <string.h>",
		"int main() {",
		"int x = 10;",
		"double y = 3.140000;",
		"char* msg = \"Hello\";",
		"x = (x * 2);",
		"printf(\"%d\\n\", x);",
		"printf(\"%f\\n\", y);",
		"printf(\"%s\\n\", msg);",
		"return 0;",
		"}",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(code, expected) {
			t.Errorf("Expected '%s' in generated code, got: %s", expected, code)
		}
	}
}

func TestCodeGen_ErrorHandling(t *testing.T) {
	// Test error handling for undefined variable
	program := &Program{
		Statements: []Statement{
			&ReleaseStmt{
				Value: &Identifier{Name: "undefined_var"},
			},
		},
	}

	cg := NewCodeGen()

	// Suppress error output during this test
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	t.Logf("About to call GenerateProgram and expect a panic...")
	defer func() {
		w.Close()
		os.Stderr = oldStderr
		io.Copy(io.Discard, r) // Drain any output
		r.Close()
		t.Logf("In defer after GenerateProgram")
		if r := recover(); r == nil {
			t.Errorf("Expected panic for undefined variable, but none occurred")
		} else {
			t.Logf("Recovered from panic: %v", r)
		}
	}()

	cg.GenerateProgram(program)
	t.Logf("GenerateProgram returned (should not reach here)")
}

func TestCodeGen_AssignmentToUndefinedVariable(t *testing.T) {
	// Test error handling for assignment to undefined variable
	program := &Program{
		Statements: []Statement{
			&AssignmentStmt{
				Name:  "undefined_var",
				Value: &NumberLiteral{Value: 42},
			},
		},
	}

	cg := NewCodeGen()

	// This should panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for assignment to undefined variable, but none occurred")
		}
	}()

	cg.GenerateProgram(program)
}
