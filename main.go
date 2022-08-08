package main

import (
	"fmt"
	"os"

	"github.com/itsert/ofin/script/interpreter"
	"github.com/itsert/ofin/script/lexer"
	"github.com/itsert/ofin/script/parser"
	"github.com/itsert/ofin/script/tools"
)

func main() {
	action := os.Args[1]
	if action == "generate" {
		if len(os.Args) < 3 {
			fmt.Println("Usage: main generate <output directory>")
			os.Exit(1)
		}
		tools.GenerateAST(os.Args[2], "Expression", []string{
			"Assign   : Name token.Token, Expr Expression",
			"Binary : Left Expression, Operator token.Token, Right Expression",
			"Grouping : Expr Expression",
			"Literal : Value interface{}",
			"Logical : Left Expression, Operator token.Token, Right Expression",
			"Unary : Operator token.Token, Right Expression",
			"Variable : Name token.Token",
		})
		tools.GenerateAST(os.Args[2], "Statement", []string{
			"StmtExpression : Expr Expression",
			"If : Condition Expression, ThenBranch Statement, ElseBranch Statement",
			"Print : Expr Expression",
			"When : Expr Expression",
			"Then : Expr Expression",
			"And : Expr Expression",
			"Scenario : Label string",
			"Var : Name token.Token, Initializer Expression",
			"Block : Statements []Statement",
			"DoNoting : Name token.Token",
		})
	} else if action == "pretty" {
		dat, err := os.ReadFile("test.ac")
		_ = err
		l := lexer.NewLexer(string(dat), "main.go")

		//fmt.Printf("Lexer %+v\n", l.Tokenize())

		stmnts, err := parser.NewParser(l).ParseProgram()
		if err != nil {
			return
		}
		interpreter.NewInterpreter().Interpret(stmnts)
		_ = stmnts

	} else {
		input := `=+(){},;`
		s := lexer.NewLexer(input, "main.go")
		tokens := s.Tokenize()

		for _, v := range tokens {
			fmt.Println(v)
		}
	}
}
