package interpreter

// import (
// 	"fmt"
// 	"strings"

// 	"github.com/itsert/ofin/script/ast"
// )

// type prettyInterpreter struct {
// }

// func (p *prettyInterpreter) Print(expr ast.Expression) interface{} {
// 	return expr.Accept(p)
// }

// func NewPrettyPrinter() *prettyInterpreter {
// 	return &prettyInterpreter{}
// }

// func (p *prettyInterpreter) VisitBinaryExpression(expr *ast.Binary) interface{} {
// 	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
// }
// func (p *prettyInterpreter) VisitGroupingExpression(expr *ast.Grouping) interface{} {
// 	return p.parenthesize("group", expr.Expr)
// }
// func (p *prettyInterpreter) VisitLiteralExpression(expr *ast.Literal) interface{} {
// 	if expr.Value == nil {
// 		return "nil"
// 	}
// 	return fmt.Sprintf("%+v", expr.Value)
// }
// func (p *prettyInterpreter) VisitUnaryExpression(expr *ast.Unary) interface{} {
// 	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
// }

// func (p *prettyInterpreter) parenthesize(name string, exprs ...ast.Expression) string {
// 	var builder strings.Builder

// 	builder.WriteString("(")
// 	builder.WriteString(name)
// 	for _, expr := range exprs {
// 		builder.WriteString(" ")
// 		builder.WriteString(fmt.Sprintf("%+v", expr.Accept(p)))
// 	}
// 	builder.WriteString(")")

// 	return builder.String()
// }
