package main

import (
	"strconv"
)

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  []Token{},
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, NewToken(EOF, "", nil, s.line))
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '%':
		s.addToken(MOD)
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
	case ' ', '\r', '\t':
		// Ignore whitespace
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			Error(s.line, "Unexpected character.")
		}
	}
}

func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}

func (s *Scanner) addToken(tokenType TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) addTokenWithLiteral(tokenType TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(tokenType, text, literal, s.line))
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		Error(s.line, "Unterminated string.")
		return
	}

	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(STRING, value)
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	value, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		Error(s.line, "Invalid number.")
		return
	}
	s.addTokenWithLiteral(NUMBER, value)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType, ok := keywords[text]
	if !ok {
		tokenType = IDENTIFIER
	}
	s.addToken(tokenType)
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func (s *Scanner) isAlphaNumeric(c byte) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func Error(line int, message string) {
	Report(line, "", message)
}
