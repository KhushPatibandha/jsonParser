package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/KhushPatibandha/jsonParser/src/lexer"
	"github.com/KhushPatibandha/jsonParser/src/parser"
)

func main() {
	var buf bytes.Buffer
	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from console")
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		buf.WriteString(line)
	}
	lines := buf.String()
	// fmt.Println(lines)

	tokens := lexer.Tokenizer(lines)
	// for _, token := range tokens {
	// 	token.Debug()
	// }

	parser := parser.NewParser(tokens)
	resultMap, err := parser.Parse()
	if err != nil {
		panic(err)
	}
	fmt.Println(resultMap)
}
