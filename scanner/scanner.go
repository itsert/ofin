package scanner

import (
	"fmt"
	"strconv"

	"github.com/itsert/ofin/merror"
	"github.com/itsert/ofin/token"
)

type Scanner struct {
	input   string
	start   int
	current int
	line    int
	tokens  []token.Token
}

func NewScanner(input string) *Scanner {
	scanner := &Scanner{
		line:    1,
		start:   0,
		current: 0,
		input:   input,
	}
	return scanner
}

func (s *Scanner) Tokenize() []token.Token {
	for !s.end() {
		s.start = s.current
		s.munchToken()
	}

	s.tokens = append(s.tokens, token.Token{
		Type:    token.EOF,
		Lexeme:  "",
		Literal: nil,
		Line:    s.line,
	})
	return s.tokens
}

func (s *Scanner) addToken(tokenType token.TokenType, literal interface{}) {
	text := s.input[s.start:s.current]
	s.tokens = append(s.tokens, token.Token{
		Type:    tokenType,
		Lexeme:  text,
		Literal: literal,
		Line:    s.line,
	})
}

func (s *Scanner) advance() byte {
	current := s.input[s.current]
	s.current += 1
	return current
}

func (s *Scanner) peek() byte {
	if s.end() {
		return 0
	}
	return s.input[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.input) {
		return 0
	}
	return s.input[s.current+1]
}

func (s *Scanner) munchToken() {
	ch := s.advance()
	switch ch {
	case '(':
		s.addToken(token.LEFT_PAREN, nil)
	case ')':
		s.addToken(token.RIGHT_PAREN, nil)
	case '{':
		s.addToken(token.LEFT_BRACE, nil)
	case '}':
		s.addToken(token.RIGHT_BRACE, nil)
	case ',':
		s.addToken(token.COMMA, nil)
	case '.':
		s.addToken(token.DOT, nil)
	case '-':
		s.addToken(token.MINUS, nil)
	case '+':
		s.addToken(token.PLUS, nil)
	case ';':
		s.addToken(token.SEMICOLON, nil)
	case '*':
		s.addToken(token.ASTERISK, nil)
	case '!':
		var ch token.TokenType
		if s.match('=') {
			ch = token.BANG_EQUAL
		} else {
			ch = token.BANG
		}
		s.addToken(ch, nil)
	case '=':
		var ch token.TokenType
		if s.match('=') {
			ch = token.EQUAL
		} else {
			ch = token.ASSIGN
		}
		s.addToken(ch, nil)
	case '<':
		var ch token.TokenType
		if s.match('=') {
			ch = token.LESS_EQUAL
		} else {
			ch = token.LESS
		}
		s.addToken(ch, nil)
	case '>':
		var ch token.TokenType
		if s.match('=') {
			ch = token.GREATER_EQUAL
		} else {
			ch = token.GREATER
		}
		s.addToken(ch, nil)
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.end() {
				s.advance()
			}
		} else {
			s.addToken(token.SLASH, nil)
		}
	case '\n':
		s.line += 1
	case ' ', '\t', '\r':
		break
	case '"':
		s.munchString()
	default:
		if isDigit(ch) {
			s.munchNumbers()
		} else if isLetter(ch) {
			s.munchIdentifier()
		} else {
			merror.Error("test.file", s.line, s.start, "Unexpected character.")
		}
	}
}

func (s *Scanner) munchIdentifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}
	currentType := token.LookupIdentifier(s.input[s.start:s.current])
	s.addToken(currentType, nil)
}

func isAlphaNumeric(c byte) bool {
	return isLetter(c) || isDigit(c)
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (s *Scanner) munchNumbers() {
	for isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}
	if f, err := strconv.ParseFloat(s.input[s.start:s.current], 32); err == nil {
		s.addToken(token.NUMBER, f)
	} else {
		msg := fmt.Sprintf("Error parsing value %s", s.input[s.start:s.current])
		merror.Error("test.of", s.line, s.start, msg)
	}
}

func (s *Scanner) munchString() {
	for s.peek() != '"' && !s.end() {
		if s.peek() == '\n' {
			merror.Error("test.of", s.line, s.start, "String did not terminate before encountering newline")
		}
		s.advance()
	}

	if s.end() {
		merror.Error("test.of", s.line, s.start, "String does not terminate")
		return
	}
	s.advance()

	str := s.input[s.start+1 : s.current-1]
	s.addToken(token.STRING, str)
}

func (s *Scanner) match(expected byte) bool {
	if s.end() {
		return false
	}
	if s.input[s.current] != expected {
		return false

	}
	s.current += 1
	return true
}

func (s *Scanner) end() bool {
	return s.current >= len(s.input)
}
