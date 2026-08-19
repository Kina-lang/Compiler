package lexer

// Newline is one of: \n, \r, \r\n
func lexNewline(token byte, scanner *Scanner) Token {
	start := scanner.cursor - 1 // Include the initial newline character in the span

	// Handle \n
	if token == '\n' {
		return Token{
			Kind:  NewlineToken,
			Value: "\n",
			Span: Span{
				Start: start,
				End:   scanner.cursor,
			},
		}
	}

	// Check for \n after \r
	next := scanner.Peek()
	if next == '\n' {
		scanner.Advance() // Consume the \n

		return Token{
			Kind:  NewlineToken,
			Value: "\r\n",
			Span: Span{
				Start: start,
				End:   scanner.cursor,
			},
		}
	}

	return Token{
		Kind:  NewlineToken,
		Value: "\r",
		Span: Span{
			Start: start,
			End:   scanner.cursor,
		},
	}
}
