package sem

import (
	"martinpetr.dev/kina/compiler/internal/diagnostics"
	"martinpetr.dev/kina/compiler/internal/treebuilder"
)

// Resolve type annotations
func resolveSignatures(inputFile InputFile, ctx *FileContext) (*FileContext, error) {
	for _, node := range inputFile.Tree.Node.Children {
		switch node := node.(type) {
			case treebuilder.FunctionDeclarationNode:
				symbol, found := ctx.SymbolTable.DirectLookup(node.Name)
				if !found {
					inputFile.Reporter.Errorf(node.Base().Span.Start, node.Base().Span.End, diagnostics.SymbolNotFoundDiagnosticCode, "Symbol not found: %s", node.Name)
					continue
				}

				signature, ok := generateFunctionSignature(node, ctx.SymbolTable, inputFile.Reporter)
				if !ok {
					// TODO: Report error? Should already be reported in the called function
					continue
				}

				symbol.SetSignature(signature)
			case treebuilder.ImportNode:
				// Ignored in this step, as signatures of the imported symbols resolved in
				// the other file's table and we can later just look them up. They are potentially
				// not resolved yet, as the file might get processed later
			default:
				inputFile.Reporter.Errorf(node.Base().Span.Start, node.Base().Span.End, diagnostics.IllegalNodeInTopLevelDiagnosticCode, "Illegal node in top-level: %s", node.Base().Kind)
		}
	}

	return ctx, nil
}
