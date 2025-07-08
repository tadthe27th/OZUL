package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ozul <source.ozul> [-c -o output.c] [-debug]")
		fmt.Println("  (no flags): interpret and run the OZUL program directly")
		fmt.Println("  -c: generate C code instead of running (advanced)")
		fmt.Println("  -o output.c: write C code to output file (with -c)")
		fmt.Println("  -debug: show debug info (tokens, AST)")
		os.Exit(1)
	}

	sourceFile := os.Args[1]
	var outputFile string
	generateC := false

	// Parse command line arguments
	for i := 2; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "-o" && i+1 < len(os.Args) {
			outputFile = os.Args[i+1]
			i++ // Skip next argument
		} else if arg == "-debug" {
			// debug = true // This line is removed as per the edit hint
		} else if arg == "-c" {
			generateC = true
		}
	}

	source, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Error reading source file: %v\n", err)
		os.Exit(1)
	}

	// Lexing
	lexer := NewLexer(string(source))
	tokens := lexer.Tokenize()

	// Parsing
	parser := NewParser(tokens)
	program := parser.Parse()

	if len(parser.Errors()) > 0 {
		fmt.Println("Parser errors:")
		for _, err := range parser.Errors() {
			fmt.Println("  ", err)
		}
		os.Exit(1)
	}

	if generateC {
		// Code generation (C)
		codegen := NewCodeGen()
		codegen.GenerateProgram(program)
		cCode := codegen.GetCode()
		if outputFile != "" {
			err := ioutil.WriteFile(outputFile, []byte(cCode), 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[ERROR] Error writing output file: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("C code written to %s\n", outputFile)
		} else {
			fmt.Println("\n=== GENERATED C CODE ===")
			fmt.Println(cCode)
		}
	} else {
		// Interpret and run the program directly
		interpreter := NewInterpreter()
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "[ERROR] Interpreter panic: %v\n", r)
			}
		}()
		interpreter.Run(program)
	}
}
