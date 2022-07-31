package tools

import (
	"fmt"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func GenerateAST(outputDir string, fileName string, types []string) {
	path := fmt.Sprintf("%s/%s.go", outputDir, strings.ToLower(fileName))
	packages := strings.Split(outputDir, "/")

	packageName := packages[len(packages)-1]
	packageName = strings.Trim(packageName, "/")

	_ = path
	f, err := os.Create(path)
	check(err)
	defer f.Close()

	f.WriteString(fmt.Sprintf("package %s\n\n", packageName))
	f.WriteString("import \"github.com/itsert/ofin/script/token\"\n")
	baseType := fileName
	f.WriteString(fmt.Sprintf("type %s interface {\n", baseType))
	f.WriteString(fmt.Sprintf("\t%s()\n", baseType))
	f.WriteString("\tAccept(visitor Visitor) interface{}\n")
	f.WriteString("}\n\n")

	for _, t := range types {
		structName := strings.TrimSpace(strings.Split(t, ":")[0])
		filedName := strings.TrimSpace(strings.Split(t, ":")[1])
		defineType(f, baseType, structName, filedName)
	}

	defineVisitor(f, baseType, types)
}

func defineType(f *os.File, baseName string, structName string, filedList string) {
	f.WriteString(fmt.Sprintf("type %s struct {\n", structName))
	fileds := strings.Split(filedList, ", ")
	for _, field := range fileds {
		f.WriteString(fmt.Sprintf("\t%s\n", field))
	}
	f.WriteString("}\n")

	f.WriteString(fmt.Sprintf("\nfunc New%s(%s) *%s{\n", structName, filedList, structName))
	f.WriteString(fmt.Sprintf("\treturn &%s{\n", structName))
	for _, field := range fileds {
		allParam := strings.Split(field, " ")
		f.WriteString(fmt.Sprintf("\t\t%s:\t%s,\n", allParam[0], allParam[0]))
	}
	f.WriteString("\t}\n")
	f.WriteString("}\n")

	f.WriteString(fmt.Sprintf("\nfunc (%s *%s) %s() {}\n\n", strings.ToLower(string(structName[0])), structName, baseName))

	f.WriteString(
		fmt.Sprintf(
			"func (%s *%s) Accept(visitor Visitor) interface{} {\n\t return visitor.Visit%s%s(%s)\n}\n\n\n",
			strings.ToLower(string(structName[0])),
			structName,
			structName,
			baseName,
			strings.ToLower(string(structName[0]))))
}

func defineVisitor(f *os.File, baseName string, types []string) {
	f.WriteString("type Visitor interface {\n")
	for _, t := range types {
		structName := strings.TrimSpace(strings.Split(t, ":")[0])
		f.WriteString(fmt.Sprintf("\tVisit%s%s(%s *%s) interface{}\n", structName, baseName, strings.ToLower(baseName), structName))

	}
	f.WriteString("}\n\n")

}
