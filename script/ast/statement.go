package ast

import "github.com/itsert/ofin/script/token"
import "github.com/itsert/ofin/script/environment"
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


type If struct {
	Condition Expression
	ThenBranch Statement
	ElseBranch Statement
}

func NewIf(Condition Expression, ThenBranch Statement, ElseBranch Statement) *If{
	return &If{
		Condition:	Condition,
		ThenBranch:	ThenBranch,
		ElseBranch:	ElseBranch,
	}
}

func (i *If) Statement() {}

func (i *If) Accept(visitor StatementVisitor) interface{} {
	 return visitor.VisitIfStatement(i)
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


type Scenario struct {
	Label string
}

func NewScenario(Label string) *Scenario{
	return &Scenario{
		Label:	Label,
	}
}

func (s *Scenario) Statement() {}

func (s *Scenario) Accept(visitor StatementVisitor) interface{} {
	 return visitor.VisitScenarioStatement(s)
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


type While struct {
	Condition Expression
	Body Statement
}

func NewWhile(Condition Expression, Body Statement) *While{
	return &While{
		Condition:	Condition,
		Body:	Body,
	}
}

func (w *While) Statement() {}

func (w *While) Accept(visitor StatementVisitor) interface{} {
	 return visitor.VisitWhileStatement(w)
}


type Block struct {
	Statements []Statement
	BlockState environment.State
}

func NewBlock(Statements []Statement, BlockState environment.State) *Block{
	return &Block{
		Statements:	Statements,
		BlockState:	BlockState,
	}
}

func (b *Block) Statement() {}

func (b *Block) Accept(visitor StatementVisitor) interface{} {
	 return visitor.VisitBlockStatement(b)
}


type DoNoting struct {
	Name token.Token
}

func NewDoNoting(Name token.Token) *DoNoting{
	return &DoNoting{
		Name:	Name,
	}
}

func (d *DoNoting) Statement() {}

func (d *DoNoting) Accept(visitor StatementVisitor) interface{} {
	 return visitor.VisitDoNotingStatement(d)
}


type StatementVisitor interface {
	VisitStmtExpressionStatement(statement *StmtExpression) interface{}
	VisitIfStatement(statement *If) interface{}
	VisitPrintStatement(statement *Print) interface{}
	VisitWhenStatement(statement *When) interface{}
	VisitThenStatement(statement *Then) interface{}
	VisitAndStatement(statement *And) interface{}
	VisitScenarioStatement(statement *Scenario) interface{}
	VisitVarStatement(statement *Var) interface{}
	VisitWhileStatement(statement *While) interface{}
	VisitBlockStatement(statement *Block) interface{}
	VisitDoNotingStatement(statement *DoNoting) interface{}
}

