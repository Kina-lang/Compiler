package lexer

import "martinpetr.dev/kina/compiler/internal/diagnostics"

func ProcessFile(path string, src []byte, reporter *diagnostics.Reporter) []Token {
	scanner := NewScanner(path, src)
	var tokens []Token

	for !scanner.IsAtEOF() {
		current := scanner.Advance()

		switch {
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

					default:
						reporter.Errorf(scanner.cursor-1, scanner.cursor, diagnostics.InvalidTokenDiagnosticCode, "Unexpected character: '%c'", current)
				}

		}
	}

	return tokens
}

func isWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || isLineBreak(b)
}

func isLineBreak(b byte) bool {
	return b == '\n' || b == '\r'
}

func sequenceMatchPredicate(scanner *Scanner, sequence []byte) func(byte) bool {
	var predicate = func(b byte) bool {
		for i, seqByte := range sequence {
			if scanner.PeekAhead(i) != seqByte {
				return false
			}
		}

		return true
	}

	return predicate
}
