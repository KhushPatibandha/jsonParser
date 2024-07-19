package jsonparser

import (
	"github.com/KhushPatibandha/jsonParser/src/lexer"
	"github.com/KhushPatibandha/jsonParser/src/parser"
)

func ParseIt(jsonString string) (interface{}, error) {
	tokens := lexer.Tokenizer(jsonString)

	// for _, token := range tokens {
	// 	token.Debug()
	// }

	p := parser.NewParser(tokens)
	return p.Parse()
}
