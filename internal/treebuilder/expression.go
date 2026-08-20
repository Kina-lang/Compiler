package treebuilder

import "martinpetr.dev/kina/compiler/internal/lexer"

var LiteralTokenTypes = []lexer.TokenType{
	lexer.StringLiteralToken,
	lexer.IntLiteralToken,
	lexer.FloatLiteralToken,
	lexer.KwTrueToken,
	lexer.KwFalseToken,
	lexer.KwNullToken,
	lexer.KwVoidToken,
}
