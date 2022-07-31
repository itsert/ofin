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
			"Binary : Left Expression, Operator token.Token, Right Expression",
			"Grouping : Expr Expression",
			"Literal : Value interface{}",
			"Unary : Operator token.Token, Right Expression",
		})
		// tools.GenerateAST(os.Args[2], "Statement", []string{
		// 	"StmtExpression : Expr Expression",
		// 	"Print : Expr Expression",
		// })
	} else if action == "pretty" {
		input := `
		1 + 2 * (3 + 4)
		`
		l := lexer.NewLexer(input, "main.go")

		expression, err := parser.NewParser(l).ParseProgram()
		if err != nil {
			return
		}

		fmt.Println(interpreter.NewPrettyPrinter().Print(expression))
		fmt.Printf("%+v\n", interpreter.NewInterpreter().Interpret(expression))

	} else {
		input := `=+(){},;`
		s := lexer.NewLexer(input, "main.go")
		tokens := s.Tokenize()

		for _, v := range tokens {
			fmt.Println(v)
		}
	}
}
