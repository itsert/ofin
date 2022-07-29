package token

type TokenType string

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
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
	SEMICOLON   = ","
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
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	WHILE    = "WHILE"
	WHEN     = "WHEN"
	SCENARIO = "SCENARIO"
	THEN     = "THEN"
	GIVEN    = "GIVEN"
	STORY    = "STORY"

	STRING = "STRING"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
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
}

func LookupIdentifier(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}
