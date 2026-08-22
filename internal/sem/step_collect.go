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
				fnSymbol := NewSymbol(node.Name, NewSymbolTable(table))

				for _, param := range node.Parameters {
					paramSymbol := NewSymbol(param.Name, nil)

					ok := fnSymbol.Table.Define(paramSymbol)
					if !ok {
						reportAlreadyDefinedSymbolError(file.Reporter, param, param.Name)
					}
				}

				ok := table.Define(fnSymbol)
				if !ok {
					reportAlreadyDefinedSymbolError(file.Reporter, node, node.Name)
				}
			case treebuilder.ImportNode:
				for _, member := range node.Members {
					var name = member.Name
					if member.Alias != "" {
						name = member.Alias
					}

					symbol := NewSymbol(name, nil)

					ok := table.Define(symbol)
					if !ok {
						reportAlreadyDefinedSymbolError(file.Reporter, member, name)
					}
				}
			default:
				file.Reporter.Errorf(node.Base().Span.Start, node.Base().Span.End, diagnostics.IllegalNodeInTopLevelDiagnosticCode, "Illegal node in top-level: %s", node.Base().Kind)
		}
	}

	return table, nil
}
