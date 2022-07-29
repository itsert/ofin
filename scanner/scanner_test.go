package scanner

import (
	"testing"

	"github.com/itsert/ofin/token"
)

func TestWIthSimpleExpression(t *testing.T) {
	input := `=+(){},;`
	tests := []struct {
		expectedType   token.TokenType
		expectedLexeme string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LEFT_PAREN, "("},
		{token.RIGHT_PAREN, ")"},
		{token.LEFT_BRACE, "{"},
		{token.RIGHT_BRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	s := NewScanner(input)
	tokens := s.Tokenize()

	if len(tokens) != len(tests) {
		t.Fatalf("Length unmatching. expected=%d, got=%d",
			len(tests), len(tokens))
	}

	for i := range tests {
		if tokens[i].Type != tests[i].expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tests[i].expectedType, tokens[i].Type)
		}

		if tokens[i].Lexeme != tests[i].expectedLexeme {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tests[i].expectedLexeme, tokens[i].Lexeme)
		}
	}

}

func TestWIthCompoundExpression(t *testing.T) {
	input := `=+(){},;==!=<><=>=!`
	tests := []struct {
		expectedType   token.TokenType
		expectedLexeme string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LEFT_PAREN, "("},
		{token.RIGHT_PAREN, ")"},
		{token.LEFT_BRACE, "{"},
		{token.RIGHT_BRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EQUAL, "=="},
		{token.BANG_EQUAL, "!="},
		{token.LESS, "<"},
		{token.GREATER, ">"},
		{token.LESS_EQUAL, "<="},
		{token.GREATER_EQUAL, ">="},
		{token.BANG, "!"},
		{token.EOF, ""},
	}

	s := NewScanner(input)
	tokens := s.Tokenize()

	if len(tokens) != len(tests) {
		t.Fatalf("Length unmatching. expected=%d, got=%d",
			len(tests), len(tokens))
	}

	for i := range tests {
		if tokens[i].Type != tests[i].expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tests[i].expectedType, tokens[i].Type)
		}

		if tokens[i].Lexeme != tests[i].expectedLexeme {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tests[i].expectedLexeme, tokens[i].Lexeme)
		}
	}

}

func TestWithFreeFormExpression(t *testing.T) {
	input := `
	// this is a comment
(( )){} // grouping stuff
!*+-/=<> <= == // operators
	`
	tests := []struct {
		expectedType   token.TokenType
		expectedLexeme string
		line           int
	}{
		{token.LEFT_PAREN, "(", 2},
		{token.LEFT_PAREN, "(", 2},
		{token.RIGHT_PAREN, ")", 2},
		{token.RIGHT_PAREN, ")", 2},
		{token.LEFT_BRACE, "{", 2},
		{token.RIGHT_BRACE, "}", 2},
		{token.BANG, "!", 3},
		{token.ASTERISK, "*", 3},
		{token.PLUS, "+", 3},
		{token.MINUS, "-", 3},
		{token.SLASH, "/", 3},
		{token.ASSIGN, "=", 3},
		{token.LESS, "<", 3},
		{token.GREATER, ">", 3},
		{token.LESS_EQUAL, "<=", 3},
		{token.EQUAL, "==", 3},
		{token.EOF, "", 4},
	}

	s := NewScanner(input)
	tokens := s.Tokenize()

	if len(tokens) != len(tests) {
		t.Fatalf("Length unmatching. expected=%d, got=%d",
			len(tests), len(tokens))
	}

	for i := range tests {
		if tokens[i].Type != tests[i].expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tests[i].expectedType, tokens[i].Type)
		}

		if tokens[i].Lexeme != tests[i].expectedLexeme {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tests[i].expectedLexeme, tokens[i].Lexeme)
		}
	}

}

func TestWithStringExpression(t *testing.T) {
	input := `
	"Hello World!"
	`
	tests := []struct {
		expectedType   token.TokenType
		expectedLexeme string
		Literal        interface{}
	}{
		{token.STRING, "STRING", "Hello World!"},
		{token.EOF, "", nil},
	}

	s := NewScanner(input)
	tokens := s.Tokenize()

	if len(tokens) != len(tests) {
		t.Fatalf("Length unmatching. expected=%d, got=%d",
			len(tests), len(tokens))
	}

	for i := range tests {
		if tokens[i].Type != tests[i].expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tests[i].expectedType, tokens[i].Type)
		}

		if tokens[i].Literal != tests[i].Literal {
			t.Fatalf("tests[%d] - lexeme wrong. expected=%q, got=%q",
				i, tests[i].Literal, tokens[i].Literal)
		}
	}

}

func TestWithNumericalExpression(t *testing.T) {
	input := `
	12.5
	`
	tests := []struct {
		expectedType   token.TokenType
		expectedLexeme string
		Literal        interface{}
	}{
		{token.NUMBER, "NUMBER", 12.5},
		{token.EOF, "", nil},
	}

	s := NewScanner(input)
	tokens := s.Tokenize()

	if len(tokens) != len(tests) {
		t.Fatalf("Length unmatching. expected=%d, got=%d",
			len(tests), len(tokens))
	}

	for i := range tests {
		if tokens[i].Type != tests[i].expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tests[i].expectedType, tokens[i].Type)
		}

		if tokens[i].Literal != tests[i].Literal {
			t.Fatalf("tests[%d] - lexeme wrong. expected=%q, got=%q",
				i, tests[i].Literal, tokens[i].Literal)
		}
	}

}

func TestWithIdentifierAndKeywordsExpression(t *testing.T) {
	input := `
	and
	given
	when
	fire
	born
	`
	tests := []struct {
		expectedType   token.TokenType
		expectedLexeme string
	}{
		{token.AND, "and"},
		{token.GIVEN, "given"},
		{token.WHEN, "when"},
		{token.IDENTIFIER, "fire"},
		{token.IDENTIFIER, "born"},
		{token.EOF, ""},
	}

	s := NewScanner(input)
	tokens := s.Tokenize()

	if len(tokens) != len(tests) {
		t.Fatalf("Length unmatching. expected=%d, got=%d",
			len(tests), len(tokens))
	}

	for i := range tests {
		if tokens[i].Type != tests[i].expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tests[i].expectedType, tokens[i].Type)
		}

		if tokens[i].Lexeme != tests[i].expectedLexeme {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tests[i].expectedLexeme, tokens[i].Lexeme)
		}
	}

}
