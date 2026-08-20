package lexer

var keywordMap = map[string]TokenType{
	"func":   KwFuncToken,
	"return": KwReturnToken,
	"int":    KwIntToken,
	"bool":   KwBoolToken,
	"string": KwStringToken,
	"float":  KwFloatToken,
	"true":   KwTrueToken,
	"false":  KwFalseToken,
	"void":   KwVoidToken,
	"null":   KwNullToken,
	"any":    KwAnyToken,
}

func identifierIsKeyword(identifier string) bool {
	_, isKeyword := keywordMap[identifier]
	return isKeyword
}

func createKeywordToken(scanner *Scanner, start int, value string) Token {
	return Token{
		Kind:  keywordMap[value],
		Value: value,
		Span: Span{
			Start: start,
			End:   scanner.cursor,
		},
	}
}
