package treebuilder

import (
	"slices"

	"martinpetr.dev/kina/compiler/internal/diagnostics"
	"martinpetr.dev/kina/compiler/internal/lexer"
)

type Scanner struct {
	tokens   []lexer.Token
	filePath string
	reporter *diagnostics.Reporter

	cursor int
}

func NewScanner(filePath string, tokens []lexer.Token, reporter *diagnostics.Reporter) *Scanner {
	return &Scanner{
		tokens:   tokens,
		filePath: filePath,
		reporter: reporter,
	}
}

func (s *Scanner) isEOF(pos int) bool {
	return pos >= len(s.tokens) || s.tokens[pos].Kind == lexer.EOFToken
}

func (s *Scanner) IsAtEOF() bool {
	return s.isEOF(s.cursor)
}

func (s *Scanner) createLexerEOFToken() lexer.Token {
	return lexer.Token{
		Kind:  lexer.EOFToken,
		Value: "EOF",
		Span: lexer.Span{
			Start: s.cursor,
			End:   s.cursor,
		},
	}
}

func (s *Scanner) Peek() lexer.Token {
	if s.isEOF(s.cursor) {
		return s.createLexerEOFToken()
	}

	return s.tokens[s.cursor]
}

func (s *Scanner) PeekAhead(offset int) lexer.Token {
	if s.isEOF(s.cursor + offset) {
		return s.createLexerEOFToken()
	}

	return s.tokens[s.cursor+offset]
}

func (s *Scanner) Expect(expectedKind lexer.TokenType) (lexer.Token, bool) {
	token := s.Peek()
	if token.Kind == expectedKind {
		s.cursor++
		return token, true
	}

	s.reporter.Errorf(
		token.Span.Start,
		token.Span.End,
		diagnostics.InvalidSyntaxDiagnosticCode,
		"Expected token of kind '%s', but got '%s'",
		expectedKind,
		token.Kind,
	)

	return token, false
}

func (s *Scanner) ExpectAny(expectedKinds ...lexer.TokenType) (lexer.Token, bool) {
	token := s.Peek()

	if slices.Contains(expectedKinds, token.Kind) {
		s.cursor++
		return token, true
	}

	s.reporter.Errorf(
		token.Span.Start,
		token.Span.End,
		diagnostics.InvalidSyntaxDiagnosticCode,
		"Expected token of kind '%v', but got '%s'",
		expectedKinds,
		token.Kind,
	)

	return token, false
}

func (s *Scanner) Advance() lexer.Token {
	token := s.Peek()
	if !s.isEOF(s.cursor) {
		s.cursor++
	}

	return token
}
