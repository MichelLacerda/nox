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
	TokenType_COLON
	TokenType_SEMICOLON
	TokenType_SLASH
	TokenType_STAR
	TokenType_PERCENT

	// One or two character tokens.
	TokenType_BANG
	TokenType_BANG_EQUAL
	TokenType_EQUAL
	TokenType_EQUAL_EQUAL
	TokenType_GREATER
	TokenType_GREATER_EQUAL
	TokenType_LESS
	TokenType_LESS_EQUAL
	TokenType_QUESTION
	TokenType_DOUBLE_STAR

	// Literals.
	TokenType_IDENTIFIER
	TokenType_STRING
	TokenType_NUMBER

	// Keywords.
	TokenType_AND
	TokenType_CLASS
	TokenType_ELSE
	TokenType_FALSE
	TokenType_FUNC
	TokenType_FOR
	TokenType_IN
	TokenType_IF
	TokenType_NIL
	TokenType_OR
	TokenType_NOT
	TokenType_PRINT
	TokenType_RETURN
	TokenType_SELF
	TokenType_SUPER
	TokenType_TRUE
	TokenType_LET
	TokenType_WHILE
	TokenType_BREAK
	TokenType_CONTINUE
	TokenType_WITH
	TokenType_AS

	// Unknown or reserved keywords.
	TokenType_Unknown
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
	TokenType_DOUBLE_STAR:   "DOUBLE_STAR",
	TokenType_STAR:          "STAR",
	TokenType_PERCENT:       "PERCENT",
	TokenType_BANG:          "BANG",
	TokenType_BANG_EQUAL:    "BANG_EQUAL",
	TokenType_EQUAL:         "EQUAL",
	TokenType_EQUAL_EQUAL:   "EQUAL_EQUAL",
	TokenType_GREATER:       "GREATER",
	TokenType_GREATER_EQUAL: "GREATER_EQUAL",
	TokenType_LESS:          "LESS",
	TokenType_LESS_EQUAL:    "LESS_EQUAL",
	TokenType_QUESTION:      "QUESTION",
	TokenType_IDENTIFIER:    "IDENTIFIER",
	TokenType_STRING:        "STRING",
	TokenType_NUMBER:        "NUMBER",
	TokenType_AND:           "AND",
	TokenType_CLASS:         "CLASS",
	TokenType_ELSE:          "ELSE",
	TokenType_FALSE:         "FALSE",
	TokenType_FUNC:          "FUNC",
	TokenType_FOR:           "FOR",
	TokenType_IN:            "IN",
	TokenType_IF:            "IF",
	TokenType_NIL:           "NIL",
	TokenType_OR:            "OR",
	TokenType_NOT:           "NOT",
	TokenType_PRINT:         "PRINT",
	TokenType_RETURN:        "RETURN",
	TokenType_SELF:          "SELF",
	TokenType_SUPER:         "SUPER",
	TokenType_TRUE:          "TRUE",
	TokenType_LET:           "LET",
	TokenType_WHILE:         "WHILE",
	TokenType_BREAK:         "BREAK",
	TokenType_CONTINUE:      "CONTINUE",
	TokenType_WITH:          "WITH",
	TokenType_AS:            "AS",
	TokenType_Unknown:       "UNKNOWN",
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
