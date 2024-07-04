package main

import (
	"bytes"
	"fmt"
	"strings"
)

type AstPrinter struct{}

func (ap *AstPrinter) Print(expr Expr) (string, error) {
	result, err := expr.Accept(ap)
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

func (ap *AstPrinter) VisitAssignExpr(expr *Assign) (interface{}, error) {
	rightStr, _ := expr.Value.Accept(ap)
	return ap.parenthesize("=", &Literal{Value: expr.Name.Lexeme}, &Literal{Value: rightStr}), nil
}

func (ap *AstPrinter) VisitBinaryExpr(expr *Binary) (interface{}, error) {
	return ap.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right), nil
}

func (ap *AstPrinter) VisitBlockStmt(stmt *Block) (interface{}, error) {
	var buf bytes.Buffer
	buf.WriteString("(block")
	for _, s := range stmt.Statements {
		str, _ := s.Accept(ap)
		buf.WriteString(" ")
		buf.WriteString(str.(string))
	}
	buf.WriteString(")")
	return buf.String(), nil
}

func (ap *AstPrinter) VisitExpressionStmt(stmt *Expression) (interface{}, error) {
	return stmt.Expression.Accept(ap)
}

func (ap *AstPrinter) VisitForStmt(stmt *For) (interface{}, error) {
	var builder strings.Builder
	builder.WriteString("(for ")

	// 1. Initializer
	if stmt.Initializer != nil {
		initStr, _ := stmt.Initializer.Accept(ap)
		builder.WriteString(initStr.(string))
		builder.WriteString(" ")
	} else {
		builder.WriteString("; ") // For empty initializer
	}

	// 2. Condition
	if stmt.Condition != nil {
		condStr, _ := stmt.Condition.Accept(ap)
		builder.WriteString(condStr.(string))
	}
	builder.WriteString("; ")

	// 3. Increment
	if stmt.Increment != nil {
		incStr, _ := stmt.Increment.Accept(ap)
		builder.WriteString(incStr.(string))
	}

	builder.WriteString(") ")

	// 4. Body
	bodyStr, _ := stmt.Body.Accept(ap)
	builder.WriteString(bodyStr.(string))

	return builder.String(), nil
}

func (ap *AstPrinter) VisitGroupingExpr(expr *Grouping) (interface{}, error) {
	return ap.parenthesize("group", expr.Expression), nil
}

func (ap *AstPrinter) VisitIfStmt(stmt *If) (interface{}, error) {
	var buf bytes.Buffer
	buf.WriteString("(if ")
	condStr, _ := stmt.Condition.Accept(ap)
	buf.WriteString(condStr.(string))
	buf.WriteString(" ")
	thenStr, _ := stmt.ThenBranch.Accept(ap)
	buf.WriteString(thenStr.(string))
	if stmt.ElseBranch != nil {
		buf.WriteString(" else ")
		elseStr, _ := stmt.ElseBranch.Accept(ap)
		buf.WriteString(elseStr.(string))
	}
	buf.WriteString(")")
	return buf.String(), nil
}

func (ap *AstPrinter) VisitLiteralExpr(expr *Literal) (interface{}, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", expr.Value), nil
}

func (ap *AstPrinter) VisitPrintStmt(stmt *Print) (interface{}, error) {
	return ap.parenthesize("print", stmt.Expression), nil
}

func (ap *AstPrinter) VisitUnaryExpr(expr *Unary) (interface{}, error) {
	return ap.parenthesize(expr.Operator.Lexeme, expr.Right), nil
}

func (ap *AstPrinter) VisitVariableExpr(expr *Variable) (interface{}, error) {
	return expr.Name.Lexeme, nil
}

func (ap *AstPrinter) VisitWhileStmt(stmt *While) (interface{}, error) {
	var buf bytes.Buffer
	buf.WriteString("(while ")
	condStr, _ := stmt.Condition.Accept(ap)
	buf.WriteString(condStr.(string))
	buf.WriteString(" ")
	bodyStr, _ := stmt.Body.Accept(ap)
	buf.WriteString(bodyStr.(string))
	buf.WriteString(")")
	return buf.String(), nil
}

func (ap *AstPrinter) VisitVarStmt(stmt *Var) (interface{}, error) {
	if stmt.Initializer != nil {
		initStr, _ := stmt.Initializer.Accept(ap)
		return fmt.Sprintf("(var %s = %s)", stmt.Name.Lexeme, initStr), nil
	}
	return fmt.Sprintf("(var %s)", stmt.Name.Lexeme), nil
}

func (ap *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var buf bytes.Buffer
	buf.WriteString("(")
	buf.WriteString(name)
	for _, expr := range exprs {
		str, _ := expr.Accept(ap)
		buf.WriteString(" ")
		buf.WriteString(str.(string))
	}
	buf.WriteString(")")
	return buf.String()
}

func (i *Interpreter) VisitVariableExpr(expr *Variable) (interface{}, error) {
	return i.environment.Get(expr.Name)
}
