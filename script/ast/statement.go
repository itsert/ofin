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


type StatementVisitor interface {
	VisitStmtExpressionStatement(statement *StmtExpression) interface{}
	VisitPrintStatement(statement *Print) interface{}
	VisitVarStatement(statement *Var) interface{}
}

