package llvmbuilder

import (
	"path/filepath"

	"martinpetr.dev/kina/compiler/internal/sem"
	"martinpetr.dev/kina/compiler/internal/treebuilder"
	"martinpetr.dev/kina/llvm/llvm"
)

type InputFile struct {
	FilePath string
	Tree *treebuilder.Tree
	Table *sem.SymbolTable
}

func BuildLLVM(projectRootPath string, entryFilePath string, files []InputFile) (map[string]string, error) {
	ctx, err := llvm.NewContext("x86_64-unknown-linux-gnu") // TODO: Hardcoded triple for now, change later
	if err != nil {
		return map[string]string{}, err
	}

	var modules map[string]*llvm.Module = make(map[string]*llvm.Module)
	for _, file := range files {
		relativePath, err := filepath.Rel(projectRootPath, file.FilePath)
		if err != nil {
			return map[string]string{}, nil
		}

		builder := llvm.NewBuilder(ctx)
		module := llvm.NewModule(relativePath, ctx)

		isEntrypoint := file.FilePath == entryFilePath
		ProcessFile(&file.Tree.Node, file.Table, module, builder, isEntrypoint)

		modules[file.FilePath] = module
	}

	var irFiles map[string]string = make(map[string]string)
	for _, file := range files {
		module := modules[file.FilePath]
		irFiles[file.FilePath] = module.String()
	}

	return irFiles, nil
}
