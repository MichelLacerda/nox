package main

var Keywords = map[string]TokenType{
	"and":    TokenType_AND,
	"class":  TokenType_CLASS,
	"do":     TokenType_DO,
	"else":   TokenType_ELSE,
	"end":    TokenType_END,
	"false":  TokenType_FALSE,
	"fun":    TokenType_FUN,
	"for":    TokenType_FOR,
	"if":     TokenType_IF,
	"nil":    TokenType_NIL,
	"or":     TokenType_OR,
	"print":  TokenType_PRINT,
	"return": TokenType_RETURN,
	"super":  TokenType_SUPER,
	"then":   TokenType_THEN,
	"this":   TokenType_THIS,
	"true":   TokenType_TRUE,
	// "var":    TokenType_VAR,
	"while": TokenType_WHILE,
	"let":   TokenType_LET,
}
