package parser

import (
	"errors"
	"fmt"

	"github.com/itsert/ofin/merror"
	"github.com/itsert/ofin/script/ast"
	"github.com/itsert/ofin/script/environment"
	"github.com/itsert/ofin/script/lexer"
	"github.com/itsert/ofin/script/token"
)

const EOF_NEWLINE_MSG = "Expect NEWLINE or EOF after %s statement"
const STMT_START_ERROR_MSG = "Expect NEWLINE and Indentation for %s block statement"

type Parser struct {
	l            *lexer.Lexer
	tokens       []token.Token
	current      int
	fileName     string
	programState *environment.ProgramState
	hasError     bool
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:            l,
		current:      0,
		tokens:       l.Tokenize(),
		fileName:     l.File,
		programState: environment.NewState(),
		hasError:     false,
	}
	return p
}

func (p *Parser) ParseProgram() (stmnts []ast.Statement, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("error  encountered")
			fmt.Println("Recovered in f", r)
		}
	}()
	var statements []ast.Statement
	for !p.end() {
		decl, _ := p.declaration()
		statements = append(statements, decl)
	}
	return statements, err
}

func (p *Parser) declaration() (stmnt ast.Statement, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("error  encountered")
			fmt.Println("Recovered in f", r)
			p.synchronize()
			p.hasError = true
		}
	}()
	if p.lookAhead(token.GIVEN) {
		p.programState.Transition(environment.GIVEN)
		return p.varDeclaration(), err
	}
	return p.actionStatements(), err
}

func (p *Parser) varDeclaration() ast.Statement {
	name := p.consume("Expecting a variable name", token.IDENTIFIER)

	var initializer ast.Expression = nil
	if p.lookAhead(token.ASSIGN) {
		initializer = p.expression()
	}
	p.consume("", token.NEWLINE, token.EOF)
	return ast.NewVar(name, initializer)
}

func (p *Parser) nonActionStatements() ast.Statement {
	if p.lookAhead(token.NEWLINE, token.EOF) {
		return ast.NewDoNoting(p.peek())
	}
	if p.lookAhead(token.IF) {
		return p.ifStatement()
	}
	if p.lookAhead(token.PRINT) {
		return p.printStatement()
	}
	if p.lookAhead(token.INDENT) {
		return ast.NewBlock(p.block())
	}
	return p.expressionStatement()
}

func (p *Parser) consumeBlockStart(name string) {
	p.consume("Expect colon before start of block", token.COLON)
	p.consume(fmt.Sprintf(STMT_START_ERROR_MSG, name), token.NEWLINE)
}
func (p *Parser) ifStatement() ast.Statement {
	condition := p.expression()
	p.consumeBlockStart("if")
	var elseBranch ast.Statement = nil
	var thenBranch ast.Statement = nil
	if p.programState.IsState(environment.GLOBAL) {
		merror.Error(p.fileName, p.peek().Line, p.peek().Line, "conditional not expected in global context")
	} else if p.programState.IsState(environment.SCENARIO) {
		thenBranch = p.actionStatements()
		if p.lookAhead(token.ELSE) {
			p.consumeBlockStart("else")
			elseBranch = p.actionStatements()
		}
	} else {
		thenBranch = p.nonActionStatements()
		if p.lookAhead(token.ELSE) {
			p.consumeBlockStart("else")
			elseBranch = p.nonActionStatements()
		}
	}
	return ast.NewIf(condition, thenBranch, elseBranch)
}

func (p *Parser) actionStatements() ast.Statement {
	if p.lookAhead(token.AND) {
		return p.andStatement()
	}

	if p.lookAhead(token.WHEN) {
		return p.whenStatement()
	}

	if p.lookAhead(token.THEN) {
		return p.thenStatement()
	}

	if p.lookAhead(token.SCENARIO) {
		return p.scenarioStatement()
	}

	return p.nonActionStatements()
}

func (p *Parser) block() []ast.Statement {
	var statements []ast.Statement
	for !p.lookAhead(token.DEDENT) && !p.end() {
		declaration, err := p.declaration()
		_ = err
		statements = append(statements, declaration)
	}

	if p.peek().Type != token.EOF && p.previous().Type != token.DEDENT {
		merror.Error(p.fileName, p.previous().Line, p.previous().Line, "Expects a dedenatation or EOF after block.")
	}

	return statements
}

func (p *Parser) subBlock() []ast.Statement {
	var statements []ast.Statement
	for !p.lookAhead(token.DEDENT) && !p.end() {
		statements = append(statements, p.nonActionStatements())
	}

	if p.peek().Type != token.EOF && p.previous().Type != token.DEDENT {
		merror.Error(p.fileName, p.previous().Line, p.previous().Line, "Expects a dedenatation or EOF after block.")
	}

	return statements
}

func (p *Parser) printStatement() ast.Statement {
	value := p.expression()
	if !p.end() {
		p.consume(fmt.Sprintf(EOF_NEWLINE_MSG, "Print"), token.NEWLINE)
	}

	return ast.NewPrint(value)
}

func (p *Parser) andStatement() ast.Statement {
	if p.programState.IsState(environment.GIVEN) {
		return p.varDeclaration()
	} else {
		value := p.expression()
		if !p.end() {
			p.consume(fmt.Sprintf(EOF_NEWLINE_MSG, "And"), token.NEWLINE)
		}
		return ast.NewAnd(value)
	}
}
func (p *Parser) whenStatement() ast.Statement {
	p.programState.Transition(environment.WHEN)
	if p.lookAhead(token.COLON) {
		p.consume(fmt.Sprintf(STMT_START_ERROR_MSG, "When"), token.NEWLINE)
		p.consume(fmt.Sprintf(STMT_START_ERROR_MSG, "When"), token.INDENT)
		return ast.NewBlock(p.subBlock())
	} else {
		value := p.expression()
		if !p.end() {
			p.consume(fmt.Sprintf(EOF_NEWLINE_MSG, "When"), token.NEWLINE)
		}
		return ast.NewWhen(value)
	}
}

func (p *Parser) thenStatement() ast.Statement {
	p.programState.Transition(environment.THEN)
	if p.lookAhead(token.COLON) {
		p.consume(fmt.Sprintf(STMT_START_ERROR_MSG, "Then"), token.NEWLINE)
		p.consume(fmt.Sprintf(STMT_START_ERROR_MSG, "Then"), token.INDENT)
		return ast.NewBlock(p.subBlock())
	} else {
		value := p.expression()
		if !p.end() {
			p.consume(fmt.Sprintf(EOF_NEWLINE_MSG, "When"), token.NEWLINE)
		}
		return ast.NewThen(value)
	}
}

func (p *Parser) scenarioStatement() ast.Statement {
	var label string
	p.programState.Transition(environment.SCENARIO)
	if p.lookAhead(token.STRING) {
		label = p.previous().Literal.(string)
	} else {
		merror.RuntimeError(p.peek(), "Expected string label")
	}
	p.consume("Expect COLON to indicate start of new block", token.COLON)
	p.consume(fmt.Sprintf(EOF_NEWLINE_MSG, "Scenario"), token.NEWLINE)
	return ast.NewScenario(label)
}

func (p *Parser) expressionStatement() ast.Statement {
	value := p.expression()
	if !p.end() {
		p.consume(fmt.Sprintf(EOF_NEWLINE_MSG, "Expression"), token.NEWLINE)
	}
	return ast.NewStmtExpression(value)
}

func (p *Parser) expression() ast.Expression {
	return p.assignment()
}

func (p *Parser) assignment() ast.Expression {
	expr := p.equality()

	if p.lookAhead(token.ASSIGN) {
		equals := p.previous()
		value := p.assignment()
		if field, ok := expr.(*ast.Variable); ok {
			name := field.Name
			return ast.NewAssign(name, value)
		}
		merror.RuntimeError(equals, "Invalid assignment target")
	}
	return expr
}

func (p *Parser) equality() ast.Expression {
	expr := p.comparison()

	for p.lookAhead(token.BANG_EQUAL, token.EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) comparison() ast.Expression {
	expr := p.addSub()

	for p.lookAhead(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right := p.addSub()
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) addSub() ast.Expression {
	expr := p.factor()

	for p.lookAhead(token.MINUS, token.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) factor() ast.Expression {
	expr := p.unary()

	for p.lookAhead(token.SLASH, token.ASTERISK) {
		operator := p.previous()
		right := p.unary()
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) unary() ast.Expression {

	if p.lookAhead(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		return ast.NewUnary(operator, right)
	}
	return p.primary()
}

func (p *Parser) primary() ast.Expression {
	if p.lookAhead(token.FALSE) {
		return ast.NewLiteral(false)
	}
	if p.lookAhead(token.TRUE) {
		return ast.NewLiteral(true)
	}
	if p.lookAhead(token.NIL) {
		return ast.NewLiteral(nil)
	}

	if p.lookAhead(token.NUMBER, token.STRING) {
		return ast.NewLiteral(p.previous().Literal)
	}

	if p.lookAhead(token.IDENTIFIER) {
		return ast.NewVariable(p.previous())
	}

	if p.lookAhead(token.LEFT_PAREN) {
		expr := p.expression()
		p.consume("Expect ')' after expression", token.RIGHT_PAREN)
		return ast.NewGrouping(expr)
	}
	merror.Error(p.fileName, p.peek().Line, p.peek().Line, "Expected expression")
	return nil
}

func (p *Parser) consume(message string, types ...token.TokenType) token.Token {
	for _, t := range types {
		if p.check(t) {
			return p.advance()
		}
	}
	merror.Error(p.fileName, p.peek().Line, p.peek().Line, message)
	return token.Token{}
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.end() {
		if p.previous().Type == token.NEWLINE {
			return
		}

		switch p.peek().Type {
		case token.SCENARIO, token.FUNCTION, token.GIVEN, token.IF, token.WHILE, token.PRINT, token.RETURN:
			return
		}
		p.advance()
	}
}

func (p *Parser) lookAhead(types ...token.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t token.TokenType) bool {
	if p.end() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) checkPrevious(t token.TokenType) bool {
	return p.previous().Type == t
}

func (p *Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p *Parser) advance() token.Token {
	currentToken := p.tokens[p.current]
	if !p.end() {
		p.current += 1
	}

	return currentToken
}

func (p *Parser) end() bool {
	return p.tokens[p.current].Type == token.EOF
}
