package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	fmt.Println("[DEBUG] Starting OZUL interpreter")
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
	debug := false
	generateC := false

	fmt.Printf("[DEBUG] Source file: %s\n", sourceFile)

	// Parse command line arguments
	for i := 2; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "-o" && i+1 < len(os.Args) {
			outputFile = os.Args[i+1]
			i++ // Skip next argument
		} else if arg == "-debug" {
			debug = true
		} else if arg == "-c" {
			generateC = true
		}
	}

	fmt.Println("[DEBUG] Reading source file...")
	source, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading source file: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("[DEBUG] Source file read successfully")

	// Lexing
	fmt.Println("[DEBUG] Starting lexing...")
	lexer := NewLexer(string(source))
	tokens := lexer.Tokenize()
	fmt.Println("[DEBUG] Lexing complete")

	// Debug: Print tokens only if debug flag is set
	if debug {
		fmt.Println("=== TOKENS ===")
		for _, token := range tokens {
			fmt.Printf("%s: %s\n", fmt.Sprint(token.Type), token.Value)
		}
	}

	// Parsing
	fmt.Println("[DEBUG] Starting parsing...")
	parser := NewParser(tokens)
	program := parser.Parse()
	fmt.Println("[DEBUG] Parsing complete")

	// Debug: Print AST only if debug flag is set
	if debug {
		fmt.Println("\n=== AST ===")
		for i, stmt := range program.Statements {
			fmt.Printf("Statement %d: %s\n", i, stmt.String())
		}
	}

	if generateC {
		// Code generation (C)
		fmt.Println("[DEBUG] Starting code generation...")
		codegen := NewCodeGen()
		codegen.GenerateProgram(program)
		cCode := codegen.GetCode()
		fmt.Println("[DEBUG] Code generation complete")
		if outputFile != "" {
			err := ioutil.WriteFile(outputFile, []byte(cCode), 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("C code written to %s\n", outputFile)
		} else {
			fmt.Println("\n=== GENERATED C CODE ===")
			fmt.Println(cCode)
		}
	} else {
		// Interpret and run the program directly
		fmt.Println("[DEBUG] Starting interpretation...")
		interpreter := NewInterpreter()
		interpreter.Run(program)
		fmt.Println("[DEBUG] Interpretation complete")
	}
}
