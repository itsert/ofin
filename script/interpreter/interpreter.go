package interpreter

import (
	"github.com/itsert/ofin/merror"
	"github.com/itsert/ofin/script/ast"
	"github.com/itsert/ofin/script/token"
)

type interpreter struct {
}

func (p *interpreter) Interpret(expr ast.Expression) interface{} {
	return p.evaluate(expr)
}

func NewInterpreter() *interpreter {
	return &interpreter{}
}

func (p *interpreter) VisitBinaryExpression(expr *ast.Binary) interface{} {
	right := p.evaluate(expr.Right)
	left := p.evaluate(expr.Left)
	_, ok1 := left.(float64)

	var leftDouble, rightDouble float64
	var leftString, rightString string
	switch i := right.(type) {
	case float64:
		rightDouble = float64(i)
	case string:
		rightString = string(i)
	default:
		merror.RuntimeError(expr.Operator, "Operands must be two numbers or two strings.")
	}

	switch i := left.(type) {
	case float64:
		leftDouble = float64(i)
	case string:
		leftString = string(i)
	default:
		merror.RuntimeError(expr.Operator, "Operands must be two numbers or two strings.")
	}
	_, ok2 := right.(float64)
	if ok1 && ok2 {
		switch expr.Operator.Type {
		case token.MINUS:
			return float64(leftDouble) - float64(rightDouble)
		case token.SLASH:
			return float64(leftDouble) / float64(rightDouble)
		case token.ASTERISK:
			return float64(leftDouble) * float64(rightDouble)
		case token.PLUS:
			return float64(leftDouble) + float64(rightDouble)
		case token.GREATER:
			return float64(leftDouble) > float64(rightDouble)
		case token.GREATER_EQUAL:
			return float64(leftDouble) >= float64(rightDouble)
		case token.LESS:
			return float64(leftDouble) < float64(rightDouble)
		case token.LESS_EQUAL:
			return float64(leftDouble) <= float64(rightDouble)
		case token.BANG_EQUAL:
			return !p.isEqual(left, right)
		case token.EQUAL:
			return p.isEqual(left, right)
		}
	}

	_, ok1 = left.(string)
	_, ok2 = right.(string)

	if ok1 && ok2 {
		switch expr.Operator.Type {
		case token.PLUS:
			return string(leftString) + string(rightString)
		}
	}

	return nil
}
func (p *interpreter) VisitGroupingExpression(expr *ast.Grouping) interface{} {
	return p.evaluate(expr.Expr)
}
func (p *interpreter) VisitLiteralExpression(expr *ast.Literal) interface{} {
	return expr.Value
}
func (p *interpreter) VisitUnaryExpression(expr *ast.Unary) interface{} {
	right := p.evaluate(expr.Right)
	switch expr.Operator.Type {
	case token.MINUS:
		switch i := right.(type) {
		case float64:
			return -i
		default:
			merror.RuntimeError(expr.Operator, "Operand must be a numbers")
		}
	case token.BANG:
		switch i := right.(type) {
		case bool:
			return bool(i)
		default:
			if i == nil {
				return false
			}
		}
		return true
	}
	return nil
}

func (p *interpreter) evaluate(expr ast.Expression) interface{} {
	return expr.Accept(p)
}

func (p *interpreter) isEqual(a interface{}, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}

	return a == b
}
