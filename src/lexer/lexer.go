package lexer

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf16"
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
			{regexp.MustCompile(`"(?:\\.|[^"])*"`), stringHandler},
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
	match := regex.FindString(lexer.remainder())
	stringLiteral := match[1 : len(match)-1]
	stringLiteral = handleEscapeCharacters(stringLiteral)
	lexer.push(NewToken(STRING, stringLiteral))
	lexer.advanceN(len(match))
}

func handleEscapeCharacters(s string) string {
	var result strings.Builder
	escape := false

	for i := 0; i < len(s); i++ {
		if escape {
			switch s[i] {
			case '"':
				result.WriteByte('"')
			case '\\':
				result.WriteByte('\\')
			case '/':
				result.WriteByte('/')
			case 'b':
				result.WriteByte('\b')
			case 'f':
				result.WriteByte('\f')
			case 'n':
				result.WriteByte('\n')
			case 'r':
				result.WriteByte('\r')
			case 't':
				result.WriteByte('\t')
			case 'u':
				if i+4 < len(s) {
					unicode, err := unicodeHandler(s[i+1 : i+5])
					if err != nil {
						panic(fmt.Sprintf("lexer error: invalid unicode escape sequence '%v'", s[i:i+5]))
					}
					if utf16.IsSurrogate(unicode) && i+10 < len(s) && s[i+5:i+7] == "\\u" {
						lowSurrogate, err := unicodeHandler(s[i+7 : i+11])
						if err != nil {
							panic(fmt.Sprintf("lexer error: invalid unicode escape sequence '%v'", s[i+7:i+11]))
						}
						unicode = utf16.DecodeRune(unicode, lowSurrogate)
						i += 6
					}
					result.WriteRune(unicode)
					i += 4
				} else {
					panic(fmt.Sprintf("lexer error: invalid unicode escape sequence '%v'", s[i:]))
				}
			default:
				result.WriteByte('\\')
				result.WriteByte(s[i])
			}
			escape = false
		} else if s[i] == '\\' {
			escape = true
		} else {
			result.WriteByte(s[i])
		}
	}
	return result.String()
}

func unicodeHandler(s string) (rune, error) {
	n, err := strconv.ParseUint(s, 16, 16)
	if err != nil {
		return 0, err
	}
	return rune(n), nil
}

func skipHandler(lexer *lexer, regex *regexp.Regexp) {
	lexer.advanceN(regex.FindStringIndex(lexer.remainder())[1])
}
