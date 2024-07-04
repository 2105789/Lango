package main

import (
	"fmt"
)

type Parser struct {
	tokens  []*Token
	current int
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{tokens, 0}
}

func (p *Parser) Parse() ([]Stmt, error) {
	statements := []Stmt{}
	for !p.isAtEnd() {
		decl, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, decl)
	}
	return statements, nil
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match(VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() (Stmt, error) {
	name, err := p.consume(IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}
	var initializer Expr
	if p.match(EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}
	return &Var{Name: name, Initializer: initializer}, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(PRINT) {
		return p.printStatement()
	} else if p.match(IF) {
		return p.ifStatement()
	} else if p.match(WHILE) {
		return p.whileStatement()
	} else if p.match(FOR) {
		return p.forStatement()
	} else if p.match(LEFT_BRACE) {
		return p.blockStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}
	return &Print{Expression: value}, nil
}

func (p *Parser) blockStatement() (Stmt, error) {
	statements := []Stmt{}
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	if _, err := p.consume(RIGHT_BRACE, "Expect '}' after block."); err != nil {
		return nil, err
	}
	return &Block{Statements: statements}, nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(SEMICOLON, "Expect ';' after expression.")
	if err != nil {
		return nil, err
	}
	return &Expression{Expression: expr}, nil
}

func (p *Parser) ifStatement() (Stmt, error) {
	p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(RIGHT_PAREN, "Expect ')' after if condition.")
	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch Stmt
	if p.match(ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &If{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}, nil
}

func (p *Parser) whileStatement() (Stmt, error) {
	p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(RIGHT_PAREN, "Expect ')' after condition.")
	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &While{Condition: condition, Body: body}, nil
}

func (p *Parser) forStatement() (Stmt, error) {
	p.consume(LEFT_PAREN, "Expect '(' after 'for'.")
	var err error

	var initializer Stmt
	if !p.check(SEMICOLON) {
		if p.match(VAR) {
			initializer, err = p.varDeclaration()
		} else {
			initializer, err = p.expressionStatement()
		}
		if err != nil {
			return nil, err
		}
	}

	var condition Expr
	if !p.check(SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	p.consume(SEMICOLON, "Expect ';' after loop condition.")

	var increment Expr
	if !p.check(RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	p.consume(RIGHT_PAREN, "Expect ')' after loop clauses.")

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &For{
		Initializer: initializer,
		Condition:   condition,
		Increment:   increment,
		Body:        body,
	}, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.match(EQUAL) {
		equals := p.previous()
		value, err := p.equality()
		if err != nil {
			return nil, err
		}

		if varExpr, ok := expr.(*Variable); ok {
			return &Assign{Name: varExpr.Name, Value: value}, nil
		}

		return nil, p.error(equals, "Invalid assignment target.")
	}

	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &Binary{expr, operator, right}
	}
	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &Binary{expr, operator, right}
	}
	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &Binary{expr, operator, right}
	}
	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.match(SLASH, STAR, MOD) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &Binary{expr, operator, right}
	}
	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &Unary{operator, right}, nil
	}
	return p.primary()
}

func (p *Parser) bitwiseAnd() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}
	return expr, nil
}

func (p *Parser) primary() (Expr, error) {
	if p.match(IDENTIFIER) {
		return &Variable{Name: p.previous()}, nil
	}
	if p.match(FALSE) {
		return &Literal{Value: false}, nil
	}
	if p.match(TRUE) {
		return &Literal{Value: true}, nil
	}
	if p.match(NIL) {
		return &Literal{Value: nil}, nil
	}
	if p.match(NUMBER, STRING) {
		return &Literal{Value: p.previous().Literal}, nil
	}
	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(RIGHT_PAREN, "Expect ')' after expression."); err != nil {
			return nil, err
		}
		return &Grouping{Expression: expr}, nil
	}
	return nil, p.error(p.peek(), "Expect expression.")
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() *Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(tokenType TokenType, message string) (*Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	return nil, p.error(p.peek(), message)
}

func (p *Parser) error(token *Token, message string) error {
	parseError(token, message)
	return fmt.Errorf("parse error at '%s': %s", token.Lexeme, message)
}

func (p *Parser) synchronize() {
}
