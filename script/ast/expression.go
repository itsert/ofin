package ast

import "github.com/itsert/ofin/script/token"
type Expression interface {
	Expression()
	Accept(visitor Visitor) interface{}
}

type Binary struct {
	Left Expression
	Operator token.Token
	Right Expression
}

func NewBinary(Left Expression, Operator token.Token, Right Expression) *Binary{
	return &Binary{
		Left:	Left,
		Operator:	Operator,
		Right:	Right,
	}
}

func (b *Binary) Expression() {}

func (b *Binary) Accept(visitor Visitor) interface{} {
	 return visitor.VisitBinaryExpression(b)
}


type Grouping struct {
	Expr Expression
}

func NewGrouping(Expr Expression) *Grouping{
	return &Grouping{
		Expr:	Expr,
	}
}

func (g *Grouping) Expression() {}

func (g *Grouping) Accept(visitor Visitor) interface{} {
	 return visitor.VisitGroupingExpression(g)
}


type Literal struct {
	Value interface{}
}

func NewLiteral(Value interface{}) *Literal{
	return &Literal{
		Value:	Value,
	}
}

func (l *Literal) Expression() {}

func (l *Literal) Accept(visitor Visitor) interface{} {
	 return visitor.VisitLiteralExpression(l)
}


type Unary struct {
	Operator token.Token
	Right Expression
}

func NewUnary(Operator token.Token, Right Expression) *Unary{
	return &Unary{
		Operator:	Operator,
		Right:	Right,
	}
}

func (u *Unary) Expression() {}

func (u *Unary) Accept(visitor Visitor) interface{} {
	 return visitor.VisitUnaryExpression(u)
}


type Visitor interface {
	VisitBinaryExpression(expression *Binary) interface{}
	VisitGroupingExpression(expression *Grouping) interface{}
	VisitLiteralExpression(expression *Literal) interface{}
	VisitUnaryExpression(expression *Unary) interface{}
}

