package treebuilder

import (
	"fmt"

	"martinpetr.dev/kina/compiler/internal/lexer"
)

func ParseStatement(scanner *Scanner) (StatementNode, bool) {
	token := scanner.Peek()

	switch token.Kind {
		case lexer.KwReturnToken:
			res, ok := ParseReturnStatement(scanner)
			return res, ok
		default:
			panic(fmt.Sprintf("Invalid token of kind '%s' in basic block", token.Kind))
	}
}

func ParseReturnStatement(scanner *Scanner) (returnStatementNode, bool) {
	returnStatement, ok := scanner.Expect(lexer.KwReturnToken)
	if !ok {
		return returnStatementNode{}, false
	}

	expression, ok := ParseExpression(scanner)
	if !ok {
		return returnStatementNode{}, false
	}

	scanner.Match(lexer.SemicolonToken) // Optional semicolon

	return NewReturnStatementNode(Span{
		Start: returnStatement.Span.Start,
	}, expression), true
}
