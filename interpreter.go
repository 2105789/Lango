package main

import (
	"fmt"
	"strconv"
)

type Interpreter struct {
	environment *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		environment: NewEnvironment(nil),
	}
}

func (i *Interpreter) Interpret(statements []Stmt) {
	for _, statement := range statements {
		_, err := i.execute(statement)
		if err != nil {
			runtimeError(err)
		}
	}
}

func (i *Interpreter) execute(stmt Stmt) (interface{}, error) {
	return stmt.Accept(i)
}

func (i *Interpreter) VisitBlockStmt(stmt *Block) (interface{}, error) {
	i.executeBlock(stmt.Statements, NewEnvironment(i.environment))
	return nil, nil
}

func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) error {
	previous := i.environment
	defer func() { i.environment = previous }()
	i.environment = environment
	for _, statement := range statements {
		_, err := i.execute(statement)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) VisitExpressionStmt(stmt *Expression) (interface{}, error) {
	return i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitForStmt(stmt *For) (interface{}, error) {
	if stmt.Initializer != nil {
		_, err := i.execute(stmt.Initializer)
		if err != nil {
			return nil, err
		}
	}

	for {
		if stmt.Condition != nil {
			cond, err := i.evaluate(stmt.Condition)
			if err != nil {
				return nil, err
			}
			if !i.isTruthy(cond) {
				break
			}
		}

		_, err := i.execute(stmt.Body)
		if err != nil {
			return nil, err
		}

		if stmt.Increment != nil {
			_, err := i.evaluate(stmt.Increment)
			if err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

func (i *Interpreter) VisitIfStmt(stmt *If) (interface{}, error) {
	cond, err := i.evaluate(stmt.Condition)
	if err != nil {
		return nil, err
	}

	if i.isTruthy(cond) {
		return i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return i.execute(stmt.ElseBranch)
	}
	return nil, nil
}

func (i *Interpreter) VisitLiteralExpr(expr *Literal) (interface{}, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitPrintStmt(stmt *Print) (interface{}, error) {
	value, err := i.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}
	fmt.Println(i.stringify(value))
	return nil, nil
}

func (i *Interpreter) VisitUnaryExpr(expr *Unary) (interface{}, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case MINUS:
		if value, ok := right.(float64); ok {
			return -value, nil
		}
		return nil, i.error(expr.Operator, "Operand must be a number.")
	case BANG:
		return !i.isTruthy(right), nil
	}

	return nil, i.error(expr.Operator, "Unexpected unary operator.")
}

func (i *Interpreter) VisitVarStmt(stmt *Var) (interface{}, error) {
	if _, ok := i.environment.values[stmt.Name.Lexeme]; !ok {
		i.environment.Define(stmt.Name.Lexeme, nil)
	}

	if stmt.Initializer != nil {
		value, err := i.evaluate(stmt.Initializer)
		if err != nil {
			return nil, err
		}
		i.environment.Assign(stmt.Name, value)
	}
	return nil, nil
}

func (i *Interpreter) VisitWhileStmt(stmt *While) (interface{}, error) {
	for {
		cond, err := i.evaluate(stmt.Condition)
		if err != nil {
			return nil, err
		}
		if !i.isTruthy(cond) {
			break
		}
		_, err = i.execute(stmt.Body)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *Interpreter) VisitBinaryExpr(expr *Binary) (interface{}, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case MINUS:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l - r, nil
			}
		}
		return nil, i.error(expr.Operator, "Operands must be numbers.")
	case SLASH:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				if r == 0 {
					return nil, i.error(expr.Operator, "Division by zero.")
				}
				return l / r, nil
			}
		}
		return nil, i.error(expr.Operator, "Operands must be numbers.")
	case STAR:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l * r, nil
			}
		}
		return nil, i.error(expr.Operator, "Operands must be numbers.")
	case MOD:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				if r == 0 {
					return nil, i.error(expr.Operator, "Modulo by zero.")
				}
				return float64(int(l) % int(r)), nil
			}
		}
		return nil, i.error(expr.Operator, "Operands of modulo must be numbers.")
	case PLUS:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l + r, nil
			}
		}
		return nil, i.error(expr.Operator, "Operands must be two numbers or two strings.")
	case GREATER:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l > r, nil
			}
		}
		return nil, i.error(expr.Operator, "Operands must be numbers.")
	case GREATER_EQUAL:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l >= r, nil
			}
		}
		return nil, i.error(expr.Operator, "Operands must be numbers.")
	case LESS:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l < r, nil
			}
		}
		return nil, i.error(expr.Operator, "Operands must be numbers.")
	case LESS_EQUAL:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l <= r, nil
			}
		}
		return nil, i.error(expr.Operator, "Operands must be numbers.")
	case BANG_EQUAL:
		return !i.isEqual(left, right), nil
	case EQUAL_EQUAL:
		return i.isEqual(left, right), nil
	}

	return nil, i.error(expr.Operator, "Unexpected binary operator.")
}

func (i *Interpreter) VisitGroupingExpr(expr *Grouping) (interface{}, error) {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitAssignExpr(expr *Assign) (interface{}, error) {
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	err = i.environment.Assign(expr.Name, value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (i *Interpreter) evaluate(expr Expr) (interface{}, error) {
	return expr.Accept(i)
}

func (i *Interpreter) isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	if b, ok := object.(bool); ok {
		return b
	}
	return true
}

func (i *Interpreter) isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}

func (i *Interpreter) stringify(object interface{}) string {
	if object == nil {
		return "nil"
	}
	if f, ok := object.(float64); ok {
		return strconv.FormatFloat(f, 'f', -1, 64)
	}
	return fmt.Sprintf("%v", object)
}

func (i *Interpreter) error(token *Token, message string) error {
	return fmt.Errorf("Runtime error at '%s': %s", token.Lexeme, message)
}
