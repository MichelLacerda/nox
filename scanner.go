package main

import (
	"fmt"
	"strconv"
)

type Scanner struct {
	source  []rune
	tokens  []*Token
	start   int
	current int
	line    int
}

func NewScanner(source []rune) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  []*Token{},
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens() []*Token {
	for !s.IsAtEnd() {
		s.start = s.current
		s.ScanToken()
	}
	s.tokens = append(s.tokens, NewToken(TokenType_EOF, "", nil, 0))
	return s.tokens
}

func (s *Scanner) ScanToken() {
	c := s.Advance()

	switch c {
	case '(':
		s.AddToken(TokenType_LEFT_PAREN)
	case ')':
		s.AddToken(TokenType_RIGHT_PAREN)
	case '{':
		s.AddToken(TokenType_LEFT_BRACE)
	case '}':
		s.AddToken(TokenType_RIGHT_BRACE)
	case '[':
		s.AddToken(TokenType_LEFT_BRACKET)
	case ']':
		s.AddToken(TokenType_RIGHT_BRACKET)
	case ',':
		s.AddToken(TokenType_COMMA)
	case '.':
		s.AddToken(TokenType_DOT)
	case '-':
		s.AddToken(TokenType_MINUS)
	case '+':
		s.AddToken(TokenType_PLUS)
	case ';':
		s.AddToken(TokenType_SEMICOLON)
	case '*':
		s.AddToken(TokenType_STAR)
	case ':':
		s.AddToken(TokenType_COLON)
	case '?':
		s.AddToken(TokenType_QUESTION)
	case '!':
		if s.Match('=') {
			s.AddToken(TokenType_BANG_EQUAL)
		} else {
			s.AddToken(TokenType_BANG)
		}
	case '=':
		if s.Match('=') {
			s.AddToken(TokenType_EQUAL_EQUAL)
		} else {
			s.AddToken(TokenType_EQUAL)
		}
	case '>':
		if s.Match('=') {
			s.AddToken(TokenType_GREATER_EQUAL)
		} else {
			s.AddToken(TokenType_GREATER)
		}
	case '<':
		if s.Match('=') {
			s.AddToken(TokenType_LESS_EQUAL)
		} else {
			s.AddToken(TokenType_LESS)
		}
	case '/':
		if s.Match('/') {
			// A comment goes until the end of the line.
			for s.Peek() != '\n' && !s.IsAtEnd() {
				s.Advance()
			}
		} else if s.Match('*') {
			// A block comment starts with /* and ends with */.
			for {
				if s.IsAtEnd() {
					runtime.ErrorAt(s.line, "Unterminated block comment.")
					break
				}

				if s.Peek() == '*' {
					s.Advance() // Consume the '*' character.
					if s.Match('/') {
						// End of block comment.
						break
					}
				} else {
					if s.Peek() == '\n' {
						s.line++
					}
					s.Advance()
				}
			}
		} else {
			s.AddToken(TokenType_SLASH)
		}
	case ' ', '\r', '\t':
		// Ignore whitespace.
	case '\n':
		s.line++
	case '"':
		s.ConsumeString()
	default:
		if s.IsDigit(c) {
			s.ConsumeNumber()
		} else if s.IsAlpha(c) {
			s.ConsumeIdentifier()
		} else {
			runtime.ErrorAt(s.line, fmt.Sprintf("Unexpected character '%c'", c))
		}
	}
}

func (s *Scanner) ConsumeIdentifier() {
	for s.IsAlphaNumeric(s.Peek()) {
		s.Advance()
	}

	text := string(s.source[s.start:s.current])
	if tokenType, ok := Keywords[text]; ok {
		s.AddToken(tokenType)
		return
	} else {
		s.AddToken(TokenType_IDENTIFIER)
	}
}

func (s *Scanner) ConsumeNumber() {
	for s.IsDigit(s.Peek()) {
		s.Advance()
	}

	// Look for a fractional part.
	if s.Peek() == '.' && s.IsDigit(s.PeekNext()) {
		s.Advance() // Consume the '.'
		for s.IsDigit(s.Peek()) {
			s.Advance()
		}
	}

	lexema := string(s.source[s.start:s.current])
	value, err := strconv.ParseFloat(lexema, 64)

	if err != nil {
		runtime.hadError = true
		runtime.ErrorAt(s.line, fmt.Sprintf("Invalid number format: %s", lexema))
	}

	s.AddTokenWithLiteral(TokenType_NUMBER, value)
}

func (s *Scanner) ConsumeString() {
	for !s.IsAtEnd() && s.Peek() != '"' {
		if s.Peek() == '\n' {
			s.line++
		}
		s.Advance()
	}

	if s.IsAtEnd() {
		runtime.ErrorAt(s.line, "Unterminated string.")
		return
	}

	// Consume the closing '"'.
	s.Advance()

	// Extract the string literal.
	literal := string(s.source[s.start+1 : s.current-1])

	s.AddTokenWithLiteral(TokenType_STRING, literal)
}

func (s *Scanner) Match(expected rune) bool {
	if s.IsAtEnd() || s.source[s.current] != expected {
		return false
	}
	s.current = s.current + 1
	return true
}

func (s *Scanner) Peek() rune {
	if s.IsAtEnd() || s.current > len(s.source) {
		return 0 // Return 0 to indicate end of input.
	}
	return s.source[s.current]
}

func (s *Scanner) PeekNext() rune {
	if s.current+1 >= len(s.source) {
		return 0 // Return 0 to indicate end of input.
	}
	return s.source[s.current+1]
}

func (s *Scanner) AddToken(tokenType TokenType) {
	s.AddTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) AddTokenWithLiteral(tokenType TokenType, literal any) {
	lexeme := string(s.source[s.start:s.current])
	s.tokens = append(s.tokens, NewToken(tokenType, lexeme, literal, s.line))
}

func (s *Scanner) Advance() rune {
	if s.IsAtEnd() {
		return 0
	}
	s.current = s.current + 1
	return s.source[s.current-1]
}

func (s *Scanner) IsDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) IsAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (s *Scanner) IsAlphaNumeric(c rune) bool {
	return s.IsAlpha(c) || s.IsDigit(c)
}

func (s *Scanner) IsAtEnd() bool {
	return s.current >= len(s.source)
}
