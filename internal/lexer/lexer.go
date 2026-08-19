package lexer

import (
	"martinpetr.dev/kina/compiler/internal/diagnostics"
)

func ProcessFile(path string, src []byte, reporter *diagnostics.Reporter) []Token {
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

				if (identifierIsKeyword(identifier)) {
					tokens = append(tokens, createKeywordToken(scanner, start, identifier))
				} else {
					tokens = append(tokens, createIdentifierToken(scanner, start, identifier))
				}

			// This case needs to be after any other case lexing anything that can start with a character
			// parsed by this case.
			case isCharacterToken(current):
				tokens = append(tokens, createCharacterToken(scanner, scanner.cursor-1, current))

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

	return tokens
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
