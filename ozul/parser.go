package main

import (
	"fmt"
)

type Parser struct {
	tokens []Token
	pos    int
	cur    Token
	errors []string
}

func NewParser(tokens []Token) *Parser {
	p := &Parser{tokens: tokens, pos: 0}
	if len(tokens) > 0 {
		p.cur = tokens[0]
	}
	return p
}

func (p *Parser) Parse() *Program {
	program := &Program{Statements: []Statement{}}
	loopCount := 0
	for p.cur.Type != EOF {
		loopCount++
		if loopCount > 1000 {
			break
		}
		if p.cur.Type == NEWLINE {
			p.nextToken()
			continue
		}
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		} else {
			p.nextToken() // Prevent infinite loop on error
		}
		// Skip any trailing newlines after a statement
		for p.cur.Type == NEWLINE {
			p.nextToken()
		}
	}
	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.cur.Type {
	case PIKACHU, PSYDUCK, EEVEE:
		return p.parseDeclaration()
	case IDENTIFIER:
		if p.peek().Type == EVOLVES_TO {
			return p.parseAssignment()
		}
		return p.parseExpressionStatement()
	case RELEASE:
		return p.parseRelease()
	case CATCH:
		return p.parseCatch()
	case NEWLINE:
		return nil
	case EOF:
		return nil
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseDeclaration() Statement {
	pokemonType := p.cur.Value
	p.nextToken() // consume Pokemon type

	if p.cur.Type != IDENTIFIER {
		p.addError("expected identifier after Pokemon type")
		p.nextToken()
		return nil
	}
	name := p.cur.Value
	p.nextToken() // consume identifier

	if p.cur.Type != IS {
		p.addError("expected 'is' after identifier")
		p.nextToken()
		return nil
	}
	p.nextToken() // consume 'is'

	value := p.parseExpression(0)

	return &DeclarationStmt{
		PokemonType: pokemonType,
		Name:        name,
		Value:       value,
	}
}

func (p *Parser) parseAssignment() Statement {
	name := p.cur.Value
	p.nextToken() // consume identifier

	if p.cur.Type != EVOLVES_TO {
		p.addError("expected 'evolves' after identifier")
		p.nextToken()
		return nil
	}
	p.nextToken() // consume 'evolves'

	if p.cur.Type != IDENTIFIER || p.cur.Value != "to" {
		p.addError("expected 'to' after 'evolves'")
		p.nextToken()
		return nil
	}
	p.nextToken() // consume 'to'

	value := p.parseExpression(0)

	return &AssignmentStmt{
		Name:  name,
		Value: value,
	}
}

func (p *Parser) parseRelease() Statement {
	p.nextToken() // consume 'release'
	value := p.parseExpression(0)

	return &ReleaseStmt{Value: value}
}

func (p *Parser) parseCatch() Statement {
	p.nextToken() // consume 'catch'

	if p.cur.Type != IDENTIFIER {
		p.addError("expected identifier after 'catch'")
		p.nextToken()
		return nil
	}
	variable := p.cur.Value
	p.nextToken() // consume identifier

	if p.cur.Type != FROM {
		p.addError("expected 'from' after identifier")
		p.nextToken()
		return nil
	}
	p.nextToken() // consume 'from'

	if p.cur.Type != TRAINER {
		p.addError("expected 'trainer' after 'from'")
		p.nextToken()
		return nil
	}
	p.nextToken() // consume 'trainer'

	return &CatchStmt{Variable: variable}
}

func (p *Parser) parseExpressionStatement() Statement {
	expr := p.parseExpression(0)
	return &ReleaseStmt{Value: expr} // Treat bare expressions as release statements
}

func (p *Parser) parseExpression(precedence int) Expression {
	left := p.parsePrimary()

	for precedence < p.getPrecedence(p.cur.Type) {
		if !p.isOperator(p.cur.Type) {
			break
		}

		operator := p.cur.Value
		p.nextToken()

		right := p.parseExpression(p.getPrecedence(TokenType(0))) // Use current operator precedence

		left = &BinaryExpr{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left
}

func (p *Parser) parsePrimary() Expression {
	switch p.cur.Type {
	case NUMBER:
		value := 0
		fmt.Sscanf(p.cur.Value, "%d", &value)
		p.nextToken()
		return &NumberLiteral{Value: value}
	case FLOAT:
		value := 0.0
		fmt.Sscanf(p.cur.Value, "%f", &value)
		p.nextToken()
		return &FloatLiteral{Value: value}
	case STRING:
		value := p.cur.Value
		p.nextToken()
		return &StringLiteral{Value: value}
	case IDENTIFIER:
		name := p.cur.Value
		p.nextToken()
		return &Identifier{Name: name}
	default:
		p.addError("unexpected token: " + p.cur.Value)
		p.nextToken()
		return nil
	}
}

func (p *Parser) nextToken() {
	p.pos++
	if p.pos < len(p.tokens) {
		p.cur = p.tokens[p.pos]
	} else {
		p.cur = Token{Type: EOF}
	}
}

func (p *Parser) peek() Token {
	if p.pos+1 < len(p.tokens) {
		return p.tokens[p.pos+1]
	}
	return Token{Type: EOF}
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Parser) getPrecedence(tokType TokenType) int {
	switch tokType {
	case MULTIPLY, DIVIDE:
		return 2
	case PLUS, MINUS:
		return 1
	default:
		return 0
	}
}

func (p *Parser) isOperator(tokType TokenType) bool {
	return tokType == PLUS || tokType == MINUS || tokType == MULTIPLY || tokType == DIVIDE
}

func (p *Parser) Errors() []string {
	return p.errors
}
