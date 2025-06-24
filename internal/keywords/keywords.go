package keywords

import "github.com/MichelLacerda/nox/internal/token"

var Keywords = map[string]token.TokenType{
	"and":      token.TokenType_AND,
	"class":    token.TokenType_CLASS,
	"else":     token.TokenType_ELSE,
	"false":    token.TokenType_FALSE,
	"for":      token.TokenType_FOR,
	"func":     token.TokenType_FUNC,
	"if":       token.TokenType_IF,
	"in":       token.TokenType_IN,
	"let":      token.TokenType_LET,
	"nil":      token.TokenType_NIL,
	"or":       token.TokenType_OR,
	"not":      token.TokenType_NOT,
	"print":    token.TokenType_PRINT,
	"return":   token.TokenType_RETURN,
	"self":     token.TokenType_SELF,
	"super":    token.TokenType_SUPER,
	"true":     token.TokenType_TRUE,
	"while":    token.TokenType_WHILE,
	"break":    token.TokenType_BREAK,
	"continue": token.TokenType_CONTINUE,
	"with":     token.TokenType_WITH,
	"as":       token.TokenType_AS,
	"import":   token.TokenType_IMPORT,
	"export":   token.TokenType_EXPORT,
}
