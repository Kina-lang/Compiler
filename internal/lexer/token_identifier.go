package lexer

func lexIdentifierValue(token byte, scanner *Scanner) string {
	restOfIdentifier := scanner.AdvanceWhile(isValidIdentifierChar)
	identifier := string(append([]byte{token}, restOfIdentifier...))

	return identifier
}

func createIdentifierToken(scanner *Scanner, start int, value string) Token {
	return Token{
		Kind:  IdentifierToken,
		Value: value,
		Span: Span{
			Start: start,
			End:   scanner.cursor,
		},
	}
}
