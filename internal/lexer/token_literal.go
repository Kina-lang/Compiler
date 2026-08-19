package lexer

import "strings"

// TODO: Make this error when there's an underscore after/before the dot
func lexNumberValue(token byte, scanner *Scanner) string {
	var includesDot bool = false

	valueString := string(token) + string(
		scanner.AdvanceWhile(func(b byte) bool {
			// Only allow one dot in the number literal
			if b == '.' {
				if includesDot {
					return false
				}

				includesDot = true
				return true
			}

			// Allow digits and underscores in the number literal
			return isDigit(b) || b == '_'
		}),
	)

	// Get the last character of the value string
	lastChar := valueString[len(valueString)-1]

	// If the last character is not a digit, revert the scanner
	// so that the next token can be lexed correctly
	if !isDigit(lastChar) {
		scanner.Revert()
		valueString = valueString[:len(valueString)-1]
	}

	return valueString
}

func createNumberToken(scanner *Scanner, start int, value string) Token {
	var isFloat bool = strings.Contains(value, ".") // Includes a dot, so it's a float

	kind := IntLiteralToken
	if isFloat {
		kind = FloatLiteralToken
	}

	return Token{
		Kind:  kind,
		Value: value,
		Span: Span{
			Start: start,
			End:   scanner.cursor,
		},
	}
}

func lexStringLiteral(token byte, scanner *Scanner) Token {
	start := scanner.cursor - 1 // Include the initial quote in the span
	wantedQuote := token

	lastCharacterWasEscape := false
	content := scanner.AdvanceWhile(func(b byte) bool {
		if b == wantedQuote && !lastCharacterWasEscape {
			return false // Stop advancing when we find the closing quote that is not escaped
		}

		if b == '\\' && !lastCharacterWasEscape {
			lastCharacterWasEscape = true // Mark that the next character is escaped
		} else {
			lastCharacterWasEscape = false // Reset the escape flag
		}

		return true
	})

	// Consume the closing quote
	scanner.Advance()

	fullContent := string(wantedQuote) + string(content) + string(wantedQuote)

	return Token{
		Kind:  StringLiteralToken,
		Value: fullContent,
		Span: Span{
			Start: start,
			End:   scanner.cursor,
		},
	}
}
