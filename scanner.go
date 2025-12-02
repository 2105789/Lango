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
	depth   int // Track brace nesting for string interpolation
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
		depth:   0,
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
		s.depth++
	case '}':
		s.addToken(RIGHT_BRACE)
		if s.depth > 0 {
			s.depth--
			// If we are inside an interpolation, we might need to resume string scanning
			// But since we don't track "inside string" state explicitly here, 
			// we rely on the parser or a more complex state machine. 
			// Actually, for simple interpolation `${...}`, the `}` ends the expression.
			// The tricky part is resuming the string *after* the expression.
			// A simple way is: if we see `}` and we were in an interpolation, check if next char is string-like?
			// No, that's ambiguous.
			// Better approach: When we hit `${`, we emit (LEFT_PAREN). 
			// Then we lex normal tokens.
			// When we hit `}`, if we are in interpolation, we emit (RIGHT_PAREN) and then (PLUS) and then RESUME string scanning.
			// However, the scanner doesn't know if `}` closes a block or an interpolation.
			// We need a stack of states. For now, let's stick to simple tokens and handle the `${` special case.
		}
	case '[':
		s.addToken(LEFT_BRACKET)
	case ']':
		s.addToken(RIGHT_BRACKET)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		if s.match('-') {
			s.addToken(MINUS_MINUS)
		} else if s.match('=') {
			s.addToken(MINUS_EQUAL)
		} else {
			s.addToken(MINUS)
		}
	case '+':
		if s.match('+') {
			s.addToken(PLUS_PLUS)
		} else if s.match('=') {
			s.addToken(PLUS_EQUAL)
		} else {
			s.addToken(PLUS)
		}
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		if s.match('=') {
			s.addToken(STAR_EQUAL)
		} else {
			s.addToken(STAR)
		}
	case '%':
		if s.match('=') {
			s.addToken(MOD_EQUAL)
		} else {
			s.addToken(MOD)
		}
	case ':':
		s.addToken(COLON)
	case '?':
		s.addToken(QUESTION)
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
		} else if s.match('=') {
			s.addToken(SLASH_EQUAL)
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
	// "Start" of the string part
	for !s.isAtEnd() {
		if s.peek() == '"' {
			// End of string
			s.advance()
			value := s.source[s.start+1 : s.current-1]
			processedValue := s.processEscapeSequences(value)
			s.addTokenWithLiteral(STRING, processedValue)
			return
		}

		if s.peek() == '$' && s.peekNext() == '{' {
			// Found interpolation start: "Hello ${
			// 1. Emit the string part before this as a STRING token
			value := s.source[s.start+1 : s.current] // exclude " and $
			processedValue := s.processEscapeSequences(value)
			s.addTokenWithLiteral(STRING, processedValue)

			// 2. Emit a PLUS token to join with the expression
			s.addToken(PLUS)

			// 3. Emit a LEFT_PAREN to group the expression (optional but safer)
			s.addToken(LEFT_PAREN)

			// 4. Consume ${
			s.advance() // $
			s.advance() // {

			// 5. Scan the expression recursively until we hit the matching }
			// We need to count braces to find the matching one.
			// Actually, we can just return to the main scan loop!
			// But we need to know we are in a string to resume string scanning after the }
			// This requires a state stack.
			// Simplified approach:
			// We can't easily jump back to string scanning from the main loop without state.
			// Let's use a recursive helper that scans tokens until it hits the closing brace of the interpolation.
			s.scanInterpolation()
			return
		}

		if s.peek() == '\n' {
			s.line++
		}
		
		if s.peek() == '\\' {
			s.advance() // skip backslash
			if !s.isAtEnd() {
				s.advance() // skip escaped char
			}
		} else {
			s.advance()
		}
	}

	Error(s.line, "Unterminated string.")
}

func (s *Scanner) scanInterpolation() {
	braceCount := 1 // We already consumed the opening {
	
	// We need to scan tokens normally until braceCount is 0
	for braceCount > 0 && !s.isAtEnd() {
		// We use a mini-loop here to pick off tokens
		// This is tricky because scanToken() resets s.start = s.current
		// We need to manage that.
		
		s.start = s.current
		c := s.advance()
		
		// Handle braces to track nesting
		if c == '{' {
			braceCount++
			s.addToken(LEFT_BRACE)
			continue
		}
		if c == '}' {
			braceCount--
			if braceCount == 0 {
				// End of interpolation expression
				// Don't emit the RIGHT_BRACE token for the interpolation closer
				// Instead, emit RIGHT_PAREN and PLUS to prepare for the rest of the string
				s.addToken(RIGHT_PAREN)
				s.addToken(PLUS)
				
				// Resume string scanning
				s.start = s.current // Start of the next string part
				// We need to "pretend" we are starting a string, so we need a dummy quote?
				// No, we can just continue scanning characters.
				// But s.string() expects a starting quote or at least s.start pointing to one?
				// Actually s.string() assumes the opening quote was consumed before calling.
				// So we can just call a helper that continues scanning.
				s.continueString()
				return
			}
			s.addToken(RIGHT_BRACE)
			continue
		}

		// For other characters, we need to delegate to the standard scanning logic
		// But scanToken() is designed to be called from the top level.
		// We can't easily reuse scanToken() because it calls advance() itself.
		// We need to put the character back and call scanToken?
		s.current-- 
		s.scanToken()
	}
	
	if braceCount > 0 {
		Error(s.line, "Unterminated interpolation.")
	}
}

func (s *Scanner) continueString() {
	// This is like string() but doesn't assume a starting quote
	// It scans until " or ${
	
	// We need to capture the content from s.start (which is right after the })
	// But s.string() logic calculates value := s.source[s.start+1 : s.current-1]
	// We need to adjust s.start to "fake" a starting quote so the math works, 
	// OR write a custom loop here.
	
	// Let's write a custom loop that behaves like the middle of string()
	// We will treat the content as a string literal.
	
	contentStart := s.current
	
	for !s.isAtEnd() {
		if s.peek() == '"' {
			// End of string
			value := s.source[contentStart : s.current]
			s.advance() // Consume "
			processedValue := s.processEscapeSequences(value)
			s.addTokenWithLiteral(STRING, processedValue)
			return
		}

		if s.peek() == '$' && s.peekNext() == '{' {
			// Another interpolation
			value := s.source[contentStart : s.current]
			processedValue := s.processEscapeSequences(value)
			s.addTokenWithLiteral(STRING, processedValue)
			
			s.addToken(PLUS)
			s.addToken(LEFT_PAREN)
			
			s.advance() // $
			s.advance() // {
			
			s.scanInterpolation()
			return
		}

		if s.peek() == '\n' {
			s.line++
		}

		if s.peek() == '\\' {
			s.advance()
			if !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.advance()
		}
	}
	
	Error(s.line, "Unterminated string.")
}

// processEscapeSequences converts escape sequences to their actual characters
func (s *Scanner) processEscapeSequences(str string) string {
	result := ""
	i := 0
	for i < len(str) {
		if str[i] == '\\' && i+1 < len(str) {
			switch str[i+1] {
			case 'n':
				result += "\n"
				i += 2
			case 't':
				result += "\t"
				i += 2
			case 'r':
				result += "\r"
				i += 2
			case '\\':
				result += "\\"
				i += 2
			case '"':
				result += "\""
				i += 2
			case '\'':
				result += "'"
				i += 2
			default:
				// Unknown escape sequence, just keep the backslash
				result += string(str[i])
				i++
			}
		} else {
			result += string(str[i])
			i++
		}
	}
	return result
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
