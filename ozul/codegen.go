package main

import (
	"fmt"
	"strings"
)

// CodeGen represents the code generator for OZUL (simplified version)
type CodeGen struct {
	// Symbol table for variables
	variables map[string]interface{}

	// Generated code (simplified representation)
	code []string
}

// NewCodeGen creates a new code generator
func NewCodeGen() *CodeGen {
	return &CodeGen{
		variables: make(map[string]interface{}),
		code:      []string{},
	}
}

// GenerateProgram generates code for the entire program
func (cg *CodeGen) GenerateProgram(program *Program) {
	cg.code = append(cg.code, "#include <stdio.h>")
	cg.code = append(cg.code, "#include <stdlib.h>")
	cg.code = append(cg.code, "#include <string.h>")
	cg.code = append(cg.code, "")
	cg.code = append(cg.code, "int main() {")

	// Generate code for each statement
	for _, stmt := range program.Statements {
		cg.generateStatement(stmt)
	}

	cg.code = append(cg.code, "    return 0;")
	cg.code = append(cg.code, "}")
}

// generateStatement generates code for a single statement
func (cg *CodeGen) generateStatement(stmt Statement) {
	switch s := stmt.(type) {
	case *DeclarationStmt:
		cg.generateDeclaration(s)
	case *AssignmentStmt:
		cg.generateAssignment(s)
	case *ReleaseStmt:
		cg.generateRelease(s)
	case *CatchStmt:
		cg.generateCatch(s)
	}
}

// generateDeclaration generates code for variable declarations
func (cg *CodeGen) generateDeclaration(stmt *DeclarationStmt) {
	value := cg.generateExpression(stmt.Value)

	switch stmt.PokemonType {
	case "Pikachu":
		cg.code = append(cg.code, fmt.Sprintf("    int %s = %s;", stmt.Name, value))
		cg.variables[stmt.Name] = "int"
	case "Psyduck":
		cg.code = append(cg.code, fmt.Sprintf("    double %s = %s;", stmt.Name, value))
		cg.variables[stmt.Name] = "double"
	case "Eevee":
		cg.code = append(cg.code, fmt.Sprintf("    char* %s = %s;", stmt.Name, value))
		cg.variables[stmt.Name] = "char*"
	default:
		panic(fmt.Sprintf("[OZUL CodeGen Error] Unknown Pokemon type: %s", stmt.PokemonType))
	}
}

// generateAssignment generates code for variable assignments
func (cg *CodeGen) generateAssignment(stmt *AssignmentStmt) {
	value := cg.generateExpression(stmt.Value)
	if typ, exists := cg.variables[stmt.Name]; exists {
		cg.code = append(cg.code, fmt.Sprintf("    %s = %s;", stmt.Name, value))
		cg.variables[stmt.Name] = typ
	} else {
		panic(fmt.Sprintf("[OZUL CodeGen Error] Variable %s not declared!", stmt.Name))
	}
}

// generateRelease generates code for output statements
func (cg *CodeGen) generateRelease(stmt *ReleaseStmt) {
	value := cg.generateExpression(stmt.Value)
	var varType string
	if id, ok := stmt.Value.(*Identifier); ok {
		t, exists := cg.variables[id.Name]
		if !exists {
			panic(fmt.Sprintf("[OZUL CodeGen Error] Variable %s not declared!", id.Name))
		}
		varType = t.(string)
	} else if _, ok := stmt.Value.(*NumberLiteral); ok {
		varType = "int"
	} else if _, ok := stmt.Value.(*FloatLiteral); ok {
		varType = "double"
	} else if _, ok := stmt.Value.(*StringLiteral); ok {
		varType = "char*"
	} else {
		// Try to infer type from expression
		varType = "int"
	}

	switch varType {
	case "int":
		cg.code = append(cg.code, fmt.Sprintf("    printf(\"%%d\\n\", %s);", value))
	case "double":
		cg.code = append(cg.code, fmt.Sprintf("    printf(\"%%f\\n\", %s);", value))
	case "char*":
		cg.code = append(cg.code, fmt.Sprintf("    printf(\"%%s\\n\", %s);", value))
	default:
		panic(fmt.Sprintf("[OZUL CodeGen Error] Unknown type for release: %s", varType))
	}
}

// generateCatch generates code for input statements
func (cg *CodeGen) generateCatch(stmt *CatchStmt) {
	cg.code = append(cg.code, fmt.Sprintf("    int %s;", stmt.Variable))
	cg.code = append(cg.code, fmt.Sprintf("    scanf(\"%%d\", &%s);", stmt.Variable))
	cg.variables[stmt.Variable] = "int"
}

// generateExpression generates code for expressions
func (cg *CodeGen) generateExpression(expr Expression) string {
	switch e := expr.(type) {
	case *NumberLiteral:
		return fmt.Sprintf("%d", e.Value)
	case *FloatLiteral:
		return fmt.Sprintf("%f", e.Value)
	case *StringLiteral:
		return fmt.Sprintf("\"%s\"", e.Value)
	case *Identifier:
		if _, exists := cg.variables[e.Name]; exists {
			return e.Name
		}
		panic(fmt.Sprintf("[OZUL CodeGen Error] Undefined variable: %s", e.Name))
	case *BinaryExpr:
		left := cg.generateExpression(e.Left)
		right := cg.generateExpression(e.Right)
		if strings.HasPrefix(left, "\"") || strings.HasPrefix(right, "\"") {
			// String concatenation - create a buffer
			bufferName := fmt.Sprintf("str_buffer_%d", len(cg.variables))
			cg.code = append(cg.code, fmt.Sprintf("    char %s[256];", bufferName))
			cg.code = append(cg.code, fmt.Sprintf("    strcpy(%s, %s);", bufferName, left))
			cg.code = append(cg.code, fmt.Sprintf("    strcat(%s, %s);", bufferName, right))
			return bufferName
		}
		return fmt.Sprintf("(%s %s %s)", left, e.Operator, right)
	default:
		panic("[OZUL CodeGen Error] Unknown expression type.")
	}
	return "" // unreachable
}

// GetCode returns the generated C code as a string
func (cg *CodeGen) GetCode() string {
	return strings.Join(cg.code, "\n")
}

// DumpCode outputs the generated code
func (cg *CodeGen) DumpCode() {
	fmt.Println(cg.GetCode())
}
