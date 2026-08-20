package treebuilder

import (
	"martinpetr.dev/kina/compiler/internal/diagnostics"
	"martinpetr.dev/kina/compiler/internal/lexer"
)

func parseFile(scanner *Scanner) fileNode {
	var children []Node = make([]Node, 0)
	start := 0

	for !scanner.IsAtEOF() {
		node, ok := parseToplevelNode(scanner)
		if !ok {
			scanner.Advance() // Skip the current token and continue parsing
			continue
		}

		children = append(children, node)
	}

	if len(children) == 0 {
		return NewFileNode(Span{
			Start: start,
			End:   start,
		}, children)
	}

	return NewFileNode(Span{
		Start: start,
		End:   children[len(children)-1].Base().Span.End,
	}, children)
}

func parseToplevelNode(scanner *Scanner) (Node, bool) {
	token := scanner.Peek()

	switch token.Kind {
	case lexer.KwFuncToken:
		res, ok := ParseFunctionDeclaration(scanner)
		return res, ok
	default:
		scanner.reporter.Errorf(token.Span.Start, token.Span.End, diagnostics.InvalidSyntaxDiagnosticCode, "Invalid token of type '%s' in top-level scope", token.Kind)
		return baseNode{}, false
	}
}
