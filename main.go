package main

import (
	"fmt"

	"github.com/itsert/ofin/scanner"
)

func main() {
	input := `=+(){},;`
	s := scanner.NewScanner(input)
	tokens := s.Tokenize()

	for _, v := range tokens {
		fmt.Println(v)
	}
}
