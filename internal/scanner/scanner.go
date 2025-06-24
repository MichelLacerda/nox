package scanner

import (
	"fmt"
	"strconv"

	"github.com/MichelLacerda/nox/internal/keywords"
	"github.com/MichelLacerda/nox/internal/token"
)

type Scanner struct {
	source  []rune
	tokens  []*token.Token
	start   int
	current int
	line    int
}

func NewScanner(source []rune) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  []*token.Token{},
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens() ([]*token.Token, error) {
	for !s.IsAtEnd() {
		s.start = s.current
		if err := s.ScanToken(); err != nil {
			return nil, fmt.Errorf("error scanning token at line %d: %w", s.line, err)
		}
	}
	s.tokens = append(s.tokens, token.NewToken(token.TokenType_EOF, "", nil, 0))
	return s.tokens, nil
}

func (s *Scanner) ScanToken() error {
	c := s.Advance()

	switch c {
	case '(':
		s.AddToken(token.TokenType_LEFT_PAREN)
	case ')':
		s.AddToken(token.TokenType_RIGHT_PAREN)
	case '{':
		s.AddToken(token.TokenType_LEFT_BRACE)
	case '}':
		s.AddToken(token.TokenType_RIGHT_BRACE)
	case '[':
		s.AddToken(token.TokenType_LEFT_BRACKET)
	case ']':
		s.AddToken(token.TokenType_RIGHT_BRACKET)
	case ',':
		s.AddToken(token.TokenType_COMMA)
	case '.':
		s.AddToken(token.TokenType_DOT)
	case '-':
		s.AddToken(token.TokenType_MINUS)
	case '+':
		s.AddToken(token.TokenType_PLUS)
	case ';':
		s.AddToken(token.TokenType_SEMICOLON)
	case '*':
		if s.Match('*') {
			s.AddToken(token.TokenType_DOUBLE_STAR)
		} else {
			s.AddToken(token.TokenType_STAR)
		}
	case '%':
		s.AddToken(token.TokenType_PERCENT)
	case ':':
		s.AddToken(token.TokenType_COLON)
	case '?':
		s.AddToken(token.TokenType_QUESTION)
	case '!':
		if s.Match('=') {
			s.AddToken(token.TokenType_BANG_EQUAL)
		} else {
			s.AddToken(token.TokenType_BANG)
		}
	case '=':
		if s.Match('=') {
			s.AddToken(token.TokenType_EQUAL_EQUAL)
		} else {
			s.AddToken(token.TokenType_EQUAL)
		}
	case '>':
		if s.Match('=') {
			s.AddToken(token.TokenType_GREATER_EQUAL)
		} else {
			s.AddToken(token.TokenType_GREATER)
		}
	case '<':
		if s.Match('=') {
			s.AddToken(token.TokenType_LESS_EQUAL)
		} else {
			s.AddToken(token.TokenType_LESS)
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
					return NewScannerError(s.line, "Unterminated block comment.")
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
			s.AddToken(token.TokenType_SLASH)
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
			panic(NewScannerError(s.line, fmt.Sprintf("Unexpected character '%c'", c)))
		}
	}
	return nil
}

func (s *Scanner) ConsumeIdentifier() {
	for s.IsAlphaNumeric(s.Peek()) {
		s.Advance()
	}

	text := string(s.source[s.start:s.current])
	if tokenType, ok := keywords.Keywords[text]; ok {
		s.AddToken(tokenType)
		return
	} else {
		s.AddToken(token.TokenType_IDENTIFIER)
	}
}

func (s *Scanner) ConsumeNumber() error {
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
		return NewScannerError(s.line, fmt.Sprintf("Invalid number format: %s", lexema))
	}

	s.AddTokenWithLiteral(token.TokenType_NUMBER, value)

	return nil
}

func (s *Scanner) ConsumeString() error {
	for !s.IsAtEnd() && s.Peek() != '"' {
		if s.Peek() == '\n' {
			s.line++
		}
		s.Advance()
	}

	if s.IsAtEnd() {
		return NewScannerError(s.line, "Unterminated string.")
	}

	// Consume the closing '"'.
	s.Advance()

	// Extract the string literal.
	literal := string(s.source[s.start+1 : s.current-1])

	s.AddTokenWithLiteral(token.TokenType_STRING, literal)

	return nil
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

func (s *Scanner) AddToken(tokenType token.TokenType) {
	s.AddTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) AddTokenWithLiteral(tokenType token.TokenType, literal any) {
	lexeme := string(s.source[s.start:s.current])
	s.tokens = append(s.tokens, token.NewToken(tokenType, lexeme, literal, s.line))
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
