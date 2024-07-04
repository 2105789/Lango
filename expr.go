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
