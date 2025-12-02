package main

import (
	"fmt"
)

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
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
	if p.match(FUN) {
		return p.function("function")
	}
	if p.match(CLASS) {
		return p.classDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() (Stmt, error) {
	nameToken, err := p.consume(IDENTIFIER, "Expect variable name.")
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
	return &Var{Name: &nameToken, Initializer: initializer}, nil
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
	} else if p.match(RETURN) {
		return p.returnStatement()
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
	expr, err := p.ternary()
	if err != nil {
		return nil, err
	}

	if p.match(EQUAL, PLUS_EQUAL, MINUS_EQUAL, STAR_EQUAL, SLASH_EQUAL, MOD_EQUAL) {
		operator := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if varExpr, ok := expr.(*Variable); ok {
			name := varExpr.Name
			if operator.Type == EQUAL {
				return &Assign{Name: name, Value: value}, nil
			}
			// Desugar compound assignment: a += b  ->  a = a + b
			// We need to create a Binary expression: a + b
			// But we need the 'a' expression again. 
			// Since 'expr' is the variable expression for 'a', we can reuse it?
			// Yes, AST nodes are immutable-ish, reusing is fine.
			// Operator mapping: PLUS_EQUAL -> PLUS
			var binOp TokenType
			switch operator.Type {
			case PLUS_EQUAL: binOp = PLUS
			case MINUS_EQUAL: binOp = MINUS
			case STAR_EQUAL: binOp = STAR
			case SLASH_EQUAL: binOp = SLASH
			case MOD_EQUAL: binOp = MOD
			}
			
			// Create a synthetic token for the binary operator
			opToken := Token{Type: binOp, Lexeme: operator.Lexeme[:1], Line: operator.Line}
			
			binaryExpr := &Binary{Left: expr, Operator: &opToken, Right: value}
			return &Assign{Name: name, Value: binaryExpr}, nil

		} else if getExpr, ok := expr.(*Get); ok {
			if operator.Type == EQUAL {
				return &Set{Object: getExpr.Object, Name: getExpr.Name, Value: value}, nil
			}
			// Desugar set: obj.prop += val -> obj.prop = obj.prop + val
			var binOp TokenType
			switch operator.Type {
			case PLUS_EQUAL: binOp = PLUS
			case MINUS_EQUAL: binOp = MINUS
			case STAR_EQUAL: binOp = STAR
			case SLASH_EQUAL: binOp = SLASH
			case MOD_EQUAL: binOp = MOD
			}
			opToken := Token{Type: binOp, Lexeme: operator.Lexeme[:1], Line: operator.Line}
			
			// We need to duplicate the Get expression for the read
			// This might evaluate Object twice if we are not careful?
			// In a tree-walk interpreter, yes, it will evaluate Object twice.
			// e.g. getObj().prop += 1  ->  getObj().prop = getObj().prop + 1
			// This is a known issue with simple desugaring. 
			// For now, we accept this limitation or we'd need a special AST node for CompoundSet.
			// Given the goal is "easy to write", this is acceptable for now.
			
			binaryExpr := &Binary{Left: expr, Operator: &opToken, Right: value}
			return &Set{Object: getExpr.Object, Name: getExpr.Name, Value: binaryExpr}, nil

		} else if indexExpr, ok := expr.(*ArrayIndex); ok {
			if operator.Type == EQUAL {
				return &ArrayAssign{Assignee: *indexExpr, Value: value}, nil
			}
			// Desugar array assign: arr[i] += val -> arr[i] = arr[i] + val
			var binOp TokenType
			switch operator.Type {
			case PLUS_EQUAL: binOp = PLUS
			case MINUS_EQUAL: binOp = MINUS
			case STAR_EQUAL: binOp = STAR
			case SLASH_EQUAL: binOp = SLASH
			case MOD_EQUAL: binOp = MOD
			}
			opToken := Token{Type: binOp, Lexeme: operator.Lexeme[:1], Line: operator.Line}
			
			binaryExpr := &Binary{Left: expr, Operator: &opToken, Right: value}
			return &ArrayAssign{Assignee: *indexExpr, Value: binaryExpr}, nil
		}

		return nil, p.error(operator, "Invalid assignment target.")
	}

	return expr, nil
}

func (p *Parser) ternary() (Expr, error) {
	expr, err := p.logic_or()
	if err != nil {
		return nil, err
	}

	if p.match(QUESTION) {
		thenBranch, err := p.expression() // Allow assignment in branches
		if err != nil {
			return nil, err
		}
		
		_, err = p.consume(COLON, "Expect ':' after then branch of ternary operator.")
		if err != nil {
			return nil, err
		}
		
		elseBranch, err := p.ternary() // Right associative? a ? b : c ? d : e -> a ? b : (c ? d : e)
		if err != nil {
			return nil, err
		}
		
		return &Ternary{Condition: expr, ThenBranch: thenBranch, ElseBranch: elseBranch}, nil
	}

	return expr, nil
}

func (p *Parser) logic_or() (Expr, error) {
	expr, err := p.logic_and()
	if err != nil {
		return nil, err
	}

	for p.match(OR) {
		operator := p.previous()
		right, err := p.logic_and()
		if err != nil {
			return nil, err
		}
		expr = &Logical{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) logic_and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = &Logical{Left: expr, Operator: operator, Right: right}
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
		expr = &Binary{expr, &operator, right}
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
		expr = &Binary{expr, &operator, right}
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
		expr = &Binary{expr, &operator, right}
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
		expr = &Binary{expr, &operator, right}
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
		return &Unary{&operator, right}, nil
	}
	return p.call()
}

func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(LEFT_PAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else if p.match(DOT) {
			name, err := p.consume(IDENTIFIER, "Expect property name after '.'.")
			if err != nil {
				return nil, err
			}
			expr = &Get{Object: expr, Name: name}
		} else if p.match(LEFT_BRACKET) {
			index, err := p.expression()
			if err != nil {
				return nil, err
			}
			bracket, err := p.consume(RIGHT_BRACKET, "Expect ']' after index.")
			if err != nil {
				return nil, err
			}
			expr = &ArrayIndex{Array: expr, Bracket: bracket, Index: index}
		} else if p.match(PLUS_PLUS, MINUS_MINUS) {
			// Postfix increment/decrement
			// Desugar: i++  ->  i = i + 1 (but return old value? No, we decided new value for simplicity)
			// Wait, if we desugar to i = i + 1, it returns the new value.
			// C/Java i++ returns OLD value. ++i returns NEW value.
			// If we want "easy to write", users might expect C behavior.
			// But implementing "return old value" via desugaring is hard without a temporary variable block.
			// For simplicity, we will make it behave like ++i (return new value).
			// Or we can just say it's a statement mostly.
			
			operator := p.previous()
			var binOp TokenType
			if operator.Type == PLUS_PLUS {
				binOp = PLUS
			} else {
				binOp = MINUS
			}
			
			// We need to check if 'expr' is a valid assignment target
			if varExpr, ok := expr.(*Variable); ok {
				// i = i + 1
				opToken := Token{Type: binOp, Lexeme: operator.Lexeme[:1], Line: operator.Line}
				one := &Literal{Value: float64(1)}
				binaryExpr := &Binary{Left: expr, Operator: &opToken, Right: one}
				expr = &Assign{Name: varExpr.Name, Value: binaryExpr}
			} else if getExpr, ok := expr.(*Get); ok {
				// obj.prop = obj.prop + 1
				opToken := Token{Type: binOp, Lexeme: operator.Lexeme[:1], Line: operator.Line}
				one := &Literal{Value: float64(1)}
				binaryExpr := &Binary{Left: expr, Operator: &opToken, Right: one}
				expr = &Set{Object: getExpr.Object, Name: getExpr.Name, Value: binaryExpr}
			} else if indexExpr, ok := expr.(*ArrayIndex); ok {
				// arr[i] = arr[i] + 1
				opToken := Token{Type: binOp, Lexeme: operator.Lexeme[:1], Line: operator.Line}
				one := &Literal{Value: float64(1)}
				binaryExpr := &Binary{Left: expr, Operator: &opToken, Right: one}
				expr = &ArrayAssign{Assignee: *indexExpr, Value: binaryExpr}
			} else {
				return nil, p.error(operator, "Invalid increment/decrement target.")
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(callee Expr) (Expr, error) {
	arguments := []Expr{}
	if !p.check(RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				p.error(p.peek(), "Can't have more than 255 arguments.")
			}
			arg, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, arg)

			if !p.match(COMMA) {
				break
			}
		}
	}

	paren, err := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return &Call{Callee: callee, Paren: paren, Arguments: arguments}, nil
}

func (p *Parser) primary() (Expr, error) {
	if p.match(IDENTIFIER) {
		nameToken := p.previous()
		return &Variable{Name: &nameToken}, nil
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
	if p.match(THIS) {
		return &This{Keyword: p.previous()}, nil
	}
	if p.match(LEFT_BRACKET) {
		bracket := p.previous()
		values := []Expr{}
		if !p.check(RIGHT_BRACKET) {
			for {
				valExpr, err := p.expression()
				if err != nil {
					return nil, err
				}
				values = append(values, valExpr)
				if !p.match(COMMA) {
					break
				}
				if p.check(RIGHT_BRACKET) {
					return nil, p.error(p.peek(), "Trailing comma in array literal not allowed.")
				}
			}
		}
		_, err := p.consume(RIGHT_BRACKET, "Expect ']' after array elements.")
		if err != nil {
			return nil, err
		}
		return &ArrayLiteral{Bracket: bracket, Values: values}, nil
	}
	if p.match(LEFT_BRACE) {
		brace := p.previous()
		keys := []Token{}
		values := []Expr{}

		// Parse key-value pairs
		if !p.check(RIGHT_BRACE) {
			for {
				// Key can be an identifier or string
				var key Token
				if p.match(IDENTIFIER, STRING) {
					key = p.previous()
				} else {
					return nil, p.error(p.peek(), "Expect property key (identifier or string).")
				}

				// Expect colon
				_, err := p.consume(COLON, "Expect ':' after property key.")
				if err != nil {
					return nil, err
				}

				// Parse value expression
				valueExpr, err := p.expression()
				if err != nil {
					return nil, err
				}

				keys = append(keys, key)
				values = append(values, valueExpr)

				// Check for comma
				if !p.match(COMMA) {
					break
				}
				// Allow trailing comma
				if p.check(RIGHT_BRACE) {
					break
				}
			}
		}

		_, err := p.consume(RIGHT_BRACE, "Expect '}' after object literal.")
		if err != nil {
			return nil, err
		}

		return &ObjectLiteral{Brace: brace, Keys: keys, Values: values}, nil
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

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(tokenType TokenType, message string) (Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	return Token{}, p.error(p.peek(), message)
}

func (p *Parser) error(token Token, message string) error {
	parseError(&token, message)
	if token.Type != EOF {
		return fmt.Errorf("parse error at '%s': %s", token.Lexeme, message)
	} else {
		return fmt.Errorf("parse error at end: %s", message)
	}
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == SEMICOLON {
			return
		}

		switch p.peek().Type {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) function(kind string) (Stmt, error) {
	nameToken, err := p.consume(IDENTIFIER, fmt.Sprintf("Expect %s name.", kind))
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LEFT_PAREN, fmt.Sprintf("Expect '(' after %s name.", kind))
	if err != nil {
		return nil, err
	}

	parameters := []Token{}
	if !p.check(RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				return nil, p.error(p.peek(), "Can't have more than 255 parameters.")
			}
			paramToken, err := p.consume(IDENTIFIER, "Expect parameter name.")
			if err != nil {
				return nil, err
			}
			parameters = append(parameters, paramToken)

			if !p.match(COMMA) {
				break
			}
		}
	}
	_, err = p.consume(RIGHT_PAREN, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LEFT_BRACE, fmt.Sprintf("Expect '{' before %s body.", kind))
	if err != nil {
		return nil, err
	}

	bodyStmt, err := p.blockStatement()
	if err != nil {
		return nil, err
	}

	block, ok := bodyStmt.(*Block)
	if !ok {
		prev := p.previous()
		return nil, p.error(prev, "Internal error parsing function body.")
	}

	return &Function{Name: nameToken, Params: parameters, Body: block.Statements}, nil
}

func (p *Parser) returnStatement() (Stmt, error) {
	keyword := p.previous()
	var value Expr = nil
	var err error

	if !p.check(SEMICOLON) {
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(SEMICOLON, "Expect ';' after return value.")
	if err != nil {
		return nil, err
	}

	return &Return{Keyword: keyword, Value: value}, nil
}

func (p *Parser) classDeclaration() (Stmt, error) {
	name, err := p.consume(IDENTIFIER, "Expect class name.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LEFT_BRACE, "Expect '{' before class body.")
	if err != nil {
		return nil, err
	}

	methods := []*Function{}
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		methodStmt, err := p.function("method")
		if err != nil {
			return nil, err
		}
		method, ok := methodStmt.(*Function)
		if !ok {
			prev := p.previous()
			return nil, p.error(prev, "Internal error parsing class method.")
		}
		methods = append(methods, method)
	}

	_, err = p.consume(RIGHT_BRACE, "Expect '}' after class body.")
	if err != nil {
		return nil, err
	}

	return &Class{Name: name, Methods: methods}, nil
}
