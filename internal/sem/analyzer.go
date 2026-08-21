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
		symbolTable, err := collect(file, nil)
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
		newCtx, err := resolveSignatures(fileContexts[file.Path])
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

// Resolve type annotations
func resolveSignatures(ctx *FileContext) (*FileContext, error) {
	return ctx, nil
}

// Validate basic blocks
func validateBodies(ctx *FileContext) (*FileContext, error) {
	return ctx, nil
}
