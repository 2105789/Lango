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
	// Ensure the result is actually a string before returning
	strResult, ok := result.(string)
	if !ok {
		// This case should ideally not happen if all visitor methods return strings or errors
		return "", fmt.Errorf("ast printer internal error: expected string result, got %T", result)
	}
	return strResult, nil
}

func (ap *AstPrinter) VisitAssignExpr(expr *Assign) (interface{}, error) {
	// Visit the value expression to get its AST string representation
	valueStr, err := expr.Value.Accept(ap)
	if err != nil {
		return nil, err
	}
	// Parenthesize using the variable name and the *result* of visiting the value expression
	return ap.parenthesize(expr.Name.Lexeme+" =", &Literal{Value: valueStr}) // Represent assignment slightly differently
}

func (ap *AstPrinter) VisitBinaryExpr(expr *Binary) (interface{}, error) {
	return ap.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (ap *AstPrinter) VisitBlockStmt(stmt *Block) (interface{}, error) {
	var buf bytes.Buffer
	buf.WriteString("(block")
	for _, s := range stmt.Statements {
		str, err := s.Accept(ap) // Check error
		if err != nil {
			return nil, err
		}
		buf.WriteString(" ")
		// Ensure the result is a string
		strResult, ok := str.(string)
		if !ok {
			return nil, fmt.Errorf("ast printer internal error: expected string result from statement, got %T", str)
		}
		buf.WriteString(strResult)
	}
	buf.WriteString(")")
	return buf.String(), nil
}

func (ap *AstPrinter) VisitExpressionStmt(stmt *Expression) (interface{}, error) {
	// Propagate result/error from expression's Accept
	return stmt.Expression.Accept(ap)
}

func (ap *AstPrinter) VisitForStmt(stmt *For) (interface{}, error) {
	var builder strings.Builder
	var err error
	var strResult string
	ok := false

	builder.WriteString("(for ")

	// 1. Initializer
	if stmt.Initializer != nil {
		initVal, err := stmt.Initializer.Accept(ap)
		if err != nil {
			return nil, err
		}
		strResult, ok = initVal.(string)
		if !ok { return nil, fmt.Errorf("ast printer: expected string from initializer, got %T", initVal) }
		builder.WriteString(strResult)
		builder.WriteString(" ")
	} else {
		builder.WriteString("nil; ") // Indicate nil initializer clearly
	}

	// 2. Condition
	if stmt.Condition != nil {
		condVal, err := stmt.Condition.Accept(ap)
		if err != nil {
			return nil, err
		}
		strResult, ok = condVal.(string)
		if !ok { return nil, fmt.Errorf("ast printer: expected string from condition, got %T", condVal) }
		builder.WriteString(strResult)
	} else {
		builder.WriteString("nil") // Indicate nil condition clearly
	}
	builder.WriteString("; ")

	// 3. Increment
	if stmt.Increment != nil {
		incVal, err := stmt.Increment.Accept(ap)
		if err != nil {
			return nil, err
		}
		strResult, ok = incVal.(string)
		if !ok { return nil, fmt.Errorf("ast printer: expected string from increment, got %T", incVal) }
		builder.WriteString(strResult)
	} else {
		builder.WriteString("nil") // Indicate nil increment clearly
	}

	builder.WriteString(") ")

	// 4. Body
	bodyVal, err := stmt.Body.Accept(ap)
	if err != nil {
		return nil, err
	}
	strResult, ok = bodyVal.(string)
	if !ok { return nil, fmt.Errorf("ast printer: expected string from body, got %T", bodyVal) }
	builder.WriteString(strResult)


	return builder.String(), nil
}

func (ap *AstPrinter) VisitGroupingExpr(expr *Grouping) (interface{}, error) {
	return ap.parenthesize("group", expr.Expression)
}

func (ap *AstPrinter) VisitIfStmt(stmt *If) (interface{}, error) {
	var buf bytes.Buffer
	var err error
	var strResult string
	ok := false

	buf.WriteString("(if ")
	condVal, err := stmt.Condition.Accept(ap)
	if err != nil { return nil, err }
	strResult, ok = condVal.(string)
    if !ok { return nil, fmt.Errorf("ast printer: expected string from condition, got %T", condVal) }
	buf.WriteString(strResult)

	buf.WriteString(" then ")
	thenVal, err := stmt.ThenBranch.Accept(ap)
	if err != nil { return nil, err }
	strResult, ok = thenVal.(string)
    if !ok { return nil, fmt.Errorf("ast printer: expected string from then branch, got %T", thenVal) }
	buf.WriteString(strResult)

	if stmt.ElseBranch != nil {
		buf.WriteString(" else ")
		elseVal, err := stmt.ElseBranch.Accept(ap)
		if err != nil { return nil, err }
		strResult, ok = elseVal.(string)
        if !ok { return nil, fmt.Errorf("ast printer: expected string from else branch, got %T", elseVal) }
		buf.WriteString(strResult)
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
	return ap.parenthesize("print", stmt.Expression)
}

func (ap *AstPrinter) VisitUnaryExpr(expr *Unary) (interface{}, error) {
	return ap.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (ap *AstPrinter) VisitVariableExpr(expr *Variable) (interface{}, error) {
	return expr.Name.Lexeme, nil
}

func (ap *AstPrinter) VisitWhileStmt(stmt *While) (interface{}, error) {
	var buf bytes.Buffer
	var err error
	var strResult string
	ok := false

	buf.WriteString("(while ")
	condVal, err := stmt.Condition.Accept(ap)
	if err != nil { return nil, err }
	strResult, ok = condVal.(string)
    if !ok { return nil, fmt.Errorf("ast printer: expected string from condition, got %T", condVal) }
	buf.WriteString(strResult)

	buf.WriteString(" ")
	bodyVal, err := stmt.Body.Accept(ap)
	if err != nil { return nil, err }
	strResult, ok = bodyVal.(string)
    if !ok { return nil, fmt.Errorf("ast printer: expected string from body, got %T", bodyVal) }
	buf.WriteString(strResult)
	buf.WriteString(")")
	return buf.String(), nil
}

func (ap *AstPrinter) VisitVarStmt(stmt *Var) (interface{}, error) {
	if stmt.Initializer != nil {
		initVal, err := stmt.Initializer.Accept(ap) // Check error
		if err != nil {
			return nil, err
		}
		// Ensure initVal is a string
		initStr, ok := initVal.(string)
		if !ok {
			return nil, fmt.Errorf("ast printer internal error: expected string result from initializer, got %T", initVal)
		}
		return fmt.Sprintf("(var %s = %s)", stmt.Name.Lexeme, initStr), nil
	}
	return fmt.Sprintf("(var %s)", stmt.Name.Lexeme), nil
}

// Add implementation for logical expressions
func (ap *AstPrinter) VisitLogicalExpr(expr *Logical) (interface{}, error) {
	// Use parenthesize, similar to VisitBinaryExpr
	return ap.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

// Add implementation for Function statement
func (ap *AstPrinter) VisitFunctionStmt(stmt *Function) (interface{}, error) {
	var buf bytes.Buffer
	buf.WriteString("(fun ")
	buf.WriteString(stmt.Name.Lexeme)
	buf.WriteString("(")
	for i, param := range stmt.Params {
		if i > 0 { buf.WriteString(" ") }
		buf.WriteString(param.Lexeme)
	}
	buf.WriteString(") ")

	// Print body as a block
	// Need to wrap the body statements in a Block node temporarily for printing
	bodyBlock := &Block{Statements: stmt.Body}
	bodyStr, err := ap.VisitBlockStmt(bodyBlock)
	if err != nil {
		return nil, err
	}
	buf.WriteString(bodyStr.(string))

	buf.WriteString(")")
	return buf.String(), nil
}

// Add implementation for Return statement
func (ap *AstPrinter) VisitReturnStmt(stmt *Return) (interface{}, error) {
	if stmt.Value == nil {
		return "(return)", nil
	}
	return ap.parenthesize("return", stmt.Value)
}

// Add implementation for Call expression
func (ap *AstPrinter) VisitCallExpr(expr *Call) (interface{}, error) {
	// Combine callee and arguments for parenthesize
	args := make([]Expr, 0, len(expr.Arguments)+1)
	args = append(args, expr.Callee)     // First element is the callee
	args = append(args, expr.Arguments...) // Then the arguments
	return ap.parenthesize("call", args...)
}

// Add implementation for Class statement
func (ap *AstPrinter) VisitClassStmt(stmt *Class) (interface{}, error) {
	var buf bytes.Buffer
	buf.WriteString("(class ")
	buf.WriteString(stmt.Name.Lexeme)

	// Print methods
	for _, method := range stmt.Methods {
		methodStr, err := ap.VisitFunctionStmt(method) // Reuse function printing
		if err != nil {
			return nil, err
		}
		buf.WriteString(" ")
		buf.WriteString(methodStr.(string))
	}

	buf.WriteString(")")
	return buf.String(), nil
}

// Add implementation for Get expression
func (ap *AstPrinter) VisitGetExpr(expr *Get) (interface{}, error) {
	return ap.parenthesize("."+expr.Name.Lexeme, expr.Object)
}

// Add implementation for Set expression
func (ap *AstPrinter) VisitSetExpr(expr *Set) (interface{}, error) {
	return ap.parenthesize("."+expr.Name.Lexeme+"=", expr.Object, expr.Value)
}

// Add implementation for This expression
func (ap *AstPrinter) VisitThisExpr(expr *This) (interface{}, error) {
	return "this", nil
}

// Add implementation for ArrayLiteral expression
func (ap *AstPrinter) VisitArrayLiteralExpr(expr *ArrayLiteral) (interface{}, error) {
	return ap.parenthesize("array", expr.Values...)
}

// Add implementation for ArrayIndex expression
func (ap *AstPrinter) VisitArrayIndexExpr(expr *ArrayIndex) (interface{}, error) {
	// Represent as (index array index_expr)
	return ap.parenthesize("index", expr.Array, expr.Index)
}

// Add implementation for ArrayAssign expression
func (ap *AstPrinter) VisitArrayAssignExpr(expr *ArrayAssign) (interface{}, error) {
	// Represent as (assign-index array index_expr value_expr)
	return ap.parenthesize("assign-index", expr.Assignee.Array, expr.Assignee.Index, expr.Value)
}

// Add implementation for ObjectLiteral expression
func (ap *AstPrinter) VisitObjectLiteralExpr(expr *ObjectLiteral) (interface{}, error) {
	// Represent as (object key1:val1 key2:val2 ...)
	return ap.parenthesize("object", expr.Values...)
}

// parenthesize now returns an error
func (ap *AstPrinter) parenthesize(name string, exprs ...Expr) (string, error) {
	var buf bytes.Buffer
	buf.WriteString("(")
	buf.WriteString(name)
	for _, expr := range exprs {
		strVal, err := expr.Accept(ap) // Check error
		if err != nil {
			return "", err // Propagate error
		}
		// Ensure strVal is a string
		strResult, ok := strVal.(string)
		if !ok {
			return "", fmt.Errorf("ast printer internal error: expected string result from expression, got %T", strVal)
		}
		buf.WriteString(" ")
		buf.WriteString(strResult)
	}
	buf.WriteString(")")
	return buf.String(), nil // Return string and nil error
}

// Add implementation for Ternary expression
func (ap *AstPrinter) VisitTernaryExpr(expr *Ternary) (interface{}, error) {
	return ap.parenthesize("ternary", expr.Condition, expr.ThenBranch, expr.ElseBranch)
}
