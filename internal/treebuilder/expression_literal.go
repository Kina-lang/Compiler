package treebuilder

import (
	"martinpetr.dev/kina/compiler/internal/diagnostics"
	"martinpetr.dev/kina/compiler/internal/lexer"
)

type LiteralType string

const (
	LiteralTypeString LiteralType = "STRING"
	LiteralTypeInt    LiteralType = "INT"
	LiteralTypeFloat  LiteralType = "FLOAT"
	LiteralTypeBool   LiteralType = "BOOL"
	LiteralTypeNull   LiteralType = "NULL"
	LiteralTypeVoid   LiteralType = "VOID"
)

func ParseLiteralExpression(scanner *Scanner) (ExpressionNode, bool) {
	token, ok := scanner.ExpectAny(LiteralTokenTypes...)
	if !ok {
		return BaseExpressionNode{}, false
	}

	var literalType LiteralType
	switch token.Kind {
	case lexer.StringLiteralToken:
		literalType = LiteralTypeString
	case lexer.IntLiteralToken:
		literalType = LiteralTypeInt
	case lexer.FloatLiteralToken:
		literalType = LiteralTypeFloat
	case lexer.KwTrueToken, lexer.KwFalseToken:
		literalType = LiteralTypeBool
	case lexer.KwNullToken:
		literalType = LiteralTypeNull
	case lexer.KwVoidToken:
		literalType = LiteralTypeVoid
	default:
		scanner.reporter.Errorf(token.Span.Start, token.Span.End, diagnostics.InvalidSyntaxDiagnosticCode, "Invalid literal token of type '%s'", token.Kind)
		return BaseExpressionNode{}, false
	}

	return NewLiteralExpressionNode(Span{
		Start: token.Span.Start,
		End:   token.Span.End,
	}, literalType, token.Value), true
}
