package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Value struct {
	Type  string // "int", "float", "string"
	Int   int
	Float float64
	Str   string
}

type Interpreter struct {
	vars map[string]Value
}

func NewInterpreter() *Interpreter {
	return &Interpreter{vars: make(map[string]Value)}
}

func (it *Interpreter) Run(program *Program) {
	const maxSteps = 10000
	steps := 0
	for _, stmt := range program.Statements {
		steps++
		if steps > maxSteps {
			panic("[OZUL Error] Execution step limit exceeded (possible infinite loop)")
		}
		it.execStatement(stmt)
	}
}

func (it *Interpreter) execStatement(stmt Statement) {
	switch s := stmt.(type) {
	case *DeclarationStmt:
		val := it.evalExpression(s.Value)
		it.vars[s.Name] = val
	case *AssignmentStmt:
		val := it.evalExpression(s.Value)
		if _, ok := it.vars[s.Name]; ok {
			it.vars[s.Name] = val
		} else {
			panic(fmt.Sprintf("[OZUL Error] Variable not declared: %s", s.Name))
		}
	case *ReleaseStmt:
		val := it.evalExpression(s.Value)
		it.printValue(val)
	case *CatchStmt:
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Enter value for %s: ", s.Variable)
		os.Stdout.Sync() // Flush output buffer so prompt is visible
		input, _ := reader.ReadString('\n')
		input = input[:len(input)-1] // remove newline
		// Default to int, could be improved
		if intVal, err := strconv.Atoi(input); err == nil {
			it.vars[s.Variable] = Value{Type: "int", Int: intVal}
		} else if floatVal, err := strconv.ParseFloat(input, 64); err == nil {
			it.vars[s.Variable] = Value{Type: "float", Float: floatVal}
		} else {
			it.vars[s.Variable] = Value{Type: "string", Str: input}
		}
	}
}

func (it *Interpreter) evalExpression(expr Expression) Value {
	switch e := expr.(type) {
	case *NumberLiteral:
		return Value{Type: "int", Int: e.Value}
	case *FloatLiteral:
		return Value{Type: "float", Float: e.Value}
	case *StringLiteral:
		return Value{Type: "string", Str: e.Value}
	case *Identifier:
		v, ok := it.vars[e.Name]
		if !ok {
			panic(fmt.Sprintf("[OZUL Error] Undefined variable: %s", e.Name))
		}
		return v
	case *BinaryExpr:
		left := it.evalExpression(e.Left)
		right := it.evalExpression(e.Right)
		if e.Operator == "+" && (left.Type == "string" || right.Type == "string") {
			return Value{Type: "string", Str: it.toString(left) + it.toString(right)}
		}
		if left.Type == "float" || right.Type == "float" {
			lf := it.toFloat(left)
			rf := it.toFloat(right)
			switch e.Operator {
			case "+":
				return Value{Type: "float", Float: lf + rf}
			case "-":
				return Value{Type: "float", Float: lf - rf}
			case "*":
				return Value{Type: "float", Float: lf * rf}
			case "/":
				if rf == 0 {
					panic("[OZUL Error] Division by zero.")
				}
				return Value{Type: "float", Float: lf / rf}
			}
		}
		li := it.toInt(left)
		ri := it.toInt(right)
		switch e.Operator {
		case "+":
			return Value{Type: "int", Int: li + ri}
		case "-":
			return Value{Type: "int", Int: li - ri}
		case "*":
			return Value{Type: "int", Int: li * ri}
		case "/":
			if ri == 0 {
				panic("[OZUL Error] Division by zero.")
			}
			return Value{Type: "int", Int: li / ri}
		}
		fmt.Fprintf(os.Stderr, "[OZUL Error] Unknown operator: %s\n", e.Operator)
		panic(fmt.Sprintf("[OZUL Error] Unknown operator: %s", e.Operator))
	default:
		fmt.Fprintln(os.Stderr, "[OZUL Error] Unknown expression type.")
		panic("[OZUL Error] Unknown expression type.")
	}
	return Value{}
}

func (it *Interpreter) printValue(val Value) {
	switch val.Type {
	case "int":
		fmt.Println(val.Int)
	case "float":
		fmt.Println(val.Float)
	case "string":
		fmt.Println(val.Str)
	}
}

func (it *Interpreter) toInt(val Value) int {
	if val.Type == "int" {
		return val.Int
	} else if val.Type == "float" {
		return int(val.Float)
	} else if val.Type == "string" {
		i, _ := strconv.Atoi(val.Str)
		return i
	}
	return 0
}

func (it *Interpreter) toFloat(val Value) float64 {
	if val.Type == "float" {
		return val.Float
	} else if val.Type == "int" {
		return float64(val.Int)
	} else if val.Type == "string" {
		f, _ := strconv.ParseFloat(val.Str, 64)
		return f
	}
	return 0.0
}

func (it *Interpreter) toString(val Value) string {
	switch val.Type {
	case "string":
		return val.Str
	case "int":
		return strconv.Itoa(val.Int)
	case "float":
		return fmt.Sprintf("%f", val.Float)
	}
	return ""
}
