package lexer

import "fmt"

type TokenKind int

const (
	STRING TokenKind = iota
	NUMBER
	BOOLEAN
	NULL

	LEFT_SQ_BRACKET
	RIGHT_SQ_BRACKET
	LEFT_CURLY_BRACKET
	RIGHT_CURLY_BRACKET
	COMMA
	COLON
	DASH
	PLUS
	EXPONENT
	EOF
)

type Token struct {
	Kind  TokenKind
	Value string
}

func (t Token) Debug() {
	if t.Kind == STRING || t.Kind == NUMBER || t.Kind == BOOLEAN {
		fmt.Printf("%s(%s)\n", TokenKindString(t.Kind), t.Value)
	} else {
		fmt.Printf("%s()\n", TokenKindString(t.Kind))
	}
}

func NewToken(k TokenKind, v string) Token {
	return Token{k, v}
}

func TokenKindString(kind TokenKind) string {
	switch kind {
	case STRING:
		return "string"
	case NUMBER:
		return "number"
	case BOOLEAN:
		return "boolean"
	case NULL:
		return "null"
	case LEFT_SQ_BRACKET:
		return "left square bracket"
	case RIGHT_SQ_BRACKET:
		return "right square bracket"
	case LEFT_CURLY_BRACKET:
		return "left curly bracket"
	case RIGHT_CURLY_BRACKET:
		return "right curly bracket"
	case COMMA:
		return "comma"
	case COLON:
		return "colon"
	case DASH:
		return "dash"
	case PLUS:
		return "plus"
	case EXPONENT:
		return "exponent"
	case EOF:
		return "EOF"
	default:
		return fmt.Sprintf("unknown(%d)", kind)
	}
}
