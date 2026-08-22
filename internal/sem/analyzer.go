package sem

import (
	"martinpetr.dev/kina/compiler/internal/diagnostics"
	"martinpetr.dev/kina/compiler/internal/treebuilder"
)

type InputFile struct {
	Path string
	Contents []byte
	Tree *treebuilder.Tree
	Reporter *diagnostics.Reporter
}

func Process(projectRootPath string, files []InputFile) (map[string]*FileContext, error) {
	fileContexts := make(map[string]*FileContext)

	// Collect
	for _, file := range files {
		symbolTable, err := collect(file, GlobalSymbolTable)
		if err != nil {
			return nil, err
		}

		fileContexts[file.Path] = &FileContext{
			FilePath: file.Path,
			SymbolTable: symbolTable,
		}
	}

	// Resolve signatures
	for _, file := range files {
		newCtx, err := resolveSignatures(file, fileContexts[file.Path])
		if err != nil {
			return nil, err
		}

		fileContexts[file.Path] = newCtx
	}

	// Validate bodies
	for _, file := range files {
		newCtx, err := validateBodies(fileContexts[file.Path])
		if err != nil {
			return nil, err
		}

		fileContexts[file.Path] = newCtx
	}

	return fileContexts, nil
}

// Validate basic blocks
func validateBodies(ctx *FileContext) (*FileContext, error) {
	return ctx, nil
}

func reportSymbolNotFoundError(reporter *diagnostics.Reporter, node treebuilder.Node, name string) {
	reporter.Errorf(node.Base().Span.Start, node.Base().Span.End, diagnostics.SymbolNotFoundDiagnosticCode, "Symbol not found: %s", name)
}

func reportAlreadyDefinedSymbolError(reporter *diagnostics.Reporter, node treebuilder.Node, name string) {
	reporter.Errorf(node.Base().Span.Start, node.Base().Span.End, diagnostics.SymbolAlreadyDefinedDiagnosticCode, "Symbol already defined: %s", name)
}

// TODO: Add support for more complex types (e.g., arrays, generics, etc.)
func generateTypeAnnotationSignature(node treebuilder.TypeAnnotationNode, table *SymbolTable, reporter *diagnostics.Reporter) (TypeSignature, bool) {
	// Check if the type is defined in current or parent symbol tables
	_, found := table.Lookup(node.TypeName)
	if !found {
		reportSymbolNotFoundError(reporter, node, node.TypeName)

		return nil, false
	};

	return NewPrimitiveTypeSignature(node.TypeName), true
}
