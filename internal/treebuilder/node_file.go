package treebuilder

import (
	"martinpetr.dev/kina/compiler/internal/diagnostics"
	"martinpetr.dev/kina/compiler/internal/lexer"
	"martinpetr.dev/kina/compiler/internal/performance"
)

func parseFile(scanner *Scanner) FileNode {
	var children = performance.NewFastArray[Node](32)
	start := 0

	for !scanner.IsAtEOF() {
		node, ok := parseToplevelNode(scanner)
		if !ok {
			scanner.Advance() // Skip the current token and continue parsing
			continue
		}

		children.Append(node)
	}

	if children.Len() == 0 {
		return NewFileNode(Span{
			Start: start,
			End:   start,
		}, children.Items())
	}

	return NewFileNode(Span{
		Start: start,
		End:   children.Last().Base().Span.End,
	}, children.Items())
}

func parseToplevelNode(scanner *Scanner) (Node, bool) {
	token := scanner.Peek()

	switch token.Kind {
	case lexer.KwFuncToken:
		res, ok := ParseFunctionDeclaration(scanner)
		return res, ok
	case lexer.KwImportToken:
		res, ok := ParseImport(scanner)
		return res, ok
	default:
		scanner.reporter.Errorf(token.Span.Start, token.Span.End, diagnostics.InvalidSyntaxDiagnosticCode, "Invalid token of type '%s' in top-level scope", token.Kind)
		return BaseNode{}, false
	}
}
