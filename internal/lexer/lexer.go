package lexer

import (
	"encoding/json"

	"martinpetr.dev/kina/compiler/internal/diagnostics"
)

type LexerResult struct {
	Tokens []Token
}

func ProcessFile(path string, src []byte, reporter *diagnostics.Reporter) LexerResult {
	scanner := NewScanner(path, src)
	var tokens []Token

	for !scanner.IsAtEOF() {
		current := scanner.Advance()

		switch {
		case isDigit(current):
			start := scanner.cursor - 1
			number := lexNumberValue(current, scanner)
			tokens = append(tokens, createNumberToken(scanner, start, number))

		case isValidIdentifierStartChar(current):
			start := scanner.cursor - 1
			identifier := lexIdentifierValue(current, scanner)

			if identifierIsKeyword(identifier) {
				tokens = append(tokens, createKeywordToken(scanner, start, identifier))
			} else {
				tokens = append(tokens, createIdentifierToken(scanner, start, identifier))
			}

		// This case needs to be after any other case lexing anything that can start with a character
		// parsed by this case.
		case isCharacterToken(current):
			tokens = append(tokens, createCharacterToken(scanner, scanner.cursor-1, current))

		case isLineBreak(current):
			tokens = append(tokens, lexNewline(current, scanner))

		case isWhitespace(current):
			// Skip whitespace

		default:
			switch current {
			case '/':
				next := scanner.Peek()
				if next == '/' || next == '*' {
					tokens = append(tokens, lexComment(current, scanner))
					break
				}

			case '"', '\'':
				tokens = append(tokens, lexStringLiteral(current, scanner))

			default:
				reporter.Errorf(scanner.cursor-1, scanner.cursor, diagnostics.InvalidTokenDiagnosticCode, "Unexpected character: '%c'", current)
			}

		}
	}

	tokens = append(tokens, Token{
		Kind:  EOFToken,
		Value: "EOF",
		Span: Span{
			Start: scanner.cursor,
			End:   scanner.cursor,
		},
	})

	return LexerResult{
		Tokens: tokens,
	}
}

func (r *LexerResult) String() (string, error) {
	json, err := json.MarshalIndent(r, "", "  ")
	return string(json), err
}

func (r *LexerResult) RemoveNonEssential() LexerResult {
	resultTokens := []Token{}

	for _, token := range r.Tokens {
		switch token.Kind {
		case LineCommentToken, BlockCommentToken, NewlineToken:
			// Skip newlines and comments
		default:
			resultTokens = append(resultTokens, token)
		}
	}

	return LexerResult{
		Tokens: resultTokens,
	}
}

func isAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func isUnderscore(b byte) bool {
	return b == '_'
}

func isValidIdentifierChar(b byte) bool {
	return isAlpha(b) || isDigit(b) || isUnderscore(b)
}

func isValidIdentifierStartChar(b byte) bool {
	return isAlpha(b) || isUnderscore(b)
}

func isWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || isLineBreak(b)
}

func isLineBreak(b byte) bool {
	return b == '\n' || b == '\r'
}

func sequenceMatchPredicate(scanner *Scanner, sequence []byte) func(byte) bool {
	return func(b byte) bool {
		for i, seqByte := range sequence {
			if scanner.PeekAhead(i) != seqByte {
				return false
			}
		}

		return true
	}
}

func isCharacterPredicate(w byte) func(byte) bool {
	return func(b byte) bool {
		return b == w
	}
}

func orPredicate(predicates ...func(byte) bool) func(byte) bool {
	return func(b byte) bool {
		for _, predicate := range predicates {
			if predicate(b) {
				return true
			}
		}

		return false
	}
}
