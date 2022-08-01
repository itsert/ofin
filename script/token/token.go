package token

type TokenType string

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func NewToken(
	Type TokenType,
	Lexeme string,
	Literal interface{},
	Line int,
) *Token {
	return &Token{
		Type:    Type,
		Lexeme:  Lexeme,
		Literal: Literal,
		Line:    Line,
	}
}

const (
	ILLEGAL = ""
	EOF     = "EOF"
	// Identifiers + literals
	IDENTIFIER = "IDENTIFIER"
	NUMBER     = "NUMBER"
	// Operators
	ASSIGN        = "="
	DOT           = "."
	PLUS          = "+"
	MINUS         = "-"
	BANG          = "!"
	ASTERISK      = "*"
	SLASH         = "/"
	LESS          = "<"
	LESS_EQUAL    = "<="
	GREATER       = ">"
	GREATER_EQUAL = ">="
	EQUAL         = "=="
	BANG_EQUAL    = "!="
	// Delimiters
	COMMA       = ","
	SEMICOLON   = ";"
	NEWLINE     = "NEWLINE"
	LEFT_PAREN  = "("
	RIGHT_PAREN = ")"
	LEFT_BRACE  = "{"
	RIGHT_BRACE = "}"
	AND         = "AND"
	OR          = "OR"
	NIL         = "NIL"
	// Keywords
	// 1343456
	FUNCTION = "FUNCTION"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	WHILE    = "WHILE"
	FOR      = "FOR"
	WHEN     = "WHEN"
	SCENARIO = "SCENARIO"
	THEN     = "THEN"
	GIVEN    = "GIVEN"
	STORY    = "STORY"
	STRING   = "STRING"

	//FUNCTIONS
	PRINT = "PRINT"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"when":   WHEN,
	"then":   THEN,
	"and":    AND,
	"given":  GIVEN,
	"story":  STORY,
	"print":  PRINT,
}

func LookupIdentifier(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}
