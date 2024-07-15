package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	jsonparser "github.com/KhushPatibandha/jsonParser"
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
	jsonString := buf.String()

	result, err := jsonparser.ParseIt(jsonString)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
