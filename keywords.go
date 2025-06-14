package main

var Keywords = map[string]TokenType{
	"and":    TokenType_AND,
	"class":  TokenType_CLASS,
	"do":     TokenType_DO,
	"else":   TokenType_ELSE,
	"end":    TokenType_END,
	"false":  TokenType_FALSE,
	"for":    TokenType_FOR,
	"func":   TokenType_FUNC,
	"if":     TokenType_IF,
	"let":    TokenType_LET,
	"nil":    TokenType_NIL,
	"or":     TokenType_OR,
	"print":  TokenType_PRINT,
	"return": TokenType_RETURN,
	"self":   TokenType_SELF,
	"super":  TokenType_SUPER,
	"then":   TokenType_THEN,
	"true":   TokenType_TRUE,
	"while":  TokenType_WHILE,
}
