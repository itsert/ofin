package merror

import (
	"fmt"
	"os"

	"github.com/itsert/ofin/script/token"
)

func Error(fileName string, line int, start int, message string) {
	fmt.Fprintf(os.Stderr, "%s:%d:%d %s\n", fileName, line, start, message)
	panic(message)
}

func RuntimeError(token token.Token, message string) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("\n[line %d]", token.Line))
	panic(message)
}
