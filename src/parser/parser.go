package parser

import (
	"fmt"
	"math"
	"strconv"

	"github.com/KhushPatibandha/jsonParser/src/lexer"
)

type Parser struct {
	tokens          []lexer.Token
	currentPosition int
}

func NewParser(t []lexer.Token) *Parser {
	return &Parser{t, 0}
}

func (p *Parser) Parse() (interface{}, error) {

	if p.isAtEnd() {
		return nil, fmt.Errorf("Empty JSON")
	}

	endBracket := p.tokens[len(p.tokens)-2]
	startBracket := p.tokens[0]
	if endBracket.Kind == lexer.COMMA {
		return nil, fmt.Errorf("Closing Parentheses Missing: Expected '}' or ']' at the end, found ','")
	}
	checkForCommaBeforeEndBracket := p.tokens[len(p.tokens)-3]
	if checkForCommaBeforeEndBracket.Kind == lexer.COMMA {
		return nil, fmt.Errorf("Trailing comma before Closing Parentheses: Expected another entry or '}' or ']'")
	}
	if endBracket.Kind != lexer.RIGHT_CURLY_BRACKET && endBracket.Kind != lexer.RIGHT_SQ_BRACKET {
		return nil, fmt.Errorf("Closing Parentheses Missing: Expected '}' or ']'")
	}
	if startBracket.Kind != lexer.LEFT_CURLY_BRACKET && startBracket.Kind != lexer.LEFT_SQ_BRACKET {
		return nil, fmt.Errorf("Opening Parentheses Missing: Expected '{' or '['")
	}

	return p.parseValue()
}

func (p *Parser) parseValue() (interface{}, error) {
	switch p.peek().Kind {
	case lexer.LEFT_CURLY_BRACKET:
		return p.parseObject()
	case lexer.LEFT_SQ_BRACKET:
		return p.parseArray()
	case lexer.STRING:
		return p.advance().Value, nil
	case lexer.NUMBER:
		return p.parseNumber()
	case lexer.BOOLEAN:
		return p.advance().Value, nil
	case lexer.NULL:
		p.advance()
		return nil, nil
	case lexer.DASH:
		return p.parseNumber()
	default:
		return nil, fmt.Errorf("unexpected token: %v", p.peek())
	}
}

func (p *Parser) parseString() (string, error) {
	if p.peek().Kind != lexer.STRING {
		return "", fmt.Errorf("expected string, found %v", p.peek())
	}
	return p.advance().Value, nil
}

func (p *Parser) parseNumber() (string, error) {
	isNegative := false
	if p.peek().Kind == lexer.DASH {
		isNegative = true
		p.advance()
	}

	if p.peek().Kind != lexer.NUMBER {
		return "", fmt.Errorf("expected number after '-', found %v", p.peek())
	}
	number := p.advance().Value

	if p.peek().Kind == lexer.EXPONENT {
		p.advance()
		expNegative := false
		if p.peek().Kind == lexer.DASH {
			expNegative = true
			p.advance()
		} else if p.peek().Kind == lexer.PLUS {
			p.advance()
		}
		if p.peek().Kind != lexer.NUMBER {
			return "", fmt.Errorf("expected number after exponent, found %v", p.peek())
		}
		expNumber, err := strconv.Atoi(p.advance().Value)
		if err != nil {
			return "", err
		}
		if expNegative {
			expNumber = -expNumber
		}

		baseNumber, err := strconv.ParseFloat(number, 64)
		if err != nil {
			return "", err
		}

		result := baseNumber * math.Pow10(expNumber)
		number = strconv.FormatFloat(result, 'f', -1, 64)
	}
	if isNegative {
		number = "-" + number
	}
	return number, nil
}

func (p *Parser) parseObject() (map[string]interface{}, error) {
	p.advance()
	mapObject := make(map[string]interface{})

	for !p.isAtEnd() && p.peek().Kind != lexer.RIGHT_CURLY_BRACKET {
		stringKey, err := p.parseString()
		if err != nil {
			return nil, err
		}

		if !p.match(lexer.COLON) {
			return nil, fmt.Errorf("expected ':', found %v", p.peek())
		}

		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		mapObject[stringKey] = value
		if !p.match(lexer.COMMA) {
			break
		}
	}

	if !p.match(lexer.RIGHT_CURLY_BRACKET) {
		return nil, fmt.Errorf("expected '}' at the end of the object, found %v", p.peek())
	}
	return mapObject, nil
}

func (p *Parser) parseArray() ([]interface{}, error) {
	p.advance()
	arraySlice := make([]interface{}, 0)

	for !p.isAtEnd() && p.peek().Kind != lexer.RIGHT_SQ_BRACKET {
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		arraySlice = append(arraySlice, value)
		if !p.match(lexer.COMMA) {
			break
		}
	}
	if !p.match(lexer.RIGHT_SQ_BRACKET) {
		return nil, fmt.Errorf("expected ']' at the end of the array, found %v", p.peek())
	}
	return arraySlice, nil
}

func (p *Parser) match(kind lexer.TokenKind) bool {
	if p.isAtEnd() || p.peek().Kind != kind {
		return false
	}
	p.advance()
	return true
}

func (p *Parser) isAtEnd() bool {
	return p.currentPosition >= len(p.tokens)
}

func (p *Parser) peek() lexer.Token {
	if p.isAtEnd() {
		return lexer.Token{Kind: lexer.EOF}
	}
	return p.tokens[p.currentPosition]
}

func (p *Parser) advance() lexer.Token {
	if !p.isAtEnd() {
		p.currentPosition++
	}
	return p.tokens[p.currentPosition-1]
}
