package ast

import "github.com/itsert/ofin/script/token"
type Statement interface {
	Statement()
	Accept(visitor StatementVisitor) interface{}
}

type StmtExpression struct {
	Expr Expression
}

func NewStmtExpression(Expr Expression) *StmtExpression{
	return &StmtExpression{
		Expr:	Expr,
	}
}

func (s *StmtExpression) Statement() {}

func (s *StmtExpression) Accept(visitor StatementVisitor) interface{} {
	 return visitor.VisitStmtExpressionStatement(s)
}


type Print struct {
	Expr Expression
}

func NewPrint(Expr Expression) *Print{
	return &Print{
		Expr:	Expr,
	}
}

func (p *Print) Statement() {}

func (p *Print) Accept(visitor StatementVisitor) interface{} {
	 return visitor.VisitPrintStatement(p)
}


type When struct {
	Expr Expression
}

func NewWhen(Expr Expression) *When{
	return &When{
		Expr:	Expr,
	}
}

func (w *When) Statement() {}

func (w *When) Accept(visitor StatementVisitor) interface{} {
	 return visitor.VisitWhenStatement(w)
}


type Then struct {
	Expr Expression
}

func NewThen(Expr Expression) *Then{
	return &Then{
		Expr:	Expr,
	}
}

func (t *Then) Statement() {}

func (t *Then) Accept(visitor StatementVisitor) interface{} {
	 return visitor.VisitThenStatement(t)
}


type And struct {
	Expr Expression
}

func NewAnd(Expr Expression) *And{
	return &And{
		Expr:	Expr,
	}
}

func (a *And) Statement() {}

func (a *And) Accept(visitor StatementVisitor) interface{} {
	 return visitor.VisitAndStatement(a)
}


type Var struct {
	Name token.Token
	Initializer Expression
}

func NewVar(Name token.Token, Initializer Expression) *Var{
	return &Var{
		Name:	Name,
		Initializer:	Initializer,
	}
}

func (v *Var) Statement() {}

func (v *Var) Accept(visitor StatementVisitor) interface{} {
	 return visitor.VisitVarStatement(v)
}


type Block struct {
	statements []Statement
}

func NewBlock(statements []Statement) *Block{
	return &Block{
		statements:	statements,
	}
}

func (b *Block) Statement() {}

func (b *Block) Accept(visitor StatementVisitor) interface{} {
	 return visitor.VisitBlockStatement(b)
}


type StatementVisitor interface {
	VisitStmtExpressionStatement(statement *StmtExpression) interface{}
	VisitPrintStatement(statement *Print) interface{}
	VisitWhenStatement(statement *When) interface{}
	VisitThenStatement(statement *Then) interface{}
	VisitAndStatement(statement *And) interface{}
	VisitVarStatement(statement *Var) interface{}
	VisitBlockStatement(statement *Block) interface{}
}

