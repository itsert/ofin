package merror

import (
	"fmt"
	"os"
)

func Error(fileName string, line int, start int, message string) {
	fmt.Fprintf(os.Stderr, "%s:%d:%d %s\n", fileName, line, start, message)
	os.Exit(1)
}
