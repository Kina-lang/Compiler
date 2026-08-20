package treebuilder

import (
	"fmt"
	"slices"

	"martinpetr.dev/kina/compiler/internal/lexer"
)

type ExpressionPrecedence int

const (
	LowestPrecedence       ExpressionPrecedence = iota
	AssignmentPrecedence                        // =
	LogicalOrPrecedence                         // ||
	LogicalAndPrecedence                        // &&
	EqualityPrecedence                          // == !=
	ComparisonPrecedence                        // < <= > >=
	SumPrecedence                               // + -
	ProductPrecedence                           // * /
	PrefixPrecedence                            // unary + - !
	CallPrecedence                              // myFunction()
	MemberAccessPrecedence                      // myObject.myProperty
)

func ParseExpression(scanner *Scanner) (ExpressionNode, bool) {
	expression, ok := parseExpressionWithPrecedence(scanner, LowestPrecedence)
	return expression, ok
}

func parseExpressionWithPrecedence(scanner *Scanner, precedence ExpressionPrecedence) (ExpressionNode, bool) {
	leftExpression, ok := parsePrefix(scanner)
	if !ok {
		return nil, false
	}

	// Until we reach EOF, semicolon or greater precedence
	for !scanner.IsAtEOF() && scanner.Peek().Kind != lexer.SemicolonToken && precedence < getInfixPrecedence(scanner.Peek().Kind) {
		// Check if the next token is an infix operator
		if !hasInfixPrecedence(scanner.Peek().Kind) {
			return leftExpression, true
		}

		leftExpression, ok = parseInfix(scanner, leftExpression)
		if !ok {
			return nil, false
		}
	}

	return leftExpression, true
}

func parsePrefix(scanner *Scanner) (ExpressionNode, bool) {
	token := scanner.Peek()

	switch {
	case slices.Contains(LiteralTokenTypes, token.Kind):
		return ParseLiteralExpression(scanner)
	default:
		panic(fmt.Sprintf("Invalid token kind of '%s' in prefix of an expression", token.Kind))
	}
}

func parseInfix(scanner *Scanner, left ExpressionNode) (ExpressionNode, bool) {
	token := scanner.Peek()

	switch token.Kind {
	default:
		panic(fmt.Sprintf("Invalid token kind of '%s' in infix of an expression", token.Kind))
	}
}

func getInfixPrecedence(tokenType lexer.TokenType) ExpressionPrecedence {
	switch tokenType {
	default:
		return LowestPrecedence
	}
}

func hasInfixPrecedence(tokenType lexer.TokenType) bool {
	// If precedence is not lowest, then it has infix precedence
	return getInfixPrecedence(tokenType) != LowestPrecedence
}
