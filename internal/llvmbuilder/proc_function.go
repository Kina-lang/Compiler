package llvmbuilder

import (
	"fmt"

	"martinpetr.dev/kina/compiler/internal/sem"
	"martinpetr.dev/kina/compiler/internal/treebuilder"
	"martinpetr.dev/kina/llvm/llvm"
)

func processFunction(node *treebuilder.FunctionDeclarationNode, table *sem.SymbolTable, module *llvm.Module, builder *llvm.Builder) *llvm.Function {
	var parameters []*llvm.Parameter
	for _, paramNode := range node.Parameters {
		p := processFunctionParameter(&paramNode, table)
		parameters = append(parameters, p)
	}

	fn := module.NewFunction(node.Name, TypeNodeToLLVMType(&node.ReturnType, table), parameters...)
	entryBlock := fn.NewBlock("entry")
	builder.SetInsertionPoint(entryBlock)

	// TODO: Process body

	if !entryBlock.IsTerminated() {
		if fn.ReturnType != llvm.Void {
			panic(fmt.Sprintf("Function %q is missing return statement and it cannot be added automatically (return type is not void)", fn.Name))
		}

		builder.CreateRet(llvm.Void.Const())
	}

	return fn
}

func processFunctionParameter(node *treebuilder.FunctionParameterNode, table *sem.SymbolTable) *llvm.Parameter {
	return llvm.NewParameter(node.Name, TypeNodeToLLVMType(&node.Type, table))
}
