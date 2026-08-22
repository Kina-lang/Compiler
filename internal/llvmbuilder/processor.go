package llvmbuilder

import (
	"martinpetr.dev/kina/compiler/internal/sem"
	"martinpetr.dev/kina/compiler/internal/treebuilder"
	"martinpetr.dev/kina/llvm/llvm"
)

func ProcessFile(node *treebuilder.FileNode, table *sem.SymbolTable, module *llvm.Module, builder *llvm.Builder, isEntrypoint bool) {
	for _, child := range node.Children {
		switch child := child.(type) {
		case treebuilder.FunctionDeclarationNode:
			fn := processFunction(&child, table, module, builder)

			if isEntrypoint && child.Name == "main" {
				module.NewFunctionAlias("kinaprog_main", fn)
			}
		default:
			panic("Unexpected node in top-level: " + child.Base().Kind)
		}
	}
}
