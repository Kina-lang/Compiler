package lexer

func lexComment(_ byte, scanner *Scanner) Token {
	next := scanner.Advance()

	switch next {
	case '/':
		return lexLineComment(scanner)
	case '*':
		return lexBlockComment(scanner)
	default:
		panic("Unexpected character after '/': " + string(next))
	}
}

func lexLineComment(scanner *Scanner) Token {
	start := scanner.cursor - 2 // Include the initial '//' in the span

	content := scanner.AdvanceUntil(isLineBreak)

	return Token{
		Kind:  LineCommentToken,
		Value: "//" + string(content),
		Span: Span{
			Start: start,
			End:   scanner.cursor,
		},
	}
}

func lexBlockComment(scanner *Scanner) Token {
	start := scanner.cursor - 2 // Include the initial '/*' in the span

	content := scanner.AdvanceUntil(sequenceMatchPredicate(scanner, []byte{'*', '/'}))

	scanner.Advance() // Consume the '*'
	scanner.Advance() // Consume the '/'

	return Token{
		Kind:  BlockCommentToken,
		Value: "/*" + string(content) + "*/",
		Span: Span{
			Start: start,
			End:   scanner.cursor,
		},
	}
}
