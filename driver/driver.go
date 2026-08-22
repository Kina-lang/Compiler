package driver

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"time"

	"martinpetr.dev/kina/compiler/internal/diagnostics"
	"martinpetr.dev/kina/compiler/internal/lexer"
	"martinpetr.dev/kina/compiler/internal/llvmbuilder"
	"martinpetr.dev/kina/compiler/internal/performance"
	"martinpetr.dev/kina/compiler/internal/sem"
	"martinpetr.dev/kina/compiler/internal/treebuilder"
	"martinpetr.dev/kina/compiler/projectConfig"
)

type Options struct {
	Out string
	EmitTokens bool
	EmitTree bool
	EmitSem bool
}

type parseProjectFilesResult struct{
	FileTrees map[string]treebuilder.Tree
	FileContents map[string][]byte
	FileReporters map[string]*diagnostics.Reporter
}

type parseFileResult struct {
	Imports []string
	Tree treebuilder.Tree
}

func Compile(projectPath string, opts Options) error {
	// Resolve and parse kina.toml project config file
	configPath := path.Join(projectPath, "kina.toml")
	config, err := projectConfig.ParseFile(configPath)
	if err != nil {
		return err
	}

	fmt.Printf("Compiling project %s...\n", projectPath)

	// Get the absolute path to the entry file
	absEntrypointPath, err := filepath.Abs(path.Join(projectPath, config.Project.Entry))
	if err != nil {
		return err
	}

	// Check if the entry file exists
	if _, err := os.Stat(absEntrypointPath); os.IsNotExist(err) {
		return fmt.Errorf("Entry file '%s' does not exist", absEntrypointPath)
	}

	diagnosticsBag := diagnostics.NewBag(true)

	// Parse (lex + ast) all project files
	// Resolved imports are also parsed recursively (results cached)
	parseResult, err := parseProjectFiles(projectPath, absEntrypointPath, diagnosticsBag, opts)
	if err != nil {
		return err
	}

	// Check if there are any diagnostics (errors/warnings) and print them
	err = diagnosticsBag.Err(os.Stderr)
	if err != nil {
		return err
	}

	var semFiles []sem.InputFile
	for filePath, tree := range parseResult.FileTrees {
		semFiles = append(semFiles, sem.InputFile{
			Path: filePath,
			Contents: parseResult.FileContents[filePath],
			Tree: &tree,
			Reporter: parseResult.FileReporters[filePath],
		})
	}

	// Semantically analyze the project files and build the symbol table
	semContexts, err := sem.Process(projectPath, absEntrypointPath, semFiles)
	if err != nil {
		return err
	}

	for filePath, ctx := range semContexts {
		EmitDebugArtifact(opts.EmitSem, "sem", ctx.SymbolTable.String, projectPath, opts.Out, filePath)
	}

	// Check if there are any diagnostics (errors/warnings) and print them
	err = diagnosticsBag.Err(os.Stderr)
	if err != nil {
		return err
	}

	var llvmFiles []llvmbuilder.InputFile
	for filePath, ctx := range semContexts {
		tree := parseResult.FileTrees[filePath]

		llvmFiles = append(llvmFiles, llvmbuilder.InputFile{
			FilePath: filePath,
			Tree: &tree,
			Table: ctx.SymbolTable,
		})
	}

	irFiles, err := llvmbuilder.BuildLLVM(projectPath, absEntrypointPath, llvmFiles)
	if err != nil {
		return err
	}

	var llFiles []string

	for filePath, ir := range irFiles {
		relFilePath, err := filepath.Rel(projectPath, filePath)
		if err != nil {
			return err
		}

		irFilePath := path.Join(opts.Out, relFilePath + ".ll")
		llFiles = append(llFiles, irFilePath)

		irFileParentDir := path.Dir(irFilePath)

		if err := os.MkdirAll(irFileParentDir, os.ModePerm); err != nil {
			return err
		}

		if err := os.WriteFile(irFilePath, []byte(ir), 0644); err != nil {
			return err
		}
	}

	return nil
}

// Parses (lex + ast) all project files starting from the entrypoint file
// Resolved imports are also parsed recursively (results cached)
func parseProjectFiles(projectRootPath string, absEntrypointPath string, diagnosticsBag *diagnostics.Bag, opts Options) (*parseProjectFilesResult, error) {
	var pathsToParse []string = []string{absEntrypointPath}                             // List of files to parse (entrypoint + imports)
	var parsedFileResults map[string]parseFileResult = make(map[string]parseFileResult) // Cache of parsed file results (key is abs file path)

	var fileContents map[string][]byte = make(map[string][]byte)
	var fileReporters map[string]*diagnostics.Reporter = make(map[string]*diagnostics.Reporter)

	var cwd, err = os.Getwd()
	if err != nil {
		return nil, err
	}

	// Loop through all files to parse
	for len(pathsToParse) > 0 {
		relativeFilePath, err := filepath.Rel(cwd, pathsToParse[0])
		if err != nil {
			return nil, err
		}

		// Read the file content
		src, err := os.ReadFile(pathsToParse[0])
		if err != nil {
			return nil, err
		}

		// Create a new diagnostics file for the current file and save it to the diagnostics files map
		diagnosticFile := diagnostics.NewFile(
			relativeFilePath,
			src,
		)
		diagnosticReporter := diagnosticsBag.For(diagnosticFile)

		fileContents[pathsToParse[0]] = src
		fileReporters[pathsToParse[0]] = diagnosticReporter

		// Parse the first file in the list
		res, err := parseFile(projectRootPath, pathsToParse[0], src, diagnosticReporter, opts)
		if err != nil {
			return nil, err
		}

		// Cache the result of the parsed file
		parsedFileResults[pathsToParse[0]] = *res

		// Add all imports of the parsed file to the list of files to parse
		// except for the ones that have already been parsed or are already in the list of files to parse
		for _, imp := range res.Imports {
			_, alreadyParsed := parsedFileResults[imp]
			alreadyInList := slices.Contains(pathsToParse, imp)

			if !alreadyParsed && !alreadyInList {
				pathsToParse = append(pathsToParse, imp)
			}
		}

		// Remove the parsed file from the list of files to parse
		pathsToParse = pathsToParse[1:]
	}

	trees := make(map[string]treebuilder.Tree)
	for filePath, parsedFileResult := range parsedFileResults {
		trees[filePath] = parsedFileResult.Tree
	}

	return &parseProjectFilesResult{
		FileTrees: trees,
		FileContents: fileContents,
		FileReporters: fileReporters,
	}, nil
}

// Parses (lex + ast) a single file and returns the result
func parseFile(projectRootPath string, filePath string, src []byte, reporter *diagnostics.Reporter, opts Options) (*parseFileResult, error) {
	fmt.Printf("Parsing file %s...\n", filePath)

	performance.ReportHeapSize("start")
	startLexer := time.Now()

	// Lex the file
	lexerResult := lexer.ProcessFile(filePath, src, reporter)
	lexerTime := time.Since(startLexer)
	performance.ReportHeapSize("lexer")

	startAsi := time.Now()

	asiResult := lexerResult.InsertSemicolons()
	asiTime := time.Since(startAsi)
	performance.ReportHeapSize("asi")

	essentialTokens := asiResult.RemoveNonEssential()

	EmitDebugArtifact(opts.EmitTokens, "tokens", essentialTokens.String, projectRootPath, opts.Out, filePath)

	startAst := time.Now()
	tree := treebuilder.BuildTree(filePath, essentialTokens.Tokens, reporter)
	astTime := time.Since(startAst)
	performance.ReportHeapSize("ast")

	EmitDebugArtifact(opts.EmitTree, "tree", tree.String, projectRootPath, opts.Out, filePath)

	fmt.Printf("Lexer: %s, ASI: %s, AST: %s\n", lexerTime, asiTime, astTime)

	resolvedImports := tree.ResolveImports(projectRootPath, filePath)

	return &parseFileResult{
		Imports: resolvedImports.GetPaths(),
		Tree: tree,
	}, nil
}

func EmitDebugArtifact(emitFlagValue bool, typeName string, dataDumpFunc func() (string, error), projectRootPath string, outRootPath string, compiledFilePath string) error {
	if !emitFlagValue {
		return nil
	}

	filePathRelativeToProjectRoot, err := filepath.Rel(projectRootPath, compiledFilePath)
	if err != nil {
		return err
	}

	fullArtifactPath := path.Join(outRootPath, "debug", typeName, filePathRelativeToProjectRoot + ".json")
	fileParentDir := path.Dir(fullArtifactPath)

	if err := os.MkdirAll(fileParentDir, os.ModePerm); err != nil {
		return err
	}

	data, err := dataDumpFunc()
	if err != nil {
		return err
	}

	if err := os.WriteFile(fullArtifactPath, []byte(data), 0644); err != nil {
		return err
	}

	return nil
}
