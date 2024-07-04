package main

type Stmt interface {
	Accept(Visitor) (interface{}, error)
}

type Visitor interface {
	VisitBinaryExpr(*Binary) (interface{}, error)
	VisitGroupingExpr(*Grouping) (interface{}, error)
	VisitLiteralExpr(*Literal) (interface{}, error)
	VisitUnaryExpr(*Unary) (interface{}, error)
	VisitVariableExpr(*Variable) (interface{}, error)
	VisitAssignExpr(*Assign) (interface{}, error)
	VisitExpressionStmt(*Expression) (interface{}, error)
	VisitPrintStmt(*Print) (interface{}, error)
	VisitVarStmt(*Var) (interface{}, error)
	VisitBlockStmt(*Block) (interface{}, error)
	VisitIfStmt(*If) (interface{}, error)
	VisitWhileStmt(*While) (interface{}, error)
	VisitForStmt(*For) (interface{}, error)
}

type Expression struct {
	Expression Expr
}

func (es *Expression) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitExpressionStmt(es)
}

type Print struct {
	Expression Expr
}

func (p *Print) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitPrintStmt(p)
}

type Var struct {
	Name        *Token
	Initializer Expr
}

func (v *Var) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitVarStmt(v)
}

type Block struct {
	Statements []Stmt
}

func (b *Block) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitBlockStmt(b)
}

type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (i *If) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitIfStmt(i)
}

type While struct {
	Condition Expr
	Body      Stmt
}

func (w *While) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitWhileStmt(w)
}

type For struct {
	Initializer Stmt
	Condition   Expr
	Increment   Expr
	Body        Stmt
}

func (f *For) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitForStmt(f)
}
