package main

import (
	"fmt"
	"os"

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
			"Unary : Operator token.Token, Right Expression",
			"Variable : Name token.Token",
		})
		tools.GenerateAST(os.Args[2], "Statement", []string{
			"StmtExpression : Expr Expression",
			"Print : Expr Expression",
			"When : Expr Expression",
			"Then : Expr Expression",
			"And : Expr Expression",
			"Var : Name token.Token, Initializer Expression",
			"Block : statements []Statement",
		})
	} else if action == "pretty" {
		// 		input := `
		// print 1 + 2 * (3 + 4)
		// print 4
		// print true
		// print 2 + 1
		// //And 1 + 2
		// //Given a = 31 + 78
		// //And b = "hello"
		// //print b
		// //a = 23
		// //print a
		// //print b
		// Then 1 + 2
		// 		`
		// 		input = `
		// Scenario:
		//    Given a = 31 + 2
		//    And b = "hello"
		//    When:
		//       a = 23
		// 	  print a
		//       print 3
		// `
		dat, err := os.ReadFile("test.ac")
		_ = err
		l := lexer.NewLexer(string(dat), "main.go")

		stmnts, err := parser.NewParser(l).ParseProgram()
		if err != nil {
			return
		}

		fmt.Printf("Lexer %+v\n", l)
		// interpreter.NewInterpreter().Interpret(stmnts)
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
