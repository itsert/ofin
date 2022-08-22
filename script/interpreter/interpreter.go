package interpreter

import (
	"fmt"
	"github.com/itsert/ofin/script/callable"

	"github.com/itsert/ofin/merror"
	"github.com/itsert/ofin/script/ast"
	"github.com/itsert/ofin/script/environment"
	"github.com/itsert/ofin/script/token"
)

type Interpreter struct {
	environment  *environment.Environment
	Global       *environment.Environment
	programState *environment.ProgramState
	label        string
}

func NewInterpreter() *Interpreter {
	globals := environment.NewEnvironment()
	defineNativeFunctions(globals)
	return &Interpreter{
		environment:  globals,
		Global:       globals,
		programState: environment.NewState(),
	}
}

func defineNativeFunctions(env *environment.Environment) {
	env.Define("clock", callable.NewClock())
}

func (p *Interpreter) VisitCallExpression(expression *ast.Call) interface{} {
	callee := p.evaluate(expression.Callee)
	var arguments []interface{}
	for _, expr := range expression.Arguments {
		arguments = append(arguments, p.evaluate(expr))
	}
	switch fn := callee.(type) {
	case callable.Callable:
		argList := len(arguments)
		if argList != fn.Arity() {
			merror.RuntimeError(
				expression.Paren,
				fmt.Sprintf("Expected %d arguments, but got %d. ",
					fn.Arity(),
					argList))
		}
		return fn.Call(p.Global, arguments)
	default:
		merror.RuntimeError(expression.Paren, "Can only call functions.")
	}
	return nil
}

func (p *Interpreter) Interpret(stmts []ast.Statement) {
	defer func() {
		if r := recover(); r != nil {
			// err = errors.New("error  encountered")
			fmt.Println("Recovered in f", r)
		}
	}()
	for _, stmt := range stmts {
		p.execute(stmt)
	}
}

func (p *Interpreter) execute(stmt ast.Statement) {
	stmt.Accept(p)
}

func (p *Interpreter) VisitLogicalExpression(expression *ast.Logical) interface{} {
	left := p.evaluate(expression.Left)

	if expression.Operator.Type == token.LOGICAL_OR {
		if p.expressBoolean(left) {
			return left
		}
	} else {
		if !p.expressBoolean(left) {
			return left
		}
	}
	return p.evaluate(expression.Right)
}

func (p *Interpreter) VisitBinaryExpression(expr *ast.Binary) interface{} {
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
func (p *Interpreter) VisitGroupingExpression(expr *ast.Grouping) interface{} {
	return p.evaluate(expr.Expr)
}
func (p *Interpreter) VisitLiteralExpression(expr *ast.Literal) interface{} {
	return expr.Value
}
func (p *Interpreter) VisitUnaryExpression(expr *ast.Unary) interface{} {
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
		return p.expressBoolean(right)
	}
	return nil
}

func (p *Interpreter) expressBoolean(expr interface{}) bool {
	switch i := expr.(type) {
	case bool:
		return i
	default:
		if i == nil {
			return false
		}
	}
	return true
}

func (p *Interpreter) VisitVariableExpression(expression *ast.Variable) interface{} {
	v, err := p.environment.Get(expression.Name)
	if err != nil {
		merror.RuntimeError(expression.Name, err.Error())
	}
	return v
}

func (p *Interpreter) VisitAssignExpression(expression *ast.Assign) interface{} {
	value := p.evaluate(expression.Expr)
	p.environment.Assign(expression.Name, value)
	return value
}

func (p *Interpreter) VisitIfStatement(statement *ast.If) interface{} {
	if p.expressBoolean(p.evaluate(statement.Condition)) {
		p.execute(statement.ThenBranch)
	} else if statement.ElseBranch != nil {
		p.execute(statement.ElseBranch)
	}
	return nil
}

func (p *Interpreter) VisitStmtExpressionStatement(statement *ast.StmtExpression) interface{} {
	p.evaluate(statement.Expr)
	return nil
}
func (p *Interpreter) VisitPrintStatement(statement *ast.Print) interface{} {
	value := p.evaluate(statement.Expr)
	fmt.Printf("%+v\n", value)
	return nil
}

func (p *Interpreter) VisitWhileStatement(statement *ast.While) interface{} {
	for p.expressBoolean(p.evaluate(statement.Condition)) {
		p.execute(statement.Body)
	}
	return nil
}

func (p *Interpreter) VisitVarStatement(statement *ast.Var) interface{} {
	var value interface{} = nil
	if statement.Initializer != nil {
		value = p.evaluate(statement.Initializer)
	}
	p.environment.Define(statement.Name.Lexeme, value)
	_, err := p.programState.Transition(environment.GIVEN)
	_ = err
	return nil
}

func (p *Interpreter) VisitWhenStatement(statement *ast.When) interface{} {
	p.executeWhen(statement)
	_, err := p.programState.Transition(environment.WHEN)
	_ = err
	return nil
}

func (p *Interpreter) executeWhen(statement *ast.When) {
	value := p.evaluate(statement.Expr)
	fmt.Printf("%+v\n", value)
}
func (p *Interpreter) VisitThenStatement(statement *ast.Then) interface{} {
	p.executeThen(statement)
	_, err := p.programState.Transition(environment.THEN)
	_ = err
	return nil
}

func (p *Interpreter) executeThen(statement *ast.Then) {
	value := p.evaluate(statement.Expr)
	result := p.expressBoolean(value)
	fmt.Printf("%+v\n", value)
	fmt.Printf("%+v\n", result)
}

func (p *Interpreter) VisitAndStatement(statement *ast.And) interface{} {
	if p.programState.IsState(environment.WHEN) {
		p.executeWhen(&ast.When{Expr: statement.Expr})
	} else if p.programState.IsState(environment.THEN) {
		p.executeThen(&ast.Then{Expr: statement.Expr})
	} else if p.programState.IsState(environment.GIVEN) {
		var varExpr interface{} = statement.Expr
		fmt.Printf("current state %v\n", p.programState.CurrentState())
		p.VisitAssignExpression(varExpr.(*ast.Assign))
	} else {
		fmt.Printf("AND: Program in invalid State:%+v\n", p.programState)
	}
	return nil
}

func (p *Interpreter) VisitScenarioStatement(statement *ast.Scenario) interface{} {
	_, err := p.programState.Transition(environment.SCENARIO)
	_ = err
	p.label = statement.Label
	return nil
}

func (p *Interpreter) VisitBlockStatement(statement *ast.Block) interface{} {
	_, err := p.programState.Transition(statement.BlockState)
	_ = err
	p.executeBlock(statement.Statements, environment.NewEnvironmentWithParent(p.environment))
	return nil
}

func (p *Interpreter) executeBlock(statements []ast.Statement, environment *environment.Environment) {
	previous := p.environment
	defer func() {
		p.environment = previous
	}()
	p.environment = environment
	for _, statement := range statements {
		p.execute(statement)
	}
}

//Convenience function to silently ignore newlines
func (p *Interpreter) VisitDoNotingStatement(statement *ast.DoNoting) interface{} {
	return nil
}

func (p *Interpreter) evaluate(expr ast.Expression) interface{} {
	return expr.Accept(p)
}

func (p *Interpreter) isEqual(a interface{}, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}

	return a == b
}
