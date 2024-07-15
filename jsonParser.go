package jsonparser

import (
	"github.com/KhushPatibandha/jsonParser/src/lexer"
	"github.com/KhushPatibandha/jsonParser/src/parser"
)

func ParseIt(jsonString string) (interface{}, error) {
	tokens := lexer.Tokenizer(jsonString)
	p := parser.NewParser(tokens)
	return p.Parse()
}
