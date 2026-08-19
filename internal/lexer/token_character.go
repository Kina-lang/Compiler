package lexer

var characterTokens = map[byte]TokenType{
	'(': ParenOpenToken,
	')': ParenCloseToken,
	'{': BraceOpenToken,
	'}': BraceCloseToken,
	'[': BracketOpenToken,
	']': BracketCloseToken,
	',': CommaToken,
	'.': DotToken,
	';': SemicolonToken,
	':': ColonToken,
}

func isCharacterToken(token byte) bool {
	_, isCharacter := characterTokens[token]
	return isCharacter
}

func createCharacterToken(scanner *Scanner, start int, value byte) Token {
	return Token{
		Kind:  characterTokens[value],
		Value: string(value),
		Span: Span{
			Start: start,
			End:   scanner.cursor,
		},
	}
}
