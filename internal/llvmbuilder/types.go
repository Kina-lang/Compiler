package llvmbuilder

import (
	"fmt"

	"martinpetr.dev/kina/compiler/internal/sem"
	"martinpetr.dev/kina/compiler/internal/treebuilder"
	"martinpetr.dev/kina/llvm/llvm"
)

// TODO: Add more
var primitiveTypeMap = map[string]llvm.Type{
	"int": llvm.Int32,
	"void": llvm.Void,
}

func TypeNodeToLLVMType(node *treebuilder.TypeAnnotationNode, table *sem.SymbolTable) llvm.Type {
	symbol, found := table.Lookup(node.TypeName)
	if !found {
		panic(fmt.Sprintf("Unable to resolve type %q", node.TypeName))
	}

	if symbol.Signature == nil {
		panic(fmt.Sprintf("Symbol %q has no signature", node.TypeName))
	}

	switch sig := (*symbol.Signature).(type) {
	case sem.PrimitiveTypeSignature:
		// TODO: Add support for custom types
		return BuiltinTypeNameToLLVMType(sig.Name)
	default:
		panic(fmt.Sprintf("Symbol %q has unexpected signature type %T", node.TypeName, sig))
	}
}

func BuiltinTypeNameToLLVMType(typeName string) llvm.Type {
	typ, ok := primitiveTypeMap[typeName]
	if ok {
		return typ
	}

	panic(fmt.Sprintf("Unable to resolve type %q", typeName))
}
