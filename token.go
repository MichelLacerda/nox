package main

import "fmt"

type TokenType int

const (
	TokenType_EOF TokenType = iota

	// Single-character tokens.
	TokenType_LEFT_PAREN
	TokenType_RIGHT_PAREN
	TokenType_LEFT_BRACE
	TokenType_RIGHT_BRACE
	TokenType_LEFT_BRACKET
	TokenType_RIGHT_BRACKET
	TokenType_COMMA
	TokenType_DOT
	TokenType_MINUS
	TokenType_PLUS
	TokenType_SEMICOLON
	TokenType_SLASH
	TokenType_STAR

	// One or two character tokens.
	TokenType_BANG
	TokenType_BANG_EQUAL
	TokenType_EQUAL
	TokenType_EQUAL_EQUAL
	TokenType_GREATER
	TokenType_GREATER_EQUAL
	TokenType_LESS
	TokenType_LESS_EQUAL

	// Literals.
	TokenType_IDENTIFIER
	TokenType_STRING
	TokenType_NUMBER

	// Keywords.
	TokenType_AND
	TokenType_CLASS
	TokenType_ELSE
	TokenType_FALSE
	TokenType_FUN
	TokenType_FOR
	TokenType_IF
	TokenType_NIL
	TokenType_OR
	TokenType_PRINT
	TokenType_RETURN
	TokenType_SUPER
	TokenType_THIS
	TokenType_TRUE
	// TokenType_VAR
	TokenType_LET
	TokenType_WHILE
	TokenType_THEN
	TokenType_END
	TokenType_DO
)

var TokenTypeNames = map[TokenType]string{
	TokenType_EOF:           "EOF",
	TokenType_LEFT_PAREN:    "LEFT_PAREN",
	TokenType_RIGHT_PAREN:   "RIGHT_PAREN",
	TokenType_LEFT_BRACE:    "LEFT_BRACE",
	TokenType_RIGHT_BRACE:   "RIGHT_BRACE",
	TokenType_LEFT_BRACKET:  "LEFT_BRACKET",
	TokenType_RIGHT_BRACKET: "RIGHT_BRACKET",
	TokenType_COMMA:         "COMMA",
	TokenType_DOT:           "DOT",
	TokenType_MINUS:         "MINUS",
	TokenType_PLUS:          "PLUS",
	TokenType_SEMICOLON:     "SEMICOLON",
	TokenType_SLASH:         "SLASH",
	TokenType_STAR:          "STAR",
	TokenType_BANG:          "BANG",
	TokenType_BANG_EQUAL:    "BANG_EQUAL",
	TokenType_EQUAL:         "EQUAL",
	TokenType_EQUAL_EQUAL:   "EQUAL_EQUAL",
	TokenType_GREATER:       "GREATER",
	TokenType_GREATER_EQUAL: "GREATER_EQUAL",
	TokenType_LESS:          "LESS",
	TokenType_LESS_EQUAL:    "LESS_EQUAL",
	TokenType_IDENTIFIER:    "IDENTIFIER",
	TokenType_STRING:        "STRING",
	TokenType_NUMBER:        "NUMBER",
	TokenType_AND:           "AND",
	TokenType_CLASS:         "CLASS",
	TokenType_ELSE:          "ELSE",
	TokenType_FALSE:         "FALSE",
	TokenType_FUN:           "FUN",
	TokenType_FOR:           "FOR",
	TokenType_IF:            "IF",
	TokenType_NIL:           "NIL",
	TokenType_OR:            "OR",
	TokenType_PRINT:         "PRINT",
	TokenType_RETURN:        "RETURN",
	TokenType_SUPER:         "SUPER",
	TokenType_THIS:          "THIS",
	TokenType_TRUE:          "TRUE",
	// TokenType_VAR:           "VAR",
	TokenType_LET:   "LET",
	TokenType_THEN:  "THEN",
	TokenType_END:   "END",
	TokenType_DO:    "DO",
	TokenType_WHILE: "WHILE",
}

func (t TokenType) String() string {
	if name, ok := TokenTypeNames[t]; ok {
		return name
	}
	return fmt.Sprintf("UNKNOWN_TOKEN_TYPE(%d)", t)
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	line    int
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s %v", t.Type, t.Lexeme, t.Literal)
}

func NewToken(tokenType TokenType, lexeme string, literal any, line int) *Token {
	return &Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		line:    line,
	}
}
