package sem

import (
	"martinpetr.dev/kina/compiler/internal/diagnostics"
	"martinpetr.dev/kina/compiler/internal/treebuilder"
)

// Create a symbol for every top-level name, type nil
func collect(file InputFile, parent *SymbolTable) (*SymbolTable, error) {
	var table = NewSymbolTable(parent)

	for _, node := range file.Tree.Node.Children {
		switch node := node.(type) {
			case treebuilder.FunctionDeclarationNode:
				symbol := NewSymbol(node.Name)

				table.Define(symbol)
			case treebuilder.ImportNode:
				for _, member := range node.Members {
					var name = member.Name
					if member.Alias != "" {
						name = member.Alias
					}

					symbol := NewSymbol(name)

					table.Define(symbol)
				}
			default:
				file.Reporter.Errorf(node.Base().Span.Start, node.Base().Span.End, diagnostics.IllegalNodeInTopLevelDiagnosticCode, "Illegal node in top-level: %s", node.Base().Kind)
		}
	}

	return table, nil
}
