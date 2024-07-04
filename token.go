package main

import "fmt"

type TokenType int

const (
	// Single-character tokens.
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR
	MOD

	// One or two character tokens.
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals.
	IDENTIFIER
	STRING
	NUMBER

	// Keywords.
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func (t Token) String() string {
	return fmt.Sprintf("{%v %s %v %d}", t.Type, t.Lexeme, t.Literal, t.Line)
}

func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) Token {
	return Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}
