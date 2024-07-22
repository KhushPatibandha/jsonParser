package lexer

import (
	"fmt"
	"regexp"
)

type regexHandler func(lex *lexer, regex *regexp.Regexp)

type regexPattern struct {
	regex   *regexp.Regexp
	handler regexHandler
}

type lexer struct {
	patterns []regexPattern
	Tokens   []Token
	source   string
	position int
	line     int
}

func Tokenizer(source string) []Token {
	lexer := createLexer(source)
	for !lexer.atEOF() {
		matched := false
		for _, pattern := range lexer.patterns {
			lineOfCode := pattern.regex.FindStringIndex(lexer.remainder())
			if lineOfCode != nil && lineOfCode[0] == 0 {
				pattern.handler(lexer, pattern.regex)
				matched = true
				break
			}
		}
		if !matched {
			panic(fmt.Sprintf("lexer error: unrecognized token '%v' near --> '%v'", lexer.remainder()[:1], lexer.remainder()))
		}
	}
	lexer.push(NewToken(EOF, "EOF"))
	return lexer.Tokens
}

func (lexer *lexer) advanceN(n int) {
	lexer.position += n
}

func (lexer *lexer) at() byte {
	return lexer.source[lexer.position]
}

func (lexer *lexer) advance() {
	lexer.position += 1
}

func (lexer *lexer) remainder() string {
	return lexer.source[lexer.position:]
}

func (lexer *lexer) push(token Token) {
	lexer.Tokens = append(lexer.Tokens, token)
}

func (lexer *lexer) atEOF() bool {
	return lexer.position >= len(lexer.source)
}

func createLexer(source string) *lexer {
	return &lexer{
		position: 0,
		line:     1,
		source:   source,
		Tokens:   make([]Token, 0),
		patterns: []regexPattern{
			{regexp.MustCompile(`[0-9]+(\.[0-9]+)?`), numberHandler},
			{regexp.MustCompile(`\s+`), skipHandler},
			{regexp.MustCompile(`"[^"]*"`), stringHandler},
			{regexp.MustCompile(`\[`), defaultHandler(LEFT_SQ_BRACKET, "[")},
			{regexp.MustCompile(`\]`), defaultHandler(RIGHT_SQ_BRACKET, "]")},
			{regexp.MustCompile(`\{`), defaultHandler(LEFT_CURLY_BRACKET, "{")},
			{regexp.MustCompile(`\}`), defaultHandler(RIGHT_CURLY_BRACKET, "}")},
			{regexp.MustCompile(`,`), defaultHandler(COMMA, ",")},
			{regexp.MustCompile(`:`), defaultHandler(COLON, ":")},
			{regexp.MustCompile(`-`), defaultHandler(DASH, "-")},
			{regexp.MustCompile(`\+`), defaultHandler(PLUS, "+")},
			{regexp.MustCompile(`[eE]`), defaultHandler(EXPONENT, "e")},
			{regexp.MustCompile(`true`), defaultHandler(BOOLEAN, "true")},
			{regexp.MustCompile(`false`), defaultHandler(BOOLEAN, "false")},
			{regexp.MustCompile(`null`), defaultHandler(NULL, "null")},
		},
	}
}

func defaultHandler(k TokenKind, v string) regexHandler {
	return func(lex *lexer, regex *regexp.Regexp) {
		lex.push(NewToken(k, v))
		lex.advanceN(len(v))
	}
}

func numberHandler(lexer *lexer, regex *regexp.Regexp) {
	lexer.push(NewToken(NUMBER, regex.FindString(lexer.remainder())))
	lexer.advanceN(len(regex.FindString(lexer.remainder())))
}

func stringHandler(lexer *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lexer.remainder())
	stringLiteral := lexer.remainder()[match[0]+1 : match[1]-1]
	lexer.push(NewToken(STRING, stringLiteral))
	lexer.advanceN(len(stringLiteral) + 2)
}

func skipHandler(lexer *lexer, regex *regexp.Regexp) {
	lexer.advanceN(regex.FindStringIndex(lexer.remainder())[1])
}
