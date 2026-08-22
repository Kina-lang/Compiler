package sem

import (
	"martinpetr.dev/kina/compiler/internal/diagnostics"
	"martinpetr.dev/kina/compiler/internal/treebuilder"
)

func generateFunctionSignature(node treebuilder.FunctionDeclarationNode, table *SymbolTable, reporter *diagnostics.Reporter) (FunctionSignature, bool) {
	functionSymbol, found := table.DirectLookup(node.Name)
	if !found {
		return NewFunctionSignature(nil, nil), false
	}

	parameterSignatures, ok := generateFunctionParameterSignatures(node.Parameters, functionSymbol.Table, reporter)
	if !ok {
		return NewFunctionSignature(nil, nil), false
	}

	returnSignature, ok := generateTypeAnnotationSignature(node.ReturnType, table, reporter)
	if !ok {
		return NewFunctionSignature(nil, nil), false
	}

	fnSig := NewFunctionSignature(parameterSignatures, returnSignature)
	functionSymbol.SetSignature(fnSig)

	return fnSig, true
}

func generateFunctionParameterSignatures(parameters []treebuilder.FunctionParameterNode, table *SymbolTable, reporter *diagnostics.Reporter) ([]TypeSignature, bool) {
	sigs := make([]TypeSignature, 0, len(parameters))

	for _, param := range parameters {
		name := param.Name

		symbol, found := table.DirectLookup(name)
		if !found {
			return nil, false
		}

		typeSignature, ok := generateTypeAnnotationSignature(param.Type, table, reporter)
		if !ok {
			return nil, false
		}

		symbol.SetSignature(typeSignature)
		sigs = append(sigs, typeSignature)
	}

	return sigs, true
}
