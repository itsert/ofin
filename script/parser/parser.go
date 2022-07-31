package parser

import (
	"errors"
	"fmt"

	"github.com/itsert/ofin/merror"
	"github.com/itsert/ofin/script/ast"
	"github.com/itsert/ofin/script/lexer"
	"github.com/itsert/ofin/script/token"
)

type Parser struct {
	l        *lexer.Lexer
	tokens   []token.Token
	current  int
	fileName string
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, current: 0, tokens: l.Tokenize(), fileName: l.File}
	return p
}

func (p *Parser) ParseProgram() (expr ast.Expression, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("error  encountered")
			fmt.Println("Recovered in f", r)
		}
	}()
	return p.expression(), err
}

func (p *Parser) expression() ast.Expression {
	return p.equality()
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

	if p.lookAhead(token.LEFT_PAREN) {
		expr := p.expression()
		p.consume(token.RIGHT_PAREN, "Expect ')' after expression")
		return ast.NewGrouping(expr)
	}
	merror.Error(p.fileName, p.peek().Line, p.peek().Line, "Expec expression")
	return nil
}

func (p *Parser) consume(t token.TokenType, message string) token.Token {
	if p.check(t) {
		return p.advance()
	}
	merror.Error(p.fileName, p.peek().Line, p.peek().Line, "Expecting an expression")
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
