package treebuilder

import (
	"martinpetr.dev/kina/compiler/internal/lexer"
	"martinpetr.dev/kina/compiler/internal/performance"
)

func ParseBasicBlock(scanner *Scanner) (BasicBlockNode, bool) {
	var statements = performance.NewFastArray[StatementNode](8)

	// If the next token is a closing brace, we have an empty block
	if scanner.Peek().Kind == lexer.BraceCloseToken {
		// Don't consume the closing brace, as it will be consumed by the caller
		return NewBasicBlockNode(Span{
			Start: scanner.Peek().Span.Start,
			End:   scanner.Peek().Span.End,
		}, statements.Items()), true
	}

	for !scanner.IsAtEOF() {
		statement, ok := ParseStatement(scanner)
		if !ok {
			return BasicBlockNode{}, false
		}

		statements.Append(statement)

		// If the next token is a closing brace, we have reached the end of the block
		if scanner.Peek().Kind == lexer.BraceCloseToken {
			// Don't consume the closing brace, as it will be consumed by the caller
			break
		}
	}

	return NewBasicBlockNode(Span{
		Start: statements.First().Base().Span.Start,
		End:   statements.Last().Base().Span.End,
	}, statements.Items()), true
}
