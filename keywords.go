package main

var Keywords = map[string]TokenType{
	"and":    TokenType_AND,
	"class":  TokenType_CLASS,
	"else":   TokenType_ELSE,
	"false":  TokenType_FALSE,
	"for":    TokenType_FOR,
	"func":   TokenType_FUNC,
	"if":     TokenType_IF,
	"in":     TokenType_IN,
	"let":    TokenType_LET,
	"nil":    TokenType_NIL,
	"or":     TokenType_OR,
	"print":  TokenType_PRINT,
	"return": TokenType_RETURN,
	"self":   TokenType_SELF,
	"super":  TokenType_SUPER,
	"true":   TokenType_TRUE,
	"while":  TokenType_WHILE,
}
