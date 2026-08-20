package treebuilder

import "martinpetr.dev/kina/compiler/internal/lexer"

var tokenTypesValidForTypeAnnotation = []lexer.TokenType{
	lexer.IdentifierToken,
	lexer.KwIntToken,
	lexer.KwFloatToken,
	lexer.KwBoolToken,
	lexer.KwStringToken,
	lexer.KwVoidToken,
	lexer.KwNullToken,
	lexer.KwAnyToken,
}

// TODO: Add support for complex types (arrays, generics, function signatures, ...)
func ParseTypeAnnotation(scanner *Scanner) (typeAnnotationNode, bool) {
	typeToken, ok := scanner.ExpectAny(tokenTypesValidForTypeAnnotation...)
	if !ok {
		return typeAnnotationNode{}, false
	}

	return NewTypeAnnotationNode(Span{
		Start: typeToken.Span.Start,
		End:   typeToken.Span.End,
	}, typeToken.Value), true
}
