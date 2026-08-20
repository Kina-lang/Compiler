package treebuilder

import (
	"martinpetr.dev/kina/compiler/internal/lexer"
)

func ParseBasicBlock(scanner *Scanner) (basicBlockNode, bool) {
	var statements []StatementNode = make([]StatementNode, 0)

	// If the next token is a closing brace, we have an empty block
	if scanner.Peek().Kind == lexer.BraceCloseToken {
		// Don't consume the closing brace, as it will be consumed by the caller
		return NewBasicBlockNode(Span{
			Start: scanner.Peek().Span.Start,
			End:   scanner.Peek().Span.End,
		}, statements), true
	}

	for !scanner.IsAtEOF() {
		statement, ok := ParseStatement(scanner)
		if !ok {
			return basicBlockNode{}, false
		}

		statements = append(statements, statement)

		// If the next token is a closing brace, we have reached the end of the block
		if scanner.Peek().Kind == lexer.BraceCloseToken {
			// Don't consume the closing brace, as it will be consumed by the caller
			break
		}
	}

	return NewBasicBlockNode(Span{
		Start: statements[0].Base().Span.Start,
		End:   statements[len(statements)-1].Base().Span.End,
	}, statements), true
}
