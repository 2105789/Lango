package main

type Expr interface {
	Accept(Visitor) (interface{}, error)
}

type Assign struct {
	Name  *Token
	Value Expr
}

func (a *Assign) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitAssignExpr(a)
}

type Binary struct {
	Left     Expr
	Operator *Token
	Right    Expr
}

func (b *Binary) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitBinaryExpr(b)
}

type Grouping struct {
	Expression Expr
}

func (g *Grouping) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitGroupingExpr(g)
}

type Literal struct {
	Value interface{}
}

func (l *Literal) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitLiteralExpr(l)
}

type Unary struct {
	Operator *Token
	Right    Expr
}

func (u *Unary) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitUnaryExpr(u)
}

type Variable struct {
	Name *Token
}

func (v *Variable) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitVariableExpr(v)
}

// New Logical expression type
type Logical struct {
	Left     Expr
	Operator Token // Will be AND or OR
	Right    Expr
}

func (l *Logical) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitLogicalExpr(l)
}

// Call Expression
type Call struct {
	Callee    Expr   // The expression that evaluates to the function
	Paren     Token  // The closing parenthesis (for error reporting)
	Arguments []Expr // List of argument expressions
}

func (c *Call) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitCallExpr(c)
}

// Get Expression (for property access, e.g., object.property)
type Get struct {
	Object Expr  // The object whose property is being accessed
	Name   Token // The property name (identifier)
}

func (g *Get) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitGetExpr(g)
}

// Set Expression (for property assignment, e.g., object.property = value)
type Set struct {
	Object Expr  // The object whose property is being set
	Name   Token // The property name (identifier)
	Value  Expr  // The value being assigned
}

func (s *Set) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitSetExpr(s)
}

// This Expression (the 'this' keyword)
type This struct {
	Keyword Token // The 'this' token
}

func (t *This) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitThisExpr(t)
}

// Array Literal Expression (e.g., [1, 2, 3])
type ArrayLiteral struct {
	Bracket Token  // The opening bracket token '['
	Values  []Expr // Expressions for the array elements
}

func (a *ArrayLiteral) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitArrayLiteralExpr(a)
}

// Array Index Expression (e.g., myArray[index])
type ArrayIndex struct {
	Array   Expr  // Expression evaluating to the array
	Bracket Token // The closing bracket token ']'
	Index   Expr  // Expression evaluating to the index
}

func (a *ArrayIndex) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitArrayIndexExpr(a)
}

// Array Assignment Expression (e.g., myArray[index] = value)
type ArrayAssign struct {
	Assignee ArrayIndex // Uses ArrayIndex structure for array and index
	Value    Expr       // Expression for the value to assign
}

func (a *ArrayAssign) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitArrayAssignExpr(a)
}

// Object Literal Expression (e.g., {name: "John", age: 30})
type ObjectLiteral struct {
	Brace  Token   // The opening brace token '{'
	Keys   []Token // Property keys (identifiers or strings)
	Values []Expr  // Property values (expressions)
}

func (o *ObjectLiteral) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitObjectLiteralExpr(o)
}

type Ternary struct {
	Condition Expr
	ThenBranch Expr
	ElseBranch Expr
}

func (t *Ternary) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitTernaryExpr(t)
}
