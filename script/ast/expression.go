package ast

import "github.com/itsert/ofin/script/token"

type Expression interface {
	Expression()
	Accept(visitor ExpressionVisitor) interface{}
}

type Assign struct {
	Name token.Token
	Expr Expression
}

func NewAssign(Name token.Token, Expr Expression) *Assign {
	return &Assign{
		Name: Name,
		Expr: Expr,
	}
}

func (a *Assign) Expression() {}

func (a *Assign) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitAssignExpression(a)
}

type Binary struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func NewBinary(Left Expression, Operator token.Token, Right Expression) *Binary {
	return &Binary{
		Left:     Left,
		Operator: Operator,
		Right:    Right,
	}
}

func (b *Binary) Expression() {}

func (b *Binary) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitBinaryExpression(b)
}

type Call struct {
	Callee    Expression
	Paren     token.Token
	Arguments []Expression
}

func NewCall(Callee Expression, Paren token.Token, Arguments []Expression) *Call {
	return &Call{
		Callee:    Callee,
		Paren:     Paren,
		Arguments: Arguments,
	}
}

func (c *Call) Expression() {}

func (c *Call) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitCallExpression(c)
}

type Grouping struct {
	Expr Expression
}

func NewGrouping(Expr Expression) *Grouping {
	return &Grouping{
		Expr: Expr,
	}
}

func (g *Grouping) Expression() {}

func (g *Grouping) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitGroupingExpression(g)
}

type Literal struct {
	Value interface{}
}

func NewLiteral(Value interface{}) *Literal {
	return &Literal{
		Value: Value,
	}
}

func (l *Literal) Expression() {}

func (l *Literal) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitLiteralExpression(l)
}

type Logical struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func NewLogical(Left Expression, Operator token.Token, Right Expression) *Logical {
	return &Logical{
		Left:     Left,
		Operator: Operator,
		Right:    Right,
	}
}

func (l *Logical) Expression() {}

func (l *Logical) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitLogicalExpression(l)
}

type Unary struct {
	Operator token.Token
	Right    Expression
}

func NewUnary(Operator token.Token, Right Expression) *Unary {
	return &Unary{
		Operator: Operator,
		Right:    Right,
	}
}

func (u *Unary) Expression() {}

func (u *Unary) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitUnaryExpression(u)
}

type Variable struct {
	Name token.Token
}

func NewVariable(Name token.Token) *Variable {
	return &Variable{
		Name: Name,
	}
}

func (v *Variable) Expression() {}

func (v *Variable) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitVariableExpression(v)
}

type ExpressionVisitor interface {
	VisitAssignExpression(expression *Assign) interface{}
	VisitBinaryExpression(expression *Binary) interface{}
	VisitCallExpression(expression *Call) interface{}
	VisitGroupingExpression(expression *Grouping) interface{}
	VisitLiteralExpression(expression *Literal) interface{}
	VisitLogicalExpression(expression *Logical) interface{}
	VisitUnaryExpression(expression *Unary) interface{}
	VisitVariableExpression(expression *Variable) interface{}
}
